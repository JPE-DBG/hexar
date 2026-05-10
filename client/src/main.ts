import { Connection } from './net/connection';
import { applySnapshot, applyDelta, GameState, SnapshotMsg, DeltaMsg, HexDTO, PlayerDTO } from './state/state';
import { Renderer } from './render/renderer';
import { setupInput } from './input/input';
import { updateHUD, calcMaintenance } from './ui/hud';
import { BuildMenu } from './ui/buildmenu';
import { TechTreePanel } from './ui/techtree';
import { AutoDropPanel } from './ui/autodrop';
import { LobbyUI } from './ui/lobby';
import { neighbors } from './hexmath';
import {
  BUILDING_GOLD, BUILDING_POWER, BUILDING_RESEARCH,
  BASE_INCOME_PER_SEC, GOLD_PER_LEVEL, GOLD_BONUS_MULTIPLIER,
  CAPITAL_POWER, RESEARCH_PER_LEVEL,
  COUNTER_SPEND_COST, COUNTER_SPEND_CAP,
  GOLD_BUILD_COST, POWER_BUILD_COST, RESEARCH_BUILD_COST,
  PROSPERITY_BONUS, COMPOUND_GROWTH_MULTIPLIER,
  TECH_PROSPERITY, TECH_COMPOUND_GROWTH, TECH_IRON_GRIP,
} from './constants';

const canvas = document.getElementById('game') as HTMLCanvasElement;
const hud = document.getElementById('hud')!;
const renderer = new Renderer(canvas);

let state: GameState | null = null;
let myPlayerId = 0;
let selectedHex: HexDTO | null = null;
let dropMap = new Set<string>();
let connection: Connection | null = null;

const buildMenu = new BuildMenu(document.body, {
  onUpgrade: (building?) => {
    if (!selectedHex || !connection) return;
    const msg: Record<string, unknown> = { type: 'action', action: 'upgrade', q: selectedHex.q, r: selectedHex.r };
    if (building) msg.building = building;
    connection.send(msg);
  },
  onDemolish: () => {
    if (!selectedHex || !connection) return;
    connection.send({ type: 'action', action: 'demolish', q: selectedHex.q, r: selectedHex.r });
  },
  onAttack: () => {
    if (!selectedHex || !connection) return;
    connection.send({ type: 'action', action: 'attack', q: selectedHex.q, r: selectedHex.r });
  },
  onDropHex: () => {
    if (!selectedHex || !connection) return;
    connection.send({ type: 'action', action: 'drop-hex', q: selectedHex.q, r: selectedHex.r });
  },
  onFortify: () => {
    if (!selectedHex || !connection) return;
    connection.send({ type: 'action', action: 'fortify', q: selectedHex.q, r: selectedHex.r });
  },
});

