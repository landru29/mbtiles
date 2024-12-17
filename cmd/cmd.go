package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
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
			os.Exit(1)
		}
	case err := <-errc:
		if err != nil {
			os.Exit(1)
		}
	}
}

func initCommands() *cobra.Command {
	var databaseFilename string

	cmdRoot := &cobra.Command{
		Use:   "mbtiles",
		Short: "manage MbTiles from OACI",
	}

	cmdRoot.AddCommand(
		processCommand(&databaseFilename),
		metadataCommand(&databaseFilename),
		tileCommand(&databaseFilename),
	)

	cmdRoot.PersistentFlags().StringVarP(&databaseFilename, "database", "d", "oaci.mbtiles", "database filename")

	return cmdRoot
}
