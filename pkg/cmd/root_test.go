package cmd

import (
	"testing"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func TestPluginName(t *testing.T) {
	tests := []struct {
		name   string
		argv0  string
		expect string
	}{
		{"bare plugin name", "kubectl-eso", "kubectl eso"},
		{"full path", "/usr/local/bin/kubectl-eso", "kubectl eso"},
		{"non-plugin binary", "eso", "kubectl-eso"},
		{"non-plugin path", "/usr/local/bin/eso", "kubectl-eso"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PluginName(tt.argv0)
			if got != tt.expect {
				t.Errorf("PluginName(%q) = %q, want %q", tt.argv0, got, tt.expect)
			}
		})
	}
}

func TestNewRootCmd(t *testing.T) {
	streams := genericclioptions.IOStreams{}
	cmd := NewRootCmd(streams)

	expected := []string{"get", "describe", "sync", "version", "completion"}
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
