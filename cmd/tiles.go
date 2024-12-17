package main

import (
	"errors"
	"mbtiles/internal/database"

	"github.com/spf13/cobra"
)

func tileCommand(databaseFilename *string) *cobra.Command {
	output := &cobra.Command{
		Use:   "tile",
		Short: "manage tiles on MbTiles",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := database.New(*databaseFilename)
			if err != nil {
				return err
			}

			defer func() {
				_ = db.Close()
			}()

			count, err := db.TilesCount()
			if err != nil {
				return err
			}

			cmd.Printf("Tiles count: %d\n", count)

			return nil
		},
	}

	output.AddCommand(
		tileGetCommand(databaseFilename),
	)

	return output
}

func tileGetCommand(databaseFilename *string) *cobra.Command {
	var (
		index int
		col   int
		row   int
	)

	output := &cobra.Command{
		Use:   "get",
		Short: "get tile from MbTiles",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := database.New(*databaseFilename)
			if err != nil {
				return err
			}

			defer func() {
				_ = db.Close()
			}()

			var tile *database.TileSample

			switch {
			case index > 0:
				tile, err = db.Tile(index)
				if err != nil {
					return err
				}

			case col > 0 && row > 0:
				tile, err = db.TileByCoordinate(col, row)
				if err != nil {
					return err
				}

			default:
				return errors.New("missing index, col or row")
			}

			cmd.Printf("Row: %d\n", tile.Row)
			cmd.Printf("Col: %d\n", tile.Col)
			cmd.Printf("Zoom: %d\n", tile.ZoomLevel)
			cmd.Printf("Type: %s\n", tile.Type)
			cmd.Printf("Width: %d\n", tile.Image.Bounds().Max.X)
			cmd.Printf("Height: %d\n", tile.Image.Bounds().Max.Y)

			return nil
		},
	}

	output.Flags().IntVarP(&index, "index", "i", -1, "tile index in database")
	output.Flags().IntVarP(&col, "col", "c", -1, "tile column in database")
	output.Flags().IntVarP(&row, "row", "r", -1, "tile row in database")

	output.AddCommand(
		metadataRewriteCommand(databaseFilename),
	)

	return output
}
