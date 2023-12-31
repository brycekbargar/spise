package domain

// Game orchestrates and owns the various state containers.
type Game struct {
	LeadingAdversary         Adversary
	LeadingAdversaryLevel    int
	SupportingAdversary      Adversary
	SupportingAdversaryLevel int
}

// Initialized Game is a domain.Game with initialized state containers.
type InitializedGame struct {
	*Game

	invadercardpool *InvaderCardpool
	invaderdeck     *InvaderDeck
}

// Init initialized the given game.
func (g *Game) Init() *InitializedGame {
	init := &InitializedGame{
		Game: g,

		invadercardpool: NewInvaderCardpool(g),
		invaderdeck:     NewInvaderDeck(g),
	}

	return init
}