const techTreePanel = new TechTreePanel(document.body, (techId) => {
  connection?.send({ type: 'action', action: 'unlock-tech', techId });
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
  sessionStorage.removeItem('hexarSession');
  history.replaceState(null, '', '/');
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

let disconnectOverlay: HTMLElement | null = null;

let reconnectBanner: HTMLElement | null = null;

function showReconnecting(attempt: number, max: number) {
  if (!reconnectBanner) {
    reconnectBanner = document.createElement('div');
    reconnectBanner.style.cssText = `
      position:fixed;top:0;left:0;width:100%;padding:6px;text-align:center;
      background:#c0392b;color:#fff;font-family:monospace;font-size:13px;z-index:150;
    `;
    document.body.appendChild(reconnectBanner);
  }
  reconnectBanner.textContent = `Reconnecting… (${attempt}/${max})`;
}

function hideReconnecting() {
  reconnectBanner?.remove();
  reconnectBanner = null;
}


function showAlreadyConnectedOverlay() {
  sessionStorage.removeItem('hexarSession');
  history.replaceState(null, '', '/');
  const overlay = document.createElement('div');
  overlay.style.cssText = `
    position:fixed;top:0;left:0;width:100%;height:100%;
    display:flex;flex-direction:column;align-items:center;justify-content:center;
    background:rgba(0,0,0,0.85);z-index:300;font-family:monospace;color:#e0e0e0;
  `;
  overlay.innerHTML = `
    <div style="font-size:28px;color:#f39c12;margin-bottom:16px">Already Connected</div>
    <div style="color:#aaa;font-size:16px;margin-bottom:24px">This game is already open in another tab.</div>
    <button onclick="location.reload()" style="padding:10px 28px;font-size:16px;background:#4ecdc4;color:#1a1a2e;border:none;border-radius:4px;cursor:pointer;font-family:monospace">Back to Lobby</button>
  `;
  document.body.appendChild(overlay);
}

function showDisconnectOverlay() {
  if (disconnectOverlay) return;
  disconnectOverlay = document.createElement('div');
  disconnectOverlay.style.cssText = `
    position:fixed;top:0;left:0;width:100%;height:100%;
    display:flex;flex-direction:column;align-items:center;justify-content:center;
    background:rgba(0,0,0,0.85);z-index:300;font-family:monospace;color:#e0e0e0;
  `;
  disconnectOverlay.innerHTML = `
    <div style="font-size:32px;color:#ff6b6b;margin-bottom:16px">Disconnected</div>
    <div style="color:#aaa;font-size:16px;margin-bottom:24px">Could not reconnect to server.</div>
    <button onclick="location.reload()" style="padding:10px 28px;font-size:16px;background:#4ecdc4;color:#1a1a2e;border:none;border-radius:4px;cursor:pointer;font-family:monospace">Back to Lobby</button>
  `;
  document.body.appendChild(disconnectOverlay);
  sessionStorage.removeItem('hexarSession');
  history.replaceState(null, '', '/');
}

function updateState(newState: GameState) {
  state = newState;

  if (state.over && myPlayerId > 0) {
    showVictory(state.winner === myPlayerId);
  }

  if (myPlayerId > 0) {
    let hexCount = 0;
    let income = 0;
    let tpRate = 0;
    const player = state.players.get(String(myPlayerId));
    for (const [, hex] of state.hexes) {
      if (hex.owner === myPlayerId) {
        hexCount++;
        income += hexIncome(hex, player);
        if (hex.building === BUILDING_RESEARCH) {
          tpRate += RESEARCH_PER_LEVEL * hex.level;
        }
      }
    }
    const maintenance = calcMaintenance(hexCount, player);
    const gold = player?.gold ?? 0;
    const tp = player?.tp ?? 0;
    updateHUD(hud, myPlayerId, gold, hexCount, income, maintenance, tp, tpRate, player?.vanguardTimer ?? 0);

    if (player) {
      techTreePanel.update(player);
    }
    autoDropPanel.update(player ?? null);

    dropMap = new Set<string>();
    if (player?.autoDropActive) {
      const droppable = [...state.hexes.values()]
        .filter(h => h.owner === myPlayerId && !h.capital &&
                     !state!.battles.some(b => b.dq === h.q && b.dr === h.r));
      if (droppable.length > 0) {
        const minIncome = Math.min(...droppable.map(h => hexIncome(h, player)));
        const incomeGroup = droppable.filter(h => hexIncome(h, player) === minIncome);
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
        const atkPwr = bestAdjacentPower(state, myPlayerId, current, player);
        const battle = state.battles.find(b => b.dq === current.q && b.dr === current.r) ?? null;
        buildMenu.updateWithActions(current, gold, isOwn, isEnemy, atkPwr, battle, player ?? null, state);
      }
    }
  }

  renderer.setState(state);
}

function onSnapshot(msg: SnapshotMsg) {
  updateState(applySnapshot(msg));
}

function onDelta(msg: DeltaMsg) {
  if (!state) return;
  updateState(applyDelta(state, msg));
}

function hexIncome(hex: HexDTO, player?: PlayerDTO | null): number {
  if (hex.building !== BUILDING_GOLD) {
    return BASE_INCOME_PER_SEC;
  }
  let income = (BASE_INCOME_PER_SEC + GOLD_PER_LEVEL * hex.level) * GOLD_BONUS_MULTIPLIER;
  if (player?.tech?.[TECH_COMPOUND_GROWTH]) {
    income *= COMPOUND_GROWTH_MULTIPLIER;
  }
  if (player?.tech?.[TECH_PROSPERITY]) {
    income += PROSPERITY_BONUS;
  }
  return income;
}

function totalInvested(hex: HexDTO): number {
  if (hex.building === 0 || hex.level === 0) return 0;
  const base = hex.building === BUILDING_RESEARCH ? RESEARCH_BUILD_COST
             : hex.building === BUILDING_GOLD     ? GOLD_BUILD_COST
             :                                      POWER_BUILD_COST;
  return base * (Math.pow(2, hex.level) - 1);
}

function bestAdjacentPower(gs: GameState, playerId: number, target: HexDTO, player?: PlayerDTO | null): number {
  let best = 0;
  for (const n of neighbors({ q: target.q, r: target.r })) {
    const key = `${n.q},${n.r}`;
    const hs = gs.hexes.get(key);
    if (!hs || hs.owner !== playerId) continue;
    let p = 0;
    if (hs.capital) p = CAPITAL_POWER;
    if (hs.building === BUILDING_POWER) p += hs.level;
    if (player?.tech?.[TECH_IRON_GRIP]) p++;
    if (p > best) best = p;
  }
  return best;
}

function startGame(code: string, token: string) {
  history.replaceState(null, '', `/#${code}:${token}`);
  sessionStorage.setItem('hexarSession', JSON.stringify({ code, token }));

  connection = new Connection(code, token, {
    onSnapshot,
    onDelta,
    onWelcome: (msg) => {
      myPlayerId = msg.playerId;
      hideReconnecting();
      console.log(`assigned player ${myPlayerId}`);
    },
    onReconnecting: showReconnecting,
    onAlreadyConnected: showAlreadyConnectedOverlay,
    onDisconnect: showDisconnectOverlay,
  });
}

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
      connection?.send({ type: 'action', action: 'claim', q, r });
      selectedHex = null;
      renderer.setSelected(null);
      buildMenu.hide();
      return;
    }

    if (dropMap.has(key) && hex.owner === myPlayerId) {
      connection?.send({ type: 'action', action: 'drop-hex', q, r });
      return;
    }

    const battle = state.battles.find(b => b.dq === q && b.dr === r) ?? null;
    if (battle && hex.owner === myPlayerId) {
      const cap = Math.min(COUNTER_SPEND_CAP, Math.floor(battle.timeLeft));
      const player = state.players.get(String(myPlayerId));
      if ((player?.gold ?? 0) >= COUNTER_SPEND_COST && battle.counterBoost < cap) {
        connection?.send({ type: 'action', action: 'counter-spend', q, r });
        return;
      }
    }

    selectedHex = hex;
    renderer.setSelected({ q, r });

    const player = state.players.get(String(myPlayerId));
    const gold = player?.gold ?? 0;
    const isOwn = hex.owner === myPlayerId;
    const isEnemy = hex.owner !== 0 && hex.owner !== myPlayerId;
    const atkPwr = bestAdjacentPower(state, myPlayerId, hex, player);
    buildMenu.updateWithActions(hex, gold, isOwn, isEnemy, atkPwr, battle, player ?? null, state);
  }
);

// Try to reconnect from sessionStorage first, then URL hash, else show lobby
function tryHashOrLobby() {
  const hash = window.location.hash.slice(1);
  const [code, token] = hash.split(':');
  if (code && token) {
    startGame(code, token);
  } else {
    new LobbyUI(document.body, ({ code, token }) => startGame(code, token));
  }
}

const saved = sessionStorage.getItem('hexarSession');
if (saved) {
  try {
    const { code, token } = JSON.parse(saved) as { code: string; token: string };
    startGame(code, token);
  } catch {
    sessionStorage.removeItem('hexarSession');
    tryHashOrLobby();
  }
} else {
  tryHashOrLobby();
}
