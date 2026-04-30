import { Connection } from './net/connection';
import { applySnapshot, GameState, SnapshotMsg } from './state/state';
import { Renderer } from './render/renderer';

const canvas = document.getElementById('game') as HTMLCanvasElement;
const hud = document.getElementById('hud')!;
const renderer = new Renderer(canvas);

let state: GameState | null = null;

function onMessage(msg: SnapshotMsg) {
  state = applySnapshot(msg);
  renderer.render(state);
  hud.textContent = `Hexes: ${state.hexes.size} | Time: ${state.elapsed.toFixed(1)}s`;
}

const wsUrl = `ws://${window.location.host}/ws`;
new Connection(wsUrl, onMessage);
