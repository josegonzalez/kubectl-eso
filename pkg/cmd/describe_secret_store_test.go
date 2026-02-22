package cmd

import (
	"bytes"
	"strings"
	"testing"

	esv1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func TestDescribeSecretStoreNoArgs(t *testing.T) {
	var buf bytes.Buffer
	streams := genericclioptions.IOStreams{Out: &buf, ErrOut: &buf}
	configFlags := genericclioptions.NewConfigFlags(true)
	cmd := NewDescribeSecretStoreCmd(streams, configFlags)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for no args, got nil")
	}
	if !strings.Contains(err.Error(), "SecretStore") {
		t.Errorf("error should contain resource type, got: %v", err)
	}
}

func TestDescribeClusterSecretStoreNoArgs(t *testing.T) {
	var buf bytes.Buffer
	streams := genericclioptions.IOStreams{Out: &buf, ErrOut: &buf}
	configFlags := genericclioptions.NewConfigFlags(true)
	cmd := NewDescribeClusterSecretStoreCmd(streams, configFlags)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for no args, got nil")
	}
	if !strings.Contains(err.Error(), "ClusterSecretStore") {
		t.Errorf("error should contain resource type, got: %v", err)
	}
}

func TestPrintStoreDetail(t *testing.T) {
	tests := []struct {
		name         string
		storeName    string
		namespace    string
		provider     *esv1.SecretStoreProvider
		conditions   []esv1.SecretStoreStatusCondition
		wantContains []string
	}{
		{
			name:      "ready AWS store (namespaced)",
			storeName: "aws-store",
			namespace: "default",
			provider:  &esv1.SecretStoreProvider{AWS: &esv1.AWSProvider{}},
			conditions: []esv1.SecretStoreStatusCondition{
				{
					Type:    esv1.SecretStoreReady,
					Status:  corev1.ConditionTrue,
					Reason:  "Valid",
					Message: "store is valid",
				},
			},
			wantContains: []string{
				"Name:", "aws-store",
				"Namespace:", "default",
				"Provider:", "AWS",
				"Ready:", "True",
				"Conditions:",
				"Ready", "True", "Valid", "store is valid",
			},
		},
		{
			name:      "unhealthy Vault store (namespaced)",
			storeName: "vault-store",
			namespace: "production",
			provider:  &esv1.SecretStoreProvider{Vault: &esv1.VaultProvider{}},
			conditions: []esv1.SecretStoreStatusCondition{
				{
					Type:    esv1.SecretStoreReady,
					Status:  corev1.ConditionFalse,
					Reason:  "ConfigError",
					Message: "unable to connect",
				},
			},
			wantContains: []string{
				"Vault",
				"Ready:", "False",
				"ConfigError", "unable to connect",
			},
		},
		{
			name:      "ready cluster store (no namespace)",
			storeName: "global-store",
			namespace: "",
			provider:  &esv1.SecretStoreProvider{GCPSM: &esv1.GCPSMProvider{}},
			conditions: []esv1.SecretStoreStatusCondition{
				{
					Type:   esv1.SecretStoreReady,
					Status: corev1.ConditionTrue,
					Reason: "Valid",
				},
			},
			wantContains: []string{
				"Name:", "global-store",
				"Scope:", "Cluster",
				"Provider:", "GCPSM",
				"Ready:", "True",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := printStoreDetail(&buf, tt.storeName, tt.namespace, tt.provider, tt.conditions)
			if err != nil {
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
