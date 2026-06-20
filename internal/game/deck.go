package game

// InitStartingDeck returns the standard 7-card starting deck.
func InitStartingDeck() Deck {
	cards := []Card{
		{Type: CardHexClaim},
		{Type: CardHexClaim},
		{Type: CardBarBurst},
		{Type: CardBarBurst},
		{Type: CardBarSpeed},
		{Type: CardBarSpeed},
		{Type: CardBasicSoldier},
	}
	return Deck{Cards: cards}
}

// CardCost returns the bar cost of a card type.
func CardCost(ct CardType) float64 {
	switch ct {
	case CardHexClaim:
		return CostHexClaim
	case CardBarBurst:
		return CostBarBurst
	case CardBarSpeed:
		return CostBarSpeed
	case CardBasicSoldier:
		return CostBasicSoldier
	default:
		return 0
	}
}

func RunDeck(state *GameState, dt float64) {
	for _, p := range state.Players {
		if len(p.Deck.Cards) == 0 {
			continue
		}
		top := p.Deck.Cards[0]
		if isCardPlayable(state, p, top) {
			p.Deck.AutoTimer = 0
			continue
		}
		p.Deck.AutoTimer += dt
		if p.Deck.AutoTimer >= DeckAutoCycleTimer {
			p.Deck.AutoTimer = 0
			cycleTo := top
			p.Deck.Cards = append(p.Deck.Cards[1:], cycleTo)
		}
	}
}

func isCardPlayable(state *GameState, p *Player, c Card) bool {
	if p.Bar < CardCost(c.Type) {
		return false
	}
	if c.Type == CardHexClaim {
		return hasFreeAdjacentHex(state, p.ID)
	}
	return true
}

func hasFreeAdjacentHex(state *GameState, pid PlayerID) bool {
	for hex, hs := range state.Hexes {
		if hs.Owner != NoPlayer {
			continue
		}
		for _, n := range hex.Neighbors() {
			if nhs, ok := state.Hexes[n]; ok && nhs.Owner == pid {
				return true
			}
		}
	}
	return false
}

func applyPlayTopCard(state *GameState, p *Player) {
	if len(p.Deck.Cards) == 0 {
		return
	}
	top := p.Deck.Cards[0]
	cost := CardCost(top.Type)
	if p.Bar < cost {
		return
	}
	if top.Type == CardHexClaim {
		// Hex Claim requires a follow-up ActionClaimHex to resolve; just deduct cost here.
		// The actual hex assignment happens in applyClaimHex.
		p.Bar -= cost
		cycleTopCard(p)
		return
	}
	p.Bar -= cost
	resolveCardEffect(state, p, top)
	cycleTopCard(p)
}

func applyPlayAsideCard(state *GameState, p *Player) {
	if p.Deck.AsideCard == nil {
		return
	}
	c := *p.Deck.AsideCard
	cost := CardCost(c.Type)
	if p.Bar < cost {
		return
	}
	// Hex Claim from aside slot: deduct cost, return card to bottom (no hex target yet).
	p.Bar -= cost
	p.Deck.Cards = append(p.Deck.Cards, c)
	p.Deck.AsideCard = nil
}

func applyPushAside(p *Player) {
	if p.Deck.AsideCard != nil {
		return
	}
	if len(p.Deck.Cards) == 0 {
		return
	}
	top := p.Deck.Cards[0]
	aside := top
	p.Deck.AsideCard = &aside
	p.Deck.Cards = p.Deck.Cards[1:]
	p.Deck.AutoTimer = 0
}

func applyClaimHex(state *GameState, p *Player, target Hex) {
	hs, ok := state.Hexes[target]
	if !ok || hs.Owner != NoPlayer {
		return
	}
	if !hasAdjacentOwned(state, p.ID, target) {
		return
	}
	hs.Owner = p.ID
}

func applyBuyCard(state *GameState, p *Player, ct CardType) {
	entry := findShopEntry(state, ct)
	if entry == nil || entry.Quantity <= 0 {
		return
	}
	if p.Bar < entry.BarCost {
		return
	}
	p.Bar -= entry.BarCost
	entry.Quantity--
	p.Deck.Cards = append(p.Deck.Cards, Card{Type: ct})
}

func applyRemoveTopCard(state *GameState, p *Player) {
	if len(p.Deck.Cards) == 0 {
		return
	}
	if p.Bar < CardRemovalCost {
		return
	}
	p.Bar -= CardRemovalCost
	p.Deck.Cards = p.Deck.Cards[1:]
}

func applyRemoveAsideCard(state *GameState, p *Player) {
	if p.Deck.AsideCard == nil {
		return
	}
	if p.Bar < CardRemovalCost {
		return
	}
	p.Bar -= CardRemovalCost
	p.Deck.AsideCard = nil
}

func resolveCardEffect(state *GameState, p *Player, c Card) {
	switch c.Type {
	case CardBarBurst:
		p.Bar += 1.0
		if p.Bar > BarMax {
			p.Bar = BarMax
		}
	case CardBarSpeed:
		p.BarBoosts = append(p.BarBoosts, BarBoost{
			Bonus:    BarSpeedBonus,
			TimeLeft: BarSpeedDuration,
		})
	case CardBasicSoldier:
		spawnUnit(state, p.ID, UnitBasicSoldier)
	}
}

func cycleTopCard(p *Player) {
	if len(p.Deck.Cards) == 0 {
		return
	}
	top := p.Deck.Cards[0]
	p.Deck.Cards = append(p.Deck.Cards[1:], top)
	p.Deck.AutoTimer = 0
}

func findShopEntry(state *GameState, ct CardType) *ShopEntry {
	for i := range state.Shop {
		if state.Shop[i].CardType == ct {
			return &state.Shop[i]
		}
	}
	return nil
}
