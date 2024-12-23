package app

import (
	"context"

	"github.com/landru29/mbtiles/internal/model"
)

// Tiles reads all tiles.
func (a Application) Tiles(ctx context.Context) (*model.TilesDescription, error) { //nolint: funlen,cyclop
	output := &model.TilesDescription{
		Col:          map[uint64][]uint64{},
		Row:          map[uint64][]uint64{},
		CountPerZoom: map[uint64]uint64{},
	}

	count, err := a.database.TilesCount(ctx)
	if err != nil {
		return nil, err
	}

	output.Count = count

	allTiles, err := a.database.AllTiles(ctx)
	if err != nil {
		return nil, err
	}

	if len(allTiles) == 0 {
		return output, nil
	}

	maxZoom := allTiles[0].ZoomLevel
	minZoom := allTiles[0].ZoomLevel
	minCol := map[uint64]uint64{}
	maxCol := map[uint64]uint64{}
	minRow := map[uint64]uint64{}
	maxRow := map[uint64]uint64{}

	for _, tile := range allTiles {
		output.CountPerZoom[tile.ZoomLevel]++

		if maxZoom < tile.ZoomLevel {
			maxZoom = tile.ZoomLevel
		}

		if minZoom > tile.ZoomLevel {
			minZoom = tile.ZoomLevel
		}

		if _, found := minCol[tile.ZoomLevel]; !found {
			minCol[tile.ZoomLevel] = tile.Col
		}

		if _, found := maxCol[tile.ZoomLevel]; !found {
			maxCol[tile.ZoomLevel] = tile.Col
		}

		if _, found := minRow[tile.ZoomLevel]; !found {
			minRow[tile.ZoomLevel] = tile.Row
		}

		if _, found := maxRow[tile.ZoomLevel]; !found {
			maxRow[tile.ZoomLevel] = tile.Row
		}

		if minCol[tile.ZoomLevel] > tile.Col {
			minCol[tile.ZoomLevel] = tile.Col
		}

		if maxCol[tile.ZoomLevel] < tile.Col {
			maxCol[tile.ZoomLevel] = tile.Col
		}

		if minRow[tile.ZoomLevel] > tile.Row {
			minRow[tile.ZoomLevel] = tile.Row
		}

		if maxRow[tile.ZoomLevel] < tile.Row {
			maxRow[tile.ZoomLevel] = tile.Row
		}
	}

	output.Zoom = []uint64{
		minZoom,
		maxZoom,
	}

	for idx := minZoom; idx <= maxZoom; idx++ {
		output.Col[idx] = []uint64{
			minCol[idx],
			maxCol[idx],
		}

		output.Row[idx] = []uint64{
			minRow[idx],
			maxRow[idx],
		}
	}

	return output, nil
}

// TileByIndex picks one tile with the specified index.
func (a Application) TileByIndex(ctx context.Context, index uint64) (*model.Tile, error) {
	tile, err := a.database.Tile(ctx, index)
	if err != nil {
		return nil, err
	}

	return tile, nil
}

// TileByCoordinates picks one tile with the specified coordinates.
func (a Application) TileByCoordinates(
	ctx context.Context,
	zoomLevel uint64,
	col uint64,
	row uint64,
) (*model.Tile, error) {
	tile, err := a.database.TileByCoordinate(ctx, model.TileRequest{
		ZoomLevel: zoomLevel,
		Col:       col,
		Row:       row,
	})
	if err != nil {
		return nil, err
	}

	return tile, nil
}

// TileRewrite revrites the tile to PNG format.
func (a Application) TileRewrite(ctx context.Context) error {
	allTiles, err := a.database.AllTiles(ctx)
	if err != nil {
		return err
	}

	return a.database.TileToPNG(ctx, a.display, allTiles)
}
