package domain_test

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/brycekbargar/spise/domain"
	"gotest.tools/v3/assert"
)

//nolint:exhaustruct
func TestDraw(t *testing.T) {
	t.Parallel()

	t.Run("EmptyDeck", func(t *testing.T) {
		t.Parallel()

		deck := &domain.InvaderDeck{}
		err := deck.Draw(domain.InvaderCard{})
		assert.ErrorIs(t, domain.ErrNoInvaderCard, err)
	})

	t.Run("Impossible", func(t *testing.T) {
		t.Parallel()

		deck := domain.NewInvaderDeck(&domain.Game{})
		err := deck.Draw(domain.StageThreeUnknown)
		assert.ErrorIs(t, domain.ErrInvalidInvaderCard, err)
	})

	t.Run("FullDeck", func(t *testing.T) {
		t.Parallel()

		deck := domain.NewInvaderDeck(&domain.Game{})
		for _, dic := range []domain.InvaderCard{
			{Stage: 1},
			{Stage: 1},
			{Stage: 1},
			{Stage: 2},
			{Stage: 2},
			{Stage: 2},
			{Stage: 2},
			{Stage: 3},
			{Stage: 3},
			{Stage: 3},
			{Stage: 3},
			{Stage: 3},
		} {
			err := deck.Draw(dic)
			assert.NilError(t, err, dic)
		}

		assert.Equal(t, "111-2222-33333", deckToString(t, deck.Drawn))
	})

	t.Run("Returnable", func(t *testing.T) {
		t.Parallel()

		cases := []struct {
			name    string
			drawn   []domain.InvaderCard
			inDeck  []domain.InvaderCard
			discard string
		}{
			{
				"=Same",
				[]domain.InvaderCard{},
				[]domain.InvaderCard{
					domain.StageOneUnknown,
					domain.StageOneUnknown,
				},
				"1*",
			},
			{
				"+Stage",
				[]domain.InvaderCard{},
				[]domain.InvaderCard{
					domain.StageTwoUnknown,
					domain.StageOneUnknown,
				},
				"2*",
			},
			{
				"-Stage",
				[]domain.InvaderCard{},
				[]domain.InvaderCard{
					domain.StageTwoUnknown,
					domain.StageThreeUnknown,
				},
				"2*",
			},
			{
				"--Stage",
				[]domain.InvaderCard{},
				[]domain.InvaderCard{
					domain.StageOneUnknown,
					domain.StageThreeUnknown,
				},
				"1",
			},
			{
				"++Stage",
				[]domain.InvaderCard{},
				[]domain.InvaderCard{
					domain.StageThreeUnknown,
					domain.StageOneUnknown,
				},
				"3",
			},
			{
				"Empty",
				[]domain.InvaderCard{},
				[]domain.InvaderCard{
					domain.StageThreeUnknown,
				},
				"3",
			},
			{
				"Disparate",
				[]domain.InvaderCard{
					domain.StageTwoUnknown,
					domain.StageThreeUnknown,
					domain.StageTwoUnknown,
					domain.StageOneUnknown,
				},
				[]domain.InvaderCard{
					domain.StageOneUnknown,
					domain.StageThreeUnknown,
				},
				"2*-3*-2*-11",
			},
		}

		for _, tc := range cases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				deck := domain.InvaderDeck{}
				deck.Drawn = make([]domain.InvaderCardDrawn, 0, len(tc.drawn))
				for _, c := range tc.drawn {
					deck.Drawn = append(
						deck.Drawn,
						domain.InvaderCardDrawn{InvaderCard: c},
					)
				}
				deck.InDeck = make(
					[]domain.InvaderCardInDeck,
					0,
					len(tc.inDeck),
				)
				for _, c := range tc.inDeck {
					deck.InDeck = append(
						deck.InDeck,
						domain.InvaderCardInDeck{InvaderCard: c},
					)
				}

				err := deck.Draw(tc.inDeck[0])
				assert.NilError(t, err)
				assert.Equal(t, tc.discard, deckToString(t, deck.Drawn))
			})
		}
	})
}

