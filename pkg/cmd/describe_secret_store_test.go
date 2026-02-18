package cmd

import (
	"bytes"
	"strings"
	"testing"

	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	corev1 "k8s.io/api/core/v1"
)

func TestPrintStoreDetail(t *testing.T) {
	tests := []struct {
		name         string
		storeName    string
		namespace    string
		provider     *esv1beta1.SecretStoreProvider
		conditions   []esv1beta1.SecretStoreStatusCondition
		wantContains []string
	}{
		{
			name:      "ready AWS store (namespaced)",
			storeName: "aws-store",
			namespace: "default",
			provider:  &esv1beta1.SecretStoreProvider{AWS: &esv1beta1.AWSProvider{}},
			conditions: []esv1beta1.SecretStoreStatusCondition{
				{
					Type:    esv1beta1.SecretStoreReady,
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
			provider:  &esv1beta1.SecretStoreProvider{Vault: &esv1beta1.VaultProvider{}},
			conditions: []esv1beta1.SecretStoreStatusCondition{
				{
					Type:    esv1beta1.SecretStoreReady,
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
			provider:  &esv1beta1.SecretStoreProvider{GCPSM: &esv1beta1.GCPSMProvider{}},
			conditions: []esv1beta1.SecretStoreStatusCondition{
				{
					Type:   esv1beta1.SecretStoreReady,
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
