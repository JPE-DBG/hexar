import { BASE_INCOME_PER_SEC, MAINTENANCE_TIER1, MAINTENANCE_TIER2, MAINTENANCE_TIER3, MAINTENANCE_TIER1_CAP, MAINTENANCE_TIER2_CAP, SUPPLY_LINES_TIER1, SUPPLY_LINES_TIER2, SUPPLY_LINES_TIER3, TECH_SUPPLY_LINES, COLORS } from '../constants';
import { PlayerDTO } from '../state/state';

const PLAYER_COLORS: Record<number, string> = {
  1: COLORS.player1,
  2: COLORS.player2,
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
  const netColor = net >= 0 ? '#7bed9f' : '#ff6b81';

  let vanguardBox = '';
  if (vanguardTimer > 0) {
    vanguardBox = `<div class="hud-vanguard">⚡ Vanguard ${vanguardTimer.toFixed(1)}s</div>`;
  }

  el.innerHTML = `
    <div class="hud-player" style="color:${color}">Player ${playerId}</div>
    <div class="hud-row">
      <span class="hud-label">Gold:</span>
      <span class="hud-value">${gold.toFixed(0)}</span>
      <span class="hud-rate" style="color:${netColor}">(${netSign}${net.toFixed(1)}/s)</span>
    </div>
    <div class="hud-row">
      <span class="hud-label">Research:</span>
      <span class="hud-value">${tp.toFixed(0)}</span>
      <span class="hud-rate">(${tpSign}${tpRate.toFixed(2)}/s)</span>
    </div>
    <div class="hud-row">
      <span class="hud-label">Hexes:</span>
      <span class="hud-value">${hexCount}</span>
    </div>
    ${vanguardBox}
  `;
}

export function calcIncome(hexCount: number): number {
  return hexCount * BASE_INCOME_PER_SEC;
}

export function calcMaintenance(hexCount: number, player?: PlayerDTO | null): number {
  const tier1 = Math.min(hexCount, MAINTENANCE_TIER1_CAP);
  const tier2 = Math.min(Math.max(hexCount - MAINTENANCE_TIER1_CAP, 0), MAINTENANCE_TIER2_CAP - MAINTENANCE_TIER1_CAP);
  const tier3 = Math.max(hexCount - MAINTENANCE_TIER2_CAP, 0);

  if (player?.tech?.[TECH_SUPPLY_LINES]) { // TechSupplyLines
    return tier1 * SUPPLY_LINES_TIER1 + tier2 * SUPPLY_LINES_TIER2 + tier3 * SUPPLY_LINES_TIER3;
  }
  return tier1 * MAINTENANCE_TIER1 + tier2 * MAINTENANCE_TIER2 + tier3 * MAINTENANCE_TIER3;
}
