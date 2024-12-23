package model

import (
	"errors"
	"fmt"
	"strings"

	pkgerrors "github.com/pkg/errors"
)

// LatLng is a coordinate.
type LatLng struct {
	Lat Coordinate
	Lng Coordinate
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

	if err := l.Lat.Set(parts[0]); err != nil {
		return pkgerrors.WithMessage(err, "cannot parse latitude")
	}

	if err := l.Lng.Set(parts[1]); err != nil {
		return pkgerrors.WithMessage(err, "cannot parse longitude")
	}

	return nil
}

// Type implements the pflag.Value interface.
func (l *LatLng) Type() string {
	return "LatLng"
}
