package game

const (
	TickRate = 100 // milliseconds per tick
	TickDt   = 0.1 // seconds per tick

	// Bar
	BarFillRate = 0.1
	BarMax      = 10.0

	// Card bar costs
	CostHexClaim     = 1.0
	CostBarBurst     = 0.0
	CostBarSpeed     = 1.0
	CostBasicSoldier = 2.0

	// +15% Bar Speed card
	BarSpeedBonus    = 0.15
	BarSpeedDuration = 15.0

	// Deck
	DeckAutoCycleTimer = 5.0 // seconds before a stuck top card auto-cycles

	// Soldier stats
	SoldierHP           = 2
	SoldierAttackPower  = 1
	SoldierAttackPeriod = 1.0 // seconds between attacks
	SoldierMovePeriod   = 2.0 // seconds between hex moves (0.5 hex/sec)

	// Capital
	CapitalStartHP       = 20
	CapitalDamagePerUnit = 1

	// Shop
	CardRemovalCost = 5.0
)
