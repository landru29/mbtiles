package database

import (
	"context"
	"database/sql"
	"errors"
	"os"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

type Connection struct {
	db       *sql.DB
	filename string
}

func createDBfile(databaseName string) error {
	file, err := os.Create(databaseName) // Create SQLite file
	if err != nil {
		return err
	}

	return file.Close()
}

// Open opens the database.
func (c *Connection) Open() error {
	var err error

	c.db, err = sql.Open("sqlite3", c.filename+"?_auto_vacuum=1")

	return err
}

func New(ctx context.Context, databaseName string) (*Connection, error) {
	conn := &Connection{
		filename: databaseName,
	}

	_, err := os.Stat(databaseName)
	if err == nil {
		if err := conn.Open(); err != nil {
			return nil, err
		}

		return conn, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		if err := createDBfile(databaseName); err != nil {
			return nil, err
		}

		if err := conn.Open(); err != nil {
			return nil, err
		}

		for _, query := range []string{
			`CREATE TABLE metadata (name text, value text);`,
			`CREATE TABLE tiles (zoom_level INTEGER NOT NULL,tile_column INTEGER NOT NULL,tile_row INTEGER NOT NULL,tile_data BLOB NOT NULL,UNIQUE (zoom_level, tile_column, tile_row) );`,
			`CREATE TABLE android_metadata (locale TEXT);`,
		} {
			statement, err := conn.db.Prepare(query)
			if err != nil {
				return nil, err
			}

			if _, err := statement.ExecContext(ctx); err != nil {
				return nil, err
			}
		}

		if err := conn.populateAndroidMetadata(ctx); err != nil {
			return nil, err
		}

		if err := conn.MetadataRewrite(ctx); err != nil {
			return nil, err
		}

		return conn, nil
	}

	return nil, err
}

// Close implements the Closer interface.
func (c Connection) Close() error {
	return c.db.Close()
}
