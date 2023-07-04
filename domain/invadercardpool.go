package domain

type Terrain int

const (
	// The jungle land type
	Jungle Terrain = 1 << iota
	// The sands land type
	Sands Terrain = 1 << iota
	// The mountain land type
	Mountain Terrain = 1 << iota
	// The wetland land type
	Wetland Terrain = 1 << iota
	// The costal land type (specifically Stage ii)
	Coastal Terrain = 1 << iota
)

// InvaderCardpool is the collection of all invader cards used for terrain predictions
type InvaderCardpool struct {
	revealed map[int][]Terrain
}

// NewInvaderCardpool initializes a new pool with no revealed cards.
// Scotland is a special case in that the Coastal terrain card starts revealed.
func NewInvaderCardpool(scotland bool) *InvaderCardpool {
	l2 := make([]Terrain, 0, 5)
	if scotland {
		l2 = append(l2, Coastal)
	}
	return &InvaderCardpool{
		revealed: map[int][]Terrain{
			1: make([]Terrain, 0, 4),
			2: l2,
			3: make([]Terrain, 0, 12),
		},
	}
}

// Predict determines the probability of terrain types for a given invader Stage.
// It is solely based on the invader cards that have been revealed.
func (ip InvaderCardpool) Predict(level int) (pcts map[Terrain]float32) {
	pcts = map[Terrain]float32{
		Jungle:   0,
		Mountain: 0,
		Sands:    0,
		Wetland:  0,
	}
	if level == 2 {
		pcts[Coastal] = 0.0
	}

	if level != 3 {
		// DANGER! check the array index
		for _, t := range ip.revealed[level] {
			delete(pcts, t)
		}
		for k := range pcts {
			pcts[k] = 100.0 / float32(len(pcts))
		}
		return
	}

	revt := map[Terrain]int{
		Jungle:   3,
		Mountain: 3,
		Sands:    3,
		Wetland:  3,
	}
	for _, t := range ip.revealed[3] {
		revt[t] = revt[t] - 1
		if revt[t] == 0 {
			delete(pcts, t)
		}
	}

	trem := 0
	for _, c := range revt {
		trem += c
	}
	for t := range pcts {
		pcts[t] = 100.0 * (float32(3-revt[t]) / float32(trem))
	}

	return
}

// Reveal marks a terrain as revealed for a given Stage.
// For Stages I/II each terrain can be revelead once.
// Stage III can have a terrain revelead up to three times.
func (ip *InvaderCardpool) Reveal(stage int, terrain ...Terrain) {
	for _, t := range terrain {
		// DANGER! check the array index
		ip.revealed[stage] = append(ip.revealed[stage], t)
	}
}
