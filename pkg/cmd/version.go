package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// Version information set via ldflags.
var (
	Version = "dev"
	Commit  = "unknown"
	Date    = "unknown"
)

// NewVersionCmd creates the version subcommand.
func NewVersionCmd(streams genericclioptions.IOStreams) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(streams.Out, "kubectl-eso version %s (commit: %s, built: %s)\n", Version, Commit, Date)
			return nil
		},
	}
}
