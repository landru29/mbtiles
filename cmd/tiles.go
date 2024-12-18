package main

import (
	"errors"
	"image/png"
	"mbtiles/internal/database"
	"os"

	"github.com/spf13/cobra"
)

func tileCommand(databaseFilename *string) *cobra.Command {
	output := &cobra.Command{
		Use:   "tile",
		Short: "manage tiles on MbTiles",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := database.New(cmd.Context(), *databaseFilename)
			if err != nil {
				return err
			}

			defer func() {
				_ = db.Close()
			}()

			count, err := db.TilesCount(cmd.Context())
			if err != nil {
				return err
			}

			cmd.Printf("Tiles count: %d\n", count)

			allTiles, err := db.AllTiles(cmd.Context())
			if err != nil {
				return err
			}

			if len(allTiles) == 0 {
				return nil
			}

			maxZoom := allTiles[0].ZoomLevel
			minZoom := allTiles[0].ZoomLevel
			minCol := map[int]int{}
			maxCol := map[int]int{}
			minRow := map[int]int{}
			maxRow := map[int]int{}
			tileCount := map[int]int{}

			for _, tile := range allTiles {
				tileCount[tile.ZoomLevel]++

				if maxZoom < tile.ZoomLevel {
					maxZoom = tile.ZoomLevel
				}

				if minZoom > tile.ZoomLevel {
					minZoom = tile.ZoomLevel
				}

				if _, found := minCol[tile.ZoomLevel]; !found {
					minCol[tile.ZoomLevel] = tile.Col
				}

				if _, found := maxCol[tile.ZoomLevel]; !found {
					maxCol[tile.ZoomLevel] = tile.Col
				}

				if _, found := minRow[tile.ZoomLevel]; !found {
					minRow[tile.ZoomLevel] = tile.Row
				}

				if _, found := maxRow[tile.ZoomLevel]; !found {
					maxRow[tile.ZoomLevel] = tile.Row
				}

				if minCol[tile.ZoomLevel] > tile.Col {
					minCol[tile.ZoomLevel] = tile.Col
				}

				if maxCol[tile.ZoomLevel] < tile.Col {
					maxCol[tile.ZoomLevel] = tile.Col
				}

				if minRow[tile.ZoomLevel] > tile.Row {
					minRow[tile.ZoomLevel] = tile.Row
				}

				if maxRow[tile.ZoomLevel] < tile.Row {
					maxRow[tile.ZoomLevel] = tile.Row
				}
			}

			cmd.Printf("Min zoom: %d\n", minZoom)
			cmd.Printf("Max zoom: %d\n", maxZoom)

			for idx := minZoom; idx <= maxZoom; idx++ {
				cmd.Printf("\nZoom: %d (%d)\n", idx, tileCount[idx])
				cmd.Printf(" - Col bounds: %d - %d\n", minCol[idx], maxCol[idx])
				cmd.Printf(" - Row bounds: %d - %d\n", minRow[idx], maxRow[idx])
			}

			return nil
		},
	}

	output.AddCommand(
		tileGetCommand(databaseFilename),
		tileRewriteCommand(databaseFilename),
	)

	return output
}

func tileGetCommand(databaseFilename *string) *cobra.Command {
	var (
		index      int
		col        int
		row        int
		zoom       int
		outputFile string
	)

	output := &cobra.Command{
		Use:   "get",
		Short: "get tile from MbTiles",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := database.New(cmd.Context(), *databaseFilename)
			if err != nil {
				return err
			}

			defer func() {
				_ = db.Close()
			}()

			var tile *database.TileSample

			switch {
			case index > 0:
				tile, err = db.Tile(cmd.Context(), index)
				if err != nil {
					return err
				}

			case col > 0 && row > 0 && zoom > 0:
				tile, err = db.TileByCoordinate(cmd.Context(), zoom, col, row)
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
				file, err := os.Create(outputFile)
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

	output.Flags().IntVarP(&index, "index", "i", -1, "tile index in database")
	output.Flags().IntVarP(&col, "col", "c", -1, "tile column in database")
	output.Flags().IntVarP(&row, "row", "r", -1, "tile row in database")
	output.Flags().IntVarP(&zoom, "zoom", "z", -1, "tile zoom in database")
	output.Flags().StringVarP(&outputFile, "output", "o", "", "out filename")

	return output
}

func tileRewriteCommand(databaseFilename *string) *cobra.Command {
	output := &cobra.Command{
		Use:   "rewrite",
		Short: "rewrite tile (PNG) to MbTiles",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := database.New(cmd.Context(), *databaseFilename)
			if err != nil {
				return err
			}

			defer func() {
				_ = db.Close()
			}()

			allTiles, err := db.AllTiles(cmd.Context())
			if err != nil {
				return err
			}

			cmd.Printf("Rewriting %d tiles\n", len(allTiles))

			return db.TileToPNG(cmd.Context(), cmd.OutOrStdout(), allTiles)
		},
	}

	return output
}
