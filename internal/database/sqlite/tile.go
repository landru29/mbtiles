package sqlite

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"io"

	"github.com/landru29/mbtiles/internal/database/sqlite/sqlc"
	"github.com/landru29/mbtiles/internal/model"
	pkgerrors "github.com/pkg/errors"
)

// TilesCount counts all the tiles.
func (c Connection) TilesCount(ctx context.Context) (uint64, error) {
	count, err := c.sqlc.TileCount(ctx)
	if err != nil {
		return 0, pkgerrors.WithMessage(err, "cannot count tiles")
	}

	return count, nil
}

// Tile picks on tile with its index.
func (c Connection) Tile(ctx context.Context, index uint64) (*model.TileSample, error) {
	tile, err := c.sqlc.TileByIndex(ctx, index)
	if err != nil {
		return nil, pkgerrors.WithMessage(err, "cannot read tile")
	}

	out, format, err := image.Decode(bytes.NewBuffer(tile.TileData))
	if err != nil {
		return nil, err
	}

	return &model.TileSample{
		Image:     out,
		Type:      format,
		ZoomLevel: tile.ZoomLevel,
		Row:       tile.TileRow,
		Col:       tile.TileColumn,
	}, err
}

// TileByCoordinate picks one tile with the specified coordinates.
func (c Connection) TileByCoordinate(ctx context.Context, request model.TileRequest) (*model.TileSample, error) {
	tile, err := c.sqlc.TileByCoordinate(ctx, sqlc.TileByCoordinateParams{
		Col:       request.Col,
		Row:       request.Row,
		ZoomLevel: request.ZoomLevel,
	})
	if err != nil {
		return nil, pkgerrors.WithMessage(err, "cannot read tile with specified coordinates")
	}

	out, format, err := image.Decode(bytes.NewBuffer(tile.TileData))
	if err != nil {
		return nil, err
	}

	return &model.TileSample{
		Image:     out,
		Type:      format,
		ZoomLevel: tile.ZoomLevel,
		Row:       tile.TileRow,
		Col:       tile.TileColumn,
	}, err
}

// AllTiles picks all tiles.
func (c Connection) AllTiles(ctx context.Context) ([]model.TileSample, error) {
	allTiles, err := c.sqlc.Tiles(ctx)
	if err != nil {
		return nil, pkgerrors.WithMessage(err, "cannot read all tiles")
	}

	output := make([]model.TileSample, len(allTiles))

	for idx, element := range allTiles {
		out, format, err := image.Decode(bytes.NewBuffer(element.TileData))
		if err != nil {
			return nil, err
		}

		output[idx] = model.TileSample{
			ZoomLevel: element.ZoomLevel,
			Row:       element.TileRow,
			Col:       element.TileColumn,
			Image:     out,
			Type:      format,
		}
	}

	return output, err
}

// TileToPNG rewrites a tile in PNG format.
func (c Connection) TileToPNG(ctx context.Context, display io.Writer, tiles []model.TileSample) error {
	for idx, tile := range tiles {
		var buffer bytes.Buffer

		if err := png.Encode(&buffer, tile.Image); err != nil {
			return pkgerrors.WithMessage(err, "cannot encode image")
		}

		_, _ = fmt.Fprintf(display, "#%d Zoom: %d - Row: %d - Col: %d\n", idx, tile.ZoomLevel, tile.Row, tile.Col)

		if err := c.sqlc.TileDataUpdate(ctx, sqlc.TileDataUpdateParams{
			TileData:  buffer.Bytes(),
			Col:       tile.Col,
			Row:       tile.Row,
			ZoomLevel: tile.ZoomLevel,
		}); err != nil {
			return pkgerrors.WithMessage(err, "cannot update tile data")
		}
	}

	return nil
}

// InsertTile adds a new tile.
func (c Connection) InsertTile(ctx context.Context, tile model.TileSample) error {
	statement, err := c.db.Prepare(`INSERT INTO tiles(zoom_level, tile_column, tile_row, tile_data) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return err
	}

	defer func() {
		_ = statement.Close()
	}()

	var imageBuf bytes.Buffer
	if err := png.Encode(&imageBuf, tile.Image); err != nil {
		return err
	}

	_, err = statement.ExecContext(ctx, tile.ZoomLevel, tile.Col, tile.Row, imageBuf.Bytes())
	if err != nil {
		return err
	}

	return nil
}
