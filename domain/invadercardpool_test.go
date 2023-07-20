package domain_test

import (
	"testing"

	"github.com/brycekbargar/spise/domain"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/google/go-cmp/cmp/cmpopts"
	"gotest.tools/v3/assert"
)

//nolint:exhaustruct
func TestInvaderCardPool_NewInvaderCardpool(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		game     *domain.Game
		revealed *domain.InvaderCard
	}{
		{"NoAdversary", &domain.Game{}, nil},
		{
			"LowScotland",
			&domain.Game{
				LeadingAdversaryLevel: 1,
				LeadingAdversary:      domain.Scotland,
			},
			nil,
		},
		{
			"Scotland",
			&domain.Game{
				LeadingAdversaryLevel: 2,
				LeadingAdversary:      domain.Scotland,
			},
			&domain.StageTwoCoastal,
		},
		{
			"HighScotland",
			&domain.Game{
				LeadingAdversaryLevel: 6,
				LeadingAdversary:      domain.Scotland,
			},
			&domain.StageTwoCoastal,
		},
		{
			"LowSupportingScotland",
			&domain.Game{
				SupportingAdversaryLevel: 1,
				SupportingAdversary:      domain.Scotland,
			},
			nil,
		},
		{
			"SupportingScotland",
			&domain.Game{
				SupportingAdversaryLevel: 2,
				SupportingAdversary:      domain.Scotland,
			},
			&domain.StageTwoCoastal,
		},
		{
			"HighSupportingScotland",
			&domain.Game{
				SupportingAdversaryLevel: 6,
				SupportingAdversary:      domain.Scotland,
			},
			&domain.StageTwoCoastal,
		},
		{
			"LowHabsburgMines",
			&domain.Game{
				LeadingAdversaryLevel: 3,
				LeadingAdversary:      domain.HabsburgMines,
			},
			nil,
		},
		{
			"HabsburgMines",
			&domain.Game{
				LeadingAdversaryLevel: 4,
				LeadingAdversary:      domain.HabsburgMines,
			},
			&domain.StageTwoCoastal,
		},
		{
			"HighHabsburgMines",
			&domain.Game{
				LeadingAdversaryLevel: 6,
				LeadingAdversary:      domain.HabsburgMines,
			},
			&domain.StageTwoCoastal,
		},
		{
			"LowSupportingHabsburgMines",
			&domain.Game{
				SupportingAdversaryLevel: 3,
				SupportingAdversary:      domain.HabsburgMines,
			},
			nil,
		},
		{
			"SupportingHabsburgMines",
			&domain.Game{
				SupportingAdversaryLevel: 4,
				SupportingAdversary:      domain.HabsburgMines,
			},
			&domain.StageTwoCoastal,
		},
		{
			"HighSupportingHabsburgMines",
			&domain.Game{
				SupportingAdversaryLevel: 6,
				SupportingAdversary:      domain.HabsburgMines,
			},
			&domain.StageTwoCoastal,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			icp := domain.NewInvaderCardpool(tc.game)
			assert.Assert(t, icp.Revealed[1].Cardinality() == 0)
			if tc.revealed == nil {
				assert.Assert(t, icp.Revealed[2].Cardinality() == 0)
			} else {
				assert.Assert(t, icp.Revealed[2].Contains(*tc.revealed))
			}
			assert.Assert(t, icp.Revealed[3].Cardinality() == 0)
		})
	}
}

//nolint:exhaustruct
func TestInvaderCardPool_Reveal(t *testing.T) {
	t.Parallel()

	t.Run("Same card", func(t *testing.T) {
		t.Parallel()

		icp := domain.NewInvaderCardpool(&domain.Game{})

		err := icp.Reveal(domain.StageOneJungle)
		assert.NilError(t, err)
		err = icp.Reveal(domain.StageOneJungle)
		assert.NilError(t, err)

		assert.Assert(t, icp.Revealed[1].Contains(domain.StageOneJungle))
	})
	t.Run("Invalid card", func(t *testing.T) {
		t.Parallel()

		cases := []struct {
			name string
			card domain.InvaderCard
		}{
			{"Stage1", domain.InvaderCard{1, domain.Jungle, domain.Wetland}},
			{"Terrain2", domain.InvaderCard{1, domain.Jungle, "not-a-terrain"}},
			{"Terrain", domain.InvaderCard{2, "not-a-terrain", domain.Jungle}},
			{"-Stage", domain.InvaderCard{-1, domain.Jungle, "not-a-terrain"}},
			{"Stage", domain.InvaderCard{4, domain.Jungle, domain.Wetland}},
			{"Terrainx2", domain.InvaderCard{3, domain.Jungle, domain.Jungle}},
		}

		for _, tc := range cases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				icp := domain.NewInvaderCardpool(&domain.Game{})

				err := icp.Reveal(tc.card)
				assert.ErrorIs(t, err, domain.ErrInvalidInvaderCard)
			})
		}
	})
}