//nolint:exhaustruct
func TestEntrenched(t *testing.T) {
	t.Parallel()

	t.Run("Entrenchable", func(t *testing.T) {
		t.Parallel()

		cases := []struct {
			name string
			game *domain.Game
			card domain.InvaderCard
			err  error
		}{
			{
				"None",
				&domain.Game{},
				domain.StageTwoMountain,
				domain.ErrNotEntrenched,
			},
			{
				"OneNonRussia",
				&domain.Game{
					LeadingAdversary:      domain.Scotland,
					LeadingAdversaryLevel: 5,
				},
				domain.StageTwoMountain,
				domain.ErrNotEntrenched,
			},
			{
				"TwoNonRussia",
				&domain.Game{
					LeadingAdversary:         domain.Scotland,
					LeadingAdversaryLevel:    6,
					SupportingAdversary:      domain.England,
					SupportingAdversaryLevel: 6,
				},
				domain.StageTwoMountain,
				domain.ErrNotEntrenched,
			},
			{
				"LowRussiaLead",
				&domain.Game{LeadingAdversary: domain.Russia},
				domain.StageTwoMountain,
				domain.ErrNotEntrenched,
			},
			{
				"LowRussiaSupport",
				&domain.Game{
					LeadingAdversary:      domain.Scotland,
					LeadingAdversaryLevel: 5,
					SupportingAdversary:   domain.Russia,
				},
				domain.StageTwoMountain,
				domain.ErrNotEntrenched,
			},
			{
				"LeadingRussia",
				&domain.Game{
					LeadingAdversary:      domain.Russia,
					LeadingAdversaryLevel: 5,
				},
				domain.StageTwoMountain,
				nil,
			},
			{
				"SupportingRussia",
				&domain.Game{
					LeadingAdversary:         domain.Scotland,
					SupportingAdversary:      domain.Russia,
					SupportingAdversaryLevel: 5,
				},
				domain.StageTwoMountain,
				nil,
			},
			{
				"StageTooLow",
				&domain.Game{
					LeadingAdversary:      domain.Russia,
					LeadingAdversaryLevel: 5,
				},
				domain.StageOneMountain,
				domain.ErrInvalidInvaderCard,
			},
		}

		for _, tc := range cases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				deck := domain.NewInvaderDeck(tc.game)
				err := deck.Entrenched(tc.card)
				assert.ErrorIs(t, err, tc.err)
			})
		}
	})
}

