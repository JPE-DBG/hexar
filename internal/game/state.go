package game

type PlayerID int

const (
	NoPlayer PlayerID = 0
	Player1  PlayerID = 1
	Player2  PlayerID = 2
)

type CardType int

const (
	CardHexClaim     CardType = iota // 1 bar — claim adjacent hex
	CardBarBurst                     // 0 bar — +1 bar instantly
	CardBarSpeed                     // 1 bar — +15% fill rate for 15s
	CardBasicSoldier                 // 2 bar — deploy soldier
)

type Card struct {
	Type CardType
}

// BarBoost tracks an active +15% bar speed effect.
type BarBoost struct {
	Bonus    float64
	TimeLeft float64
}

// Deck is the player's card queue: top = index 0, bottom = last.
type Deck struct {
	Cards     []Card
	AsideCard *Card
	AutoTimer float64 // counts up; card auto-cycles when it exceeds DeckAutoCycleTimer
}

type Player struct {
	ID        PlayerID
	Bar       float64
	BarBoosts []BarBoost
	Deck      Deck
	CapitalHP int
}

type UnitType int

const (
	UnitBasicSoldier UnitType = iota
)

type Unit struct {
	ID          int
	Owner       PlayerID
	Type        UnitType
	Pos         Hex
	HP          int
	MoveTimer   float64 // counts up; unit moves when it reaches SoldierMovePeriod
	AttackTimer float64 // counts up; unit attacks when it reaches SoldierAttackPeriod
}

type HexState struct {
	Owner       PlayerID
	Capital     bool
	HasBuilding bool
}

type ShopEntry struct {
	CardType CardType
	BarCost  float64
	Quantity int
}

type GameState struct {
	Hexes      map[Hex]*HexState
	Players    map[PlayerID]*Player
	Units      []*Unit
	Shop       []ShopEntry
	NextUnitID int

	Elapsed       float64
	Over          bool
	Winner        PlayerID
	WinReason     string
	Waiting       bool
	Paused        bool
	PauseTimeLeft float64
}

func NewGameState() *GameState {
	return &GameState{
		Hexes:   make(map[Hex]*HexState),
		Players: make(map[PlayerID]*Player),
	}
}

func NewPlayer(id PlayerID) *Player {
	p := &Player{
		ID:        id,
		Bar:       0,
		CapitalHP: CapitalStartHP,
	}
	p.Deck = InitStartingDeck()
	return p
}

func opponentOf(id PlayerID) PlayerID {
	if id == Player1 {
		return Player2
	}
	return Player1
}
