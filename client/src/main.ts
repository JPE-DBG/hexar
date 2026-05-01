import { Connection } from './net/connection';
import { applySnapshot, GameState, SnapshotMsg } from './state/state';
import { Renderer } from './render/renderer';
import { setupInput } from './input/input';
import { updateHUD, calcIncome, calcMaintenance } from './ui/hud';

const canvas = document.getElementById('game') as HTMLCanvasElement;
const hud = document.getElementById('hud')!;
const renderer = new Renderer(canvas);

let state: GameState | null = null;
let myPlayerId = 0;

function onSnapshot(msg: SnapshotMsg) {
  state = applySnapshot(msg);
  renderer.render(state);

  if (myPlayerId > 0) {
    let hexCount = 0;
    for (const [, hex] of state.hexes) {
      if (hex.owner === myPlayerId) hexCount++;
    }
    const income = calcIncome(hexCount);
    const maintenance = calcMaintenance(hexCount);
    const player = state.players.get(String(myPlayerId));
    const gold = player?.gold ?? 0;
    updateHUD(hud, myPlayerId, gold, hexCount, income, maintenance);
  }
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
    connection.send({ type: 'action', action: 'claim', q, r });
  }
);
