import { pixelToHex } from '../hexmath';

const DRAG_THRESHOLD = 5;

export function setupInput(
  canvas: HTMLCanvasElement,
  getOffset: () => { x: number; y: number },
  onHexClick: (q: number, r: number) => void,
  onPan: (dx: number, dy: number) => void,
) {
  let dragStart: { x: number; y: number } | null = null;
  let isDragging = false;

  canvas.addEventListener('pointerdown', (ev) => {
    dragStart = { x: ev.clientX, y: ev.clientY };
    isDragging = false;
    canvas.setPointerCapture(ev.pointerId);
  });

  canvas.addEventListener('pointermove', (ev) => {
    if (!dragStart) return;
    const dx = ev.clientX - dragStart.x;
    const dy = ev.clientY - dragStart.y;
    if (!isDragging && (Math.abs(dx) > DRAG_THRESHOLD || Math.abs(dy) > DRAG_THRESHOLD)) {
      isDragging = true;
    }
    if (isDragging) {
      onPan(dx, dy);
      dragStart = { x: ev.clientX, y: ev.clientY };
    }
  });

  canvas.addEventListener('pointerup', (ev) => {
    if (dragStart && !isDragging) {
      const rect = canvas.getBoundingClientRect();
      const offset = getOffset();
      const hex = pixelToHex(ev.clientX - rect.left - offset.x, ev.clientY - rect.top - offset.y);
      onHexClick(hex.q, hex.r);
    }
    dragStart = null;
    isDragging = false;
  });

  canvas.addEventListener('pointercancel', () => {
    dragStart = null;
    isDragging = false;
  });
}
