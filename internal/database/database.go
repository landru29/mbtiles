// Package database defines the database requirements.
package database

import (
	"context"
	"io"

	"github.com/landru29/mbtiles/internal/model"
)

// Connection is a database connection.
type Connection interface {
	io.Closer
	MetadataRewrite(ctx context.Context, minCoord model.LatLng, maxCood model.LatLng) error
	Metadata(ctx context.Context) (map[string]string, error)
	TilesCount(ctx context.Context) (uint64, error)
	Tile(ctx context.Context, index uint64) (*model.TileSample, error)
	TileByCoordinate(ctx context.Context, request model.TileRequest) (*model.TileSample, error)
	AllTiles(ctx context.Context) ([]model.TileSample, error)
	TileToPNG(ctx context.Context, display io.Writer, tiles []model.TileSample) error
	InsertTile(ctx context.Context, tile model.TileSample) error
}
