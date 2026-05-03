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

	GoldBuildCost     = 60.0
	PowerBuildCost    = 60.0
	ResearchBuildCost = 80.0

	GoldBonusMultiplier = 1.5
	GoldPerLevel        = 0.6
	ResearchPerLevel    = 0.2
	PowerPerLevel       = 1

	DemolishRefund = 0.5
	AutoDropRefund = 0.5 // separate from DemolishRefund; Resilience tech changes this in M5

	CounterSpendCostPerSec = 50.0
	CounterSpendCap        = 3

	InstantTakeoverMinDiff = 3 // power diff must exceed this for instant takeover

	ConquestThreshold = 0.60
	ConquestHoldTime  = 10.0

	AutoDropGracePeriod = 10.0

	GameDuration = 30 * 60 // seconds

	TechCostBlitz          = 20.0
	TechCostFortify        = 20.0
	TechCostProsperity     = 25.0
	TechCostReclamation    = 25.0
	TechCostVanguard       = 30.0
	TechCostGarrison       = 30.0
	TechCostSupplyLines    = 40.0
	TechCostDominion       = 40.0
	TechCostResilience     = 45.0
	TechCostIronGrip       = 55.0
	TechCostCompoundGrowth = 65.0
	TechCostSiegeMastery   = 75.0
)
