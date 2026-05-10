import { HexDTO, BattleDTO, PlayerDTO, GameState } from '../state/state';
import { neighbors } from '../hexmath';
import {
  GOLD_BUILD_COST, GOLD_PER_LEVEL, GOLD_BONUS_MULTIPLIER,
  POWER_BUILD_COST, POWER_PER_LEVEL,
  RESEARCH_BUILD_COST, RESEARCH_PER_LEVEL,
  DEMOLISH_REFUND, ATTACK_COST, CAPITAL_POWER, BASE_INCOME_PER_SEC,
  BUILDING_GOLD, BUILDING_POWER, BUILDING_RESEARCH,
  FORTIFY_COST, COUNTER_SPEND_COST, COUNTER_SPEND_CAP,
  PROSPERITY_BONUS, COMPOUND_GROWTH_MULTIPLIER,
  RECLAMATION_ATTACK_COST, VANGUARD_ATTACK_COST,
  GARRISON_MAX_BOOST,
  TECH_COMPOUND_GROWTH, TECH_PROSPERITY, TECH_RECLAMATION,
  TECH_VANGUARD, TECH_GARRISON, TECH_FORTIFY, TECH_IRON_GRIP,
} from '../constants';

export interface SidebarCallbacks {
  onToolSelect: (tool: string | null) => void;
  onUpgrade: () => void;
  onDemolish: () => void;
  onAttack: () => void;
  onDropHex: () => void;
  onFortify: () => void;
  onCounterSpend: () => void;
}

// Enhanced SVG icons (32x32, color-coded)
const ICON_ECONOMY = `<svg width="32" height="32" viewBox="0 0 32 32">
  <circle cx="16" cy="16" r="13" fill="#f9ca2422" stroke="#f9ca24" stroke-width="2.5"/>
  <text x="16" y="21" text-anchor="middle" fill="#f9ca24" font-size="14" font-weight="bold" font-family="monospace">G</text>
</svg>`;

const ICON_POWER = `<svg width="32" height="32" viewBox="0 0 32 32">
  <path d="M16 4 L16 28 M10 20 L22 20" stroke="#e74c3c" stroke-width="3" stroke-linecap="round"/>
  <polygon points="16,2 12,10 20,10" fill="#e74c3c"/>
</svg>`;

const ICON_RESEARCH = `<svg width="32" height="32" viewBox="0 0 32 32">
  <path d="M10,26 L16,6 L22,26" stroke="#45b7d1" stroke-width="2.5" fill="none"/>
  <ellipse cx="16" cy="25" rx="8" ry="3.5" stroke="#45b7d1" stroke-width="2" fill="none"/>
  <line x1="12" y1="17" x2="20" y2="17" stroke="#45b7d1" stroke-width="2"/>
</svg>`;

const BUILD_COSTS: Record<number, number> = {
  [BUILDING_GOLD]: GOLD_BUILD_COST,
  [BUILDING_POWER]: POWER_BUILD_COST,
  [BUILDING_RESEARCH]: RESEARCH_BUILD_COST,
};

function upgradeCost(building: number, level: number): number {
  const base = BUILD_COSTS[building] ?? RESEARCH_BUILD_COST;
  return base * Math.pow(2, level);
}

function demolishRefund(building: number, level: number): number {
  const base = BUILD_COSTS[building] ?? RESEARCH_BUILD_COST;
  return base * (Math.pow(2, level) - 1) * DEMOLISH_REFUND;
}

function effectiveAttackCost(player: PlayerDTO | null, targetHex: HexDTO): number {
  if (!player) return ATTACK_COST;
  let cost = ATTACK_COST;
  if (player.tech?.[TECH_RECLAMATION] && targetHex.previousOwner === player.id) {
    cost -= RECLAMATION_ATTACK_COST;
  }
  if (player.tech?.[TECH_VANGUARD] && player.vanguardTimer > 0) {
    cost -= VANGUARD_ATTACK_COST;
  }
  return Math.max(0, cost);
}

function countAdjacentOwned(hex: HexDTO, state: GameState, owner: number): number {
  let count = 0;
  for (const n of neighbors({ q: hex.q, r: hex.r })) {
    const key = `${n.q},${n.r}`;
    const hs = state.hexes.get(key);
    if (hs && hs.owner === owner) count++;
  }
  return count;
}

function calcPower(hex: HexDTO, state?: GameState): number {
  let p = 0;
  if (hex.capital) p = CAPITAL_POWER;
  if (hex.building === BUILDING_POWER) p += hex.level;
  if (state && hex.owner > 0) {
    const ownerPlayer = state.players.get(String(hex.owner));
    if (ownerPlayer?.tech?.[TECH_IRON_GRIP]) p++;
  }
  return p;
}

