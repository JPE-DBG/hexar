import { Connection } from './net/connection';
import { applySnapshot, GameState, SnapshotMsg, HexDTO } from './state/state';
import { Renderer } from './render/renderer';
import { setupInput } from './input/input';
import { updateHUD, calcMaintenance } from './ui/hud';
import { BuildMenu } from './ui/buildmenu';
import { TechTreePanel } from './ui/techtree';
import { AutoDropPanel } from './ui/autodrop';
import { neighbors } from './hexmath';
import {
  BUILDING_GOLD, BUILDING_POWER, BUILDING_RESEARCH,
  BASE_INCOME_PER_SEC, GOLD_PER_LEVEL, GOLD_BONUS_MULTIPLIER,
  CAPITAL_POWER, RESEARCH_PER_LEVEL,
  COUNTER_SPEND_COST, COUNTER_SPEND_CAP,
  GOLD_BUILD_COST, POWER_BUILD_COST, RESEARCH_BUILD_COST,
} from './constants';

const canvas = document.getElementById('game') as HTMLCanvasElement;
const hud = document.getElementById('hud')!;
const renderer = new Renderer(canvas);

let state: GameState | null = null;
let myPlayerId = 0;
let selectedHex: HexDTO | null = null;
let dropMap = new Set<string>();

const buildMenu = new BuildMenu(document.body, {
  onUpgrade: (building?) => {
    if (!selectedHex) return;
    const msg: Record<string, unknown> = { type: 'action', action: 'upgrade', q: selectedHex.q, r: selectedHex.r };
    if (building) msg.building = building;
    connection.send(msg);
  },
  onDemolish: () => {
    if (!selectedHex) return;
    connection.send({ type: 'action', action: 'demolish', q: selectedHex.q, r: selectedHex.r });
  },
  onAttack: () => {
    if (!selectedHex) return;
    connection.send({ type: 'action', action: 'attack', q: selectedHex.q, r: selectedHex.r });
  },
  onDropHex: () => {
    if (!selectedHex) return;
    connection.send({ type: 'action', action: 'drop-hex', q: selectedHex.q, r: selectedHex.r });
  },
  onFortify: () => {
    if (!selectedHex) return;
    connection.send({ type: 'action', action: 'fortify', q: selectedHex.q, r: selectedHex.r });
  },
});

const techTreePanel = new TechTreePanel(document.body, (techId) => {
  connection.send({ type: 'action', action: 'unlock-tech', techId });
});

const autoDropPanel = new AutoDropPanel(document.body);

window.addEventListener('keydown', (e) => {
  if (e.key === 't' || e.key === 'T') {
    techTreePanel.toggle();
    if (state && myPlayerId > 0) {
      const player = state.players.get(String(myPlayerId));
      if (player) techTreePanel.update(player);
    }
  }
});

let victoryOverlay: HTMLElement | null = null;

function showVictory(isWinner: boolean) {
  if (victoryOverlay) return;
  victoryOverlay = document.createElement('div');
  victoryOverlay.style.cssText = `
    position:fixed;top:0;left:0;width:100%;height:100%;
    display:flex;flex-direction:column;align-items:center;justify-content:center;
    background:rgba(0,0,0,0.75);z-index:200;font-family:monospace;
  `;
  const msg = isWinner ? 'You win!' : 'You lose!';
  const color = isWinner ? '#4ecdc4' : '#ff6b6b';
  victoryOverlay.innerHTML = `<div style="font-size:48px;font-weight:bold;color:${color}">${msg}</div>
    <div style="color:#aaa;margin-top:8px;font-size:16px">Capital captured</div>`;
  document.body.appendChild(victoryOverlay);
}

