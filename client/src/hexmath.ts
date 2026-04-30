import { HEX_SIZE } from './constants';

export interface Hex {
  q: number;
  r: number;
}

const DIRECTIONS: Hex[] = [
  { q: 1, r: 0 }, { q: 1, r: -1 }, { q: 0, r: -1 },
  { q: -1, r: 0 }, { q: -1, r: 1 }, { q: 0, r: 1 },
];

export function neighbors(h: Hex): Hex[] {
  return DIRECTIONS.map(d => ({ q: h.q + d.q, r: h.r + d.r }));
}

export function hexToPixel(h: Hex): { x: number; y: number } {
  const x = HEX_SIZE * (Math.sqrt(3) * h.q + (Math.sqrt(3) / 2) * h.r);
  const y = HEX_SIZE * (1.5 * h.r);
  return { x, y };
}

export function pixelToHex(px: number, py: number): Hex {
  const q = (Math.sqrt(3) / 3 * px - 1 / 3 * py) / HEX_SIZE;
  const r = (2 / 3 * py) / HEX_SIZE;
  return hexRound(q, r);
}

function hexRound(q: number, r: number): Hex {
  const s = -q - r;
  let rq = Math.round(q);
  let rr = Math.round(r);
  const rs = Math.round(s);

  const qDiff = Math.abs(rq - q);
  const rDiff = Math.abs(rr - r);
  const sDiff = Math.abs(rs - s);

  if (qDiff > rDiff && qDiff > sDiff) {
    rq = -rr - rs;
  } else if (rDiff > sDiff) {
    rr = -rq - rs;
  }

  return { q: rq, r: rr };
}

export function hexKey(h: Hex): string {
  return `${h.q},${h.r}`;
}
