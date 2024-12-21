// Package model is the data model.
package model

import (
	"image"

	"golang.org/x/exp/constraints"
)

// TileSample represents one tile.
type TileSample struct {
	ZoomLevel uint64
	Row       uint64
	Col       uint64
	Image     image.Image
	Type      string
}

// TilesDescription is the set of tiles description.
type TilesDescription struct {
	Count        uint64
	Zoom         []uint64
	Col          map[uint64][]uint64
	Row          map[uint64][]uint64
	CountPerZoom map[uint64]uint64
}

// TileRequest is a request for a tile.
type TileRequest struct {
	ZoomLevel uint64
	Row       uint64
	Col       uint64
}

// Min is the generic min function.
func Min[P constraints.Ordered](first P, second P) P { //nolint: ireturn
	if first > second {
		return second
	}

	return first
}

// Max is the generic miax function.
func Max[P constraints.Ordered](first P, second P) P { //nolint: ireturn
	if first < second {
		return second
	}

	return first
}
