package main

import (
	"os"

	"github.com/josegonzalez/kubectl-eso/pkg/cmd"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func main() {
	streams := genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}
	rootCmd := cmd.NewRootCmd(streams)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
