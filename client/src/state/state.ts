export interface HexDTO {
  q: number;
  r: number;
  owner: number;
  building: number;
  level: number;
  capital: boolean;
  fortifyTimer: number;
  previousOwner: number;
}

export interface PlayerDTO {
  id: number;
  gold: number;
  tp: number;
  tech: boolean[];
  autoDropActive: boolean;
  autoDropGrace: number;
  vanguardTimer: number;
}

export interface BattleDTO {
  aq: number;
  ar: number;
  dq: number;
  dr: number;
  timeLeft: number;
  attacker: number;
  defender: number;
  counterBoost: number;
}

export interface SnapshotMsg {
  type: 'snapshot';
  hexes: Record<string, HexDTO>;
  players: Record<string, PlayerDTO>;
  battles: BattleDTO[];
  elapsed: number;
  over: boolean;
  winner: number;
  waiting: boolean;
}

export interface DeltaMsg {
  type: 'delta';
  players: PlayerDTO[];
  battles: BattleDTO[];
  elapsed: number;
  over: boolean;
  winner: number;
  waiting: boolean;
  hexChanges?: HexDTO[];
}

export interface GameState {
  hexes: Map<string, HexDTO>;
  players: Map<string, PlayerDTO>;
  battles: BattleDTO[];
  elapsed: number;
  over: boolean;
  winner: number;
  waiting: boolean;
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

  return { hexes, players, battles: msg.battles || [], elapsed: msg.elapsed, over: msg.over ?? false, winner: msg.winner ?? 0, waiting: msg.waiting ?? false };
}

export function applyDelta(state: GameState, msg: DeltaMsg): GameState {
  const players = new Map(state.players);
  for (const p of msg.players) {
    players.set(String(p.id), p);
  }

  const hexes = new Map(state.hexes);
  for (const hex of (msg.hexChanges ?? [])) {
    hexes.set(`${hex.q},${hex.r}`, hex);
  }

  return { hexes, players, battles: msg.battles, elapsed: msg.elapsed, over: msg.over, winner: msg.winner, waiting: msg.waiting };
}
