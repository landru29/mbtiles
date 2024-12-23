package main

import (
	"github.com/spf13/cobra"
)

func processCommand() *cobra.Command {
	workerCount := 5

	output := &cobra.Command{
		Use:   "generate",
		Short: "generate MbTiles from OACI",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return appli(cmd.Context()).Generate(
				cmd.Context(),
				options(cmd.Context()),
				workerCount,
			)
		},
	}

	output.Flags().IntVarP(&workerCount, "workers", "w", 5, "number of simultaneous http requests")

	return output
}
