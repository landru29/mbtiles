package main

import (
	"github.com/spf13/cobra"
)

func metadataCommand() *cobra.Command {
	output := &cobra.Command{
		Use:   "metadata",
		Short: "manage metadata on MbTiles",
		RunE: func(cmd *cobra.Command, _ []string) error {
			metadata, err := appli(cmd.Context()).Metadata(cmd.Context())
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
		metadataRewriteCommand(),
	)

	return output
}

func metadataRewriteCommand() *cobra.Command {
	output := &cobra.Command{
		Use:   "rewrite",
		Short: "rewrite metadata on MbTiles",
		RunE: func(cmd *cobra.Command, _ []string) error {
			coord := coordinates(cmd.Context())

			if err := appli(cmd.Context()).MetadataRewrite(cmd.Context(), coord[0], coord[1]); err != nil {
				return err
			}

			return nil
		},
	}

	return output
}
