package app

import (
	"context"
	"fmt"
	"strconv"

	"github.com/landru29/mbtiles/internal/database"
	"github.com/landru29/mbtiles/internal/model"
	"github.com/landru29/mbtiles/internal/tile"
	"github.com/landru29/mbtiles/internal/tile/oaci"
	pkgerrors "github.com/pkg/errors"
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
			func(index int, tile model.Tile) error {
				_, _ = fmt.Fprintf(
					a.display,
					"#%d | 🔍%d - ↓%d/%d - →%d/%d (%d, %d)\n",
					index,
					tile.ZoomLevel,
					tile.Row,
					currentLayer.RowMax(),
					tile.Col,
					currentLayer.ColMax(),
					tile.Image.Bounds().Max.X,
					tile.Image.Bounds().Max.Y,
				)

				if err := a.database.InsertTile(ctx, tile.TMS()); err != nil {
					return pkgerrors.WithMessage(err, "cannot insert tile")
				}

				a.zoomDetectionLock.Lock()
				defer a.zoomDetectionLock.Unlock()

				if a.maxDetectedZoom == 0 || a.maxDetectedZoom < tile.ZoomLevel {
					a.maxDetectedZoom = tile.ZoomLevel
				}

				if a.minDetectedZoom == 0 || a.minDetectedZoom > tile.ZoomLevel {
					a.minDetectedZoom = tile.ZoomLevel
				}

				if a.detectedFormat == "" {
					_, _ = fmt.Fprintf(a.display, "detected format: %s\n", tile.OriginalFormat)
				}

				a.detectedFormat = options.Format
				if options.Format == model.FormatNoTransform {
					a.detectedFormat = tile.OriginalFormat
				}

				return nil
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

	if err := a.database.UpdateMetadata(
		ctx,
		database.MetadataMaxzoom,
		strconv.FormatUint(a.maxDetectedZoom, 10),
	); err != nil {
		return pkgerrors.WithMessage(err, "cannot set max zoom")
	}

	if err := a.database.UpdateMetadata(
		ctx,
		database.MetadataMinzoom,
		strconv.FormatUint(a.minDetectedZoom, 10),
	); err != nil {
		return pkgerrors.WithMessage(err, "cannot set min zoom")
	}

	if err := a.database.UpdateMetadata(
		ctx,
		database.MetadataFormat,
		a.detectedFormat.String(),
	); err != nil {
		return pkgerrors.WithMessage(err, "cannot set min zoom")
	}

	return nil
}
