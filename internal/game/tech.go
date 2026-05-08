package game

func TechCost(id TechID) float64 {
	switch id {
	case TechBlitz:
		return TechCostBlitz
	case TechFortify:
		return TechCostFortify
	case TechProsperity:
		return TechCostProsperity
	case TechReclamation:
		return TechCostReclamation
	case TechVanguard:
		return TechCostVanguard
	case TechGarrison:
		return TechCostGarrison
	case TechSupplyLines:
		return TechCostSupplyLines
	case TechWarChest:
		return TechCostWarChest
	case TechResilience:
		return TechCostResilience
	case TechIronGrip:
		return TechCostIronGrip
	case TechCompoundGrowth:
		return TechCostCompoundGrowth
	case TechSiegeMastery:
		return TechCostSiegeMastery
	default:
		return 0
	}
}
