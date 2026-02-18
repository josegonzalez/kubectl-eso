package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/josegonzalez/kubectl-eso/pkg/eso"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func TestRunGetSecret(t *testing.T) {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-secret",
			Namespace: "default",
			Labels:    map[string]string{eso.LabelManaged: "true"},
		},
		Type: corev1.SecretTypeOpaque,
	}

	tests := []struct {
		name         string
		args         []string
		wantContains []string
	}{
		{
			name:         "table output",
			args:         []string{"secret"},
			wantContains: []string{"my-secret", "Yes"},
		},
		{
			name:         "json output",
			args:         []string{"secret", "-o", "json"},
			wantContains: []string{"my-secret"},
		},
		{
			name:         "yaml output",
			args:         []string{"secret", "-o", "yaml"},
			wantContains: []string{"my-secret"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := setupFakeClients("default", nil, secret)
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

func TestRunGetSecretHelmFiltering(t *testing.T) {
	secrets := []corev1.Secret{
		{
			ObjectMeta: metav1.ObjectMeta{Name: "app-secret", Namespace: "default"},
			Type:       corev1.SecretTypeOpaque,
		},
		{
			ObjectMeta: metav1.ObjectMeta{Name: "helm-release", Namespace: "default"},
			Type:       corev1.SecretType(eso.HelmSecretType),
		},
	}

	cleanup := setupFakeClients("default", nil, &secrets[0], &secrets[1])
	defer cleanup()

	var buf bytes.Buffer
	streams := genericclioptions.IOStreams{Out: &buf}
	configFlags := genericclioptions.NewConfigFlags(true)
	cmd := NewGetCmd(streams, configFlags)
	cmd.SetArgs([]string{"secret"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "app-secret") {
		t.Errorf("expected app-secret in output:\n%s", output)
	}
	if strings.Contains(output, "helm-release") {
		t.Errorf("expected helm-release to be filtered out:\n%s", output)
	}
}

func TestRunGetSecretAllNamespaces(t *testing.T) {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-secret",
			Namespace: "production",
		},
		Type: corev1.SecretTypeOpaque,
	}

	cleanup := setupFakeClients("default", nil, secret)
	defer cleanup()

	var buf bytes.Buffer
	streams := genericclioptions.IOStreams{Out: &buf}
	configFlags := genericclioptions.NewConfigFlags(true)
	cmd := NewGetCmd(streams, configFlags)
	cmd.SetArgs([]string{"secret", "-A"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "NAMESPACE") {
		t.Errorf("expected NAMESPACE header:\n%s", output)
	}
}
