package cmd

import (
	"testing"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func TestNewRootCmd(t *testing.T) {
	streams := genericclioptions.IOStreams{}
	cmd := NewRootCmd(streams)

	expected := []string{"get", "describe", "annotate", "sync", "version", "completion"}
	subs := cmd.Commands()
	names := make(map[string]bool, len(subs))
	for _, c := range subs {
		names[c.Name()] = true
	}

	for _, name := range expected {
		if !names[name] {
			t.Errorf("expected subcommand %q not found", name)
		}
	}
}
