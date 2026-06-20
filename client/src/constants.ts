export const HEX_SIZE = 30;
export const MAP_RADIUS = 4;

export const CARD_TYPES = {
  UNIT: 'unit',
  BUILDING: 'building',
  BAR: 'bar',
} as const;

export type CardType = typeof CARD_TYPES[keyof typeof CARD_TYPES];