export class Sidebar {
  private el: HTMLElement;
  private selectedTool: string | null = null;
  private callbacks: SidebarCallbacks;
  private lastKey = '';

  constructor(parent: HTMLElement, callbacks: SidebarCallbacks) {
    this.callbacks = callbacks;
    this.el = document.createElement('div');
    this.el.id = 'sidebar';
    this.el.innerHTML = `
      <div class="sidebar-section">
        <div class="section-label">BUILD</div>
        <button class="sidebar-btn" data-tool="economy" data-hotkey="Q">
          ${ICON_ECONOMY}
          <span class="btn-label">Economy</span>
          <span class="btn-cost">60g</span>
          <span class="btn-hotkey">Q</span>
        </button>
        <button class="sidebar-btn" data-tool="power" data-hotkey="W">
          ${ICON_POWER}
          <span class="btn-label">Power</span>
          <span class="btn-cost">60g</span>
          <span class="btn-hotkey">W</span>
        </button>
        <button class="sidebar-btn" data-tool="research" data-hotkey="E">
          ${ICON_RESEARCH}
          <span class="btn-label">Research</span>
          <span class="btn-cost">80g</span>
          <span class="btn-hotkey">E</span>
        </button>
      </div>

      <div class="sidebar-section">
        <div class="section-label">MANAGE</div>
        <button class="sidebar-btn" data-tool="demolish" data-hotkey="D">
          <span class="btn-icon">🗑</span>
          <span class="btn-label">Demolish</span>
          <span class="btn-hotkey">D</span>
        </button>
        <button class="sidebar-btn" data-tool="sell-hex" data-hotkey="X">
          <span class="btn-icon">❌</span>
          <span class="btn-label">Sell Hex</span>
          <span class="btn-hotkey">X</span>
        </button>
      </div>

      <div class="sidebar-section sidebar-context" id="sidebar-context">
        <div class="section-label">SELECT HEX</div>
        <div class="context-hint">Click a hex to see actions</div>
      </div>
    `;
    parent.appendChild(this.el);
    this.bindEvents();
  }

  private bindEvents() {
    this.el.addEventListener('pointerdown', (e) => {
      const btn = (e.target as HTMLElement).closest('.sidebar-btn') as HTMLElement | null;
      if (!btn) return;
      e.preventDefault();

      const tool = btn.getAttribute('data-tool');
      const action = btn.getAttribute('data-action');

      if (tool) {
        this.selectTool(tool);
      } else if (action) {
        // Context actions
        switch (action) {
          case 'upgrade': this.callbacks.onUpgrade(); break;
          case 'demolish': this.callbacks.onDemolish(); break;
          case 'attack': this.callbacks.onAttack(); break;
          case 'drop-hex': this.callbacks.onDropHex(); break;
          case 'fortify': this.callbacks.onFortify(); break;
          case 'counter-spend': this.callbacks.onCounterSpend(); break;
        }
      }
    });
  }

  selectTool(tool: string | null) {
    // Toggle: if clicking already-selected tool, deselect it
    if (this.selectedTool === tool) {
      tool = null;
    }

    this.selectedTool = tool;

    // Update visual state
    this.el.querySelectorAll('.sidebar-btn[data-tool]').forEach(btn => {
      btn.classList.toggle('selected', btn.getAttribute('data-tool') === tool);
    });

    this.callbacks.onToolSelect(tool);
  }

  getSelectedTool(): string | null {
    return this.selectedTool;
  }

