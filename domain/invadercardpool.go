package domain

import (
	"fmt"

	mapset "github.com/deckarep/golang-set/v2"
)

// InvaderCardpool is used for terrain predictions over the course of a game.
type InvaderCardpool struct {
	revealed map[int]mapset.Set[InvaderCard]
}

// NewInvaderCardpool initializes a new pool with no revealed cards.
// Scotland and Habsburg Mining expedition are a special cases
// where the "Coastal Lands" terrain card starts revealed.
func NewInvaderCardpool(scotland bool, habsmine bool) *InvaderCardpool {
	icp := &InvaderCardpool{
		revealed: map[int]mapset.Set[InvaderCard]{
			1: mapset.NewSetWithSize[InvaderCard](4),
			2: mapset.NewSetWithSize[InvaderCard](5),
			3: mapset.NewSetWithSize[InvaderCard](6),
		},
	}
	if scotland || habsmine {
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
		for t := range icp.revealed[stage].Iter() {
			delete(pcts, t.terrain)
		}
		for t := range pcts {
			pcts[t] = 1.0 / float64(len(pcts))
		}
	case 3:
		revt := make(map[Terrain]int)
		for _, t := range StandardTerrains {
			revt[t] = 0
		}
		for t := range icp.revealed[3].Iter() {
			revt[t.terrain]++
			revt[t.terrain2]++
		}
		rem := 6 - icp.revealed[3].Cardinality()
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
func (icp *InvaderCardpool) Reveal(icard InvaderCard) error {
	if 1 > icard.stage || icard.stage > 3 {
		return fmt.Errorf(
			"ErrInvalidInvaderCard %w : %d is not a stage, expected I, II, or III",
			ErrInvalidInvaderCard,
			icard.stage,
		)
	}

	icp.revealed[icard.stage].Add(icard)

	return nil
}
