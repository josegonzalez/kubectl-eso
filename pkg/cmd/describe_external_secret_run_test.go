package cmd

import (
	"bytes"
	"strings"
	"testing"

	esv1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestRunDescribeExternalSecret(t *testing.T) {
	es := &esv1.ExternalSecret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-es",
			Namespace: "default",
		},
		Spec: esv1.ExternalSecretSpec{
			SecretStoreRef: esv1.SecretStoreRef{
				Name: "my-store",
				Kind: "SecretStore",
			},
			Target: esv1.ExternalSecretTarget{
				Name: "target-secret",
			},
		},
	}

	tests := []struct {
		name         string
		args         []string
		wantContains []string
	}{
		{
			name:         "table output",
			args:         []string{"external-secret", "my-es"},
			wantContains: []string{"my-es", "my-store", "target-secret"},
		},
		{
			name:         "json output",
			args:         []string{"external-secret", "my-es", "-o", "json"},
			wantContains: []string{"my-es", "my-store"},
		},
		{
			name:         "yaml output",
			args:         []string{"external-secret", "my-es", "-o", "yaml"},
			wantContains: []string{"my-es", "my-store"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := setupFakeClients("default", []client.Object{es})
			defer cleanup()

			var buf bytes.Buffer
			streams := genericclioptions.IOStreams{Out: &buf}
			configFlags := genericclioptions.NewConfigFlags(true)
			cmd := NewDescribeCmd(streams, configFlags)
			cmd.SetArgs(tt.args)

			if err := cmd.Execute(); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			output := buf.String()
			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("output missing %q:\n%s", want, output)
				}
			}
		})
	}
}

func TestRunDescribeExternalSecretNotFound(t *testing.T) {
	cleanup := setupFakeClients("default", nil)
	defer cleanup()

	var buf bytes.Buffer
	streams := genericclioptions.IOStreams{Out: &buf, ErrOut: &buf}
	configFlags := genericclioptions.NewConfigFlags(true)
	cmd := NewDescribeCmd(streams, configFlags)
	cmd.SetArgs([]string{"external-secret", "nonexistent"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for nonexistent ExternalSecret, got nil")
	}
}
