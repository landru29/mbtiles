package main

import (
	"errors"
	"image/png"
	"os"
	"path/filepath"

	"github.com/landru29/mbtiles/internal/model"
	"github.com/spf13/cobra"
)

func tileCommand() *cobra.Command {
	output := &cobra.Command{
		Use:   "tile",
		Short: "manage tiles on MbTiles",
		RunE: func(cmd *cobra.Command, _ []string) error {
			description, err := appli(cmd.Context()).Tiles(cmd.Context())
			if err != nil {
				return err
			}

			cmd.Printf("Tiles count: %d\n", description.Count)

			cmd.Printf("Min zoom: %d\n", description.Zoom[0])
			cmd.Printf("Max zoom: %d\n", description.Zoom[1])

			for idx := description.Zoom[0]; idx <= description.Zoom[1]; idx++ {
				cmd.Printf("\nZoom: %d (%d)\n", idx, description.CountPerZoom[idx])
				cmd.Printf(" - Col bounds: %d - %d\n", description.Col[idx][0], description.Col[idx][1])
				cmd.Printf(" - Row bounds: %d - %d\n", description.Row[idx][0], description.Row[idx][1])
			}

			return nil
		},
	}

	output.AddCommand(
		tileGetCommand(),
		tileRewriteCommand(),
	)

	return output
}

func tileGetCommand() *cobra.Command {
	var (
		index      int64
		col        int64
		row        int64
		zoom       int64
		outputFile string
	)

	output := &cobra.Command{
		Use:   "get",
		Short: "get tile from MbTiles",
		RunE: func(cmd *cobra.Command, _ []string) error {
			app := appli(cmd.Context())

			var tile *model.TileSample

			switch {
			case index > 0:
				var err error

				tile, err = app.TileByIndex(cmd.Context(), uint64(index))
				if err != nil {
					return err
				}

			case col > 0 && row > 0 && zoom > 0:
				var err error

				tile, err = app.TileByCoordinates(cmd.Context(), uint64(zoom), uint64(col), uint64(row))
				if err != nil {
					return err
				}

			default:
				return errors.New("missing index, col or row")
			}

			if tile.Image == nil {
				return errors.New("no tile found")
			}

			cmd.Printf("Row: %d\n", tile.Row)
			cmd.Printf("Col: %d\n", tile.Col)
			cmd.Printf("Zoom: %d\n", tile.ZoomLevel)
			cmd.Printf("Type: %s\n", tile.Type)
			cmd.Printf("Width: %d\n", tile.Image.Bounds().Max.X)
			cmd.Printf("Height: %d\n", tile.Image.Bounds().Max.Y)

			if outputFile != "" {
				file, err := os.Create(filepath.Clean(outputFile))
				if err != nil {
					return err
				}

				defer func() {
					_ = file.Close()
				}()

				return png.Encode(file, tile.Image)
			}

			return nil
		},
	}

	output.Flags().Int64VarP(&index, "index", "i", -1, "tile index in database")
	output.Flags().Int64VarP(&col, "col", "c", -1, "tile column in database")
	output.Flags().Int64VarP(&row, "row", "r", -1, "tile row in database")
	output.Flags().Int64VarP(&zoom, "zoom", "z", -1, "tile zoom in database")
	output.Flags().StringVarP(&outputFile, "output", "o", "", "out filename")

	return output
}

func tileRewriteCommand() *cobra.Command {
	output := &cobra.Command{
		Use:   "rewrite",
		Short: "rewrite tile (PNG) to MbTiles",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return appli(cmd.Context()).TileRewrite(cmd.Context())
		},
	}

	return output
}
