// Package main is the main command line.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/landru29/mbtiles/internal/app"
	"github.com/landru29/mbtiles/internal/database/sqlite"
	"github.com/landru29/mbtiles/internal/model"
	"github.com/spf13/cobra"
)

type (
	appContext    struct{}
	optionContext struct{}
)

func main() {
	cmdRoot := initCommands()

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()

	errc := make(chan error)
	go func() {
		errc <- cmdRoot.ExecuteContext(ctx)
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigChan:
		cancel()

		if err := <-errc; err != nil {
			os.Exit(1) //nolint: gocritic
		}
	case err := <-errc:
		if err != nil {
			os.Exit(1)
		}
	}
}

func appli(ctx context.Context) *app.Application {
	application, found := ctx.Value(appContext{}).(*app.Application)
	if !found {
		return &app.Application{} // avoid null pointer.
	}

	return application
}

func options(ctx context.Context) model.Option {
	zoom, found := ctx.Value(optionContext{}).(model.Option)
	if !found {
		return model.Option{} // avoid null pointer.
	}

	return zoom
}

func initCommands() *cobra.Command {
	databaseFilename := ""
	maxZoom := uint64(10)
	minZoom := uint64(4)
	minCoord := model.LatLng{
		Lat: 41.990226,
		Lng: -5.593299,
	}
	format := "png"

	maxCoord := model.LatLng{
		Lat: 51.251834,
		Lng: 8.561345,
	}

	cmdRoot := &cobra.Command{
		Use:   "mbtiles",
		Short: "manage MbTiles from OACI",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			globalOptions := model.Option{
				CoordinateMin: minCoord,
				CoordinateMax: maxCoord,
				ZoomMin:       minZoom,
				ZoomMax:       maxZoom,
				Format:        format,
			}

			database, err := sqlite.New(cmd.Context(), databaseFilename, globalOptions)
			if err != nil {
				return err
			}

			application := app.New(database, cmd.OutOrStdout())

			cmd.SetContext(
				context.WithValue(
					context.WithValue(cmd.Context(), appContext{}, application),
					optionContext{},
					globalOptions,
				),
			)

			return nil
		},
		PersistentPostRunE: func(cmd *cobra.Command, _ []string) error {
			return appli(cmd.Context()).Close()
		},
	}

	cmdRoot.AddCommand(
		processCommand(),
		metadataCommand(),
		tileCommand(),
		sourceCommand(),
	)

	cmdRoot.PersistentFlags().StringVarP(&databaseFilename, "database", "d", "oaci.mbtiles", "database filename")
	cmdRoot.PersistentFlags().VarP(&minCoord, "min", "", "minimum coordinate")
	cmdRoot.PersistentFlags().VarP(&maxCoord, "max", "", "minimum coordinate")
	cmdRoot.PersistentFlags().Uint64VarP(&maxZoom, "max-zoom", "", 10, "max zoom")
	cmdRoot.PersistentFlags().Uint64VarP(&minZoom, "min-zoom", "", 4, "min zoom")
	cmdRoot.PersistentFlags().StringVarP(&format, "format", "f", "png", "tile format")

	return cmdRoot
}
