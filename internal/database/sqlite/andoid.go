package sqlite

import (
	"context"

	pkgerrors "github.com/pkg/errors"
)

func (c Connection) populateAndroidMetadata(ctx context.Context) error {
	if err := c.sqlc.PopulateAndroidMetadata(ctx, "fr_FR"); err != nil {
		return pkgerrors.WithMessage(err, "cannot populate android metadata")
	}

	return nil
}
