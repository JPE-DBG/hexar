import { HexDTO, BattleDTO, PlayerDTO, GameState } from '../state/state';
import { neighbors } from '../hexmath';
import {
  GOLD_BUILD_COST, GOLD_PER_LEVEL, GOLD_BONUS_MULTIPLIER,
  POWER_BUILD_COST, POWER_PER_LEVEL,
  RESEARCH_BUILD_COST, RESEARCH_PER_LEVEL,
  DEMOLISH_REFUND, ATTACK_COST, CAPITAL_POWER, BASE_INCOME_PER_SEC,
  BUILDING_GOLD, BUILDING_POWER, BUILDING_RESEARCH,
  FORTIFY_COST,
  PROSPERITY_BONUS, COMPOUND_GROWTH_MULTIPLIER,
  RECLAMATION_ATTACK_COST, VANGUARD_ATTACK_COST,
  GARRISON_MAX_BOOST,
} from '../constants';

export interface BuildMenuCallbacks {
  onUpgrade: (building?: 'gold' | 'power' | 'research') => void;
  onDemolish: () => void;
  onAttack: () => void;
  onDropHex: () => void;
  onFortify: () => void;
}

const BUILD_COSTS: Record<number, number> = {
  [BUILDING_GOLD]: GOLD_BUILD_COST,
  [BUILDING_POWER]: POWER_BUILD_COST,
  [BUILDING_RESEARCH]: RESEARCH_BUILD_COST,
};
const BUILDING_NAMES: Record<number, string> = { [BUILDING_GOLD]: 'Gold', [BUILDING_POWER]: 'Power', [BUILDING_RESEARCH]: 'Research' };

