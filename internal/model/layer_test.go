package model_test

import (
	"testing"

	"github.com/landru29/mbtiles/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestLatLngToTile(t *testing.T) {
	t.Run("top-left", func(t *testing.T) {
		x, y := model.LatLngToTile(51.251834, -5.593299, 11)
		assert.EqualValues(t, 992, x)
		assert.EqualValues(t, 683, y)
	})

	t.Run("bottom-right", func(t *testing.T) {
		x, y := model.LatLngToTile(41.990226, 8.561345, 11)
		assert.EqualValues(t, 1072, x)
		assert.EqualValues(t, 760, y)
	})
}
