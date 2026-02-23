package cmd

import (
	"bytes"
	"strings"
	"testing"

	esv1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1"
	"github.com/josegonzalez/kubectl-eso/pkg/eso"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestRunDescribeSecret(t *testing.T) {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-secret",
			Namespace: "default",
			Labels:    map[string]string{eso.LabelManaged: "true"},
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			"username": []byte("admin"),
		},
	}

	es := &esv1.ExternalSecret{
		ObjectMeta: metav1.ObjectMeta{Name: "my-es", Namespace: "default"},
		Spec: esv1.ExternalSecretSpec{
			SecretStoreRef: esv1.SecretStoreRef{Name: "my-store", Kind: "SecretStore"},
			Target:         esv1.ExternalSecretTarget{Name: "my-secret"},
		},
	}

	tests := []struct {
		name         string
		args         []string
		wantContains []string
	}{
		{
			name:         "table output",
			args:         []string{"secret", "my-secret"},
			wantContains: []string{"my-secret", "Yes", "my-store", "username"},
		},
		{
			name:         "table with decode",
			args:         []string{"secret", "my-secret", "--decode"},
			wantContains: []string{"admin"},
		},
		{
			name:         "json output",
			args:         []string{"secret", "my-secret", "-o", "json"},
			wantContains: []string{"my-secret"},
		},
		{
			name:         "yaml output",
			args:         []string{"secret", "my-secret", "-o", "yaml"},
			wantContains: []string{"my-secret"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := setupFakeClients("default", []client.Object{es}, secret)
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

func TestRunDescribeSecretHelmRejection(t *testing.T) {
	helmSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "helm-secret",
			Namespace: "default",
		},
		Type: corev1.SecretType(eso.HelmSecretType),
	}

	cleanup := setupFakeClients("default", nil, helmSecret)
	defer cleanup()

	var buf bytes.Buffer
	streams := genericclioptions.IOStreams{Out: &buf, ErrOut: &buf}
	configFlags := genericclioptions.NewConfigFlags(true)
	cmd := NewDescribeCmd(streams, configFlags)
	cmd.SetArgs([]string{"secret", "helm-secret"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for helm secret, got nil")
	}
	if !strings.Contains(err.Error(), "Helm-managed") {
		t.Errorf("expected Helm-managed error, got %v", err)
	}
}

func TestRunDescribeSecretNotFound(t *testing.T) {
	cleanup := setupFakeClients("default", nil)
	defer cleanup()

	var buf bytes.Buffer
	streams := genericclioptions.IOStreams{Out: &buf, ErrOut: &buf}
	configFlags := genericclioptions.NewConfigFlags(true)
	cmd := NewDescribeCmd(streams, configFlags)
	cmd.SetArgs([]string{"secret", "nonexistent"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for nonexistent secret, got nil")
	}
}

func TestRunDescribeSecretStoreFromExternalSecret(t *testing.T) {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "eso-secret",
			Namespace: "default",
			Labels:    map[string]string{eso.LabelManaged: "true"},
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{"key": []byte("val")},
	}

	es := &esv1.ExternalSecret{
		ObjectMeta: metav1.ObjectMeta{Name: "my-es", Namespace: "default"},
		Spec: esv1.ExternalSecretSpec{
			SecretStoreRef: esv1.SecretStoreRef{Name: "vault", Kind: "ClusterSecretStore"},
			Target:         esv1.ExternalSecretTarget{Name: "eso-secret"},
		},
	}

	cleanup := setupFakeClients("default", []client.Object{es}, secret)
	defer cleanup()

	var buf bytes.Buffer
	streams := genericclioptions.IOStreams{Out: &buf}
	configFlags := genericclioptions.NewConfigFlags(true)
	cmd := NewDescribeCmd(streams, configFlags)
	cmd.SetArgs([]string{"secret", "eso-secret"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	for _, want := range []string{"Store:", "vault", "Store Kind:", "ClusterSecretStore"} {
		if !strings.Contains(output, want) {
			t.Errorf("output missing %q:\n%s", want, output)
		}
	}
}
