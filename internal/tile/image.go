package tile

import (
	"context"
	"image"
)

type Loader interface {
	LoadImage(ctx context.Context, zoomLevel uint64, col uint64, row uint64) (image.Image, error)
}
