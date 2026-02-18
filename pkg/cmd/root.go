package cmd

import (
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// PluginName returns the command name to display in help text.
// When invoked as a kubectl plugin (argv[0] starts with "kubectl-"),
// it returns "kubectl eso" (with a space) per Krew best practices.
func PluginName(argv0 string) string {
	if strings.HasPrefix(filepath.Base(argv0), "kubectl-") {
		return "kubectl eso"
	}
	return "kubectl-eso"
}

// NewRootCmd creates the root command for kubectl-eso.
func NewRootCmd(streams genericclioptions.IOStreams) *cobra.Command {
	configFlags := genericclioptions.NewConfigFlags(true)

	cmd := &cobra.Command{
		Use:           "kubectl-eso",
		Short:         "kubectl plugin for External Secrets Operator",
		Long:          "A kubectl plugin to manage Kubernetes Secrets in the context of the External Secrets Operator (external-secrets.io).",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	configFlags.AddFlags(cmd.PersistentFlags())

	cmd.AddCommand(NewVersionCmd(streams))
	cmd.AddCommand(NewCompletionCmd(streams))
	cmd.AddCommand(NewGetCmd(streams, configFlags))
	cmd.AddCommand(NewDescribeCmd(streams, configFlags))
	cmd.AddCommand(NewAnnotateCmd(streams, configFlags))
	cmd.AddCommand(NewSyncCmd(streams, configFlags))

	return cmd
}
