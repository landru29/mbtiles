// Package sqlite is the SQLite3 implementation.
package sqlite

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"os"
	"path/filepath"

	"github.com/landru29/mbtiles/internal/database/sqlite/sqlc"
	"github.com/landru29/mbtiles/internal/model"
	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
	pkgerrors "github.com/pkg/errors"
	migrate "github.com/rubenv/sql-migrate"
)

//go:generate sqlc generate --file ./sqlc.yaml

//go:embed migrations/*.sql
var migrationsFS embed.FS

// Connection is the sqlite implementation of the database.
type Connection struct {
	db       *sql.DB
	sqlc     *sqlc.Queries
	filename string
}

func createDBfile(databaseName string) error {
	file, err := os.Create(filepath.Clean(databaseName)) // Create SQLite file
	if err != nil {
		return err
	}

	return file.Close()
}

// Open opens the database.
func (c *Connection) Open() error {
	var err error

	c.db, err = sql.Open("sqlite3", filepath.Clean(c.filename)+"?_auto_vacuum=1")

	c.sqlc = sqlc.New(c.db)

	return err
}

// New creates the database.
func New(ctx context.Context, databaseName string, minCoord model.LatLng, maxCood model.LatLng) (*Connection, error) {
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

		if err := conn.initDatabase(ctx, minCoord, maxCood); err != nil {
			return nil, err
		}

		return conn, nil
	}

	return nil, err
}

func (c Connection) initDatabase(ctx context.Context, minCoord model.LatLng, maxCood model.LatLng) error {
	migrations := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: migrationsFS,
		Root:       "migrations",
	}

	_, err := migrate.Exec(c.db, "sqlite3", migrations, migrate.Up)
	if err != nil {
		return pkgerrors.WithMessage(err, "cannot setup MbTiles")
	}

	if err := c.populateAndroidMetadata(ctx); err != nil {
		return err
	}

	if err := c.MetadataRewrite(ctx, minCoord, maxCood); err != nil {
		return err
	}

	return nil
}

// Close implements the Closer interface.
func (c Connection) Close() error {
	return c.db.Close()
}