function onSnapshot(msg: SnapshotMsg) {
  state = applySnapshot(msg);

  if (state.over && myPlayerId > 0) {
    showVictory(state.winner === myPlayerId);
  }

  if (myPlayerId > 0) {
    let hexCount = 0;
    let income = 0;
    let tpRate = 0;
    for (const [, hex] of state.hexes) {
      if (hex.owner === myPlayerId) {
        hexCount++;
        income += hexIncome(hex);
        if (hex.building === BUILDING_RESEARCH) {
          tpRate += RESEARCH_PER_LEVEL * hex.level;
        }
      }
    }
    const maintenance = calcMaintenance(hexCount);
    const player = state.players.get(String(myPlayerId));
    const gold = player?.gold ?? 0;
    const tp = player?.tp ?? 0;
    updateHUD(hud, myPlayerId, gold, hexCount, income, maintenance, tp, tpRate);

    if (player) {
      techTreePanel.update(player);
    }
    autoDropPanel.update(player ?? null);

    // Compute drop map: all hexes at minimum (income, totalInvested) score
    dropMap = new Set<string>();
    if (player?.autoDropActive) {
      const droppable = [...state.hexes.values()]
        .filter(h => h.owner === myPlayerId && !h.capital &&
                     !state!.battles.some(b => b.dq === h.q && b.dr === h.r));
      if (droppable.length > 0) {
        const minIncome = Math.min(...droppable.map(h => hexIncome(h)));
        const incomeGroup = droppable.filter(h => hexIncome(h) === minIncome);
        const minInvested = Math.min(...incomeGroup.map(h => totalInvested(h)));
        const candidates = incomeGroup.filter(h => totalInvested(h) === minInvested);
        for (const h of candidates) dropMap.add(`${h.q},${h.r}`);
      }
    }
    renderer.setDropMap(dropMap);

    if (selectedHex) {
      const key = `${selectedHex.q},${selectedHex.r}`;
      const current = state.hexes.get(key);
      if (current) {
        selectedHex = current;
        const isOwn = current.owner === myPlayerId;
        const isEnemy = current.owner !== 0 && current.owner !== myPlayerId;
        const atkPwr = bestAdjacentPower(state, myPlayerId, current);
        const battle = state.battles.find(b => b.dq === current.q && b.dr === current.r) ?? null;
        buildMenu.updateWithActions(current, gold, isOwn, isEnemy, atkPwr, battle, player ?? null);
      }
    }
  }

  renderer.setState(state);
}

function hexIncome(hex: HexDTO): number {
  if (hex.building === BUILDING_GOLD) {
    return (BASE_INCOME_PER_SEC + GOLD_PER_LEVEL * hex.level) * GOLD_BONUS_MULTIPLIER;
  }
  return BASE_INCOME_PER_SEC;
}

function totalInvested(hex: HexDTO): number {
  if (hex.building === 0 || hex.level === 0) return 0;
  const base = hex.building === BUILDING_RESEARCH ? RESEARCH_BUILD_COST
             : hex.building === BUILDING_GOLD     ? GOLD_BUILD_COST
             :                                      POWER_BUILD_COST;
  return base * (Math.pow(2, hex.level) - 1);
}

function bestAdjacentPower(gs: GameState, playerId: number, target: HexDTO): number {
  let best = 0;
  for (const n of neighbors({ q: target.q, r: target.r })) {
    const key = `${n.q},${n.r}`;
    const hs = gs.hexes.get(key);
    if (!hs || hs.owner !== playerId) continue;
    let p = 0;
    if (hs.capital) p = CAPITAL_POWER;
    if (hs.building === BUILDING_POWER) p += hs.level;
    if (p > best) best = p;
  }
  return best;
}

const wsUrl = `ws://${window.location.host}/ws`;
const connection = new Connection(wsUrl, {
  onSnapshot,
  onWelcome: (msg) => {
    myPlayerId = msg.playerId;
    console.log(`assigned player ${myPlayerId}`);
  },
});

setupInput(
  canvas,
  () => renderer.getOffset(),
  (q, r) => {
    if (!state) return;
    const key = `${q},${r}`;
    const hex = state.hexes.get(key);

    if (!hex) {
      selectedHex = null;
      renderer.setSelected(null);
      buildMenu.hide();
      return;
    }

    if (hex.owner === 0) {
      connection.send({ type: 'action', action: 'claim', q, r });
      selectedHex = null;
      renderer.setSelected(null);
      buildMenu.hide();
      return;
    }

    // Priority 1: crisis drop (red pulsing hex) — emergency recovery
    if (dropMap.has(key) && hex.owner === myPlayerId) {
      connection.send({ type: 'action', action: 'drop-hex', q, r });
      return;
    }

    // Priority 2: counter-spend (amber pulsing, still actionable)
    const battle = state.battles.find(b => b.dq === q && b.dr === r) ?? null;
    if (battle && hex.owner === myPlayerId) {
      const cap = Math.min(COUNTER_SPEND_CAP, Math.floor(battle.timeLeft));
      const player = state.players.get(String(myPlayerId));
      if ((player?.gold ?? 0) >= COUNTER_SPEND_COST && battle.counterBoost < cap) {
        connection.send({ type: 'action', action: 'counter-spend', q, r });
        return;
      }
    }

    // Priority 3: normal selection → show build menu
    selectedHex = hex;
    renderer.setSelected({ q, r });

    const player = state.players.get(String(myPlayerId));
    const gold = player?.gold ?? 0;
    const isOwn = hex.owner === myPlayerId;
    const isEnemy = hex.owner !== 0 && hex.owner !== myPlayerId;
    const atkPwr = bestAdjacentPower(state, myPlayerId, hex);
    buildMenu.updateWithActions(hex, gold, isOwn, isEnemy, atkPwr, battle, player ?? null);
  }
);
