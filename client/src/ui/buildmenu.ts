import { HexDTO } from '../state/state';

export interface BuildMenuCallbacks {
  onBuild: (building: 'economy' | 'defense') => void;
  onUpgrade: () => void;
  onDemolish: () => void;
  onAttack: () => void;
}

const BUILD_COSTS: Record<number, number> = { 1: 80, 2: 60, 3: 80 };

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
        case 'build-economy': this.callbacks.onBuild('economy'); break;
        case 'build-defense': this.callbacks.onBuild('defense'); break;
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
      gold >= 80, gold >= 60, gold >= upgCost,
      canAttack, attackerPower, defPower, hasBattle,
    ].join('|');

    if (key === this.lastKey) {
      this.el.style.display = 'flex';
      return;
    }
    this.lastKey = key;

    let html = `<span style="margin-right:4px">[${hex.q},${hex.r}] Pwr:${defPower}</span>`;

    if (isOwn) {
      if (hex.building === 0) {
        html += this.makeBtn('Economy (80g)', gold >= 80, 'build-economy');
        html += this.makeBtn('Defense (60g)', gold >= 60, 'build-defense');
      } else {
        const refund = demolishRefund(hex.building, hex.level);
        html += this.makeBtn(`Demolish (+${refund}g)`, true, 'demolish');
        html += this.makeBtn(`Upgrade (${upgCost}g)`, gold >= upgCost, 'upgrade');
      }
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
