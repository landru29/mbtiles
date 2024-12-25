// Package app is the main application.
package app

import (
	"io"
	"sync"

	"github.com/landru29/mbtiles/internal/database"
	"github.com/landru29/mbtiles/internal/model"
)

// Application is the main application.
type Application struct {
	database          database.Connection
	display           io.Writer
	maxDetectedZoom   uint64
	minDetectedZoom   uint64
	detectedFormat    model.Format
	zoomDetectionLock *sync.Mutex
}

// New creates the application.
func New(database database.Connection, display io.Writer) *Application {
	if display == nil {
		display = io.Discard
	}

	return &Application{
		database:          database,
		display:           display,
		zoomDetectionLock: &sync.Mutex{},
	}
}

// Close implements the io.Closer interface.
func (a Application) Close() error {
	return a.database.Close()
}
