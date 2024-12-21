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
	appContext   struct{}
	coordContext struct{}
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

func coordinates(ctx context.Context) []model.LatLng {
	coord, found := ctx.Value(coordContext{}).([]model.LatLng)
	if !found {
		return []model.LatLng{{}, {}} // avoid null pointer.
	}

	return coord
}

func initCommands() *cobra.Command {
	databaseFilename := ""
	minCoord := model.LatLng{
		Lat: 41.990226,
		Lng: -5.593299,
	}

	maxCoord := model.LatLng{
		Lat: 51.251834,
		Lng: 8.561345,
	}

	cmdRoot := &cobra.Command{
		Use:   "mbtiles",
		Short: "manage MbTiles from OACI",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			database, err := sqlite.New(cmd.Context(), databaseFilename, minCoord, maxCoord)
			if err != nil {
				return err
			}

			application := app.New(database, cmd.OutOrStdout())

			cmd.SetContext(
				context.WithValue(
					context.WithValue(cmd.Context(), appContext{}, application),
					coordContext{},
					[]model.LatLng{minCoord, maxCoord},
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
	)

	cmdRoot.PersistentFlags().StringVarP(&databaseFilename, "database", "d", "oaci.mbtiles", "database filename")
	cmdRoot.PersistentFlags().VarP(&minCoord, "min", "", "minimum coordinate")
	cmdRoot.PersistentFlags().VarP(&maxCoord, "max", "", "minimum coordinate")

	return cmdRoot
}
