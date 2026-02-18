package cmd

import (
	"bytes"
	"strings"
	"testing"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func TestGetClientsError(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"annotate", []string{"annotate", "my-secret"}},
		{"sync", []string{"sync", "my-es"}},
		{"get external-secret", []string{"get", "external-secret"}},
		{"get secret", []string{"get", "secret"}},
		{"get secret-store", []string{"get", "secret-store"}},
		{"get cluster-secret-store", []string{"get", "cluster-secret-store"}},
		{"describe external-secret", []string{"describe", "external-secret", "my-es"}},
		{"describe secret", []string{"describe", "secret", "my-secret"}},
		{"describe secret-store", []string{"describe", "secret-store", "my-ss"}},
		{"describe cluster-secret-store", []string{"describe", "cluster-secret-store", "my-css"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := setupFailingClients()
			defer cleanup()

			var buf bytes.Buffer
			streams := genericclioptions.IOStreams{Out: &buf, ErrOut: &buf}
			root := NewRootCmd(streams)
			root.SetArgs(tt.args)

			err := root.Execute()
			if err == nil {
				t.Fatal("expected error from failing clients, got nil")
			}
			if !strings.Contains(err.Error(), "client error") {
				t.Errorf("expected 'client error', got %v", err)
			}
		})
	}
}

func TestGetNamespaceError(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"annotate", []string{"annotate", "my-secret"}},
		{"sync", []string{"sync", "my-es"}},
		{"get external-secret", []string{"get", "external-secret"}},
		{"get secret", []string{"get", "secret"}},
		{"get secret-store", []string{"get", "secret-store"}},
		{"describe external-secret", []string{"describe", "external-secret", "my-es"}},
		{"describe secret", []string{"describe", "secret", "my-secret"}},
		{"describe secret-store", []string{"describe", "secret-store", "my-ss"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := setupFailingNamespace()
			defer cleanup()

			var buf bytes.Buffer
			streams := genericclioptions.IOStreams{Out: &buf, ErrOut: &buf}
			root := NewRootCmd(streams)
			root.SetArgs(tt.args)

			err := root.Execute()
			if err == nil {
				t.Fatal("expected error from failing namespace, got nil")
			}
			if !strings.Contains(err.Error(), "namespace error") {
				t.Errorf("expected 'namespace error', got %v", err)
			}
		})
	}
}
