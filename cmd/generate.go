package main

import (
	"context"
	"fmt"
	"image"
	"mbtiles/internal/database"
	"mbtiles/internal/tile"
	"mbtiles/internal/tile/oaci"

	"github.com/spf13/cobra"
)

func processCommand(databaseFilename *string) *cobra.Command {
	output := &cobra.Command{
		Use:   "generate",
		Short: "generate MbTiles from OACI",
		RunE: func(cmd *cobra.Command, args []string) error {
			currentBox := tile.New(
				6,  // ZoomLevel
				21, // RowMin
				24, // RowMax
				30, // ColMin
				33, // ColMax
			)

			database, err := database.New(cmd.Context(), *databaseFilename)
			if err != nil {
				return err
			}

			defer func() {
				_ = database.Close()
			}()

			for currentBox.ZoomLevel < 12 {
				if err := currentBox.Loop(context.Background(), oaci.Client{}, func(img image.Image, zoomLevel uint64, col uint64, row uint64) error {
					fmt.Printf("zoom:%d - row: %d - col: %d (%d, %d)\n", zoomLevel, row, col, img.Bounds().Max.X, img.Bounds().Max.Y)

					return database.InsertTile(cmd.Context(), img, zoomLevel, col, row)
				}); err != nil {
					return err
				}

				nextBox, err := currentBox.ToZoom(currentBox.ZoomLevel + 1)
				if err != nil {
					return err
				}

				currentBox = *nextBox
			}

			return nil
		},
	}

	return output
}
