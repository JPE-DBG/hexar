package game

const (
	TickRate = 100 // milliseconds per tick
	TickDt   = 0.1 // seconds per tick

	BaseIncomePerSec    = 2.0
	MaintenanceTier1    = 1.0 // hexes 1-10
	MaintenanceTier2    = 2.0 // hexes 11-20
	MaintenanceTier3    = 3.0 // hexes 21+
	MaintenanceTier1Cap = 10
	MaintenanceTier2Cap = 20

	ClaimCost  = 10.0
	AttackCost = 100.0

	EconomyBuildCost  = 60.0
	DefenseBuildCost  = 60.0
	ResearchBuildCost = 80.0

	EconomyBonusMultiplier = 1.5
	EconomyPerLevel        = 0.6
	ResearchPerLevel       = 0.1
	DefensePerLevel        = 1

	DemolishRefund = 0.5

	CounterSpendCostPerSec = 50.0
	CounterSpendCap        = 3

	InstantTakeoverMinDiff = 3 // power diff must exceed this for instant takeover

	ConquestThreshold    = 0.60
	ConquestHoldTime     = 10.0
	TechDomThreshold     = 0.35
	TechDomHoldTime      = 10.0
	TechDomRequiredLevel = 4

	AutoDropGracePeriod = 10.0

	GameDuration = 30 * 60 // seconds

	TechCostIronGrip          = 50.0
	TechCostProductionBoom    = 40.0
	TechCostEfficientConquest = 35.0
	TechCostGarrison          = 60.0

	EfficientConquestAttackCost = 75.0
	ProductionBoomBonus         = 0.30
	IronGripPowerBonus          = 1
	GarrisonMaxBonus            = 3
)
