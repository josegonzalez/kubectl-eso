package cmd

import (
	"bytes"
	"strings"
	"testing"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func TestVersionCmd(t *testing.T) {
	var buf bytes.Buffer
	streams := genericclioptions.IOStreams{Out: &buf}
	cmd := NewVersionCmd(streams)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "kubectl-eso version") {
		t.Errorf("expected output to contain 'kubectl-eso version', got %q", output)
	}
}
