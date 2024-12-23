package app

import (
	"context"
	"fmt"

	"github.com/landru29/mbtiles/internal/model"
	"github.com/landru29/mbtiles/internal/tile"
	"github.com/landru29/mbtiles/internal/tile/oaci"
)

// Generate launch the MbTiles file generation.
func (a Application) Generate(
	ctx context.Context,
	options model.Option,
	workerCount int,
) error {
	currentLayer := model.NewLayer(
		options.ZoomMax,
		options.CoordinateMin,
		options.CoordinateMax,
	)

	defer func() {
		_ = a.database.Close()
	}()

	for currentLayer.ZoomLevel >= options.ZoomMin {
		if err := tile.Loop(
			ctx,
			currentLayer,
			oaci.Client{},
			func(tile model.Tile) error {
				_, _ = fmt.Fprintf(
					a.display,
					"🔍%d - ↓%d/%d - →%d/%d (%d, %d)\n",
					tile.ZoomLevel,
					tile.Row,
					currentLayer.RowMax(),
					tile.Col,
					currentLayer.ColMax(),
					tile.Image.Bounds().Max.X,
					tile.Image.Bounds().Max.Y,
				)

				return a.database.InsertTile(ctx, tile.TMS())
			},
			workerCount,
			a.display,
		); err != nil {
			return err
		}

		nextBox, err := currentLayer.ToZoom(currentLayer.ZoomLevel - 1)
		if err != nil {
			return err
		}

		currentLayer = *nextBox
	}

	return nil
}
