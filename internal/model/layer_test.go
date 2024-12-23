package model_test

import (
	"testing"

	"github.com/landru29/mbtiles/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestLatLngToTile(t *testing.T) {
	// Zoom 10
	// X: [989, 1087] a convertir
	// Y: [681, 772]
	t.Run("top-left zoom 11", func(t *testing.T) {
		layer := model.NewLayer(11, model.LatLng{}, model.LatLng{})

		x, y := layer.LatLngToTile(51.251834, -5.593299)

		assert.EqualValues(t, 992, x)

		assert.EqualValues(t, 683, y)
	})

	t.Run("bottom-right zoom 11", func(t *testing.T) {
		layer := model.NewLayer(11, model.LatLng{}, model.LatLng{})

		x, y := layer.LatLngToTile(41.990226, 8.561345)

		assert.EqualValues(t, 1072, x)

		assert.EqualValues(t, 760, y)
	})

	// Zoom 10
	// X: [496, 536]
	// Y: [643, 682]

	t.Run("top-left zoom 10", func(t *testing.T) {
		layer := model.NewLayer(10, model.LatLng{}, model.LatLng{})

		x, y := layer.LatLngToTile(51.251834, -5.593299)

		assert.EqualValues(t, 496, x)

		assert.EqualValues(t, 341, y)
	})

	t.Run("bottom-right zoom 10", func(t *testing.T) {
		layer := model.NewLayer(10, model.LatLng{}, model.LatLng{})

		x, y := layer.LatLngToTile(41.990226, 8.561345)

		assert.EqualValues(t, 536, x)

		assert.EqualValues(t, 380, y)
	})
}
