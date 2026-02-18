package cmd

import (
	"bytes"
	"strings"
	"testing"
	"time"

	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestPrintExternalSecretTable(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name          string
		items         []esv1beta1.ExternalSecret
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
			wantContains:  []string{"NAME", "STORE", "REFRESH INTERVAL", "READY", "AGE"},
		},
		{
			name:          "empty list no headers",
			items:         nil,
			allNamespaces: false,
			noHeaders:     true,
			wantMissing:   []string{"NAME"},
		},
		{
			name:          "single item",
			allNamespaces: false,
			noHeaders:     false,
			items: []esv1beta1.ExternalSecret{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:              "my-es",
						Namespace:         "default",
						CreationTimestamp: metav1.NewTime(now.Add(-1 * time.Hour)),
					},
					Spec: esv1beta1.ExternalSecretSpec{
						SecretStoreRef: esv1beta1.SecretStoreRef{
							Name: "my-store",
						},
						RefreshInterval: &metav1.Duration{Duration: 1 * time.Hour},
					},
					Status: esv1beta1.ExternalSecretStatus{
						Conditions: []esv1beta1.ExternalSecretStatusCondition{
							{
								Type:   esv1beta1.ExternalSecretReady,
								Status: corev1.ConditionTrue,
							},
						},
					},
				},
			},
			wantContains: []string{"my-es", "my-store", "1h0m0s", "True"},
		},
		{
			name:          "all namespaces shows namespace column",
			allNamespaces: true,
			noHeaders:     false,
			items: []esv1beta1.ExternalSecret{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:              "my-es",
						Namespace:         "production",
						CreationTimestamp: metav1.NewTime(now),
					},
					Spec: esv1beta1.ExternalSecretSpec{
						SecretStoreRef: esv1beta1.SecretStoreRef{Name: "store"},
					},
				},
			},
			wantContains: []string{"NAMESPACE", "production"},
		},
		{
			name:          "not ready condition",
			allNamespaces: false,
			noHeaders:     false,
			items: []esv1beta1.ExternalSecret{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:              "failing-es",
						CreationTimestamp: metav1.NewTime(now),
					},
					Spec: esv1beta1.ExternalSecretSpec{
						SecretStoreRef: esv1beta1.SecretStoreRef{Name: "store"},
					},
					Status: esv1beta1.ExternalSecretStatus{
						Conditions: []esv1beta1.ExternalSecretStatusCondition{
							{
								Type:   esv1beta1.ExternalSecretReady,
								Status: corev1.ConditionFalse,
							},
						},
					},
				},
			},
			wantContains: []string{"failing-es", "False"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := printExternalSecretTable(&buf, tt.items, tt.allNamespaces, tt.noHeaders)
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
