package domain

import (
	"fmt"
)

// InvaderCardpool is used for terrain predictions over the course of a game.
type InvaderCardpool struct {
	revealed map[int][]InvaderCard
}

// NewInvaderCardpool initializes a new pool with no revealed cards.
// Scotland and Habsburg Mining expedition are a special cases
// where the "Coastal Lands" terrain card starts revealed.
func NewInvaderCardpool(scotland bool, habsmine bool) *InvaderCardpool {
	icp := &InvaderCardpool{
		revealed: map[int][]InvaderCard{
			1: make([]InvaderCard, 0, 4),
			2: make([]InvaderCard, 0, 5),
			3: make([]InvaderCard, 0, 6),
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
		for _, t := range icp.revealed[stage] {
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
		for _, t := range icp.revealed[3] {
			rt[t.terrain] = rt[t.terrain] + 1
			rt[t.terrain2] = rt[t.terrain2] + 1
		}
		rem := 6 - len(icp.revealed[3])
		for t := range pcts {
			// Each terrain appears on only 3 Stage III cards
			if rt[t] == 3 {
				delete(pcts, t)
				continue
			}
			pcts[t] = float32(rem) / float32(3-rt[t])
		}
	} else {
		return nil, fmt.Errorf("%q is not a stage, expected I, II, or III", stage)
	}

	return pcts, nil
}

// Reveal excludes cards from future predictions
func (ip *InvaderCardpool) Reveal(ic InvaderCard) {
	// DANGER! check the array index
	ip.revealed[ic.stage] = append(ip.revealed[ic.stage], ic)
}

// InvaderCard has a phase and one or more terrain types
type InvaderCard struct {
	stage    int
	terrain  Terrain
	terrain2 Terrain
}

// All possible invader cards
var (
	// Stage I Invader Card w/ Jungle Terrain
	StageOneJungle = InvaderCard{1, Jungle, None}
	// Stage I Invader Card w/ Sands Terrain
	StageOneSands = InvaderCard{1, Sands, None}
	// Stage I Invader Card w/ Mountain Terrain
	StageOneMountain = InvaderCard{1, Mountain, None}
	// Stage I Invader Card w/ Wetland Terrain
	StageOneWetland = InvaderCard{1, Wetland, None}
	// Stage II Invader Card w/ Jungle Terrain
	StageTwoJungle = InvaderCard{2, Jungle, None}
	// Stage II Invader Card w/ Sands Terrain
	StageTwoSands = InvaderCard{2, Sands, None}
	// Stage II Invader Card w/ Mountain Terrain
	StageTwoMountain = InvaderCard{2, Mountain, None}
	// Stage II Invader Card w/ Wetland Terrain
	StageTwoWetland = InvaderCard{2, Wetland, None}
	// Unique Stage II Invader Card w/ "Coastal Lands" Terrain
	StageTwoCoastal = InvaderCard{2, CoastalLands, None}
	// Stage III Invader Card w/ Jungle + Sands Terrain
	StageThreeJungleSands = InvaderCard{3, Jungle, Sands}
	// Stage III Invader Card w/ Jungle + Mountain Terrain
	StageThreeJungleMountain = InvaderCard{3, Jungle, Mountain}
	// Stage III Invader Card w/ Jungle + Wetland Terrain
	StageThreeJungleWetland = InvaderCard{3, Jungle, Wetland}
	// Stage III Invader Card w/ Sands + Mountain Terrain
	StageThreeSandsMountain = InvaderCard{3, Sands, Mountain}
	// Stage III Invader Card w/ Sands + Wetland Terrain
	StageThreeSandsWetland = InvaderCard{3, Sands, Wetland}
	// Stage III Invader Card w/ Mountain + Wetland Terrain
	StageThreeMountainWetland = InvaderCard{3, Mountain, Wetland}
)
