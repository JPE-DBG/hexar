export interface HexDTO {
  q: number;
  r: number;
  owner: number;
  building: number;
  level: number;
  capital: boolean;
}

export interface PlayerDTO {
  id: number;
  gold: number;
  tp: number;
}

export interface SnapshotMsg {
  type: 'snapshot';
  hexes: Record<string, HexDTO>;
  players: Record<string, PlayerDTO>;
  elapsed: number;
}

export interface GameState {
  hexes: Map<string, HexDTO>;
  players: Map<string, PlayerDTO>;
  elapsed: number;
}

export function applySnapshot(msg: SnapshotMsg): GameState {
  const hexes = new Map<string, HexDTO>();
  for (const [key, hex] of Object.entries(msg.hexes)) {
    hexes.set(key, hex);
  }

  const players = new Map<string, PlayerDTO>();
  for (const [key, player] of Object.entries(msg.players)) {
    players.set(key, player);
  }

  return { hexes, players, elapsed: msg.elapsed };
}
