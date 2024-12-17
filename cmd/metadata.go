package main

import (
	"mbtiles/internal/database"

	"github.com/spf13/cobra"
)

func metadataCommand(databaseFilename *string) *cobra.Command {
	output := &cobra.Command{
		Use:   "metadata",
		Short: "manage metadata on MbTiles",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := database.New(*databaseFilename)
			if err != nil {
				return err
			}

			defer func() {
				_ = db.Close()
			}()

			metadata, err := db.Metadata()
			if err != nil {
				return err
			}

			for key, value := range metadata {
				cmd.Printf("%s: %s\n", key, value)
			}

			return nil
		},
	}

	output.AddCommand(
		metadataRewriteCommand(databaseFilename),
	)

	return output
}

func metadataRewriteCommand(databaseFilename *string) *cobra.Command {
	output := &cobra.Command{
		Use:   "rewrite",
		Short: "rewrite metadata on MbTiles",
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := database.New(*databaseFilename)
			if err != nil {
				return err
			}

			defer func() {
				_ = db.Close()
			}()

			if err := db.MetadataRewrite(); err != nil {
				return err
			}

			return nil
		},
	}

	return output
}