//nolint:exhaustruct
func TestReturn(t *testing.T) {
	t.Parallel()

	t.Run("NotReturnable", func(t *testing.T) {
		t.Parallel()

		deck := domain.NewInvaderDeck(&domain.Game{})
		for _, dic := range []domain.InvaderCard{
			{Stage: 1},
			{Stage: 1},
			{Stage: 1},
			{Stage: 2},
			{Stage: 2},
			{Stage: 2},
			{Stage: 2},
			{Stage: 3},
		} {
			err := deck.Draw(dic)
			assert.NilError(t, err, dic)
		}

		err := deck.Return(domain.InvaderCard{Stage: 1})
		assert.ErrorIs(t, err, domain.ErrInvaderCardNotReturnable)
	})

	t.Run("Not Drawn", func(t *testing.T) {
		t.Parallel()

		deck := domain.NewInvaderDeck(&domain.Game{})
		for _, dic := range []domain.InvaderCard{
			domain.StageOneWetland,
		} {
			err := deck.Draw(dic)
			assert.NilError(t, err, dic)
		}

		err := deck.Return(domain.StageOneJungle)
		assert.ErrorIs(t, err, domain.ErrInvalidInvaderCard)
	})

	t.Run("Returnable", func(t *testing.T) {
		t.Parallel()

		cases := []struct {
			name     string
			drawn    []domain.InvaderCard
			inDeck   []domain.InvaderCard
			returned domain.InvaderCard
			discard  string
		}{
			{
				"Basic",
				[]domain.InvaderCard{},
				[]domain.InvaderCard{
					domain.StageOneJungle,
					domain.StageOneUnknown,
				},
				domain.StageOneJungle,
				"1*",
			},
			{
				"Many",
				[]domain.InvaderCard{
					domain.StageOneJungle,
					domain.StageOneWetland,
					domain.StageOneMountain,
				},
				[]domain.InvaderCard{
					domain.StageOneSands,
					domain.StageTwoUnknown,
				},
				domain.StageOneMountain,
				"1*1*-2*-1*",
			},
			{
				"TwoAway",
				[]domain.InvaderCard{
					domain.StageOneJungle,
					domain.StageOneWetland,
					domain.StageTwoMountain,
					domain.StageThreeMountainWetland,
				},
				[]domain.InvaderCard{
					domain.StageTwoSands,
					domain.StageThreeUnknown,
				},
				domain.StageThreeMountainWetland,
				"11-2*-3*-2*",
			},
			{
				"TwoAwayBridge",
				[]domain.InvaderCard{
					domain.StageOneJungle,
					domain.StageOneWetland,
					domain.StageTwoMountain,
					domain.StageThreeMountainWetland,
				},
				[]domain.InvaderCard{
					domain.StageTwoSands,
					domain.StageThreeUnknown,
				},
				domain.StageTwoMountain,
				"1*1*-3*3*-2*",
			},
		}

		for _, tc := range cases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				deck := domain.InvaderDeck{}
				deck.Drawn = make([]domain.InvaderCardDrawn, 0, len(tc.drawn))
				for _, c := range tc.drawn {
					deck.Drawn = append(
						deck.Drawn,
						domain.InvaderCardDrawn{InvaderCard: c},
					)
				}
				deck.InDeck = make(
					[]domain.InvaderCardInDeck,
					0,
					len(tc.inDeck),
				)
				for _, c := range tc.inDeck {
					deck.InDeck = append(
						deck.InDeck,
						domain.InvaderCardInDeck{InvaderCard: c},
					)
				}

				err := deck.Draw(tc.inDeck[0])
				assert.NilError(t, err)

				err = deck.Return(tc.returned)
				assert.NilError(t, err)

				assert.Equal(t, tc.returned, deck.InDeck[0].InvaderCard)
				assert.Equal(t, tc.discard, deckToString(t, deck.Drawn))
			})
		}
	})
}

func deckToString[IC any](t *testing.T, cards []IC) string {
	t.Helper()

	ostg := -1
	var actual bytes.Buffer
	for _, card := range cards {
		// TODO: Fix this when generics suck less
		var actcard domain.InvaderCard
		switch dc := any(card).(type) {
		case domain.InvaderCardInDeck:
			actcard = dc.InvaderCard
		case domain.InvaderCardDrawn:
			actcard = dc.InvaderCard
		}

		sym := ""
		var stg int
		switch actcard {
		case domain.StageTwoCoastal:
			sym = "C"
			stg = 2
		case domain.StageTwoSaltDeposits:
			sym = "S"
			stg = 2
		default:
			if 1 > actcard.Stage || actcard.Stage > 3 {
				t.Fatalf("invalidcard: %v", actcard)
			}

			sym = strconv.Itoa(actcard.Stage)
			stg = actcard.Stage
		}

		idc, ok := any(card).(domain.InvaderCardInDeck)
		if ok && idc.SpeciallyPlaced == true {
			sym += "*"
		}
		dc, ok := any(card).(domain.InvaderCardDrawn)
		if ok && dc.PastReturnable == true {
			sym += "*"
		}

		if ostg != -1 && ostg != stg {
			actual.WriteString("-")
		}
		actual.WriteString(sym)
		ostg = stg
	}

	return actual.String()
}

