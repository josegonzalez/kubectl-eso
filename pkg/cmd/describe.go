package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// NewDescribeCmd creates the describe parent command.
func NewDescribeCmd(streams genericclioptions.IOStreams, configFlags *genericclioptions.ConfigFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Show details of a resource",
	}

	cmd.PersistentFlags().StringP("output", "o", "table", "Output format: table, json, yaml")

	cmd.AddCommand(NewDescribeExternalSecretCmd(streams, configFlags))
	cmd.AddCommand(NewDescribeSecretCmd(streams, configFlags))
	cmd.AddCommand(NewDescribeSecretStoreCmd(streams, configFlags))
	cmd.AddCommand(NewDescribeClusterSecretStoreCmd(streams, configFlags))

	return cmd
}
