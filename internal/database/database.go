package database

import (
	"bytes"
	"database/sql"
	"errors"
	"image"
	"image/jpeg"
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

	c.db, err = sql.Open("sqlite3", c.filename)

	return err
}

func New(databaseName string) (*Connection, error) {
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

		sqliteDatabase, err := sql.Open("sqlite3", databaseName)
		if err != nil {
			return nil, err
		}

		for _, query := range []string{
			`CREATE TABLE metadata (name text, value text);`,
			`CREATE TABLE tiles (zoom_level integer, tile_column integer, tile_row integer, tile_data blob);`,
			`CREATE UNIQUE INDEX tile_index on tiles (zoom_level, tile_column, tile_row);`,
		} {
			statement, err := sqliteDatabase.Prepare(query)
			if err != nil {
				return nil, err
			}

			if _, err := statement.Exec(); err != nil {
				return nil, err
			}
		}

		if err := conn.MetadataRewrite(); err != nil {
			return nil, err
		}

		return conn, nil
	}

	return nil, err
}

func (c Connection) insertMetadata(name string, value string) error {
	statement, err := c.db.Prepare(`INSERT INTO metadata(name, value) VALUES (?, ?)`)
	if err != nil {
		return err
	}
	_, err = statement.Exec(name, value)
	if err != nil {
		return err
	}

	return nil
}

func (c Connection) InsertTile(img image.Image, zoomLevel uint64, col uint64, row uint64) error {
	statement, err := c.db.Prepare(`INSERT INTO tiles(zoom_level, tile_column, tile_row, tile_data) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return err
	}

	var imageBuf bytes.Buffer
	if err := jpeg.Encode(&imageBuf, img, &jpeg.Options{
		Quality: 100,
	}); err != nil {
		return err
	}

	_, err = statement.Exec(zoomLevel, col, row, imageBuf.Bytes())
	if err != nil {
		return err
	}

	return nil
}

// Close implements the Closer interface.
func (c Connection) Close() error {
	return c.db.Close()
}
