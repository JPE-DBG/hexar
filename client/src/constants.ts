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
export const RESEARCH_PER_LEVEL = 0.1;

export const DEMOLISH_REFUND = 0.5;

export const MAINTENANCE_TIER1 = 1.0;
export const MAINTENANCE_TIER2 = 2.0;
export const MAINTENANCE_TIER3 = 3.0;
export const MAINTENANCE_TIER1_CAP = 10;
export const MAINTENANCE_TIER2_CAP = 20;

export const CLAIM_COST = 10;
export const ATTACK_COST = 100;
export const CAPITAL_POWER = 1;

export const BUILDING_GOLD = 1;
export const BUILDING_POWER = 2;
export const BUILDING_RESEARCH = 3;
