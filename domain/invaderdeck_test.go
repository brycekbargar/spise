package domain_test

import (
	"strings"
	"testing"

	"github.com/brycekbargar/spise/domain"
	"gotest.tools/v3/assert"
)

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
		}, "111-3-2222-3333"},
		{"BP3", &domain.Game{
			LeadingAdversary:      domain.BrandenburgPrussia,
			LeadingAdversaryLevel: 3,
		}, "111-3-222-3333"},
		{"BP4", &domain.Game{
			LeadingAdversary:      domain.BrandenburgPrussia,
			LeadingAdversaryLevel: 4,
		}, "11-3-222-3333"},
		{"BP5", &domain.Game{
			LeadingAdversary:      domain.BrandenburgPrussia,
			LeadingAdversaryLevel: 5,
		}, "1-3-222-3333"},
		{"BP6", &domain.Game{
			LeadingAdversary:      domain.BrandenburgPrussia,
			LeadingAdversaryLevel: 6,
		}, "3-222-3333"},
		{"HLC3", &domain.Game{
			LeadingAdversary:      domain.HabsburgLivestock,
			LeadingAdversaryLevel: 3,
		}, "11-2222-33333"},
		{"HME4", &domain.Game{
			LeadingAdversary:      domain.HabsburgMines,
			LeadingAdversaryLevel: 4,
		}, "111-2S22-33333"},
		{"R4", &domain.Game{
			LeadingAdversary:      domain.Russia,
			LeadingAdversaryLevel: 4,
		}, "111-2-3-2-3-2-3-2-33"},
		{"S2", &domain.Game{
			LeadingAdversary:      domain.Scotland,
			LeadingAdversaryLevel: 2,
		}, "11-22-1-C2-33333"},
		{"S4", &domain.Game{
			LeadingAdversary:      domain.Scotland,
			LeadingAdversaryLevel: 4,
		}, "11-22-3-C2-3333"},
		{"BP5HLC3", &domain.Game{
			LeadingAdversary:         domain.BrandenburgPrussia,
			LeadingAdversaryLevel:    5,
			SupportingAdversary:      domain.HabsburgLivestock,
			SupportingAdversaryLevel: 3,
		}, "3-222-3333"},
		{"BP5HME4", &domain.Game{
			LeadingAdversary:         domain.BrandenburgPrussia,
			LeadingAdversaryLevel:    5,
			SupportingAdversary:      domain.HabsburgMines,
			SupportingAdversaryLevel: 4,
		}, "1-3-S22-3333"},
		{"BP5R4", &domain.Game{
			LeadingAdversary:         domain.BrandenburgPrussia,
			LeadingAdversaryLevel:    5,
			SupportingAdversary:      domain.Russia,
			SupportingAdversaryLevel: 4,
		}, "1-33-2-3-2-3-2-3"},
		{"BP5S4", &domain.Game{
			LeadingAdversary:         domain.BrandenburgPrussia,
			LeadingAdversaryLevel:    5,
			SupportingAdversary:      domain.Scotland,
			SupportingAdversaryLevel: 4,
		}, "3-2-3-C2-333"},
		{"HLC3BP5", &domain.Game{
			LeadingAdversary:         domain.HabsburgLivestock,
			LeadingAdversaryLevel:    3,
			SupportingAdversary:      domain.BrandenburgPrussia,
			SupportingAdversaryLevel: 5,
		}, "3-222-3333"},
		{"HLC3HME4", &domain.Game{
			LeadingAdversary:         domain.HabsburgLivestock,
			LeadingAdversaryLevel:    3,
			SupportingAdversary:      domain.HabsburgMines,
			SupportingAdversaryLevel: 4,
		}, "11-2S22-33333"},
		{"HLC3R4", &domain.Game{
			LeadingAdversary:         domain.HabsburgLivestock,
			LeadingAdversaryLevel:    3,
			SupportingAdversary:      domain.Russia,
			SupportingAdversaryLevel: 4,
		}, "11-2-3-2-3-2-3-2-33"},
		{"HLC3S4", &domain.Game{
			LeadingAdversary:         domain.HabsburgLivestock,
			LeadingAdversaryLevel:    3,
			SupportingAdversary:      domain.Scotland,
			SupportingAdversaryLevel: 4,
		}, "1-22-3-C2-3333"},
		{"HME4BP5", &domain.Game{
			LeadingAdversary:         domain.HabsburgMines,
			LeadingAdversaryLevel:    4,
			SupportingAdversary:      domain.BrandenburgPrussia,
			SupportingAdversaryLevel: 5,
		}, "1-3-2S2-3333"},
		{"HME4HLC3", &domain.Game{
			LeadingAdversary:         domain.HabsburgMines,
			LeadingAdversaryLevel:    4,
			SupportingAdversary:      domain.HabsburgLivestock,
			SupportingAdversaryLevel: 3,
		}, "11-2S22-33333"},
		{"HME4R4", &domain.Game{
			LeadingAdversary:         domain.HabsburgMines,
			LeadingAdversaryLevel:    4,
			SupportingAdversary:      domain.Russia,
			SupportingAdversaryLevel: 4,
		}, "111-2-3-S-3-2-3-2-33"},
		{"HME4S4", &domain.Game{
			LeadingAdversary:         domain.HabsburgMines,
			LeadingAdversaryLevel:    4,
			SupportingAdversary:      domain.Scotland,
			SupportingAdversaryLevel: 4,
		}, "11-2S-3-C2-3333"},
		{"R4BP5", &domain.Game{
			LeadingAdversary:         domain.Russia,
			LeadingAdversaryLevel:    4,
			SupportingAdversary:      domain.BrandenburgPrussia,
			SupportingAdversaryLevel: 5,
		}, "1-3-2-3-2-3-2-33"},
		{"R4HLC3", &domain.Game{
			LeadingAdversary:         domain.Russia,
			LeadingAdversaryLevel:    4,
			SupportingAdversary:      domain.HabsburgLivestock,
			SupportingAdversaryLevel: 3,
		}, "11-2-3-2-3-2-3-2-33"},
		{"R4HME4", &domain.Game{
			LeadingAdversary:         domain.Russia,
			LeadingAdversaryLevel:    4,
			SupportingAdversary:      domain.HabsburgMines,
			SupportingAdversaryLevel: 4,
		}, "111-2-3-S-3-2-3-2-33"},
		{"R4S4", &domain.Game{
			LeadingAdversary:         domain.Russia,
			LeadingAdversaryLevel:    4,
			SupportingAdversary:      domain.Scotland,
			SupportingAdversaryLevel: 4,
		}, "11-2-3-2-3-3-C-3-2-3"},
		{"S4BP5", &domain.Game{
			LeadingAdversary:         domain.Scotland,
			LeadingAdversaryLevel:    4,
			SupportingAdversary:      domain.BrandenburgPrussia,
			SupportingAdversaryLevel: 5,
		}, "3-22-3-C-333"},
		{"S4HLC3", &domain.Game{
			LeadingAdversary:         domain.Scotland,
			LeadingAdversaryLevel:    4,
			SupportingAdversary:      domain.HabsburgLivestock,
			SupportingAdversaryLevel: 3,
		}, "1-22-3-C2-3333"},
		{"S4HME4", &domain.Game{
			LeadingAdversary:         domain.Scotland,
			LeadingAdversaryLevel:    4,
			SupportingAdversary:      domain.HabsburgMines,
			SupportingAdversaryLevel: 4,
		}, "112S-3-C2-3333"},
		{"S4R4", &domain.Game{
			LeadingAdversary:         domain.Scotland,
			LeadingAdversaryLevel:    4,
			SupportingAdversary:      domain.Russia,
			SupportingAdversaryLevel: 4,
		}, "11-2-3-2-3-3-C-3-2-3"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			icp := domain.NewInvaderDeck(tc.game)
			stages := strings.Split(strings.ReplaceAll(tc.initial, "-", ""), "")
			for i, s := range stages {
				card := icp.InDeck[i]

				switch s {
				case "C":
					assert.Equal(t, domain.StageTwoCoastal, card.Card)
				case "S":
					assert.Equal(t, domain.StageTwoSaltDeposits, card.Card)
				case "1":
					assert.Equal(t, domain.StageOneUnknown, card.Card)
				case "2":
					assert.Equal(t, domain.StageTwoUnknown, card.Card)
				case "3":
					assert.Equal(t, domain.StageThreeUnknown, card.Card)
				}
			}
		})
	}
}
