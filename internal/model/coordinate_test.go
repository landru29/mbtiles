package model_test

import (
	"testing"

	"github.com/landru29/mbtiles/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCoordinate(t *testing.T) {
	t.Run("full", func(t *testing.T) {
		var coordinate model.Coordinate

		require.NoError(t, coordinate.Set("35°25'33.22"))

		assert.InDelta(t, 35.425894, float64(coordinate), 0.0001)

		assert.Equal(t, "35°25'33.220000", coordinate.String())
	})

	t.Run("negative", func(t *testing.T) {
		var coordinate model.Coordinate

		require.NoError(t, coordinate.Set("-35°25'33.22"))

		assert.InDelta(t, -35.425894, float64(coordinate), 0.0001)

		assert.Equal(t, "-35°25'33.220000", coordinate.String())
	})

	t.Run("no seconds", func(t *testing.T) {
		var coordinate model.Coordinate

		require.NoError(t, coordinate.Set("35°25'"))

		assert.InDelta(t, 35.416666, float64(coordinate), 0.0001)

		assert.Equal(t, "35°25'", coordinate.String())
	})

	t.Run("no minutes", func(t *testing.T) {
		var coordinate model.Coordinate

		require.NoError(t, coordinate.Set("35°"))

		assert.InDelta(t, 35.0, float64(coordinate), 0.0001)

		assert.Equal(t, "35°", coordinate.String())
	})
}
