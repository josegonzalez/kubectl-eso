package cmd

import (
	"bytes"
	"strings"
	"testing"

	esv1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1"
)

func TestPrintStoreTable(t *testing.T) {
	tests := []struct {
		name          string
		items         []storeItem
		allNamespaces bool
		noHeaders     bool
		wantContains  []string
		wantMissing   []string
	}{
		{
			name:          "empty list with headers",
			items:         nil,
			allNamespaces: false,
			noHeaders:     false,
			wantContains:  []string{"NAME", "READY", "PROVIDER", "AGE"},
		},
		{
			name:          "single store ready",
			allNamespaces: false,
			noHeaders:     false,
			items: []storeItem{
				{Name: "my-store", Namespace: "default", Ready: "True", Provider: "AWS", Age: "1h"},
			},
			wantContains: []string{"my-store", "True", "AWS"},
		},
		{
			name:          "all namespaces",
			allNamespaces: true,
			noHeaders:     false,
			items: []storeItem{
				{Name: "store", Namespace: "staging", Ready: "Unknown", Provider: "Vault", Age: "0s"},
			},
			wantContains: []string{"NAMESPACE", "staging", "Vault"},
		},
		{
			name:      "cluster store (no namespace column)",
			noHeaders: false,
			items: []storeItem{
				{Name: "global-store", Ready: "True", Provider: "AWS", Age: "1d"},
			},
			wantContains: []string{"global-store", "True", "AWS"},
		},
		{
			name:      "no headers",
			noHeaders: true,
			items: []storeItem{
				{Name: "store", Ready: "Unknown", Provider: "Vault", Age: "0s"},
			},
			wantContains: []string{"store", "Vault"},
			wantMissing:  []string{"NAME", "READY"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := printStoreTable(&buf, tt.items, tt.allNamespaces, tt.noHeaders)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			output := buf.String()
			for _, want := range tt.wantContains {
				if !strings.Contains(output, want) {
					t.Errorf("output missing %q:\n%s", want, output)
				}
			}
			for _, missing := range tt.wantMissing {
				if strings.Contains(output, missing) {
					t.Errorf("output should not contain %q:\n%s", missing, output)
				}
			}
		})
	}
}

func TestGetProviderName(t *testing.T) {
	tests := []struct {
		name     string
		provider *esv1.SecretStoreProvider
		want     string
	}{
		{"nil provider", nil, "<none>"},
		{"AWS", &esv1.SecretStoreProvider{AWS: &esv1.AWSProvider{}}, "AWS"},
		{"AzureKV", &esv1.SecretStoreProvider{AzureKV: &esv1.AzureKVProvider{}}, "AzureKV"},
		{"GCPSM", &esv1.SecretStoreProvider{GCPSM: &esv1.GCPSMProvider{}}, "GCPSM"},
		{"Vault", &esv1.SecretStoreProvider{Vault: &esv1.VaultProvider{}}, "Vault"},
		{"Kubernetes", &esv1.SecretStoreProvider{Kubernetes: &esv1.KubernetesProvider{}}, "Kubernetes"},
		{"Fake", &esv1.SecretStoreProvider{Fake: &esv1.FakeProvider{}}, "Fake"},
		{"unknown", &esv1.SecretStoreProvider{}, "<unknown>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getProviderName(tt.provider)
			if got != tt.want {
				t.Errorf("getProviderName() = %q, want %q", got, tt.want)
			}
		})
	}
}
