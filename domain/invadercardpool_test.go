package domain_test

import (
	"testing"

	"github.com/brycekbargar/spise/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvaderCardPool_Reveal(t *testing.T) {
	t.Parallel()

	t.Run("Same card", func(t *testing.T) {
		t.Parallel()

		icp := domain.NewInvaderCardpool(false, false)

		err := icp.Reveal(domain.StageOneJungle)
		require.NoError(t, err)

		err = icp.Reveal(domain.StageOneJungle)
		assert.Error(t, err)
	})
	t.Run("Scotland coastal", func(t *testing.T) {
		t.Parallel()

		icp := domain.NewInvaderCardpool(true, false)

		err := icp.Reveal(domain.StageTwoCoastal)
		assert.Error(t, err)
	})
	t.Run("Habsburg mining coastal", func(t *testing.T) {
		t.Parallel()

		icp := domain.NewInvaderCardpool(false, true)

		err := icp.Reveal(domain.StageTwoCoastal)
		assert.Error(t, err)
	})
	t.Run("Invalid card", func(t *testing.T) {
		t.Parallel()

		icp := domain.NewInvaderCardpool(false, false)

		var ic domain.InvaderCard
		err := icp.Reveal(ic)
		assert.Error(t, err)
	})
	t.Run("All the cards", func(t *testing.T) {
		t.Parallel()

		icp := domain.NewInvaderCardpool(false, false)

		for _, ic := range []domain.InvaderCard{
			domain.StageOneJungle,
			domain.StageOneSands,
			domain.StageOneMountain,
			domain.StageOneWetland,
			domain.StageTwoJungle,
			domain.StageTwoSands,
			domain.StageTwoMountain,
			domain.StageTwoWetland,
			domain.StageTwoCoastal,
			domain.StageThreeJungleMountain,
			domain.StageThreeJungleSands,
			domain.StageThreeJungleWetland,
			domain.StageThreeMountainSands,
			domain.StageThreeMountainWetland,
			domain.StageThreeSandsWetland,
		} {
			err := icp.Reveal(ic)
			assert.NoError(t, err)
		}
	})
}
