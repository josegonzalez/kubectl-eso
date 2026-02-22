package cmd

import (
	"bytes"
	"strings"
	"testing"
	"time"

	esv1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestRunDescribeSecretStore(t *testing.T) {
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
		Status: esv1.SecretStoreStatus{
			Conditions: []esv1.SecretStoreStatusCondition{
				{
					Type:   esv1.SecretStoreReady,
					Status: corev1.ConditionTrue,
					Reason: "Valid",
				},
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
			args:         []string{"secret-store", "my-store"},
			wantContains: []string{"my-store", "AWS", "True"},
		},
		{
			name:         "json output",
			args:         []string{"secret-store", "my-store", "-o", "json"},
			wantContains: []string{"my-store"},
		},
		{
			name:         "yaml output",
			args:         []string{"secret-store", "my-store", "-o", "yaml"},
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

func TestRunDescribeClusterSecretStore(t *testing.T) {
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
		Status: esv1.SecretStoreStatus{
			Conditions: []esv1.SecretStoreStatusCondition{
				{
					Type:   esv1.SecretStoreReady,
					Status: corev1.ConditionTrue,
					Reason: "Valid",
				},
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
			args:         []string{"cluster-secret-store", "global-store"},
			wantContains: []string{"global-store", "Vault", "Cluster"},
		},
		{
			name:         "json output",
			args:         []string{"cluster-secret-store", "global-store", "-o", "json"},
			wantContains: []string{"global-store"},
		},
		{
			name:         "yaml output",
			args:         []string{"cluster-secret-store", "global-store", "-o", "yaml"},
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

func TestRunDescribeSecretStoreNotFound(t *testing.T) {
	cleanup := setupFakeClients("default", nil)
	defer cleanup()

	var buf bytes.Buffer
	streams := genericclioptions.IOStreams{Out: &buf, ErrOut: &buf}
	configFlags := genericclioptions.NewConfigFlags(true)
	cmd := NewDescribeCmd(streams, configFlags)
	cmd.SetArgs([]string{"secret-store", "nonexistent"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for nonexistent SecretStore, got nil")
	}
	if !strings.Contains(err.Error(), "failed to get SecretStore") {
		t.Errorf("expected 'failed to get SecretStore' error, got %v", err)
	}
}

func TestRunDescribeClusterSecretStoreNotFound(t *testing.T) {
	cleanup := setupFakeClients("default", nil)
	defer cleanup()

	var buf bytes.Buffer
	streams := genericclioptions.IOStreams{Out: &buf, ErrOut: &buf}
	configFlags := genericclioptions.NewConfigFlags(true)
	cmd := NewDescribeCmd(streams, configFlags)
	cmd.SetArgs([]string{"cluster-secret-store", "nonexistent"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for nonexistent ClusterSecretStore, got nil")
	}
	if !strings.Contains(err.Error(), "failed to get ClusterSecretStore") {
		t.Errorf("expected 'failed to get ClusterSecretStore' error, got %v", err)
	}
}
