package cmd

import (
	"bytes"
	"strings"
	"testing"
	"time"

	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestRunGetExternalSecret(t *testing.T) {
	es := &esv1beta1.ExternalSecret{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "my-es",
			Namespace:         "default",
			CreationTimestamp: metav1.NewTime(time.Now().Add(-1 * time.Hour)),
		},
		Spec: esv1beta1.ExternalSecretSpec{
			SecretStoreRef:  esv1beta1.SecretStoreRef{Name: "my-store"},
			RefreshInterval: &metav1.Duration{Duration: 1 * time.Hour},
		},
	}

	tests := []struct {
		name         string
		args         []string
		wantContains []string
	}{
		{
			name:         "table output",
			args:         []string{"external-secret"},
			wantContains: []string{"my-es", "my-store"},
		},
		{
			name:         "json output",
			args:         []string{"external-secret", "-o", "json"},
			wantContains: []string{"my-es", "my-store"},
		},
		{
			name:         "yaml output",
			args:         []string{"external-secret", "-o", "yaml"},
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
			cmd := NewGetCmd(streams, configFlags)
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

func TestRunGetExternalSecretAllNamespaces(t *testing.T) {
	es := &esv1beta1.ExternalSecret{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "my-es",
			Namespace:         "production",
			CreationTimestamp: metav1.NewTime(time.Now()),
		},
		Spec: esv1beta1.ExternalSecretSpec{
			SecretStoreRef: esv1beta1.SecretStoreRef{Name: "store"},
		},
	}

	cleanup := setupFakeClients("default", []client.Object{es})
	defer cleanup()

	var buf bytes.Buffer
	streams := genericclioptions.IOStreams{Out: &buf}
	configFlags := genericclioptions.NewConfigFlags(true)
	cmd := NewGetCmd(streams, configFlags)
	cmd.SetArgs([]string{"external-secret", "-A"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "NAMESPACE") {
		t.Errorf("expected NAMESPACE header in output:\n%s", output)
	}
	if !strings.Contains(output, "production") {
		t.Errorf("expected production namespace in output:\n%s", output)
	}
}
