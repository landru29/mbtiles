package model

import (
	"errors"
	"math"
)

// Layer is a layer on a specific zoom level.
type Layer struct {
	ZoomLevel uint64
	RowMin    uint64
	RowMax    uint64
	ColMin    uint64
	ColMax    uint64
}

// New creates a layer.
func New(
	zoomLevel uint64,
	rowMin uint64,
	rowMax uint64,
	colMin uint64,
	colMax uint64,
) Layer {
	return Layer{
		ZoomLevel: zoomLevel,
		RowMin:    rowMin,
		RowMax:    rowMax,
		ColMin:    colMin,
		ColMax:    colMax,
	}
}

// NewFromLatLng creates a layer from lat-lng coordinates.
func NewFromLatLng(zoomLevel uint64, minCoord LatLng, maxCoord LatLng) Layer {
	colMin, rowMin := LatLngToTile(minCoord.Lat, minCoord.Lng, zoomLevel)
	colMax, rowMax := LatLngToTile(maxCoord.Lat, maxCoord.Lng, zoomLevel)

	return Layer{
		ZoomLevel: zoomLevel,
		RowMin:    Min(rowMin, rowMax),
		RowMax:    Max(rowMin, rowMax),
		ColMin:    Min(colMin, colMax),
		ColMax:    Max(colMin, colMax),
	}
}

// ToZoom redefine the layer to the specific zoom level.
func (l Layer) ToZoom(zoomLevel uint64) (*Layer, error) {
	if zoomLevel < l.ZoomLevel {
		return nil, errors.New("cannot decrease zoom")
	}

	coeficient := uint64(1)
	for range zoomLevel - l.ZoomLevel {
		coeficient *= 2
	}

	return &Layer{
		ZoomLevel: zoomLevel,
		RowMin:    l.RowMin * coeficient,
		RowMax:    l.RowMax * coeficient,
		ColMin:    l.ColMin * coeficient,
		ColMax:    l.ColMax * coeficient,
	}, nil
}

// Columns extrats all the columns.
func (l Layer) Columns() []uint64 {
	output := []uint64{}
	for idx := l.ColMin; idx <= l.ColMax; idx++ {
		output = append(output, idx)
	}

	return output
}

// Rows extracts all the rows.
func (l Layer) Rows() []uint64 {
	output := []uint64{}
	for idx := l.RowMin; idx <= l.RowMax; idx++ {
		output = append(output, idx)
	}

	return output
}

// LatLngToTile converts lat-lng coordinates to tile.
func LatLngToTile(lat float64, lng float64, zoomLevel uint64) (uint64, uint64) {
	coef := float64(1)
	for range zoomLevel {
		coef *= 2
	}

	latRad := lat * math.Pi / 180.0

	xTile := coef * (lng + 180) / 360
	yTile := coef * (1 - (math.Log(math.Tan(latRad)+1/math.Cos(latRad)) / math.Pi)) / 2

	return uint64(xTile), uint64(yTile)
}
