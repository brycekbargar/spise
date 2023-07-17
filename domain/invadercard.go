package domain

import "errors"

var ErrInvalidInvaderCard = errors.New("invalid invader card")

// InvaderCard has a phase and one or more terrain types.
type InvaderCard struct {
	stage    int
	terrain  Terrain
	terrain2 Terrain
}

// All possible invader cards.
var (
	// Stage I Invader Card w/ Jungle Terrain.
	StageOneJungle = InvaderCard{1, Jungle, UnknownTerrain}
	// Stage I Invader Card w/ Mountain Terrain.
	StageOneMountain = InvaderCard{1, Mountain, UnknownTerrain}
	// Stage I Invader Card w/ Sands Terrain.
	StageOneSands = InvaderCard{1, Sands, UnknownTerrain}
	// Stage I Invader Card w/ Wetland Terrain.
	StageOneWetland = InvaderCard{1, Wetland, UnknownTerrain}
	// Stage II Invader Card w/ Jungle Terrain.
	StageTwoJungle = InvaderCard{2, Jungle, UnknownTerrain}
	// Stage II Invader Card w/ Mountain Terrain.
	StageTwoMountain = InvaderCard{2, Mountain, UnknownTerrain}
	// Stage II Invader Card w/ Sands Terrain.
	StageTwoSands = InvaderCard{2, Sands, UnknownTerrain}
	// Stage II Invader Card w/ Wetland Terrain.
	StageTwoWetland = InvaderCard{2, Wetland, UnknownTerrain}
	// Unique Stage II Invader Card w/ "Coastal Lands" Terrain.
	StageTwoCoastal = InvaderCard{2, CoastalLands, UnknownTerrain}
	// Stage III Invader Card w/ Jungle + Mountain Terrain.
	StageThreeJungleMountain = InvaderCard{3, Jungle, Mountain}
	// Stage III Invader Card w/ Jungle + Sands Terrain.
	StageThreeJungleSands = InvaderCard{3, Jungle, Sands}
	// Stage III Invader Card w/ Jungle + Wetland Terrain.
	StageThreeJungleWetland = InvaderCard{3, Jungle, Wetland}
	// Stage III Invader Card w/ Sands + Mountain Terrain.
	StageThreeMountainSands = InvaderCard{3, Mountain, Sands}
	// Stage III Invader Card w/ Mountain + Wetland Terrain.
	StageThreeMountainWetland = InvaderCard{3, Mountain, Wetland}
	// Stage III Invader Card w/ Sands + Wetland Terrain.
	StageThreeSandsWetland = InvaderCard{3, Sands, Wetland}
)

var (
	// All Stage I Invader Cards.
	StageOneInvaderCards = []InvaderCard{
		StageOneJungle,
		StageOneMountain,
		StageOneSands,
		StageOneWetland,
	}
	// All Stage II Invader Cards.
	StageTwoInvaderCards = []InvaderCard{
		StageTwoJungle,
		StageTwoMountain,
		StageTwoSands,
		StageTwoWetland,
		StageTwoCoastal,
	}
	// All Stage III Invader Cards.
	StageThreeInvaderCards = []InvaderCard{
		StageThreeJungleMountain,
		StageThreeJungleSands,
		StageThreeJungleWetland,
		StageThreeMountainSands,
		StageThreeMountainWetland,
		StageThreeSandsWetland,
	}
	// All Invader Cards.
	AllInvaderCards = append(
		append(StageOneInvaderCards, StageTwoInvaderCards...),
		StageThreeInvaderCards...)
)
