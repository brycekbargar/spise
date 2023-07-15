package domain_test

import (
	"testing"

	"github.com/brycekbargar/spise/domain"
	"github.com/google/go-cmp/cmp/cmpopts"
	"gotest.tools/v3/assert"
)

func TestInvaderCardPool_Reveal(t *testing.T) {
	t.Parallel()

	t.Run("Same card", func(t *testing.T) {
		t.Parallel()

		icp := domain.NewInvaderCardpool(false, false)

		err := icp.Reveal(domain.StageOneJungle)
		assert.NilError(t, err)

		err = icp.Reveal(domain.StageOneJungle)
		assert.NilError(t, err)
	})
	t.Run("Invalid card", func(t *testing.T) {
		t.Parallel()

		icp := domain.NewInvaderCardpool(false, false)

		var ic domain.InvaderCard
		err := icp.Reveal(ic)
		assert.ErrorIs(t, err, domain.ErrInvalidInvaderCard)
	})
	t.Run("All the cards", func(t *testing.T) {
		t.Parallel()

		icp := domain.NewInvaderCardpool(false, false)

		for _, ic := range domain.AllInvaderCards {
			err := icp.Reveal(ic)
			assert.NilError(t, err)
		}
	})
}

func TestInvaderCardPool_Predict(t *testing.T) {
	t.Parallel()

	t.Run("Invalid Stages", func(t *testing.T) {
		t.Parallel()

		icp := domain.NewInvaderCardpool(false, false)
		_, err := icp.Predict(0)
		assert.ErrorIs(t, err, domain.ErrInvalidInvaderCard)

		_, err = icp.Predict(4)
		assert.ErrorIs(t, err, domain.ErrInvalidInvaderCard)
	})
	t.Run("Base Stage I", func(t *testing.T) {
		t.Parallel()

		icp := domain.NewInvaderCardpool(false, false)
		pct, err := icp.Predict(1)
		assert.NilError(t, err)
		assert.Equal(t, pct[domain.Jungle], .25)
		assert.Equal(t, pct[domain.Mountain], .25)
		assert.Equal(t, pct[domain.Sands], .25)
		assert.Equal(t, pct[domain.Wetland], .25)
		assert.Equal(t, pct[domain.CoastalLands], float64(0))
	})
	t.Run("Base Stage II", func(t *testing.T) {
		t.Parallel()

		icp := domain.NewInvaderCardpool(false, false)
		pct, err := icp.Predict(2)
		assert.NilError(t, err)
		assert.Equal(t, pct[domain.Jungle], .20)
		assert.Equal(t, pct[domain.Mountain], .20)
		assert.Equal(t, pct[domain.Sands], .20)
		assert.Equal(t, pct[domain.Wetland], .20)
		assert.Equal(t, pct[domain.CoastalLands], .20)
	})
	t.Run("Base Stage III", func(t *testing.T) {
		t.Parallel()

		icp := domain.NewInvaderCardpool(false, false)
		pct, err := icp.Predict(3)
		assert.NilError(t, err)
		assert.Equal(t, pct[domain.Jungle], .50)
		assert.Equal(t, pct[domain.Mountain], .50)
		assert.Equal(t, pct[domain.Sands], .50)
		assert.Equal(t, pct[domain.Wetland], .50)
		assert.Equal(t, pct[domain.CoastalLands], float64(0))
	})

	t.Run("Scotland", func(t *testing.T) {
		t.Parallel()

		icp := domain.NewInvaderCardpool(true, false)
		pct, err := icp.Predict(1)
		assert.NilError(t, err)
		assert.Equal(t, pct[domain.Jungle], .25)

		pct, err = icp.Predict(2)
		assert.NilError(t, err)
		assert.Equal(t, pct[domain.Jungle], .25)
		assert.Equal(t, pct[domain.CoastalLands], float64(0))
	})
	t.Run("Habsburg Mines", func(t *testing.T) {
		t.Parallel()

		icp := domain.NewInvaderCardpool(false, true)
		pct, err := icp.Predict(1)
		assert.NilError(t, err)
		assert.Equal(t, pct[domain.Jungle], .25)

		pct, err = icp.Predict(2)
		assert.NilError(t, err)
		assert.Equal(t, pct[domain.Jungle], .25)
		assert.Equal(t, pct[domain.CoastalLands], float64(0))
	})

	t.Run("During Stage I", func(t *testing.T) {
		t.Parallel()

		icp := domain.NewInvaderCardpool(false, false)

		err := icp.Reveal(domain.StageOneJungle)
		assert.NilError(t, err)
		pct, err := icp.Predict(1)
		assert.NilError(t, err)
		assert.Equal(t, pct[domain.Jungle], float64(0))

		err = icp.Reveal(domain.StageTwoWetland)
		assert.NilError(t, err)
		pct, err = icp.Predict(1)
		assert.NilError(t, err)
		assert.DeepEqual(
			t,
			pct[domain.Wetland],
			float64(.33),
			cmpopts.EquateApprox(0, .004),
		)

		err = icp.Reveal(domain.StageOneWetland)
		assert.NilError(t, err)
		pct, err = icp.Predict(1)
		assert.NilError(t, err)
		assert.Equal(t, pct[domain.Wetland], float64(0))

		for _, ic := range domain.StageOneInvaderCards {
			err = icp.Reveal(ic)
			assert.NilError(t, err)
		}

		pct, err = icp.Predict(1)
		assert.NilError(t, err)
		assert.Equal(t, pct[domain.Sands], float64(0))
		assert.Equal(t, pct[domain.Wetland], float64(0))
	})

	t.Run("All of Stage III Terrain", func(t *testing.T) {
		t.Parallel()

		icp := domain.NewInvaderCardpool(false, false)

		err := icp.Reveal(domain.StageThreeJungleMountain)
		assert.NilError(t, err)
		err = icp.Reveal(domain.StageThreeJungleSands)
		assert.NilError(t, err)
		err = icp.Reveal(domain.StageThreeJungleWetland)
		assert.NilError(t, err)
		pct, err := icp.Predict(3)
		assert.NilError(t, err)
		assert.Equal(t, pct[domain.Jungle], float64(0))
	})

	t.Run("None of Stage III Terrain", func(t *testing.T) {
		t.Parallel()

		icp := domain.NewInvaderCardpool(false, false)

		err := icp.Reveal(domain.StageThreeJungleMountain)
		assert.NilError(t, err)
		err = icp.Reveal(domain.StageThreeJungleSands)
		assert.NilError(t, err)
		err = icp.Reveal(domain.StageThreeMountainSands)
		assert.NilError(t, err)
		pct, err := icp.Predict(3)
		assert.NilError(t, err)
		assert.Equal(t, pct[domain.Wetland], float64(1))
	})

	t.Run("During Stage III", func(t *testing.T) {
		t.Parallel()

		icp := domain.NewInvaderCardpool(false, false)

		err := icp.Reveal(domain.StageThreeJungleMountain)
		assert.NilError(t, err)
		pct, err := icp.Predict(3)
		assert.NilError(t, err)
		assert.Equal(t, pct[domain.Jungle], float64(.4))

		err = icp.Reveal(domain.StageThreeMountainSands)
		assert.NilError(t, err)
		pct, err = icp.Predict(3)
		assert.NilError(t, err)
		assert.Equal(t, pct[domain.Jungle], float64(.5))

		err = icp.Reveal(domain.StageThreeMountainWetland)
		assert.NilError(t, err)
		pct, err = icp.Predict(3)
		assert.NilError(t, err)
		assert.DeepEqual(
			t,
			pct[domain.Jungle],
			float64(.66),
			cmpopts.EquateApprox(0, .007),
		)

		err = icp.Reveal(domain.StageThreeSandsWetland)
		assert.NilError(t, err)
		pct, err = icp.Predict(3)
		assert.NilError(t, err)
		assert.Equal(t, pct[domain.Jungle], float64(1))

		err = icp.Reveal(domain.StageThreeJungleMountain)
		assert.NilError(t, err)
		pct, err = icp.Predict(3)
		assert.NilError(t, err)
		assert.Equal(t, pct[domain.Jungle], float64(1))
	})
}
