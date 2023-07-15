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
		rt := make(map[Terrain]int)
		for _, t := range StandardTerrains {
			rt[t] = 0
		}
		for t := range icp.revealed[3].Iter() {
			rt[t.terrain]++
			rt[t.terrain2]++
		}
		rem := 6 - icp.revealed[3].Cardinality()
		for t := range pcts {
			// Each terrain appears on only 3 Stage III cards
			if rt[t] == 3 {
				delete(pcts, t)
				continue
			}
			pcts[t] = float64(3-rt[t]) / float64(rem)
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
func (icp *InvaderCardpool) Reveal(ic InvaderCard) error {
	if 1 > ic.stage || ic.stage > 3 {
		return fmt.Errorf(
			"ErrInvalidInvaderCard %w : %d is not a stage, expected I, II, or III",
			ErrInvalidInvaderCard,
			ic.stage,
		)
	}

	icp.revealed[ic.stage].Add(ic)
	return nil
}
