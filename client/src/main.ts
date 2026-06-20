import './style.css';
import { Renderer, MockHex } from './render/renderer.js';
import { renderShop, ShopCard } from './ui/shop.js';
import { renderDeckPanel, MockCard } from './ui/deck-panel.js';
import { renderBarMeter } from './ui/bar-meter.js';
import { MAP_RADIUS } from './constants.js';

// ── Mock state ─────────────────────────────────────────────────────────────────

function hexDistance(aq: number, ar: number, bq: number, br: number): number {
  return (Math.abs(aq - bq) + Math.abs(ar - br) + Math.abs(aq + ar - bq - br)) / 2;
}

function buildMockMap(): MockHex[] {
  const hexes: MockHex[] = [];
  for (let q = -MAP_RADIUS; q <= MAP_RADIUS; q++) {
    const r1 = Math.max(-MAP_RADIUS, -q - MAP_RADIUS);
    const r2 = Math.min(MAP_RADIUS, -q + MAP_RADIUS);
    for (let r = r1; r <= r2; r++) {
      const dP1 = hexDistance(q, r, -MAP_RADIUS, 0);
      const dP2 = hexDistance(q, r, MAP_RADIUS, 0);
      const isP1Capital = q === -MAP_RADIUS && r === 0;
      const isP2Capital = q === MAP_RADIUS && r === 0;
      const owner = isP1Capital || dP1 <= 2 ? 1
        : isP2Capital || dP2 <= 2 ? 2
        : 0;
      hexes.push({ q, r, owner, capital: isP1Capital || isP2Capital });
    }
  }
  // Add a mock soldier for P1 near center
  const soldier = hexes.find(h => h.q === -1 && h.r === 0);
  if (soldier) soldier.unit = 1;
  return hexes;
}

const mockDeck: MockCard[] = [
  { name: 'Basic Soldier',  type: 'unit',     cost: 2, effect: 'Deploy soldier toward capital' },
  { name: 'Hex Claim',      type: 'building', cost: 1, effect: 'Claim adjacent unclaimed hex' },
  { name: 'Hex Claim',      type: 'building', cost: 1, effect: 'Claim adjacent unclaimed hex' },
  { name: '+1 Bar Burst',   type: 'bar',      cost: 0, effect: 'Add +1 bar instantly' },
  { name: '+1 Bar Burst',   type: 'bar',      cost: 0, effect: 'Add +1 bar instantly' },
  { name: '+15% Bar Speed', type: 'bar',      cost: 1, effect: 'Bar rate ×1.15 for 15s' },
  { name: '+15% Bar Speed', type: 'bar',      cost: 1, effect: 'Bar rate ×1.15 for 15s' },
];

const mockShop: ShopCard[] = [
  { name: 'Tower',         type: 'building', cost: 4, qty: 3 },
  { name: 'Bar Boost',     type: 'bar',      cost: 3, qty: 5 },
  { name: 'Spawner',       type: 'building', cost: 6, qty: 2 },
  { name: 'Heavy Soldier', type: 'unit',     cost: 3, qty: 4 },
];

const mockBar = 6.2;
const mockAside: MockCard | null = null;

// ── Init ───────────────────────────────────────────────────────────────────────

const canvas = document.getElementById('game') as HTMLCanvasElement;
const shopEl = document.getElementById('shop-strip')!;
const barEl = document.getElementById('bar-meter')!;
const deckEl = document.getElementById('deck-panel')!;

const renderer = new Renderer(canvas);
renderer.setHexes(buildMockMap());

renderShop(shopEl, mockShop);
renderBarMeter(barEl, mockBar);
renderDeckPanel(deckEl, mockDeck, mockAside, mockDeck.length);
