// Package app is the main application.
package app

import (
	"io"

	"github.com/landru29/mbtiles/internal/database"
)

// Application is the main application.
type Application struct {
	database database.Connection
	display  io.Writer
}

// New creates the application.
func New(database database.Connection, display io.Writer) *Application {
	if display == nil {
		display = io.Discard
	}

	return &Application{
		database: database,
		display:  display,
	}
}

// Close implements the io.Closer interface.
func (a Application) Close() error {
	return a.database.Close()
}