//nolint:exhaustruct
func TestInvaderDeck_NewInvaderDeck(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		game    *domain.Game
		initial string
	}{
		{"Base", &domain.Game{}, "111-2222-33333"},
		{"BP2", &domain.Game{
			LeadingAdversary:      domain.BrandenburgPrussia,
			LeadingAdversaryLevel: 2,
		}, "111-3*-2222-3333"},
		{"BP3", &domain.Game{
			LeadingAdversary:      domain.BrandenburgPrussia,
			LeadingAdversaryLevel: 3,
		}, "11-3*-2222-3333"},
		{"BP4", &domain.Game{
			LeadingAdversary:      domain.BrandenburgPrussia,
			LeadingAdversaryLevel: 4,
		}, "11-3*-222-3333"},
		{"BP5", &domain.Game{
			LeadingAdversary:      domain.BrandenburgPrussia,
			LeadingAdversaryLevel: 5,
		}, "1-3*-222-3333"},
		{"BP6", &domain.Game{
			LeadingAdversary:      domain.BrandenburgPrussia,
			LeadingAdversaryLevel: 6,
		}, "3*-222-3333"},
		{"HLC3", &domain.Game{
			LeadingAdversary:      domain.HabsburgLivestock,
			LeadingAdversaryLevel: 3,
		}, "11-2222-33333"},
		{"HME4", &domain.Game{
			LeadingAdversary:      domain.HabsburgMines,
			LeadingAdversaryLevel: 4,
		}, "111-2S*22-33333"},
		{"R4", &domain.Game{
			LeadingAdversary:      domain.Russia,
			LeadingAdversaryLevel: 4,
		}, "111-2-3*-2-3*-2-3*-2-3*3"},
		{"S2", &domain.Game{
			LeadingAdversary:      domain.Scotland,
			LeadingAdversaryLevel: 2,
		}, "11-2*2*-1-C*2-33333"},
		{"S4", &domain.Game{
			LeadingAdversary:      domain.Scotland,
			LeadingAdversaryLevel: 4,
		}, "11-2*2*-3*-C*2-3333"},
		{"BP5HLC3", &domain.Game{
			LeadingAdversary:         domain.BrandenburgPrussia,
			LeadingAdversaryLevel:    5,
			SupportingAdversary:      domain.HabsburgLivestock,
			SupportingAdversaryLevel: 3,
		}, "3*-222-3333"},
		{"BP5HME4", &domain.Game{
			LeadingAdversary:         domain.BrandenburgPrussia,
			LeadingAdversaryLevel:    5,
			SupportingAdversary:      domain.HabsburgMines,
			SupportingAdversaryLevel: 4,
		}, "1-3*-S*22-3333"},
		{"BP5R4", &domain.Game{
			LeadingAdversary:         domain.BrandenburgPrussia,
			LeadingAdversaryLevel:    5,
			SupportingAdversary:      domain.Russia,
			SupportingAdversaryLevel: 4,
		}, "1-3*3*-2-3*-2-3*-2-3*"},
		{"BP5S4", &domain.Game{
			LeadingAdversary:         domain.BrandenburgPrussia,
			LeadingAdversaryLevel:    5,
			SupportingAdversary:      domain.Scotland,
			SupportingAdversaryLevel: 4,
		}, "3*-2*-3*-C*2-333"},
		{"HLC3BP5", &domain.Game{
			LeadingAdversary:         domain.HabsburgLivestock,
			LeadingAdversaryLevel:    3,
			SupportingAdversary:      domain.BrandenburgPrussia,
			SupportingAdversaryLevel: 5,
		}, "3*-222-3333"},
		{"HLC3HME4", &domain.Game{
			LeadingAdversary:         domain.HabsburgLivestock,
			LeadingAdversaryLevel:    3,
			SupportingAdversary:      domain.HabsburgMines,
			SupportingAdversaryLevel: 4,
		}, "11-2S*22-33333"},
		{"HLC3R4", &domain.Game{
			LeadingAdversary:         domain.HabsburgLivestock,
			LeadingAdversaryLevel:    3,
			SupportingAdversary:      domain.Russia,
			SupportingAdversaryLevel: 4,
		}, "11-2-3*-2-3*-2-3*-2-3*3"},
		{"HLC3S4", &domain.Game{
			LeadingAdversary:         domain.HabsburgLivestock,
			LeadingAdversaryLevel:    3,
			SupportingAdversary:      domain.Scotland,
			SupportingAdversaryLevel: 4,
		}, "1-2*2*-3*-C*2-3333"},
		{"HME4BP5", &domain.Game{
			LeadingAdversary:         domain.HabsburgMines,
			LeadingAdversaryLevel:    4,
			SupportingAdversary:      domain.BrandenburgPrussia,
			SupportingAdversaryLevel: 5,
		}, "1-3*-2S*2-3333"},
		{"HME4HLC3", &domain.Game{
			LeadingAdversary:         domain.HabsburgMines,
			LeadingAdversaryLevel:    4,
			SupportingAdversary:      domain.HabsburgLivestock,
			SupportingAdversaryLevel: 3,
		}, "11-2S*22-33333"},
		{"HME4R4", &domain.Game{
			LeadingAdversary:         domain.HabsburgMines,
			LeadingAdversaryLevel:    4,
			SupportingAdversary:      domain.Russia,
			SupportingAdversaryLevel: 4,
		}, "111-2-3*-S*-3*-2-3*-2-3*3"},
		{"HME4S4", &domain.Game{
			LeadingAdversary:         domain.HabsburgMines,
			LeadingAdversaryLevel:    4,
			SupportingAdversary:      domain.Scotland,
			SupportingAdversaryLevel: 4,
		}, "11-2*S*-3*-C*2-3333"},
		{"R4BP5", &domain.Game{
			LeadingAdversary:         domain.Russia,
			LeadingAdversaryLevel:    4,
			SupportingAdversary:      domain.BrandenburgPrussia,
			SupportingAdversaryLevel: 5,
		}, "1-3*-2-3*-2-3*-2-3*3"},
		{"R4HLC3", &domain.Game{
			LeadingAdversary:         domain.Russia,
			LeadingAdversaryLevel:    4,
			SupportingAdversary:      domain.HabsburgLivestock,
			SupportingAdversaryLevel: 3,
		}, "11-2-3*-2-3*-2-3*-2-3*3"},
		{"R4HME4", &domain.Game{
			LeadingAdversary:         domain.Russia,
			LeadingAdversaryLevel:    4,
			SupportingAdversary:      domain.HabsburgMines,
			SupportingAdversaryLevel: 4,
		}, "111-2-3*-S*-3*-2-3*-2-3*3"},
		{"R4S4", &domain.Game{
			LeadingAdversary:         domain.Russia,
			LeadingAdversaryLevel:    4,
			SupportingAdversary:      domain.Scotland,
			SupportingAdversaryLevel: 4,
		}, "11-2*-3*-2*-3*3*-C*-3*-2-3*"},
		{"S4BP5", &domain.Game{
			LeadingAdversary:         domain.Scotland,
			LeadingAdversaryLevel:    4,
			SupportingAdversary:      domain.BrandenburgPrussia,
			SupportingAdversaryLevel: 5,
		}, "3*-2*2*-3*-C*-333"},
		{"S4HLC3", &domain.Game{
			LeadingAdversary:         domain.Scotland,
			LeadingAdversaryLevel:    4,
			SupportingAdversary:      domain.HabsburgLivestock,
			SupportingAdversaryLevel: 3,
		}, "1-2*2*-3*-C*2-3333"},
		{"S4HME4", &domain.Game{
			LeadingAdversary:         domain.Scotland,
			LeadingAdversaryLevel:    4,
			SupportingAdversary:      domain.HabsburgMines,
			SupportingAdversaryLevel: 4,
		}, "11-2*S*-3*-C*2-3333"},
		{"S4R4", &domain.Game{
			LeadingAdversary:         domain.Scotland,
			LeadingAdversaryLevel:    4,
			SupportingAdversary:      domain.Russia,
			SupportingAdversaryLevel: 4,
		}, "11-2*-3*-2*-3*3*-C*-3*-2-3*"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			deck := domain.NewInvaderDeck(tc.game)
			assert.Equal(t, tc.initial, deckToString(t, deck.InDeck))
		})
	}
}
