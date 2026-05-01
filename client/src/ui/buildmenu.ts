import { HexDTO } from '../state/state';

export interface BuildMenuCallbacks {
  onBuild: (building: 'economy' | 'defense') => void;
  onUpgrade: () => void;
  onDemolish: () => void;
  onAttack: () => void;
}

const BUILDING_NAMES: Record<number, string> = { 1: 'Economy', 2: 'Defense', 3: 'Research' };
const BUILD_COSTS: Record<string, number> = { economy: 80, defense: 60 };
const UPGRADE_BASE: Record<number, number> = { 1: 40, 2: 30, 3: 40 };

function upgradeCost(building: number, level: number): number {
  const base = UPGRADE_BASE[building] ?? 40;
  return base * Math.pow(2, level);
}

export class BuildMenu {
  private el: HTMLElement;
  private callbacks: BuildMenuCallbacks;

  constructor(parent: HTMLElement, callbacks: BuildMenuCallbacks) {
    this.callbacks = callbacks;
    this.el = document.createElement('div');
    this.el.id = 'build-menu';
    this.el.style.cssText = `
      position: fixed; bottom: 20px; left: 50%; transform: translateX(-50%);
      background: rgba(0,0,0,0.8); padding: 12px 16px; border-radius: 8px;
      display: none; gap: 8px; align-items: center; font-family: monospace;
      color: #e0e0e0; font-size: 13px;
    `;
    parent.appendChild(this.el);
  }

  update(hex: HexDTO | null, gold: number, isOwn: boolean, isEnemy: boolean) {
    if (!hex) {
      this.el.style.display = 'none';
      return;
    }

    this.el.style.display = 'flex';
    let html = `<span>[${hex.q},${hex.r}]</span>`;

    if (isOwn) {
      if (hex.building === 0) {
        html += this.btn('Economy (80g)', gold >= 80, () => this.callbacks.onBuild('economy'));
        html += this.btn('Defense (60g)', gold >= 60, () => this.callbacks.onBuild('defense'));
      } else {
        const name = BUILDING_NAMES[hex.building] ?? '?';
        const cost = upgradeCost(hex.building, hex.level);
        html += `<span>${name} L${hex.level}</span>`;
        html += this.btn(`Upgrade (${cost}g)`, gold >= cost, () => this.callbacks.onUpgrade());
        html += this.btn('Demolish', true, () => this.callbacks.onDemolish());
      }
    } else if (isEnemy) {
      html += this.btn('Attack (100g)', gold >= 100, () => this.callbacks.onAttack());
    }

    this.el.innerHTML = html;
    this.bindButtons();
  }

  private btn(label: string, enabled: boolean, _action: () => void): string {
    const cls = enabled ? 'bm-btn' : 'bm-btn disabled';
    return `<button class="${cls}" ${enabled ? '' : 'disabled'}>${label}</button>`;
  }

  private bindButtons() {
    const buttons = this.el.querySelectorAll('button.bm-btn:not(.disabled)');
    const actions = this.collectActions();
    buttons.forEach((btn, i) => {
      if (actions[i]) btn.addEventListener('click', actions[i]);
    });
  }

  private collectActions(): (() => void)[] {
    return this._pendingActions;
  }

  private _pendingActions: (() => void)[] = [];

  updateWithActions(hex: HexDTO | null, gold: number, isOwn: boolean, isEnemy: boolean) {
    if (!hex) {
      this.el.style.display = 'none';
      return;
    }

    this.el.style.display = 'flex';
    this._pendingActions = [];
    let html = `<span style="margin-right:8px">[${hex.q},${hex.r}] Pwr:${this.calcPower(hex)}</span>`;

    if (isOwn) {
      if (hex.building === 0) {
        html += this.makeBtn('Economy (80g)', gold >= 80);
        this._pendingActions.push(() => this.callbacks.onBuild('economy'));
        html += this.makeBtn('Defense (60g)', gold >= 60);
        this._pendingActions.push(() => this.callbacks.onBuild('defense'));
      } else {
        const name = BUILDING_NAMES[hex.building] ?? '?';
        const cost = upgradeCost(hex.building, hex.level);
        html += `<span style="margin-right:8px">${name} L${hex.level}</span>`;
        html += this.makeBtn(`Upgrade (${cost}g)`, gold >= cost);
        this._pendingActions.push(() => this.callbacks.onUpgrade());
        html += this.makeBtn('Demolish', true);
        this._pendingActions.push(() => this.callbacks.onDemolish());
      }
    } else if (isEnemy) {
      html += this.makeBtn('Attack (100g)', gold >= 100);
      this._pendingActions.push(() => this.callbacks.onAttack());
    }

    this.el.innerHTML = html;
    const buttons = this.el.querySelectorAll('button:not([disabled])');
    buttons.forEach((btn, i) => {
      if (this._pendingActions[i]) {
        btn.addEventListener('click', this._pendingActions[i]);
      }
    });
  }

  private makeBtn(label: string, enabled: boolean): string {
    const style = enabled
      ? 'background:#4a4a6a;color:#fff;border:1px solid #6a6a8a;padding:4px 10px;border-radius:4px;cursor:pointer;margin:0 4px'
      : 'background:#2a2a3a;color:#666;border:1px solid #3a3a4a;padding:4px 10px;border-radius:4px;margin:0 4px';
    return `<button style="${style}" ${enabled ? '' : 'disabled'}>${label}</button>`;
  }

  private calcPower(hex: HexDTO): number {
    let p = 0;
    if (hex.capital) p = 1;
    if (hex.building === 2) p += hex.level + 1;
    return p;
  }

  hide() {
    this.el.style.display = 'none';
  }
}
