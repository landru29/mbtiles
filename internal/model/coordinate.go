package model

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	pkgerrors "github.com/pkg/errors"
)

const (
	coordinateRegexp = `(-?\d+)°((\d+)')?([\d.]+)?([NSEW])?`
)

// Coordinate is a GNSS coordinate.
type Coordinate float64

// String implements the pflag.Value interface.
func (l *Coordinate) String() string {
	if l == nil {
		return "nil"
	}

	value := math.Round(float64(*l)*float64(1e9)) / float64(1e9)

	output := ""

	sign := math.Signbit(value)
	value = math.Abs(value)

	deg := int(value)
	minutes := int((value - float64(deg)) * 60.0)
	seconds := math.Round((value-float64(deg)-float64(minutes)/60.0)*3600.0*float64(1e4)) / float64(1e4)

	if deg != 0 {
		output = fmt.Sprintf("%d°", deg)
	}

	if minutes != 0 {
		output = fmt.Sprintf("%s%d'", output, minutes)
	}

	if seconds != 0 {
		output = fmt.Sprintf("%s%f", output, seconds)
	}

	return map[bool]string{true: "-"}[sign] + output
}

// Set implements the pflag.Value interface.
func (l *Coordinate) Set(data string) error {
	if value, err := strconv.ParseFloat(data, 64); err == nil {
		*l = Coordinate(value)

		return nil
	}

	regxp, err := regexp.Compile(coordinateRegexp) //nolint: gocritic
	if err != nil {
		return pkgerrors.WithMessage(err, "cannot build coordinate regexp")
	}

	matcher := regxp.FindAllStringSubmatch(data, -1)
	if len(matcher) != 1 || len(matcher[0]) != 6 {
		return fmt.Errorf("cannot parse coordinate %s", data)
	}

	sign := 1.0

	if strings.ToUpper(matcher[0][5]) == "W" || strings.ToUpper(matcher[0][5]) == "S" {
		sign = -1.0
	}

	if matcher[0][1] != "" {
		deg, err := strconv.ParseInt(matcher[0][1], 10, 64)
		if err != nil {
			return pkgerrors.WithMessage(err, "cannot parse degree")
		}

		if deg < 0 {
			sign *= -1.0

			deg = -deg
		}

		*l = Coordinate(float64(deg))
	}

	if matcher[0][3] != "" {
		minutes, err := strconv.ParseInt(matcher[0][3], 10, 64)
		if err != nil {
			return pkgerrors.WithMessage(err, "cannot parse minutes")
		}

		*l += Coordinate(float64(minutes) / 60.0)
	}

	if matcher[0][4] != "" {
		seconds, err := strconv.ParseFloat(matcher[0][4], 64)
		if err != nil {
			return pkgerrors.WithMessage(err, "cannot parse seconds")
		}

		*l += Coordinate(float64(seconds) / 3600.0)
	}

	*l *= Coordinate(sign)

	return nil
}

// Type implements the pflag.Value interface.
func (l *Coordinate) Type() string {
	return "Coordinate"
}
