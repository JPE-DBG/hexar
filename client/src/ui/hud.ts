import { BASE_INCOME_PER_SEC, MAINTENANCE_TIER1, MAINTENANCE_TIER2, MAINTENANCE_TIER3, MAINTENANCE_TIER1_CAP, MAINTENANCE_TIER2_CAP, SUPPLY_LINES_TIER1, SUPPLY_LINES_TIER2, SUPPLY_LINES_TIER3 } from '../constants';
import { PlayerDTO } from '../state/state';

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
  maintenance: number,
  tp: number,
  tpRate: number,
  vanguardTimer: number
) {
  const net = income - maintenance;
  const netSign = net >= 0 ? '+' : '';
  const tpSign = tpRate >= 0 ? '+' : '';
  const color = PLAYER_COLORS[playerId] ?? '#e0e0e0';

  let vanguardIndicator = '';
  if (vanguardTimer > 0) {
    vanguardIndicator = ` | <span style="color:#fa0">⚡Vanguard ${vanguardTimer.toFixed(1)}s</span>`;
  }

  el.innerHTML =
    `<span style="color:${color}">Player ${playerId}</span>` +
    ` | Gold: ${gold.toFixed(0)} (${netSign}${net.toFixed(1)}/s)` +
    ` | TP: ${tp.toFixed(0)} (${tpSign}${tpRate.toFixed(2)}/s)` +
    ` | Hexes: ${hexCount}` +
    vanguardIndicator;
}

export function calcIncome(hexCount: number): number {
  return hexCount * BASE_INCOME_PER_SEC;
}

export function calcMaintenance(hexCount: number, player?: PlayerDTO | null): number {
  const tier1 = Math.min(hexCount, MAINTENANCE_TIER1_CAP);
  const tier2 = Math.min(Math.max(hexCount - MAINTENANCE_TIER1_CAP, 0), MAINTENANCE_TIER2_CAP - MAINTENANCE_TIER1_CAP);
  const tier3 = Math.max(hexCount - MAINTENANCE_TIER2_CAP, 0);

  if (player?.tech?.[6]) { // TechSupplyLines = 6
    return tier1 * SUPPLY_LINES_TIER1 + tier2 * SUPPLY_LINES_TIER2 + tier3 * SUPPLY_LINES_TIER3;
  }
  return tier1 * MAINTENANCE_TIER1 + tier2 * MAINTENANCE_TIER2 + tier3 * MAINTENANCE_TIER3;
}
