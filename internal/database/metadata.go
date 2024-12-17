package database

import "context"

func (c Connection) insertMetadata(ctx context.Context, metadata map[string]string) error {
	for name, value := range metadata {
		statement, err := c.db.Prepare(`INSERT INTO metadata(name, value) VALUES (?, ?)`)
		if err != nil {
			return err
		}
		_, err = statement.ExecContext(ctx, name, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c Connection) MetadataRewrite(ctx context.Context) error {
	statement, err := c.db.Prepare(`DELETE FROM metadata;`)
	if err != nil {
		return err
	}
	_, err = statement.ExecContext(ctx)
	if err != nil {
		return err
	}

	if err := c.insertMetadata(ctx, map[string]string{
		"bounds":      "-5.68558502538914379,41.8265329407147419,9.0168685878921071,51.2910711518300815",
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
	rows, err := c.db.QueryContext(ctx, "SELECT name, value FROM metadata")
	if err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	output := map[string]string{}

	for rows.Next() {
		var (
			key   string
			value string
		)

		if err := rows.Scan(
			&key,
			&value,
		); err != nil {
			return nil, err
		}

		output[key] = value
	}

	return output, nil
}
