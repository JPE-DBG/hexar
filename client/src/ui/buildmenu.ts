import { HexDTO } from '../state/state';
import {
  GOLD_BUILD_COST, GOLD_PER_LEVEL, GOLD_BONUS_MULTIPLIER,
  POWER_BUILD_COST, POWER_PER_LEVEL,
  RESEARCH_BUILD_COST, RESEARCH_PER_LEVEL,
  DEMOLISH_REFUND, ATTACK_COST, CAPITAL_POWER, BASE_INCOME_PER_SEC,
  BUILDING_GOLD, BUILDING_POWER, BUILDING_RESEARCH,
} from '../constants';

export interface BuildMenuCallbacks {
  onUpgrade: (building?: 'gold' | 'power' | 'research') => void;
  onDemolish: () => void;
  onAttack: () => void;
}

const BUILD_COSTS: Record<number, number> = {
  [BUILDING_GOLD]: GOLD_BUILD_COST,
  [BUILDING_POWER]: POWER_BUILD_COST,
  [BUILDING_RESEARCH]: RESEARCH_BUILD_COST,
};
const BUILDING_NAMES: Record<number, string> = { [BUILDING_GOLD]: 'Gold', [BUILDING_POWER]: 'Power', [BUILDING_RESEARCH]: 'Research' };

function upgradeDelta(building: number, currentLevel: number): string {
  if (building === BUILDING_GOLD) {
    const next = (BASE_INCOME_PER_SEC + GOLD_PER_LEVEL * (currentLevel + 1)) * GOLD_BONUS_MULTIPLIER;
    const curr = currentLevel === 0 ? BASE_INCOME_PER_SEC : (BASE_INCOME_PER_SEC + GOLD_PER_LEVEL * currentLevel) * GOLD_BONUS_MULTIPLIER;
    return `+${(next - curr).toFixed(1)}/s`;
  }
  if (building === BUILDING_POWER) return `+${POWER_PER_LEVEL} Pwr`;
  return `+${RESEARCH_PER_LEVEL} TP/s`;
}

function upgradeCost(building: number, level: number): number {
  const base = BUILD_COSTS[building] ?? RESEARCH_BUILD_COST;
  return base * Math.pow(2, level);
}

function demolishRefund(building: number, level: number): number {
  const base = BUILD_COSTS[building] ?? RESEARCH_BUILD_COST;
  return base * (Math.pow(2, level) - 1) * DEMOLISH_REFUND;
}

function upgradeLabel(building: number, currentLevel: number): string {
  const name = BUILDING_NAMES[building] ?? '?';
  const targetLevel = currentLevel + 1;
  const cost = upgradeCost(building, currentLevel);
  const delta = upgradeDelta(building, currentLevel);
  return `${name} ${targetLevel} (${cost}g) ${delta}`;
}

export class BuildMenu {
  private el: HTMLElement;
  private callbacks: BuildMenuCallbacks;
  private lastKey = '';

  constructor(parent: HTMLElement, callbacks: BuildMenuCallbacks) {
    this.callbacks = callbacks;
    this.el = document.createElement('div');
    this.el.id = 'build-menu';
    this.el.style.cssText = `
      position: fixed; bottom: 20px; left: 50%; transform: translateX(-50%);
      background: rgba(0,0,0,0.8); padding: 8px 16px; border-radius: 8px;
      display: none; gap: 8px; align-items: center; font-family: monospace;
      color: #e0e0e0; font-size: 13px; height: 40px; box-sizing: border-box;
      flex-wrap: nowrap; white-space: nowrap;
    `;
    parent.appendChild(this.el);

    this.el.addEventListener('pointerdown', (e) => {
      const btn = (e.target as HTMLElement).closest('button[data-action]') as HTMLElement | null;
      if (!btn || btn.hasAttribute('disabled')) return;
      e.preventDefault();
      const action = btn.getAttribute('data-action');
      switch (action) {
        case 'upgrade-gold': this.callbacks.onUpgrade('gold'); break;
        case 'upgrade-power': this.callbacks.onUpgrade('power'); break;
        case 'upgrade-research': this.callbacks.onUpgrade('research'); break;
        case 'upgrade': this.callbacks.onUpgrade(); break;
        case 'demolish': this.callbacks.onDemolish(); break;
        case 'attack': this.callbacks.onAttack(); break;
      }
    });
  }

