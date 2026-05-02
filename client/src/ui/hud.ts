import { BASE_INCOME_PER_SEC, MAINTENANCE_TIER1, MAINTENANCE_TIER2, MAINTENANCE_TIER3, MAINTENANCE_TIER1_CAP, MAINTENANCE_TIER2_CAP } from '../constants';

const PLAYER_COLORS: Record<number, string> = {
  1: '#4ecdc4',
  2: '#ff6b6b',
};

export function updateHUD(
  el: HTMLElement,
  playerId: number,
  gold: number,
  hexCount: number,
  income: number,
  maintenance: number
) {
  const net = income - maintenance;
  const sign = net >= 0 ? '+' : '';
  const color = PLAYER_COLORS[playerId] ?? '#e0e0e0';
  el.innerHTML = `<span style="color:${color}">Player ${playerId}</span> | Gold: ${gold.toFixed(0)} (${sign}${net.toFixed(1)}/s) | Hexes: ${hexCount}`;
}

export function calcIncome(hexCount: number): number {
  return hexCount * BASE_INCOME_PER_SEC;
}

export function calcMaintenance(hexCount: number): number {
  const tier1 = Math.min(hexCount, MAINTENANCE_TIER1_CAP);
  const tier2 = Math.min(Math.max(hexCount - MAINTENANCE_TIER1_CAP, 0), MAINTENANCE_TIER2_CAP - MAINTENANCE_TIER1_CAP);
  const tier3 = Math.max(hexCount - MAINTENANCE_TIER2_CAP, 0);
  return tier1 * MAINTENANCE_TIER1 + tier2 * MAINTENANCE_TIER2 + tier3 * MAINTENANCE_TIER3;
}
