package model_test

import (
	"testing"

	"github.com/landru29/mbtiles/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLatLngToTile(t *testing.T) {
	// Zoom 10
	// X: [989, 1087] a convertir
	// Y: [681, 772]
	t.Run("top-left zoom 11", func(t *testing.T) {
		layer := model.NewLayer(11, model.LatLng{}, model.LatLng{})

		x, y := layer.LatLngToTile(model.LatLng{Lat: 51.251834, Lng: -5.593299})

		assert.EqualValues(t, 992, x)

		assert.EqualValues(t, 683, y)
	})

	t.Run("bottom-right zoom 11", func(t *testing.T) {
		layer := model.NewLayer(11, model.LatLng{}, model.LatLng{})

		x, y := layer.LatLngToTile(model.LatLng{Lat: 41.990226, Lng: 8.561345})

		assert.EqualValues(t, 1072, x)

		assert.EqualValues(t, 760, y)
	})

	// Zoom 10
	// X: [496, 536]
	// Y: [643, 682]

	t.Run("top-left zoom 10", func(t *testing.T) {
		layer := model.NewLayer(10, model.LatLng{}, model.LatLng{})

		x, y := layer.LatLngToTile(model.LatLng{Lat: 51.251834, Lng: -5.593299})

		assert.EqualValues(t, 496, x)

		assert.EqualValues(t, 341, y)
	})

	t.Run("bottom-right zoom 10", func(t *testing.T) {
		layer := model.NewLayer(10, model.LatLng{}, model.LatLng{})

		x, y := layer.LatLngToTile(model.LatLng{Lat: 41.990226, Lng: 8.561345})

		assert.EqualValues(t, 536, x)

		assert.EqualValues(t, 380, y)
	})

	t.Run("Ouessant zoom 11", func(t *testing.T) {
		layer := model.NewLayer(9, model.LatLng{}, model.LatLng{})

		coordinate := model.LatLng{}
		require.NoError(t, coordinate.Set("48°27'48N,W 005°03'49W"))

		x, y := layer.LatLngToTile(coordinate)

		assert.EqualValues(t, 248, x)

		assert.EqualValues(t, 176, y)
	})
}
