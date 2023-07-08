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
	t.Run("Invalid card", func(t *testing.T) {
		t.Parallel()

		icp := domain.NewInvaderCardpool(false, false)

		var ic domain.InvaderCard
		err := icp.Reveal(ic)
		assert.Error(t, err)
	})
}
