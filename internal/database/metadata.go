package database

func (c Connection) MetadataRewrite() error {
	statement, err := c.db.Prepare(`TRUNCATE TABLE metadata`)
	if err != nil {
		return err
	}
	_, err = statement.Exec()
	if err != nil {
		return err
	}

	if err := c.insertMetadata("name", "oaci"); err != nil {
		return err
	}

	if err := c.insertMetadata("format", "jpg"); err != nil {
		return err
	}

	if err := c.insertMetadata("minzoom", "6"); err != nil {
		return err
	}

	if err := c.insertMetadata("maxzoom", "11"); err != nil {
		return err
	}

	return nil
}

// Metadata reads all the metadata.
func (c Connection) Metadata() (map[string]string, error) {
	rows, err := c.db.Query("SELECT name, value FROM metadata")
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
