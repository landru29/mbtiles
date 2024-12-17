package database

import "context"

func (c Connection) populateAndroidMetadata(ctx context.Context) error {
	statement, err := c.db.Prepare(`INSERT INTO android_metadata(locale) VALUES (?)`)
	if err != nil {
		return err
	}
	_, err = statement.ExecContext(ctx, "fr_FR")
	if err != nil {
		return err
	}

	return nil
}