function upgradeDelta(building: number, currentLevel: number, player?: PlayerDTO | null): string {
  if (building === BUILDING_GOLD) {
    let next = (BASE_INCOME_PER_SEC + GOLD_PER_LEVEL * (currentLevel + 1)) * GOLD_BONUS_MULTIPLIER;
    let curr = currentLevel === 0 ? BASE_INCOME_PER_SEC : (BASE_INCOME_PER_SEC + GOLD_PER_LEVEL * currentLevel) * GOLD_BONUS_MULTIPLIER;

    // Apply Compound Growth multiplier
    if (player?.tech?.[10]) {
      next *= COMPOUND_GROWTH_MULTIPLIER;
      if (currentLevel > 0) curr *= COMPOUND_GROWTH_MULTIPLIER;
    }

    // Apply Prosperity bonus
    if (player?.tech?.[2]) {
      next += PROSPERITY_BONUS;
      if (currentLevel > 0) curr += PROSPERITY_BONUS;
    }

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

function effectiveAttackCost(player: PlayerDTO | null, targetHex: HexDTO): number {
  if (!player) return ATTACK_COST;

  let cost = ATTACK_COST;

  // Reclamation: -50g if previously owned by attacker
  if (player.tech?.[3] && targetHex.previousOwner === player.id) {
    cost -= RECLAMATION_ATTACK_COST;
  }

  // Vanguard: -50g if timer active (within 12s of last capture)
  if (player.tech?.[4] && player.vanguardTimer > 0) {
    cost -= VANGUARD_ATTACK_COST;
  }

  // Both techs stack: 100 - 50 - 50 = 0 (free attack when reclaiming during Vanguard)
  return Math.max(0, cost);
}

function countAdjacentOwned(hex: HexDTO, state: GameState, owner: number): number {
  let count = 0;
  for (const n of neighbors({ q: hex.q, r: hex.r })) {
    const key = `${n.q},${n.r}`;
    const hs = state.hexes.get(key);
    if (hs && hs.owner === owner) {
      count++;
    }
  }
  return count;
}

function upgradeLabel(building: number, currentLevel: number, player?: PlayerDTO | null): string {
  const name = BUILDING_NAMES[building] ?? '?';
  const targetLevel = currentLevel + 1;
  const cost = upgradeCost(building, currentLevel);
  const delta = upgradeDelta(building, currentLevel, player);
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
        case 'drop-hex': this.callbacks.onDropHex(); break;
        case 'fortify': this.callbacks.onFortify(); break;
      }
    });
  }

  updateWithActions(hex: HexDTO | null, gold: number, isOwn: boolean, isEnemy: boolean, attackerPower = 0, battle: BattleDTO | null = null, player: PlayerDTO | null = null, state: GameState | null = null) {
    if (!hex) {
      this.el.style.display = 'none';
      this.lastKey = '';
      return;
    }

    const defPower = this.calcPower(hex, player, state ?? undefined);
    // Garrison activates immediately when battle starts, so attacker must exceed defPower + garrison
    let garrisonBonus = 0;
    if (isEnemy && state && hex.owner > 0) {
      const ownerPlayer = state.players.get(String(hex.owner));
      if (ownerPlayer?.tech?.[5]) {
        garrisonBonus = Math.min(GARRISON_MAX_BOOST, countAdjacentOwned(hex, state, hex.owner));
      }
    }
    const effectiveDefPower = defPower + garrisonBonus;
    const upgCost = upgradeCost(hex.building, hex.level);
    const hasBattle = battle !== null;
    const effectiveCost = effectiveAttackCost(player, hex);
    const canAttack = !hasBattle && gold >= effectiveCost && attackerPower > effectiveDefPower;
    const key = [
      hex.q, hex.r, hex.building, hex.level, isOwn, isEnemy,
      gold >= GOLD_BUILD_COST, gold >= RESEARCH_BUILD_COST, gold >= upgCost,
      canAttack, attackerPower, defPower,
      battle ? Math.floor(battle.timeLeft) : -1,
      hex.fortifyTimer > 0 ? 1 : 0, player?.tech?.[1] ? 1 : 0,
    ].join('|');

    if (key === this.lastKey) {
      this.el.style.display = 'flex';
      return;
    }
    this.lastKey = key;

    const ownerPlayer = state?.players.get(String(hex.owner));
    let powerStr = `Pwr:${defPower}`;
    // Garrison: shown as defense-only note, not part of static power total
    if (state && hex.owner > 0 && ownerPlayer?.tech?.[5]) {
      const garrisonBonus = Math.min(GARRISON_MAX_BOOST, countAdjacentOwned(hex, state, hex.owner));
      if (garrisonBonus > 0) powerStr += ` +${garrisonBonus} def`;
    }
    let html = `<span style="margin-right:4px">[${hex.q},${hex.r}] ${powerStr}</span>`;

    if (isOwn) {
      const hasBuilding = hex.building !== 0;
      const refund = hasBuilding ? demolishRefund(hex.building, hex.level) : 0;

      if (!hasBuilding) {
        html += this.makeBtn(upgradeLabel(BUILDING_GOLD, 0, player), gold >= GOLD_BUILD_COST, 'upgrade-gold');
        html += this.makeBtn(upgradeLabel(BUILDING_POWER, 0, player), gold >= POWER_BUILD_COST, 'upgrade-power');
        html += this.makeBtn(upgradeLabel(BUILDING_RESEARCH, 0, player), gold >= RESEARCH_BUILD_COST, 'upgrade-research');
      } else {
        html += this.makeBtn(upgradeLabel(hex.building, hex.level, player), gold >= upgCost, 'upgrade');
      }

      html += `<span style="border-left:1px solid #555;height:20px;margin:0 8px;display:inline-block;vertical-align:middle"></span>`;
      html += this.makeBtn(`🗑${hasBuilding ? ` (+${refund}g)` : ''}`, hasBuilding, 'demolish');

      // Sell hex: voluntary drop at any time (no battle, non-capital)
      if (!hex.capital && !battle) {
        const sellRefund = hasBuilding ? demolishRefund(hex.building, hex.level) : 0;
        html += this.makeBtn(`Sell hex${sellRefund > 0 ? ` (+${sellRefund}g)` : ''}`, true, 'drop-hex');
      }

      // Fortify: available when Fortify tech is owned and hex is not already fortified
      if (player?.tech?.[1] && hex.fortifyTimer <= 0 && !battle) {
        html += this.makeBtn(`Fortify (${FORTIFY_COST}g)`, gold >= FORTIFY_COST, 'fortify');
      }
    } else if (isEnemy) {
      html += this.makeBtn(`Attack (${effectiveCost}g)`, canAttack, 'attack');
      if (hasBattle) {
        html += `<span style="color:#fa0;margin-left:4px">Battle in progress</span>`;
      } else if (attackerPower <= effectiveDefPower) {
        html += `<span style="color:#f66;margin-left:4px">Need Pwr > ${effectiveDefPower}</span>`;
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

  private calcPower(hex: HexDTO, _player?: PlayerDTO | null, state?: GameState): number {
    let p = 0;
    if (hex.capital) p = CAPITAL_POWER;
    if (hex.building === BUILDING_POWER) p += hex.level;
    // Iron Grip: use HEX OWNER's tech — server applies it per-owner in hexEffectivePower()
    // Garrison NOT included: server only applies it in resolveBattle, not ValidateAttack
    if (state && hex.owner > 0) {
      const ownerPlayer = state.players.get(String(hex.owner));
      if (ownerPlayer?.tech?.[9]) p++; // TechIronGrip
    }
    return p;
  }

  hide() {
    this.el.style.display = 'none';
    this.lastKey = '';
  }
}
