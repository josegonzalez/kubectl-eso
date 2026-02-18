package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// NewGetCmd creates the get parent command.
func NewGetCmd(streams genericclioptions.IOStreams, configFlags *genericclioptions.ConfigFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Display one or more resources",
	}

	cmd.PersistentFlags().StringP("output", "o", "table", "Output format: table, json, yaml")
	cmd.PersistentFlags().Bool("no-headers", false, "Omit table header row")

	cmd.AddCommand(NewGetExternalSecretCmd(streams, configFlags))
	cmd.AddCommand(NewGetSecretCmd(streams, configFlags))
	cmd.AddCommand(NewGetSecretStoreCmd(streams, configFlags))
	cmd.AddCommand(NewGetClusterSecretStoreCmd(streams, configFlags))

	return cmd
}
