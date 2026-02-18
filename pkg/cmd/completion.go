package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// NewCompletionCmd creates the completion subcommand.
func NewCompletionCmd(streams genericclioptions.IOStreams) *cobra.Command {
	return &cobra.Command{
		Use:       "completion [bash|zsh|fish|powershell]",
		Short:     "Generate shell completion script",
		Long:      "Generate shell completion script for bash, zsh, fish, or powershell.",
		Args:      cobra.ExactArgs(1),
		ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return cmd.Root().GenBashCompletionV2(streams.Out, true)
			case "zsh":
				return cmd.Root().GenZshCompletion(streams.Out)
			case "fish":
				return cmd.Root().GenFishCompletion(streams.Out, true)
			case "powershell":
				return cmd.Root().GenPowerShellCompletionWithDesc(streams.Out)
			default:
				return fmt.Errorf("unsupported shell: %s", args[0])
			}
		},
	}
}
