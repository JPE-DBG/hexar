package game

type PlayerID int

const NoPlayer PlayerID = 0

type BuildingType int

const (
	BuildingNone     BuildingType = 0
	BuildingGold     BuildingType = 1
	BuildingPower    BuildingType = 2
	BuildingResearch BuildingType = 3
)

type HexState struct {
	Owner         PlayerID     `json:"owner"`
	Building      BuildingType `json:"building"`
	Level         int          `json:"level"`
	Capital       bool         `json:"capital"`
	PreviousOwner PlayerID     `json:"previousOwner"`
	FortifyTimer  float64      `json:"fortifyTimer"`
	UpgradeTimer  float64      `json:"upgradeTimer"`
}

func (h *HexState) Power() int {
	p := 0
	if h.Capital {
		p = 1
	}
	if h.Building == BuildingPower {
		p += h.Level
	}
	return p
}

type TechID int

const (
	TechBlitz TechID = iota
	TechFortify
	TechProsperity
	TechReclamation
	TechVanguard
	TechGarrison
	TechSupplyLines
	TechWarChest
	TechResilience
	TechIronGrip
	TechCompoundGrowth
	TechSiegeMastery
	TechCount // = 12
)

type Player struct {
	ID             PlayerID        `json:"id"`
	Gold           float64         `json:"gold"`
	TP             float64         `json:"tp"`
	Tech           [TechCount]bool `json:"tech"`
	AutoDropGrace  float64         `json:"autoDropGrace"`
	AutoDropActive bool            `json:"autoDropActive"`
	VanguardTimer  float64         `json:"vanguardTimer"`
}

type Battle struct {
	AttackerHex  Hex      `json:"attackerHex"`
	DefenderHex  Hex      `json:"defenderHex"`
	Attacker     PlayerID `json:"attacker"`
	Defender     PlayerID `json:"defender"`
	TimeLeft     float64  `json:"timeLeft"`
	CounterBoost int      `json:"counterBoost"`
}

type GameState struct {
	Hexes         map[Hex]*HexState    `json:"hexes"`
	Players       map[PlayerID]*Player `json:"players"`
	Battles       []Battle             `json:"battles"`
	Elapsed       float64              `json:"elapsed"`
	Over          bool                 `json:"over"`
	Winner        PlayerID             `json:"winner"`
	WinReason     string               `json:"winReason"`
	Waiting       bool                 `json:"waiting"`
	Paused        bool                 `json:"paused"`
	PauseTimeLeft float64              `json:"pauseTimeLeft"`
}

func NewGameState() *GameState {
	return &GameState{
		Hexes:   make(map[Hex]*HexState),
		Players: make(map[PlayerID]*Player),
	}
}
