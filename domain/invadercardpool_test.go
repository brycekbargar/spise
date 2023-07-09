package domain_test

import (
	"testing"

	"github.com/brycekbargar/spise/domain"
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
		assert.ErrorContains(t, err, "not a stage")
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
			assert.NilError(t, err)
		}
	})
}
