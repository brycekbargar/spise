package domain

import (
	"fmt"

	mapset "github.com/deckarep/golang-set/v2"
)

// InvaderCardpool is used for terrain predictions over the course of a game.
type InvaderCardpool struct {
	Revealed map[int]mapset.Set[InvaderCard]
}

// NewInvaderCardpool initializes a new pool with no revealed cards.
func NewInvaderCardpool(game *Game) *InvaderCardpool {
	icp := &InvaderCardpool{
		Revealed: map[int]mapset.Set[InvaderCard]{
			1: mapset.NewSetWithSize[InvaderCard](4),
			2: mapset.NewSetWithSize[InvaderCard](5),
			3: mapset.NewSetWithSize[InvaderCard](6),
		},
	}
	if (game.LeadingAdversary == Scotland && game.LeadingAdversaryLevel >= 2) ||
		(game.SupportingAdversary == Scotland && game.SupportingAdversaryLevel >= 2) ||
		(game.LeadingAdversary == HabsburgMines && game.LeadingAdversaryLevel >= 4) ||
		(game.SupportingAdversary == HabsburgMines && game.SupportingAdversaryLevel >= 4) {

		err := icp.Reveal(StageTwoCoastal)
		if err != nil {
			panic(err)
		}
	}

	return icp
}

// Predict the probabilities of the next terrain for a given invader Stage.
// It is solely based on the invader cards that have been revealed.
func (icp InvaderCardpool) Predict(stage int) (map[Terrain]float64, error) {
	pcts := make(map[Terrain]float64)
	for _, t := range StandardTerrains {
		pcts[t] = 0.0
	}
	if stage == 2 {
		pcts[CoastalLands] = 0.0
	}

	switch stage {
	case 1:
		fallthrough
	case 2:
		for t := range icp.Revealed[stage].Iter() {
			delete(pcts, t.Terrain)
		}
		for t := range pcts {
			pcts[t] = 1.0 / float64(len(pcts))
		}
	case 3:
		revt := make(map[Terrain]int)
		for _, t := range StandardTerrains {
			revt[t] = 0
		}
		for t := range icp.Revealed[3].Iter() {
			revt[t.Terrain]++
			revt[t.Terrain2]++
		}
		rem := 6 - icp.Revealed[3].Cardinality()
		for trn := range pcts {
			// Each terrain appears on only 3 Stage III cards
			if revt[trn] == 3 {
				delete(pcts, trn)

				continue
			}
			pcts[trn] = float64(3-revt[trn]) / float64(rem)
		}
	default:
		return nil, fmt.Errorf(
			"ErrInvalidInvaderCard %w : %d is not a stage, expected I, II, or III",
			ErrInvalidInvaderCard,
			stage,
		)
	}

	return pcts, nil
}

// Reveal excludes cards from future predictions.
func (icp *InvaderCardpool) Reveal(card InvaderCard) error {
	for _, ic := range AllInvaderCards {
		if ic == card {
			icp.Revealed[card.Stage].Add(card)

			return nil
		}
	}

	return ErrInvalidInvaderCard
}
