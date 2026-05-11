import './style.css';
import { Connection } from './net/connection';
import { applySnapshot, applyDelta, GameState, SnapshotMsg, DeltaMsg, HexDTO, PlayerDTO } from './state/state';
import { Renderer } from './render/renderer';
import { setupInput } from './input/input';
import { updateHUD, calcMaintenance } from './ui/hud';
import { Sidebar } from './ui/sidebar';
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
  COLORS,
} from './constants';

const canvas = document.getElementById('game') as HTMLCanvasElement;
const hud = document.getElementById('hud')!;
const renderer = new Renderer(canvas);

let state: GameState | null = null;
let myPlayerId = 0;
let roomCode = '';
let selectedHex: HexDTO | null = null;
let selectedTool: string | null = null;
let dropMap = new Set<string>();
let connection: Connection | null = null;

const sidebar = new Sidebar(document.body, {
  onToolSelect: (tool) => {
    selectedTool = tool;
  },
  onUpgrade: () => {
    if (!selectedHex || !connection) return;
    connection.send({ type: 'action', action: 'upgrade', q: selectedHex.q, r: selectedHex.r });
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
  onCounterSpend: () => {
    if (!selectedHex || !connection) return;
    connection.send({ type: 'action', action: 'counter-spend', q: selectedHex.q, r: selectedHex.r });
  },
});

const techTreePanel = new TechTreePanel(document.body, (techId) => {
  connection?.send({ type: 'action', action: 'unlock-tech', techId });
});

const autoDropPanel = new AutoDropPanel(document.body);

window.addEventListener('keydown', (e) => {
  // Don't intercept text input
  if (e.target instanceof HTMLInputElement) return;

  // Tech tree toggle
  if (e.key === 't' || e.key === 'T') {
    techTreePanel.toggle();
    if (state && myPlayerId > 0) {
      const player = state.players.get(String(myPlayerId));
      if (player) techTreePanel.update(player);
    }
    return;
  }

  // Sidebar tool shortcuts with toggle
  const key = e.key.toUpperCase();
  switch (key) {
    case 'Q':
      sidebar.selectTool(selectedTool === 'economy' ? null : 'economy');
      break;
    case 'W':
      sidebar.selectTool(selectedTool === 'power' ? null : 'power');
      break;
    case 'E':
      sidebar.selectTool(selectedTool === 'research' ? null : 'research');
      break;
    case 'ESCAPE':
      sidebar.selectTool(null);
      break;
  }

  // Context action shortcuts (require selected hex)
  if (!selectedHex || !state) return;

  switch (key) {
    case ' ': // Space = upgrade
      e.preventDefault();
      if (selectedHex.owner === myPlayerId && selectedHex.building > 0) {
        connection?.send({ type: 'action', action: 'upgrade', q: selectedHex.q, r: selectedHex.r });
      }
      break;
    case 'A': // Attack
      if (selectedHex.owner !== myPlayerId && selectedHex.owner > 0) {
        connection?.send({ type: 'action', action: 'attack', q: selectedHex.q, r: selectedHex.r });
      }
      break;
    case 'D': // Demolish
      if (selectedHex.owner === myPlayerId && selectedHex.building > 0) {
        connection?.send({ type: 'action', action: 'demolish', q: selectedHex.q, r: selectedHex.r });
      }
      break;
    case 'X': // Sell hex
      if (selectedHex.owner === myPlayerId && !selectedHex.capital) {
        connection?.send({ type: 'action', action: 'drop-hex', q: selectedHex.q, r: selectedHex.r });
      }
      break;
    case 'F': // Fortify
      if (selectedHex.owner === myPlayerId) {
        connection?.send({ type: 'action', action: 'fortify', q: selectedHex.q, r: selectedHex.r });
      }
      break;
    case 'C': { // Counter-spend
      const hex = selectedHex;
      const battle = state.battles.find(b => b.dq === hex.q && b.dr === hex.r);
      if (battle && hex.owner === myPlayerId) {
        connection?.send({ type: 'action', action: 'counter-spend', q: hex.q, r: hex.r });
      }
      break;
    }
  }
});

let victoryOverlay: HTMLElement | null = null;
let waitingOverlay: HTMLElement | null = null;

function showWaitingOverlay() {
  if (waitingOverlay) return;
  waitingOverlay = document.createElement('div');
  waitingOverlay.className = 'overlay';
  waitingOverlay.style.background = COLORS.overlayBg;
  waitingOverlay.innerHTML = `
    <div class="overlay-title" style="font-size:36px;color:${COLORS.accent}">Waiting for opponent…</div>
    <div class="overlay-sub" style="color:${COLORS.textMuted};font-size:16px">Share this room code with player 2:</div>
    <div style="font-size:56px;font-weight:bold;color:#fff;letter-spacing:10px;margin-top:16px">${roomCode}</div>
  `;
  document.body.appendChild(waitingOverlay);
}

function hideWaitingOverlay() {
  waitingOverlay?.remove();
  waitingOverlay = null;
}

function showVictory(isWinner: boolean, winReason: string) {
  if (victoryOverlay) return;
  sessionStorage.removeItem('hexarSession');
  history.replaceState(null, '', '/');
  victoryOverlay = document.createElement('div');
  victoryOverlay.className = 'overlay';
  victoryOverlay.style.background = COLORS.overlayBg;
  const msg = isWinner ? 'You win!' : 'You lose!';
  const color = isWinner ? COLORS.accent : COLORS.player2;
  const subtitle = winReason === 'forfeit'
    ? (isWinner ? 'Opponent forfeited' : 'You forfeited')
    : 'Capital captured';
  victoryOverlay.innerHTML = `
    <div class="overlay-title" style="font-size:48px;color:${color}">${msg}</div>
    <div class="overlay-sub" style="color:${COLORS.textMuted};font-size:16px">${subtitle}</div>
  `;
  document.body.appendChild(victoryOverlay);
}

let disconnectOverlay: HTMLElement | null = null;
let reconnectBanner: HTMLElement | null = null;
let pauseBanner: HTMLElement | null = null;

function updatePauseBanner(timeLeft: number) {
  if (!pauseBanner) {
    pauseBanner = document.createElement('div');
    pauseBanner.className = 'banner';
    pauseBanner.style.background = COLORS.pause;
    document.body.appendChild(pauseBanner);
  }
  const mins = Math.floor(timeLeft / 60);
  const secs = Math.floor(timeLeft % 60);
  const countdown = `${mins}:${String(secs).padStart(2, '0')}`;
  pauseBanner.innerHTML = `<b>Game paused</b> — opponent disconnected. Reconnect within <b>${countdown}</b> or forfeit.`;
}

function hidePauseBanner() {
  pauseBanner?.remove();
  pauseBanner = null;
}

function showReconnecting(attempt: number, max: number) {
  if (!reconnectBanner) {
    reconnectBanner = document.createElement('div');
    reconnectBanner.className = 'banner';
    reconnectBanner.style.background = COLORS.reconnect;
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
  overlay.className = 'overlay';
  overlay.style.background = 'rgba(0,0,0,0.85)';
  overlay.innerHTML = `
    <div class="overlay-title" style="font-size:28px;color:${COLORS.warning};margin-bottom:16px">Already Connected</div>
    <div class="overlay-sub" style="color:${COLORS.textMuted};font-size:16px;margin-bottom:24px">This game is already open in another tab.</div>
    <button onclick="location.reload()" class="btn overlay-action" style="background:${COLORS.accent};color:${COLORS.background}">Back to Lobby</button>
  `;
  document.body.appendChild(overlay);
}

function showDisconnectOverlay() {
  if (disconnectOverlay) return;
  disconnectOverlay = document.createElement('div');
  disconnectOverlay.className = 'overlay';
  disconnectOverlay.style.background = 'rgba(0,0,0,0.85)';
  disconnectOverlay.innerHTML = `
    <div class="overlay-title" style="font-size:32px;color:${COLORS.player2};margin-bottom:16px">Disconnected</div>
    <div class="overlay-sub" style="color:${COLORS.textMuted};font-size:16px;margin-bottom:24px">Could not reconnect to server.</div>
    <button onclick="location.reload()" class="btn overlay-action" style="background:${COLORS.accent};color:${COLORS.background}">Back to Lobby</button>
  `;
  document.body.appendChild(disconnectOverlay);
  sessionStorage.removeItem('hexarSession');
  history.replaceState(null, '', '/');
}

function updateState(newState: GameState) {
  state = newState;

  if (state.waiting) {
    showWaitingOverlay();
    return;
  }
  hideWaitingOverlay();

  if (state.paused) {
    updatePauseBanner(state.pauseTimeLeft);
  } else {
    hidePauseBanner();
  }

  if (state.over && myPlayerId > 0) {
    showVictory(state.winner === myPlayerId, state.winReason);
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

    // Trigger gold floaters on economy hexes (throttled inside renderer to 1 per 2s per hex)
    for (const [, hex] of state.hexes) {
      if (hex.owner === myPlayerId && hex.building === BUILDING_GOLD) {
        renderer.addFloater(hex.q, hex.r, `+${hexIncome(hex, player).toFixed(1)}`);
      }
    }

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
        sidebar.updateContext(current, gold, isOwn, isEnemy, atkPwr, battle, player ?? null, state);
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
  // Trigger capture flash for hex ownership changes
  if (msg.hexChanges) {
    for (const hex of msg.hexChanges) {
      if (hex.owner > 0 && hex.previousOwner !== undefined && hex.previousOwner !== hex.owner) {
        const flashColor = hex.owner === 1 ? COLORS.player1 : COLORS.player2;
        renderer.addCaptureFlash(hex.q, hex.r, flashColor);
      }
    }
  }
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
  roomCode = code;
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
      sidebar.hide();
      return;
    }

    // If tool is selected and hex is valid target, apply tool
    if (selectedTool && hex.owner === myPlayerId && hex.building === 0) {
      let building: string | undefined;
      switch (selectedTool) {
        case 'economy':
          building = 'gold';
          break;
        case 'power':
          building = 'power';
          break;
        case 'research':
          building = 'research';
          break;
      }
      if (building) {
        connection?.send({ type: 'action', action: 'upgrade', q, r, building });
        // Tool stays selected for batch operations
        return;
      }
    }

    // No tool selected or invalid target → normal hex interaction
    if (hex.owner === 0) {
      connection?.send({ type: 'action', action: 'claim', q, r });
      selectedHex = null;
      renderer.setSelected(null);
      sidebar.hide();
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
    sidebar.updateContext(hex, gold, isOwn, isEnemy, atkPwr, battle, player ?? null, state);
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
