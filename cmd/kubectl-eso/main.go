package main

import (
	"os"

	"github.com/josegonzalez/kubectl-eso/pkg/cmd"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func main() {
	streams := genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}
	rootCmd := cmd.NewRootCmd(streams)

	execCmd := rootCmd
	if cmd.PluginName(os.Args[0]) == "kubectl eso" {
		kubectl := &cobra.Command{Use: "kubectl", SilenceUsage: true, SilenceErrors: true}
		rootCmd.Use = "eso"
		kubectl.AddCommand(rootCmd)
		kubectl.SetArgs(append([]string{"eso"}, os.Args[1:]...))
		execCmd = kubectl
	}

	if err := execCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
