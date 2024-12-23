package sqlite

import (
	"context"
	"fmt"
	"strconv"

	"github.com/landru29/mbtiles/internal/database/sqlite/sqlc"
	"github.com/landru29/mbtiles/internal/model"
	pkgerrors "github.com/pkg/errors"
)

func (c Connection) insertMetadata(ctx context.Context, metadata map[string]string) error {
	for name, value := range metadata {
		if err := c.sqlc.InsertMetadata(ctx, sqlc.InsertMetadataParams{
			Name:  name,
			Value: value,
		}); err != nil {
			return pkgerrors.WithMessage(err, "cannot insert metadata")
		}
	}

	return nil
}

// MetadataRewrite rewrites the correct metadata.
func (c Connection) MetadataRewrite(
	ctx context.Context,
	options model.Option,
) error {
	if err := c.sqlc.WipeAllMetadata(ctx); err != nil {
		return pkgerrors.WithMessage(err, "cannot wipe all metadata")
	}

	if err := c.insertMetadata(ctx, map[string]string{
		"bounds": fmt.Sprintf(
			"%f,%f,%f,%f",
			model.Min(options.CoordinateMin.Lng, options.CoordinateMax.Lng),
			model.Min(options.CoordinateMin.Lat, options.CoordinateMax.Lat),
			model.Max(options.CoordinateMin.Lng, options.CoordinateMax.Lng),
			model.Max(options.CoordinateMin.Lat, options.CoordinateMax.Lat),
		),
		"name":        "oaci_1_250",
		"format":      options.Format,
		"minzoom":     strconv.FormatUint(options.ZoomMin, 10),
		"maxzoom":     strconv.FormatUint(options.ZoomMax, 10),
		"type":        "overlay",
		"description": "SIA France",
		"version":     "1.3",
	}); err != nil {
		return err
	}

	return nil
}

// Metadata reads all the metadata.
func (c Connection) Metadata(ctx context.Context) (map[string]string, error) {
	metadata, err := c.sqlc.Metadata(ctx)
	if err != nil {
		return nil, pkgerrors.WithMessage(err, "cannot read metadata")
	}

	output := map[string]string{}

	for _, element := range metadata {
		output[element.Name] = element.Value
	}

	return output, nil
}
