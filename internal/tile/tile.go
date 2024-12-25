// Package tile manages tile layers.
package tile

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/landru29/mbtiles/internal/model"
)

// Loop iterates through a layer to execute a processor.
func Loop(
	ctx context.Context,
	layer model.Layer,
	loader Loader,
	processor func(workerIndex int, tile model.Tile) error,
	workerCount int,
	display io.Writer,
) error {
	var wait sync.WaitGroup

	requester := make(chan *model.TileRequest, 1)

	if display == nil {
		display = io.Discard
	}

	for idx := range workerCount {
		wait.Add(1)

		go worker(ctx, idx, &wait, requester, loader, processor, display)
	}

	for _, col := range layer.Columns() {
		for _, row := range layer.Rows() {
			requester <- &model.TileRequest{
				Col:       col,
				Row:       row,
				ZoomLevel: layer.ZoomLevel,
			}
		}
	}

	close(requester)

	wait.Wait()

	return nil
}

func worker(
	ctx context.Context,
	index int,
	wait *sync.WaitGroup,
	requester chan *model.TileRequest,
	loader Loader,
	processor func(workerIndex int, tile model.Tile) error,
	display io.Writer,
) {
	defer func() {
		wait.Done()
	}()

	for req := range requester {
		if req == nil {
			return
		}

		attempt := 0

		err := backoff.Retry(func() error {
			if attempt != 0 {
				_, _ = fmt.Fprintf(
					display,
					"#%d [%d] 🔍%d ↓%d →%d\n",
					index,
					attempt,
					req.ZoomLevel,
					req.Col,
					req.Row,
				)
			}

			img, err := loader.LoadImage(
				ctx,
				*req,
			)

			switch {
			case errors.Is(err, os.ErrNotExist):
				attempt++

				return backoff.Permanent(err)
			case err != nil:
				attempt++

				return err
			}

			if err := processor(index, model.Tile{
				Image:     img,
				ZoomLevel: req.ZoomLevel,
				Col:       req.Col,
				Row:       req.Row,
			}); err != nil {
				return backoff.Permanent(err)
			}

			return nil
		}, backoff.WithMaxRetries(
			backoff.NewConstantBackOff(200*time.Millisecond),
			3,
		))

		switch {
		case errors.Is(err, os.ErrNotExist):
			_, _ = fmt.Fprintf(
				display,
				"#%d NOT FOUND 🔍%d ↓%d →%d => %s\n",
				index,
				req.ZoomLevel,
				req.Col,
				req.Row,
				err,
			)

		case err != nil:
			_, _ = fmt.Fprintf(
				display,
				"#%d ERROR 🔍%d ↓%d →%d => %s\n",
				index,
				req.ZoomLevel,
				req.Col,
				req.Row,
				err,
			)
		}
	}
}
