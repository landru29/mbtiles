package app

import (
	"context"
	"image"

	"github.com/landru29/mbtiles/internal/model"
	"github.com/landru29/mbtiles/internal/tile/oaci"
)

// Download downloads one tile.
func (a Application) Download(ctx context.Context, coordinate model.LatLng, zoomLevel uint64) (image.Image, error) {
	col, row := model.Layer{ZoomLevel: zoomLevel}.LatLngToTile(coordinate.Lat, coordinate.Lng)

	return oaci.Client{}.LoadImage(ctx, model.TileRequest{
		ZoomLevel: zoomLevel,
		Row:       row,
		Col:       col,
	})
}
