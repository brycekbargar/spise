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
		icp.Reveal(StageTwoCoastal)
	}
	return icp
}

// Predict the probabilities of the next terrain for a given invader Stage.
// It is solely based on the invader cards that have been revealed.
func (icp InvaderCardpool) Predict(stage int) (map[Terrain]float32, error) {
	pcts := make(map[Terrain]float32)
	for _, t := range StandardTerrains {
		pcts[t] = 0.0
	}
	if stage == 2 {
		pcts[CoastalLands] = 0.0
	}

	if stage == 1 || stage == 2 {
		for t := range icp.revealed[stage].Iter() {
			delete(pcts, t.terrain)
		}
		for t := range pcts {
			pcts[t] = 100.0 / float32(len(pcts))
		}
	} else if stage == 3 {
		rt := make(map[Terrain]int)
		for _, t := range StandardTerrains {
			rt[t] = 0
		}
		for t := range icp.revealed[3].Iter() {
			rt[t.terrain] = rt[t.terrain] + 1
			rt[t.terrain2] = rt[t.terrain2] + 1
		}
		rem := 6 - icp.revealed[3].Cardinality()
		for t := range pcts {
			// Each terrain appears on only 3 Stage III cards
			if rt[t] == 3 {
				delete(pcts, t)
				continue
			}
			pcts[t] = float32(rem) / float32(3-rt[t])
		}
	} else {
		return nil, fmt.Errorf("%d is not a stage, expected I, II, or III", stage)
	}

	return pcts, nil
}

// Reveal excludes cards from future predictions
func (icp *InvaderCardpool) Reveal(ic InvaderCard) error {
	if 1 > ic.stage || ic.stage > 3 {
		return fmt.Errorf("%d is not a stage, expected I, II, or III", ic.stage)
	}

	icp.revealed[ic.stage].Add(ic)
	return nil
}
