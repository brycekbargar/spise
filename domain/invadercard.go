package domain

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
	// Stage I Invader Card w/ Mountain Terrain
	StageOneMountain = InvaderCard{1, Mountain, None}
	// Stage I Invader Card w/ Sands Terrain
	StageOneSands = InvaderCard{1, Sands, None}
	// Stage I Invader Card w/ Wetland Terrain
	StageOneWetland = InvaderCard{1, Wetland, None}
	// Stage II Invader Card w/ Jungle Terrain
	StageTwoJungle = InvaderCard{2, Jungle, None}
	// Stage II Invader Card w/ Mountain Terrain
	StageTwoMountain = InvaderCard{2, Mountain, None}
	// Stage II Invader Card w/ Sands Terrain
	StageTwoSands = InvaderCard{2, Sands, None}
	// Stage II Invader Card w/ Wetland Terrain
	StageTwoWetland = InvaderCard{2, Wetland, None}
	// Unique Stage II Invader Card w/ "Coastal Lands" Terrain
	StageTwoCoastal = InvaderCard{2, CoastalLands, None}
	// Stage III Invader Card w/ Jungle + Mountain Terrain
	StageThreeJungleMountain = InvaderCard{3, Jungle, Mountain}
	// Stage III Invader Card w/ Jungle + Sands Terrain
	StageThreeJungleSands = InvaderCard{3, Jungle, Sands}
	// Stage III Invader Card w/ Jungle + Wetland Terrain
	StageThreeJungleWetland = InvaderCard{3, Jungle, Wetland}
	// Stage III Invader Card w/ Sands + Mountain Terrain
	StageThreeMountainSands = InvaderCard{3, Mountain, Sands}
	// Stage III Invader Card w/ Mountain + Wetland Terrain
	StageThreeMountainWetland = InvaderCard{3, Mountain, Wetland}
	// Stage III Invader Card w/ Sands + Wetland Terrain
	StageThreeSandsWetland = InvaderCard{3, Sands, Wetland}
)

var (
	StageOneInvaderCards = []InvaderCard{
		StageOneJungle,
		StageOneMountain,
		StageOneSands,
		StageOneWetland,
	}
	StageTwoInvaderCards = []InvaderCard{
		StageTwoJungle,
		StageTwoMountain,
		StageTwoSands,
		StageTwoWetland,
		StageTwoCoastal,
	}
	StageThreeInvaderCards = []InvaderCard{
		StageThreeJungleMountain,
		StageThreeJungleSands,
		StageThreeJungleWetland,
		StageThreeMountainSands,
		StageThreeMountainWetland,
		StageThreeSandsWetland,
	}
	AllInvaderCards = append(
		append(StageOneInvaderCards, StageTwoInvaderCards...),
		StageThreeInvaderCards...)
)
