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
  el.textContent = `Player ${playerId} | Gold: ${gold.toFixed(0)} (${sign}${net.toFixed(1)}/s) | Hexes: ${hexCount}`;
}

export function calcIncome(hexCount: number): number {
  return hexCount * 2.0;
}

export function calcMaintenance(hexCount: number): number {
  const tier1 = Math.min(hexCount, 10);
  const tier2 = Math.min(Math.max(hexCount - 10, 0), 10);
  const tier3 = Math.max(hexCount - 20, 0);
  return tier1 * 1.0 + tier2 * 2.0 + tier3 * 3.0;
}
