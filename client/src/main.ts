import { Connection } from './net/connection';
import { applySnapshot, GameState, SnapshotMsg, HexDTO } from './state/state';
import { Renderer } from './render/renderer';
import { setupInput } from './input/input';
import { updateHUD, calcMaintenance } from './ui/hud';
import { BuildMenu } from './ui/buildmenu';
import { TechTreePanel } from './ui/techtree';
import { AutoDropPanel } from './ui/autodrop';
import { neighbors } from './hexmath';
import { BUILDING_GOLD, BUILDING_POWER, BUILDING_RESEARCH, BASE_INCOME_PER_SEC, GOLD_PER_LEVEL, GOLD_BONUS_MULTIPLIER, CAPITAL_POWER, RESEARCH_PER_LEVEL } from './constants';

const canvas = document.getElementById('game') as HTMLCanvasElement;
const hud = document.getElementById('hud')!;
const renderer = new Renderer(canvas);

let state: GameState | null = null;
let myPlayerId = 0;
let selectedHex: HexDTO | null = null;

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
  onCounterSpend: () => {
    if (!selectedHex) return;
    connection.send({ type: 'action', action: 'counter-spend', q: selectedHex.q, r: selectedHex.r });
  },
});

const techTreePanel = new TechTreePanel(document.body, (techId) => {
  connection.send({ type: 'action', action: 'unlock-tech', techId });
});

const autoDropPanel = new AutoDropPanel(document.body, (q, r) => {
  connection.send({ type: 'action', action: 'drop-hex', q, r });
});

window.addEventListener('keydown', (e) => {
  if (e.key === 't' || e.key === 'T') {
    techTreePanel.toggle();
    if (state && myPlayerId > 0) {
      const player = state.players.get(String(myPlayerId));
      if (player) techTreePanel.update(player);
    }
  }
});

function onSnapshot(msg: SnapshotMsg) {
  state = applySnapshot(msg);
  renderer.render(state);

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
    autoDropPanel.update(player ?? null, state.hexes, state.battles);

    if (selectedHex) {
      const key = `${selectedHex.q},${selectedHex.r}`;
      const current = state.hexes.get(key);
      if (current) {
        selectedHex = current;
        const isOwn = current.owner === myPlayerId;
        const isEnemy = current.owner !== 0 && current.owner !== myPlayerId;
        const atkPwr = bestAdjacentPower(state, myPlayerId, current);
        const battle = state.battles.find(b => b.dq === current.q && b.dr === current.r) ?? null;
        buildMenu.updateWithActions(current, gold, isOwn, isEnemy, atkPwr, battle);
      }
    }
  }
}

function hexIncome(hex: HexDTO): number {
  if (hex.building === BUILDING_GOLD) {
    return (BASE_INCOME_PER_SEC + GOLD_PER_LEVEL * hex.level) * GOLD_BONUS_MULTIPLIER;
  }
  return BASE_INCOME_PER_SEC;
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

    selectedHex = hex;
    renderer.setSelected({ q, r });

    const player = state.players.get(String(myPlayerId));
    const gold = player?.gold ?? 0;
    const isOwn = hex.owner === myPlayerId;
    const isEnemy = hex.owner !== 0 && hex.owner !== myPlayerId;
    const atkPwr = bestAdjacentPower(state, myPlayerId, hex);
    const battle = state.battles.find(b => b.dq === hex.q && b.dr === hex.r) ?? null;
    buildMenu.updateWithActions(hex, gold, isOwn, isEnemy, atkPwr, battle);
  }
);
