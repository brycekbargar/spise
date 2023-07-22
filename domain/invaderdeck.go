package domain

type InvaderDeck struct {
	game *Game

	Drawn  []InvaderCardDrawn
	InDeck []InvaderCardInDeck
}

func NewInvaderDeck(game *Game) *InvaderDeck {
	// var initial []InvaderCard

	return &InvaderDeck{
		game:   game,
		Drawn:  []InvaderCardDrawn{},
		InDeck: []InvaderCardInDeck{},
	}
}

// InvaderCardInDeck wraps possibilites for the invader deck.
type InvaderCardInDeck struct {
	Card InvaderCard
}

// InvaderCardDrawn wraps drawn invader cards.
type InvaderCardDrawn struct {
	Card InvaderCard
}

// Special Invader Cards.
var (
	// Special Stage II Invader Card for Habsburg Mining Expedition.
	StageTwoSaltDeposits = InvaderCard{2, "salt-deposits", UnknownTerrain}
	// Stage I Invader Card with unknown terrain.
	StageOneUnknown = InvaderCard{1, UnknownTerrain, UnknownTerrain}
	// Stage II Invader Card with unknown terrain.
	StageTwoUnknown = InvaderCard{2, UnknownTerrain, UnknownTerrain}
	// Stage III Invader Card with unknown terrain.
	StageThreeUnknown = InvaderCard{3, UnknownTerrain, UnknownTerrain}
)
