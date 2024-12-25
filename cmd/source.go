package main

import (
	"errors"
	"image/jpeg"
	"os"
	"path/filepath"

	"github.com/landru29/mbtiles/internal/model"
	pkgerrors "github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func sourceCommand() *cobra.Command {
	var (
		outputFilename string
		zoom           int64
		coordinate     model.LatLng
	)

	output := &cobra.Command{
		Use:   "source",
		Short: "download a tile from the source",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if zoom < 0 {
				return errors.New("missing zoom")
			}

			img, err := appli(cmd.Context()).Download(cmd.Context(), coordinate, uint64(zoom))
			if err != nil {
				return err
			}

			file, err := os.Create(filepath.Clean(outputFilename))
			if err != nil {
				return pkgerrors.WithMessage(err, "cannot create file")
			}

			defer func() {
				_ = file.Close()
			}()

			return jpeg.Encode(file, img, &jpeg.Options{Quality: 100})
		},
	}

	output.Flags().StringVarP(&outputFilename, "output", "o", "", "output filename for JPEG output")
	output.Flags().Int64VarP(&zoom, "zoom", "z", -1, "tile zoom")
	output.Flags().VarP(&coordinate, "coordinate", "c", "GNSS coordinate")

	return output
}
