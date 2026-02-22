package cmd

import (
	"bytes"
	"strings"
	"testing"
	"time"

	esv1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestRunGetSecretStore(t *testing.T) {
	ss := &esv1.SecretStore{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "my-store",
			Namespace:         "default",
			CreationTimestamp: metav1.NewTime(time.Now()),
		},
		Spec: esv1.SecretStoreSpec{
			Provider: &esv1.SecretStoreProvider{
				AWS: &esv1.AWSProvider{},
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
			args:         []string{"secret-store"},
			wantContains: []string{"my-store", "AWS"},
		},
		{
			name:         "json output",
			args:         []string{"secret-store", "-o", "json"},
			wantContains: []string{"my-store"},
		},
		{
			name:         "yaml output",
			args:         []string{"secret-store", "-o", "yaml"},
			wantContains: []string{"my-store"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := setupFakeClients("default", []client.Object{ss})
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

func TestRunGetSecretStoreAllNamespaces(t *testing.T) {
	ss := &esv1.SecretStore{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "my-store",
			Namespace:         "production",
			CreationTimestamp: metav1.NewTime(time.Now()),
		},
		Spec: esv1.SecretStoreSpec{
			Provider: &esv1.SecretStoreProvider{
				AWS: &esv1.AWSProvider{},
			},
		},
	}

	cleanup := setupFakeClients("default", []client.Object{ss})
	defer cleanup()

	var buf bytes.Buffer
	streams := genericclioptions.IOStreams{Out: &buf}
	configFlags := genericclioptions.NewConfigFlags(true)
	cmd := NewGetCmd(streams, configFlags)
	cmd.SetArgs([]string{"secret-store", "-A"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "NAMESPACE") {
		t.Errorf("expected NAMESPACE header:\n%s", output)
	}
	if !strings.Contains(output, "production") {
		t.Errorf("expected production namespace:\n%s", output)
	}
}

func TestRunGetClusterSecretStore(t *testing.T) {
	css := &esv1.ClusterSecretStore{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "global-store",
			CreationTimestamp: metav1.NewTime(time.Now()),
		},
		Spec: esv1.SecretStoreSpec{
			Provider: &esv1.SecretStoreProvider{
				Vault: &esv1.VaultProvider{},
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
			args:         []string{"cluster-secret-store"},
			wantContains: []string{"global-store", "Vault"},
		},
		{
			name:         "json output",
			args:         []string{"cluster-secret-store", "-o", "json"},
			wantContains: []string{"global-store"},
		},
		{
			name:         "yaml output",
			args:         []string{"cluster-secret-store", "-o", "yaml"},
			wantContains: []string{"global-store"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := setupFakeClients("default", []client.Object{css})
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
