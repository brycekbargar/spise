package domain

// Terrain is the type of land
type Terrain string

const (
	// The land type is unknown
	None Terrain = ""
	// The jungle land type
	Jungle Terrain = "jungle"
	// The sands land type
	Sands Terrain = "sands"
	// The mountain land type
	Mountain Terrain = "mountain"
	// The wetland land type
	Wetland Terrain = "wetland"
	// The coastal land type (specifically Stage ii)
	CoastalLands Terrain = "coastal-lands"
)

// AllTerrains is the complete list of Terrain Types
var (
	AllTerrains = []Terrain{Jungle, Sands, Mountain, Wetland, CoastalLands}
	// StandardTerrains is the list of Terrain types common to all Phases
	StandardTerrains = AllTerrains[:4]
)