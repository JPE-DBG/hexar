import { Connection } from './net/connection';
import { applySnapshot, GameState, SnapshotMsg, HexDTO } from './state/state';
import { Renderer } from './render/renderer';
import { setupInput } from './input/input';
import { updateHUD, calcMaintenance } from './ui/hud';
import { BuildMenu } from './ui/buildmenu';
import { neighbors } from './hexmath';

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
});

function onSnapshot(msg: SnapshotMsg) {
  state = applySnapshot(msg);
  renderer.render(state);

  if (myPlayerId > 0) {
    let hexCount = 0;
    let income = 0;
    for (const [, hex] of state.hexes) {
      if (hex.owner === myPlayerId) {
        hexCount++;
        income += hexIncome(hex);
      }
    }
    const maintenance = calcMaintenance(hexCount);
    const player = state.players.get(String(myPlayerId));
    const gold = player?.gold ?? 0;
    updateHUD(hud, myPlayerId, gold, hexCount, income, maintenance);

    if (selectedHex) {
      const key = `${selectedHex.q},${selectedHex.r}`;
      const current = state.hexes.get(key);
      if (current) {
        selectedHex = current;
        const isOwn = current.owner === myPlayerId;
        const isEnemy = current.owner !== 0 && current.owner !== myPlayerId;
        const atkPwr = bestAdjacentPower(state, myPlayerId, current);
        const hasBattle = state.battles.some(b => b.dq === current.q && b.dr === current.r);
        buildMenu.updateWithActions(current, gold, isOwn, isEnemy, atkPwr, hasBattle);
      }
    }
  }
}

function hexIncome(hex: HexDTO): number {
  if (hex.building === 1) {
    return (2.0 + 0.6 * hex.level) * 1.5;
  }
  return 2.0;
}

function bestAdjacentPower(gs: GameState, playerId: number, target: HexDTO): number {
  let best = 0;
  for (const n of neighbors({ q: target.q, r: target.r })) {
    const key = `${n.q},${n.r}`;
    const hs = gs.hexes.get(key);
    if (!hs || hs.owner !== playerId) continue;
    let p = 0;
    if (hs.capital) p = 1;
    if (hs.building === 2) p += hs.level;
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
    const hasBattle = state.battles.some(b => b.dq === hex.q && b.dr === hex.r);
    buildMenu.updateWithActions(hex, gold, isOwn, isEnemy, atkPwr, hasBattle);
  }
);
