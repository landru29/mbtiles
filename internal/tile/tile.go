package tile

import (
	"context"
	"errors"
	"fmt"
	"image"
	"time"

	"github.com/cenkalti/backoff"
)

type Box struct {
	ZoomLevel uint64
	RowMin    uint64
	RowMax    uint64
	ColMin    uint64
	ColMax    uint64
	colFail   map[uint64]uint64
	rowFail   map[uint64]uint64
}

func New(
	ZoomLevel uint64,
	RowMin uint64,
	RowMax uint64,
	ColMin uint64,
	ColMax uint64,
) Box {
	return Box{
		ZoomLevel: ZoomLevel,
		RowMin:    RowMin,
		RowMax:    RowMax,
		ColMin:    ColMin,
		ColMax:    ColMax,
		colFail:   map[uint64]uint64{},
		rowFail:   map[uint64]uint64{},
	}
}

func (b Box) ToZoom(zoomLevel uint64) (*Box, error) {
	if zoomLevel < b.ZoomLevel {
		return nil, errors.New("cannot decrease zoom")
	}

	coeficient := uint64(1)
	for range zoomLevel - b.ZoomLevel {
		coeficient *= 2
	}

	return &Box{
		ZoomLevel: zoomLevel,
		RowMin:    b.RowMin * coeficient,
		RowMax:    b.RowMax * coeficient,
		ColMin:    b.ColMin * coeficient,
		ColMax:    b.ColMax * coeficient,
		colFail:   map[uint64]uint64{},
		rowFail:   map[uint64]uint64{},
	}, nil
}

func (b Box) columns() []uint64 {
	output := []uint64{}
	for idx := b.ColMin; idx <= b.ColMax; idx++ {
		output = append(output, idx)
	}

	return output
}

func (b Box) rows() []uint64 {
	output := []uint64{}
	for idx := b.RowMin; idx <= b.RowMax; idx++ {
		output = append(output, idx)
	}

	return output
}

func (b *Box) Loop(ctx context.Context, loader Loader, processor func(img image.Image, zoomLevel uint64, col uint64, row uint64) error) error {
	for _, col := range b.columns() {
		for _, row := range b.rows() {
			attempt := 0

			if b.colFail[col] > 1 || b.rowFail[row] > 1 {
				fmt.Printf("too many errors on zoom:%d - row: %d - col: %d - skipping\n", b.ZoomLevel, row, col)

				continue
			}

			if err := backoff.Retry(func() error {
				if attempt != 0 {
					fmt.Printf("  #%d zoom:%d - row: %d - col: %d\n", attempt, b.ZoomLevel, row, col)
				}

				img, err := loader.LoadImage(ctx, b.ZoomLevel, col, row)
				if err != nil {
					attempt++

					return err
				}

				if err := processor(img, b.ZoomLevel, col, row); err != nil {
					return backoff.Permanent(err)
				}

				return nil
			}, backoff.WithMaxRetries(
				backoff.NewConstantBackOff(200*time.Millisecond),
				3,
			)); err != nil {
				fmt.Printf("ERROR zoom:%d - row: %d - col: %d => %s\n", b.ZoomLevel, row, col, err)

				b.colFail[col]++
				b.rowFail[row]++
			}
		}
	}

	return nil
}
