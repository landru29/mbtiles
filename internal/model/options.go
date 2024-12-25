package model

// Option is the application option.
type Option struct {
	CoordinateMin LatLng
	CoordinateMax LatLng
	ZoomMin       uint64
	ZoomMax       uint64
	Format        Format
	Name          string
	Description   string
}
