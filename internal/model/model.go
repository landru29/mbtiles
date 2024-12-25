// Package model is the data model.
package model

import (
	"image"

	"golang.org/x/exp/constraints"
)

// Tile represents one tile.
type Tile struct {
	ZoomLevel      uint64
	Row            uint64
	Col            uint64
	Image          image.Image
	RawImage       []byte
	OriginalFormat Format
	Type           Format
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

// TMS converts a request to (Tile Map Service coordinate.
func (t TileRequest) TMS() TileRequest {
	return TileRequest{
		ZoomLevel: t.ZoomLevel,
		Row:       uint64(int64(1)<<t.ZoomLevel) - 1 - t.Row,
		Col:       t.Col,
	}
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

// TMS converts a tile to (Tile Map Service coordinate.
func (t Tile) TMS() Tile {
	return Tile{
		ZoomLevel: t.ZoomLevel,
		Row:       uint64(int64(1)<<t.ZoomLevel) - 1 - t.Row,
		Col:       t.Col,
		Image:     t.Image,
		Type:      t.Type,
	}
}
