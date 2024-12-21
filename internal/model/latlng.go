package model

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	pkgerrors "github.com/pkg/errors"
)

// LatLng is a coordinate.
type LatLng struct {
	Lat float64
	Lng float64
}

// String implements the pflag.Value interface.
func (l *LatLng) String() string {
	if l == nil {
		return "nil"
	}

	return fmt.Sprintf("%f,%f", l.Lat, l.Lng)
}

// Set implements the pflag.Value interface.
func (l *LatLng) Set(data string) error {
	parts := strings.Split(data, ",")

	if len(parts) != 2 {
		return errors.New("should be <lat>,<lng>")
	}

	lat, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return pkgerrors.WithMessage(err, "cannot parse latitude")
	}

	lng, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return pkgerrors.WithMessage(err, "cannot parse longitude")
	}

	*l = LatLng{
		Lat: lat,
		Lng: lng,
	}

	return nil
}

// Type implements the pflag.Value interface.
func (l *LatLng) Type() string {
	return "LatLng"
}