  updateContext(
    hex: HexDTO | null,
    gold: number,
    isOwn: boolean,
    isEnemy: boolean,
    attackerPower: number,
    battle: BattleDTO | null,
    player: PlayerDTO | null,
    state: GameState | null
  ) {
    const context = this.el.querySelector('#sidebar-context') as HTMLElement;
    if (!hex) {
      context.innerHTML = `
        <div class="section-label">SELECT HEX</div>
        <div class="context-hint">Click a hex to see actions</div>
      `;
      this.lastKey = '';
      return;
    }

    const defPower = calcPower(hex, state ?? undefined);
    let garrisonBonus = 0;
    if (isEnemy && state && hex.owner > 0) {
      const ownerPlayer = state.players.get(String(hex.owner));
      if (ownerPlayer?.tech?.[TECH_GARRISON]) {
        garrisonBonus = Math.min(GARRISON_MAX_BOOST, countAdjacentOwned(hex, state, hex.owner));
      }
    }
    const effectiveDefPower = defPower + garrisonBonus;
    const effectiveCost = effectiveAttackCost(player, hex);
    const canAttack = !battle && gold >= effectiveCost && attackerPower > effectiveDefPower;

    // Cache key to avoid unnecessary re-renders
    const key = [
      hex.q, hex.r, hex.building, hex.level, isOwn, isEnemy,
      gold >= GOLD_BUILD_COST, gold >= RESEARCH_BUILD_COST,
      canAttack, attackerPower, defPower,
      battle ? Math.floor(battle.timeLeft) : -1,
      hex.fortifyTimer > 0 ? 1 : 0,
    ].join('|');

    if (key === this.lastKey) {
      return;
    }
    this.lastKey = key;

    let powerStr = `Pwr:${defPower}`;
    if (state && hex.owner > 0 && isEnemy) {
      const ownerPlayer = state.players.get(String(hex.owner));
      if (ownerPlayer?.tech?.[TECH_GARRISON]) {
        const gb = Math.min(GARRISON_MAX_BOOST, countAdjacentOwned(hex, state, hex.owner));
        if (gb > 0) powerStr += ` +${gb} def`;
      }
    }

    let html = `<div class="section-label">[${hex.q},${hex.r}] ${powerStr}</div>`;

    if (isOwn) {
      const hasBuilding = hex.building !== 0;
      const refund = hasBuilding ? demolishRefund(hex.building, hex.level) : 0;

      // Upgrade button
      if (hasBuilding) {
        const upgCost = upgradeCost(hex.building, hex.level);
        const canUpgrade = gold >= upgCost;
        html += `
          <button class="sidebar-btn ${canUpgrade ? '' : 'disabled'}"
                  data-action="upgrade"
                  ${canUpgrade ? '' : 'disabled'}>
            <span class="btn-icon">⬆</span>
            <span class="btn-label">Upgrade L${hex.level + 1}</span>
            <span class="btn-cost">${upgCost}g</span>
            <span class="btn-hotkey">Space</span>
          </button>
        `;
      }

      // Demolish button
      html += `
        <button class="sidebar-btn ${hasBuilding ? '' : 'disabled'}"
                data-action="demolish"
                ${hasBuilding ? '' : 'disabled'}>
          <span class="btn-icon">🗑</span>
          <span class="btn-label">Demolish${hasBuilding ? ` (+${Math.floor(refund)}g)` : ''}</span>
          <span class="btn-hotkey">D</span>
        </button>
      `;

      // Sell hex button (only if not capital and no battle)
      if (!hex.capital && !battle) {
        const sellRefund = hasBuilding ? demolishRefund(hex.building, hex.level) : 0;
        html += `
          <button class="sidebar-btn" data-action="drop-hex">
            <span class="btn-icon">❌</span>
            <span class="btn-label">Sell Hex${sellRefund > 0 ? ` (+${Math.floor(sellRefund)}g)` : ''}</span>
            <span class="btn-hotkey">X</span>
          </button>
        `;
      }

      // Fortify button
      if (player?.tech?.[TECH_FORTIFY] && hex.fortifyTimer <= 0 && !battle) {
        const canFortify = gold >= FORTIFY_COST;
        html += `
          <button class="sidebar-btn ${canFortify ? '' : 'disabled'}"
                  data-action="fortify"
                  ${canFortify ? '' : 'disabled'}>
            <span class="btn-icon">🛡</span>
            <span class="btn-label">Fortify (${FORTIFY_COST}g)</span>
            <span class="btn-hotkey">F</span>
          </button>
        `;
      }

      // Counter-spend button (during battle)
      if (battle) {
        const cap = Math.min(COUNTER_SPEND_CAP, Math.floor(battle.timeLeft));
        const canCounterSpend = gold >= COUNTER_SPEND_COST && battle.counterBoost < cap;
        html += `
          <button class="sidebar-btn ${canCounterSpend ? '' : 'disabled'}"
                  data-action="counter-spend"
                  ${canCounterSpend ? '' : 'disabled'}>
            <span class="btn-icon">💰</span>
            <span class="btn-label">Counter +1P (${COUNTER_SPEND_COST}g)</span>
            <span class="btn-hotkey">C</span>
          </button>
        `;
      }
    } else if (isEnemy) {
      // Attack button
      html += `
        <button class="sidebar-btn ${canAttack ? '' : 'disabled'}"
                data-action="attack"
                ${canAttack ? '' : 'disabled'}>
          <span class="btn-icon">⚔</span>
          <span class="btn-label">Attack (${effectiveCost}g)</span>
          <span class="btn-hotkey">A</span>
        </button>
      `;

      if (battle) {
        html += `<div class="context-hint" style="color:#fa0">Battle in progress...</div>`;
      } else if (attackerPower <= effectiveDefPower) {
        html += `<div class="context-hint" style="color:#f66">Need Pwr > ${effectiveDefPower}</div>`;
      }
    }

    context.innerHTML = html;
  }

  hide() {
    const context = this.el.querySelector('#sidebar-context') as HTMLElement;
    context.innerHTML = `
      <div class="section-label">SELECT HEX</div>
      <div class="context-hint">Click a hex to see actions</div>
    `;
    this.lastKey = '';
  }
}