  updateWithActions(hex: HexDTO | null, gold: number, isOwn: boolean, isEnemy: boolean, attackerPower = 0, hasBattle = false) {
    if (!hex) {
      this.el.style.display = 'none';
      this.lastKey = '';
      return;
    }

    const defPower = this.calcPower(hex);
    const upgCost = upgradeCost(hex.building, hex.level);
    const canAttack = !hasBattle && gold >= ATTACK_COST && attackerPower > defPower;
    const key = [
      hex.q, hex.r, hex.building, hex.level, isOwn, isEnemy,
      gold >= GOLD_BUILD_COST, gold >= RESEARCH_BUILD_COST, gold >= upgCost,
      canAttack, attackerPower, defPower, hasBattle,
    ].join('|');

    if (key === this.lastKey) {
      this.el.style.display = 'flex';
      return;
    }
    this.lastKey = key;

    let html = `<span style="margin-right:4px">[${hex.q},${hex.r}] Pwr:${defPower}</span>`;

    if (isOwn) {
      const hasBuilding = hex.building !== 0;
      const refund = hasBuilding ? demolishRefund(hex.building, hex.level) : 0;

      if (!hasBuilding) {
        html += this.makeBtn(upgradeLabel(BUILDING_GOLD, 0), gold >= GOLD_BUILD_COST, 'upgrade-gold');
        html += this.makeBtn(upgradeLabel(BUILDING_POWER, 0), gold >= POWER_BUILD_COST, 'upgrade-power');
        html += this.makeBtn(upgradeLabel(BUILDING_RESEARCH, 0), gold >= RESEARCH_BUILD_COST, 'upgrade-research');
      } else {
        html += this.makeBtn(upgradeLabel(hex.building, hex.level), gold >= upgCost, 'upgrade');
      }

      html += `<span style="border-left:1px solid #555;height:20px;margin:0 8px;display:inline-block;vertical-align:middle"></span>`;
      html += this.makeBtn(`🗑${hasBuilding ? ` (+${refund}g)` : ''}`, hasBuilding, 'demolish');
    } else if (isEnemy) {
      html += this.makeBtn(`Attack (${ATTACK_COST}g)`, canAttack, 'attack');
      if (hasBattle) {
        html += `<span style="color:#fa0;margin-left:4px">Battle in progress</span>`;
      } else if (attackerPower <= defPower) {
        html += `<span style="color:#f66;margin-left:4px">Need Pwr > ${defPower}</span>`;
      }
    }

    this.el.innerHTML = html;
    this.el.style.display = 'flex';
  }

  private makeBtn(label: string, enabled: boolean, action: string): string {
    const style = enabled
      ? 'background:#4a4a6a;color:#fff;border:1px solid #6a6a8a;padding:4px 10px;border-radius:4px;cursor:pointer;margin:0 2px'
      : 'background:#2a2a3a;color:#666;border:1px solid #3a3a4a;padding:4px 10px;border-radius:4px;margin:0 2px';
    return `<button style="${style}" data-action="${action}" ${enabled ? '' : 'disabled'}>${label}</button>`;
  }

  private calcPower(hex: HexDTO): number {
    let p = 0;
    if (hex.capital) p = CAPITAL_POWER;
    if (hex.building === BUILDING_POWER) p += hex.level;
    return p;
  }

  hide() {
    this.el.style.display = 'none';
    this.lastKey = '';
  }
}
