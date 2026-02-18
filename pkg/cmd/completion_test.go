package cmd

import (
	"bytes"
	"testing"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func TestCompletionCmd(t *testing.T) {
	shells := []string{"bash", "zsh", "fish", "powershell"}
	for _, shell := range shells {
		t.Run(shell, func(t *testing.T) {
			var buf bytes.Buffer
			streams := genericclioptions.IOStreams{Out: &buf}
			root := NewRootCmd(streams)
			root.SetArgs([]string{"completion", shell})

			if err := root.Execute(); err != nil {
				t.Fatalf("unexpected error for shell %s: %v", shell, err)
			}

			if buf.Len() == 0 {
				t.Errorf("expected non-empty output for shell %s", shell)
			}
		})
	}
}

func TestCompletionInvalidShell(t *testing.T) {
	var buf bytes.Buffer
	streams := genericclioptions.IOStreams{Out: &buf}
	root := NewRootCmd(streams)
	root.SetArgs([]string{"completion", "invalid"})

	if err := root.Execute(); err == nil {
		t.Error("expected error for invalid shell, got nil")
	}
}
