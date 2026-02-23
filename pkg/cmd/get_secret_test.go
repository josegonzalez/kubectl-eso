package cmd

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/josegonzalez/kubectl-eso/pkg/eso"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestPrintSecretTable(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name          string
		secrets       []corev1.Secret
		allNamespaces bool
		noHeaders     bool
		storeMap      map[string]storeRef
		wantContains  []string
		wantMissing   []string
	}{
		{
			name:          "empty list with headers",
			secrets:       nil,
			allNamespaces: false,
			noHeaders:     false,
			wantContains:  []string{"NAME", "TYPE", "ESO-MANAGED", "STORE", "AGE"},
		},
		{
			name:          "no headers",
			secrets:       nil,
			allNamespaces: false,
			noHeaders:     true,
			wantMissing:   []string{"NAME"},
		},
		{
			name:          "eso managed secret",
			allNamespaces: false,
			noHeaders:     false,
			secrets: []corev1.Secret{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:              "my-secret",
						Namespace:         "default",
						CreationTimestamp: metav1.NewTime(now.Add(-2 * time.Hour)),
						Labels: map[string]string{
							eso.LabelManaged: "true",
						},
					},
					Type: corev1.SecretTypeOpaque,
				},
			},
			storeMap:     map[string]storeRef{"default/my-secret": {name: "my-store", kind: "SecretStore"}},
			wantContains: []string{"my-secret", "Opaque", "Yes", "my-store"},
		},
		{
			name:          "unmanaged secret",
			allNamespaces: false,
			noHeaders:     false,
			secrets: []corev1.Secret{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:              "plain-secret",
						Namespace:         "default",
						CreationTimestamp: metav1.NewTime(now),
					},
					Type: corev1.SecretTypeOpaque,
				},
			},
			wantContains: []string{"plain-secret", "No"},
		},
		{
			name:          "all namespaces",
			allNamespaces: true,
			noHeaders:     false,
			secrets: []corev1.Secret{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:              "my-secret",
						Namespace:         "production",
						CreationTimestamp: metav1.NewTime(now),
					},
					Type: corev1.SecretTypeOpaque,
				},
			},
			wantContains: []string{"NAMESPACE", "production"},
		},
		{
			name:          "store from map",
			allNamespaces: false,
			noHeaders:     true,
			secrets: []corev1.Secret{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:              "eso-secret",
						Namespace:         "default",
						CreationTimestamp: metav1.NewTime(now),
						Labels:            map[string]string{eso.LabelManaged: "true"},
					},
					Type: corev1.SecretTypeOpaque,
				},
			},
			storeMap:     map[string]storeRef{"default/eso-secret": {name: "vault", kind: "SecretStore"}},
			wantContains: []string{"eso-secret", "Yes", "vault"},
		},
		{
			name:          "no match in map",
			allNamespaces: false,
			noHeaders:     true,
			secrets: []corev1.Secret{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:              "other-secret",
						Namespace:         "default",
						CreationTimestamp: metav1.NewTime(now),
						Labels:            map[string]string{eso.LabelManaged: "true"},
					},
					Type: corev1.SecretTypeOpaque,
				},
			},
			storeMap:    map[string]storeRef{"default/different": {name: "vault", kind: "SecretStore"}},
			wantMissing: []string{"vault"},
		},
		{
			name:          "non-managed secret skips map lookup",
			allNamespaces: false,
			noHeaders:     true,
			secrets: []corev1.Secret{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:              "plain-secret",
						Namespace:         "default",
						CreationTimestamp: metav1.NewTime(now),
					},
					Type: corev1.SecretTypeOpaque,
				},
			},
			storeMap:    map[string]storeRef{"default/plain-secret": {name: "vault", kind: "SecretStore"}},
			wantMissing: []string{"vault"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := printSecretTable(&buf, tt.secrets, tt.allNamespaces, tt.noHeaders, tt.storeMap)
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
