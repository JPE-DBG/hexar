import { HexDTO } from '../state/state';

export interface BuildMenuCallbacks {
  onUpgrade: (building?: 'economy' | 'defense' | 'research') => void;
  onDemolish: () => void;
  onAttack: () => void;
}

const BUILD_COSTS: Record<number, number> = { 1: 60, 2: 60, 3: 80 };
const BUILDING_NAMES: Record<number, string> = { 1: 'Economy', 2: 'Defense', 3: 'Research' };

function upgradeCost(building: number, level: number): number {
  const base = BUILD_COSTS[building] ?? 80;
  return base * Math.pow(2, level);
}

function demolishRefund(building: number, level: number): number {
  const base = BUILD_COSTS[building] ?? 80;
  return base * (Math.pow(2, level) - 1) * 0.5;
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
        case 'upgrade-economy': this.callbacks.onUpgrade('economy'); break;
        case 'upgrade-defense': this.callbacks.onUpgrade('defense'); break;
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
    const canAttack = !hasBattle && gold >= 100 && attackerPower > defPower;
    const key = [
      hex.q, hex.r, hex.building, hex.level, isOwn, isEnemy,
      gold >= 60, gold >= 80, gold >= upgCost,
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
        // Empty hex: offer all three building types as first upgrade
        html += this.makeBtn(`Economy (60g)`, gold >= 60, 'upgrade-economy');
        html += this.makeBtn(`Defense (60g)`, gold >= 60, 'upgrade-defense');
        html += this.makeBtn(`Research (80g)`, gold >= 80, 'upgrade-research');
      } else {
        // Occupied hex: upgrade existing building
        const name = BUILDING_NAMES[hex.building] ?? '?';
        html += this.makeBtn(`${name} L${hex.level}→${hex.level + 1} (${upgCost}g)`, gold >= upgCost, 'upgrade');
      }

      // Separator + Demolish
      html += `<span style="border-left:1px solid #555;height:20px;margin:0 8px;display:inline-block;vertical-align:middle"></span>`;
      html += this.makeBtn(`🗑${hasBuilding ? ` (+${refund}g)` : ''}`, hasBuilding, 'demolish');
    } else if (isEnemy) {
      html += this.makeBtn('Attack (100g)', canAttack, 'attack');
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
    if (hex.capital) p = 1;
    if (hex.building === 2) p += hex.level;
    return p;
  }

  hide() {
    this.el.style.display = 'none';
    this.lastKey = '';
  }
}
