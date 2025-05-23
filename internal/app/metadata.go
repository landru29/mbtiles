package app

import (
	"context"

	"github.com/landru29/mbtiles/internal/model"
)

// Metadata reads all the metadata.
func (a Application) Metadata(ctx context.Context) (map[string]string, error) {
	metadata, err := a.database.Metadata(ctx)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

// MetadataRewrite rewritres the correct metadata.
func (a Application) MetadataRewrite(
	ctx context.Context,
	options model.Option,
) error {
	if err := a.database.MetadataRewrite(ctx, options); err != nil {
		return err
	}

	return nil
}