//nolint:exhaustruct
func TestInvaderCardPool_Predict(t *testing.T) {
	t.Parallel()

	t.Run("Invalid Stages", func(t *testing.T) {
		cases := []struct {
			name  string
			stage int
		}{
			{"-Stage", 0},
			{"+Stage", 4},
		}

		for _, tc := range cases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				icp := domain.NewInvaderCardpool(&domain.Game{})
				_, err := icp.Predict(tc.stage)
				assert.ErrorIs(t, err, domain.ErrInvalidInvaderCard)
			})
		}
	})

	t.Run("Base", func(t *testing.T) {
		t.Parallel()
		cases := []struct {
			name        string
			stage       int
			revealed    mapset.Set[domain.InvaderCard]
			predictions map[domain.Terrain]float64
		}{
			{
				"StageI",
				1,
				mapset.NewSet[domain.InvaderCard](),
				map[domain.Terrain]float64{
					domain.Jungle:   .25,
					domain.Mountain: .25,
					domain.Sands:    .25,
					domain.Wetland:  .25,
				},
			},
			{
				"StageII",
				2,
				mapset.NewSet[domain.InvaderCard](),
				map[domain.Terrain]float64{
					domain.Jungle:       .20,
					domain.Mountain:     .20,
					domain.Sands:        .20,
					domain.Wetland:      .20,
					domain.CoastalLands: .20,
				},
			},
			{
				"StageIII",
				3,
				mapset.NewSet[domain.InvaderCard](),
				map[domain.Terrain]float64{
					domain.Jungle:   .50,
					domain.Mountain: .50,
					domain.Sands:    .50,
					domain.Wetland:  .50,
				},
			},
		}

		for _, tc := range cases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				icp := domain.NewInvaderCardpool(&domain.Game{})
				icp.Revealed[tc.stage] = tc.revealed

				p, err := icp.Predict(tc.stage)
				assert.NilError(t, err)
				assert.DeepEqual(t, p, tc.predictions)
			})
		}
	})

	t.Run("Stage I", func(t *testing.T) {
		t.Parallel()
		cases := []struct {
			name        string
			stage       int
			revealed    mapset.Set[domain.InvaderCard]
			predictions map[domain.Terrain]float64
		}{
			{
				"One",
				1,
				mapset.NewSet(domain.StageOneJungle),
				map[domain.Terrain]float64{
					domain.Mountain: .33,
					domain.Sands:    .33,
					domain.Wetland:  .33,
				},
			},
			{
				"All",
				1,
				mapset.NewSet[domain.InvaderCard](
					domain.StageOneInvaderCards...,
				),
				map[domain.Terrain]float64{},
			},
		}

		for _, tc := range cases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				icp := domain.NewInvaderCardpool(&domain.Game{})
				icp.Revealed[tc.stage] = tc.revealed

				p, err := icp.Predict(tc.stage)
				assert.NilError(t, err)
				assert.DeepEqual(t, p,
					tc.predictions,
					cmpopts.EquateApprox(0, .004))
			})
		}
	})

	t.Run("Stage III", func(t *testing.T) {
		t.Parallel()
		cases := []struct {
			name        string
			stage       int
			revealed    mapset.Set[domain.InvaderCard]
			predictions map[domain.Terrain]float64
		}{
			{
				"AllTerrain",
				3,
				mapset.NewSet(
					domain.StageThreeJungleMountain,
					domain.StageThreeMountainSands,
					domain.StageThreeMountainWetland,
				),
				map[domain.Terrain]float64{
					domain.Jungle:  .66,
					domain.Sands:   .66,
					domain.Wetland: .66,
				},
			},
			{
				"NoneTerrain",
				3,
				mapset.NewSet[domain.InvaderCard](
					domain.StageThreeJungleMountain,
					domain.StageThreeJungleWetland,
					domain.StageThreeMountainWetland,
				),
				map[domain.Terrain]float64{
					domain.Jungle:   .33,
					domain.Mountain: .33,
					domain.Sands:    1,
					domain.Wetland:  .33,
				},
			},
		}

		for _, tc := range cases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				icp := domain.NewInvaderCardpool(&domain.Game{})
				icp.Revealed[tc.stage] = tc.revealed

				p, err := icp.Predict(tc.stage)
				assert.NilError(t, err)
				assert.DeepEqual(t, p,
					tc.predictions,
					cmpopts.EquateApprox(0, .007))
			})
		}
	})
}
