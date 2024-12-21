package tile

import (
	"context"
	"image"

	"github.com/landru29/mbtiles/internal/model"
)

// Loader is a tile loader.
type Loader interface {
	LoadImage(ctx context.Context, request model.TileRequest) (image.Image, error)
}
