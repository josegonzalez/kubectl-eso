package cmd

import (
	"testing"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func TestNewDescribeCmd(t *testing.T) {
	streams := genericclioptions.IOStreams{}
	configFlags := genericclioptions.NewConfigFlags(true)
	cmd := NewDescribeCmd(streams, configFlags)

	expected := []string{"external-secret", "secret", "secret-store", "cluster-secret-store"}
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
