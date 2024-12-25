// Package database defines the database requirements.
package database

import (
	"context"
	"io"

	"github.com/landru29/mbtiles/internal/model"
)

const (
	// MetadataBounds is the Bounds metadata name.
	MetadataBounds = "bounds"
	// MetadataName is the Name metadata name.
	MetadataName = "name"
	// MetadataFormat is the Format metadata name.
	MetadataFormat = "format"
	// MetadataMinzoom is the Minzoom metadata name.
	MetadataMinzoom = "minzoom"
	// MetadataMaxzoom is the Maxzoom metadata name.
	MetadataMaxzoom = "maxzoom"
	// MetadataType is the Type metadata name.
	MetadataType = "type"
	// MetadataDescription is the Description metadata name.
	MetadataDescription = "description"
	// MetadataVersion is the Version metadata name.
	MetadataVersion = "version"
)

// Connection is a database connection.
type Connection interface {
	io.Closer
	MetadataRewrite(
		ctx context.Context,
		options model.Option,
	) error
	Metadata(ctx context.Context) (map[string]string, error)
	TilesCount(ctx context.Context) (uint64, error)
	Tile(ctx context.Context, index uint64) (*model.Tile, error)
	TileByCoordinate(ctx context.Context, request model.TileRequest) (*model.Tile, error)
	AllTiles(ctx context.Context) ([]model.Tile, error)
	InsertTile(ctx context.Context, tile model.Tile) error
	UpdateMetadata(ctx context.Context, name string, value string) error
}
