export const HEX_SIZE = 30;
export const TICK_RATE = 100;

export const COLORS = {
  unclaimed: '#2d2d44',
  player1: '#4ecdc4',
  player2: '#ff6b6b',
  grid: '#3d3d5c',
  background: '#1a1a2e',
  capital: '#ffd93d',
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

export const BUILDING_GOLD = 1;
export const BUILDING_POWER = 2;
export const BUILDING_RESEARCH = 3;

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
