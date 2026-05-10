export const HEX_SIZE = 30;
export const TICK_RATE = 100;

export const COLORS = {
  // Canvas colors — used by renderer
  background: '#16213e',
  unclaimed:  '#2a2a4a',
  player1:    '#45b7d1',
  player2:    '#e74c3c',
  grid:       '#0f3460',
  capital:    '#f9ca24',
  battle:     '#f0932b',  // orange-amber, distinct from capital gold

  // DOM colors — applied programmatically from TS, never hardcoded in CSS
  pause:      '#6c3483',
  reconnect:  '#c0392b',
  warning:    '#f39c12',
  accent:     '#45b7d1',  // same as player1 — single interactive accent
  textMuted:  '#999999',
  panelBg:    'rgba(22,33,62,0.95)',
  overlayBg:  'rgba(14,20,40,0.88)',
} as const;

// Economy constants — mirrors internal/game/constants.go
export const BASE_INCOME_PER_SEC = 2.0;

export const GOLD_BUILD_COST = 60;
export const GOLD_PER_LEVEL = 0.6;
export const GOLD_BONUS_MULTIPLIER = 1.5;

export const POWER_BUILD_COST = 60;
export const POWER_PER_LEVEL = 1;

export const RESEARCH_BUILD_COST = 80;
export const RESEARCH_PER_LEVEL = 0.2;

export const DEMOLISH_REFUND = 0.5;
export const AUTO_DROP_REFUND = 0.5;

export const MAINTENANCE_TIER1 = 1.0;
export const MAINTENANCE_TIER2 = 2.0;
export const MAINTENANCE_TIER3 = 3.0;
export const MAINTENANCE_TIER1_CAP = 10;
export const MAINTENANCE_TIER2_CAP = 20;

export const CLAIM_COST = 10;
export const ATTACK_COST = 100;
export const CAPITAL_POWER = 1;

export const COUNTER_SPEND_COST = 50;
export const COUNTER_SPEND_CAP = 3;

export const RECLAMATION_ATTACK_COST = 50;
export const VANGUARD_ATTACK_COST = 50;
export const FORTIFY_COST = 40;
export const FORTIFY_DURATION = 90;

// Tech bonus constants (mirrors internal/game/constants.go)
export const PROSPERITY_BONUS = 1.0;
export const COMPOUND_GROWTH_MULTIPLIER = 1.25;
export const SUPPLY_LINES_TIER1 = 0.9;
export const SUPPLY_LINES_TIER2 = 1.8;
export const SUPPLY_LINES_TIER3 = 2.7;
export const GARRISON_MAX_BOOST = 2;

export const BUILDING_GOLD = 1;
export const BUILDING_POWER = 2;
export const BUILDING_RESEARCH = 3;

// Tech IDs — must match TechID iota order in internal/game/state.go
export const TECH_BLITZ = 0;
export const TECH_FORTIFY = 1;
export const TECH_PROSPERITY = 2;
export const TECH_RECLAMATION = 3;
export const TECH_VANGUARD = 4;
export const TECH_GARRISON = 5;
export const TECH_SUPPLY_LINES = 6;
export const TECH_WAR_CHEST = 7;
export const TECH_RESILIENCE = 8;
export const TECH_IRON_GRIP = 9;
export const TECH_COMPOUND_GROWTH = 10;
export const TECH_SIEGE_MASTERY = 11;

export interface TechDef {
  id: number;
  name: string;
  cost: number;
  description: string;
  archetype: 'Aggressor' | 'Defender' | 'Builder' | 'Territorial';
}

export const TECH_DEFS: TechDef[] = [
  { id: 0,  name: 'Blitz',          cost: 20, description: 'Unclaimed hex claims cost 0g',                       archetype: 'Aggressor'   },
  { id: 1,  name: 'Fortify',        cost: 20, description: 'Spend 40g to prevent instant-takeover for 90s',      archetype: 'Defender'    },
  { id: 2,  name: 'Prosperity',     cost: 25, description: 'Economy buildings +1/sec additional',                archetype: 'Builder'     },
  { id: 3,  name: 'Reclamation',    cost: 25, description: 'Recapture previously-owned hex costs 50g',           archetype: 'Territorial' },
  { id: 4,  name: 'Vanguard',       cost: 30, description: 'After capture, next attack within 12s costs 50g',    archetype: 'Aggressor'   },
  { id: 5,  name: 'Garrison',       cost: 30, description: 'Adjacent owned hexes +1 Power in defense (cap +2)',  archetype: 'Defender'    },
  { id: 6,  name: 'Supply Lines',   cost: 40, description: 'Maintenance ×0.9/1.8/2.7 per tier',                  archetype: 'Builder'     },
  { id: 7,  name: 'War Chest',    cost: 30, description: 'Capture an enemy hex, recover 30g',              archetype: 'Territorial' },
  { id: 8,  name: 'Resilience',     cost: 45, description: 'Auto-drop grace 10s → 20s; drop refund 70%',         archetype: 'Defender'    },
  { id: 9,  name: 'Iron Grip',      cost: 55, description: 'All owned hexes permanently +1 Power',               archetype: 'Aggressor'   },
  { id: 10, name: 'Compound Growth',cost: 65, description: 'Economy buildings ×1.25',                            archetype: 'Builder'     },
  { id: 11, name: 'Siege Mastery',  cost: 75, description: 'Attack timers -40% (min 3s); ties → attacker wins',  archetype: 'Aggressor'   },
];
