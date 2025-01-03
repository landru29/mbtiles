package model

import (
	"math"
)

// Layer is a layer on a specific zoom level.
type Layer struct {
	ZoomLevel uint64
	LngMin    Coordinate
	LngMax    Coordinate
	LatMin    Coordinate
	LatMax    Coordinate
}

// NewLayer creates a layer.
func NewLayer(
	zoomLevel uint64,
	minLatLng LatLng,
	maxLatLng LatLng,
) Layer {
	return Layer{
		ZoomLevel: zoomLevel,
		LngMin:    Min(minLatLng.Lng, maxLatLng.Lng),
		LngMax:    Max(minLatLng.Lng, maxLatLng.Lng),
		LatMin:    Min(minLatLng.Lat, maxLatLng.Lat),
		LatMax:    Max(minLatLng.Lat, maxLatLng.Lat),
	}
}

// ToZoom redefine the layer to the specific zoom level.
func (l Layer) ToZoom(zoomLevel uint64) (*Layer, error) {
	return &Layer{
		ZoomLevel: zoomLevel,
		LngMin:    l.LngMin,
		LngMax:    l.LngMax,
		LatMin:    l.LatMin,
		LatMax:    l.LatMax,
	}, nil
}

// Columns extrats all the columns.
func (l Layer) Columns() []uint64 {
	output := []uint64{}
	for idx := l.ColMin(); idx <= l.ColMax(); idx++ {
		output = append(output, idx)
	}

	return output
}

// Rows extracts all the rows.
func (l Layer) Rows() []uint64 {
	output := []uint64{}
	for idx := l.RowMin(); idx <= l.RowMax(); idx++ {
		output = append(output, idx)
	}

	return output
}

// RowMin is the tile minimum row.
func (l Layer) RowMin() uint64 {
	return Min(l.YTile(l.LatMin), l.YTile(l.LatMax))
}

// RowMax is the tile maximum row.
func (l Layer) RowMax() uint64 {
	return Max(l.YTile(l.LatMin), l.YTile(l.LatMax))
}

// ColMin is the tile minimum column.
func (l Layer) ColMin() uint64 {
	return Min(l.XTile(l.LngMin), l.XTile(l.LngMax))
}

// ColMax is the tile maximum column.
func (l Layer) ColMax() uint64 {
	return Max(l.XTile(l.LngMin), l.XTile(l.LngMax))
}

// XTile convert longitude to tile column.
func (l Layer) XTile(lng Coordinate) uint64 {
	coef := float64(int(1) << l.ZoomLevel)

	return uint64(coef * (float64(lng) + 180.0) / 360.0)
}

// YTile converts latitude to tile row.
func (l Layer) YTile(lat Coordinate) uint64 {
	coef := float64(int(1) << (l.ZoomLevel - 1))

	latRad := float64(lat) * math.Pi / 180.0

	yTile := coef * (1.0 - (math.Log(math.Tan(latRad)+1.0/math.Cos(latRad)) / math.Pi))

	return uint64(yTile)
}

// LatLngToTile converts lat-lng coordinates to tile.
// Returns col, row.
func (l Layer) LatLngToTile(coord LatLng) (uint64, uint64) {
	return l.XTile(coord.Lng), l.YTile(coord.Lat)
}
