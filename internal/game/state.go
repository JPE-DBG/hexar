package game

type PlayerID int

const NoPlayer PlayerID = 0

type BuildingType int

const (
	BuildingNone     BuildingType = 0
	BuildingEconomy  BuildingType = 1
	BuildingDefense  BuildingType = 2
	BuildingResearch BuildingType = 3
)

type HexState struct {
	Owner    PlayerID     `json:"owner"`
	Building BuildingType `json:"building"`
	Level    int          `json:"level"`
	Capital  bool         `json:"capital"`
}

func (h *HexState) Power() int {
	p := 0
	if h.Capital {
		p = 1
	}
	if h.Building == BuildingDefense {
		p += h.Level
	}
	return p
}

type TechID int

const (
	TechIronGrip TechID = iota
	TechProductionBoom
	TechEfficientConquest
	TechGarrison
	TechCount
)

type Player struct {
	ID   PlayerID        `json:"id"`
	Gold float64         `json:"gold"`
	TP   float64         `json:"tp"`
	Tech [TechCount]bool `json:"tech"`
}

func (p *Player) TechLevel() int {
	count := 0
	for _, unlocked := range p.Tech {
		if unlocked {
			count++
		}
	}
	return count
}

type Battle struct {
	AttackerHex Hex      `json:"attackerHex"`
	DefenderHex Hex      `json:"defenderHex"`
	Attacker    PlayerID `json:"attacker"`
	Defender    PlayerID `json:"defender"`
	TimeLeft    float64  `json:"timeLeft"`
}

type GameState struct {
	Hexes   map[Hex]*HexState    `json:"hexes"`
	Players map[PlayerID]*Player `json:"players"`
	Battles []Battle             `json:"battles"`
	Elapsed float64              `json:"elapsed"`
	Over    bool                 `json:"over"`
	Winner  PlayerID             `json:"winner"`
}

func NewGameState() *GameState {
	return &GameState{
		Hexes:   make(map[Hex]*HexState),
		Players: make(map[PlayerID]*Player),
	}
}
