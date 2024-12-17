package database

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"io"
)

type TileSample struct {
	ZoomLevel int
	Row       int
	Col       int
	Image     image.Image
	Type      string
}

func (c Connection) TilesCount(ctx context.Context) (int, error) {
	rows, err := c.db.QueryContext(ctx, "SELECT count(*) FROM tiles")
	if err != nil {
		return 0, err
	}

	if err := rows.Err(); err != nil {
		return 0, err
	}

	defer func() {
		_ = rows.Close()
	}()

	var count int

	if rows.Next() {
		if err := rows.Scan(
			&count,
		); err != nil {
			return 0, err
		}

	}

	return count, nil
}

func (c Connection) Tile(ctx context.Context, index int) (*TileSample, error) {
	output := TileSample{}

	rows, err := c.db.QueryContext(ctx, fmt.Sprintf("SELECT zoom_level, tile_column, tile_row, tile_data FROM tiles LIMIT 1 OFFSET %d", index))
	if err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	if rows.Next() {
		var data []byte

		if err := rows.Scan(
			&output.ZoomLevel,
			&output.Col,
			&output.Row,
			&data,
		); err != nil {
			return nil, err
		}

		out, format, err := image.Decode(bytes.NewBuffer(data))
		if err != nil {
			return nil, err
		}

		output.Image = out
		output.Type = format
	}

	return &output, err
}

func (c Connection) TileByCoordinate(ctx context.Context, col int, row int) (*TileSample, error) {
	output := TileSample{}

	rows, err := c.db.QueryContext(ctx, "SELECT zoom_level, tile_column, tile_row, tile_data FROM tiles WHERE tile_column=? AND tile_row=?", col, row)
	if err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	if rows.Next() {
		var data []byte

		if err := rows.Scan(
			&output.ZoomLevel,
			&output.Col,
			&output.Row,
			&data,
		); err != nil {
			return nil, err
		}

		out, format, err := image.Decode(bytes.NewBuffer(data))
		if err != nil {
			return nil, err
		}

		output.Image = out
		output.Type = format
	}

	return &output, err
}

func (c Connection) AllTiles(ctx context.Context) ([]TileSample, error) {
	output := []TileSample{}

	rows, err := c.db.QueryContext(ctx, "SELECT zoom_level, tile_column, tile_row, tile_data FROM tiles")
	if err != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var (
			data    []byte
			element TileSample
		)

		if err := rows.Scan(
			&element.ZoomLevel,
			&element.Col,
			&element.Row,
			&data,
		); err != nil {
			return nil, err
		}

		out, format, err := image.Decode(bytes.NewBuffer(data))
		if err != nil {
			return nil, err
		}

		element.Image = out
		element.Type = format

		output = append(output, element)
	}

	return output, err
}

func (c Connection) TileToPNG(ctx context.Context, display io.Writer, tiles []TileSample) error {
	rows, err := c.db.Query("SELECT zoom_level, tile_column, tile_row, tile_data FROM tiles")
	if err != nil {
		return err
	}

	if err := rows.Err(); err != nil {
		return err
	}

	defer func() {
		_ = rows.Close()
	}()

	for idx, tile := range tiles {
		var buffer bytes.Buffer

		png.Encode(&buffer, tile.Image)

		fmt.Fprintf(display, "#%d Zoom: %d - Row: %d - Col: %d\n", idx, tile.ZoomLevel, tile.Row, tile.Col)

		if _, err := c.db.ExecContext(ctx, "UPDATE tiles SET tile_data=? WHERE zoom_level = ? AND tile_column = ? AND tile_row = ?", buffer.Bytes(), tile.ZoomLevel, tile.Col, tile.Row); err != nil {
			return err
		}
	}

	return nil
}
