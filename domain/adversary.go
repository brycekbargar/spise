package domain

// Adversaries are optional personalities for the faceless Invaders.
type Adversary string

const (
	// The adversary is unknown.
	UnknownAdversary Adversary = ""
	// The Kingdom of Brandenburg-Prussia.
	BrandenburgPrussia Adversary = "brandenburg-prussia"
	// The Kingdom of England.
	England Adversary = "england"
	// The Kingdom of France (Plantation Colony).
	France Adversary = "france-plantation-colony"
	// Habsburg Monarchy (Livestock Colony).
	HabsburgLivestock Adversary = "habsburg-livestock-colony"
	// Habsburg Mining Expedition.
	HabsburgMines Adversary = "habsburg-mining-expedition"
	// The Tsardom of Russia.
	Russia Adversary = "russia"
	// The Kingdom of Scotland.
	Sweden Adversary = "sweden"
)
