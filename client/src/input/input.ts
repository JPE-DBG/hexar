import { pixelToHex } from '../hexmath';

export function setupInput(
  canvas: HTMLCanvasElement,
  getOffset: () => { x: number; y: number },
  onHexClick: (q: number, r: number) => void
) {
  canvas.addEventListener('click', (ev) => {
    const rect = canvas.getBoundingClientRect();
    const offset = getOffset();
    const px = ev.clientX - rect.left - offset.x;
    const py = ev.clientY - rect.top - offset.y;
    const hex = pixelToHex(px, py);
    onHexClick(hex.q, hex.r);
  });
}
