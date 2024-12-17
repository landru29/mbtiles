package database

import (
	"bytes"
	"fmt"
	"image"
)

type TileSample struct {
	ZoomLevel int
	Row       int
	Col       int
	Image     image.Image
	Type      string
}

func (c Connection) TilesCount() (int, error) {
	rows, err := c.db.Query("SELECT count(*) FROM tiles")
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

func (c Connection) Tile(index int) (*TileSample, error) {
	output := TileSample{}

	rows, err := c.db.Query(fmt.Sprintf("SELECT zoom_level, tile_column, tile_row, tile_data FROM tiles LIMIT 1 OFFSET %d", index))
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

func (c Connection) TileByCoordinate(col int, row int) (*TileSample, error) {
	output := TileSample{}

	rows, err := c.db.Query("SELECT zoom_level, tile_column, tile_row, tile_data FROM tiles WHERE tile_column=? AND tile_row=?", col, row)
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
