package sqlite

import (
	"context"
	"fmt"

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
func (c Connection) MetadataRewrite(ctx context.Context, minCoord model.LatLng, maxCood model.LatLng) error {
	if err := c.sqlc.WipeAllMetadata(ctx); err != nil {
		return pkgerrors.WithMessage(err, "cannot wipe all metadata")
	}

	if err := c.insertMetadata(ctx, map[string]string{
		"bounds": fmt.Sprintf(
			"%f,%f,%f,%f",
			model.Min(minCoord.Lng, maxCood.Lng),
			model.Min(minCoord.Lat, maxCood.Lat),
			model.Max(minCoord.Lng, maxCood.Lng),
			model.Max(minCoord.Lat, maxCood.Lat),
		),
		"name":        "oaci_1_250",
		"format":      "png",
		"minzoom":     "6",
		"maxzoom":     "11",
		"type":        "overlay",
		"description": "SIA France",
		"version":     "1.1",
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
