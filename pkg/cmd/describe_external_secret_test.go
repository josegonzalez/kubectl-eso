package cmd

import (
	"bytes"
	"strings"
	"testing"

	esv1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func TestDescribeExternalSecretNoArgs(t *testing.T) {
	var buf bytes.Buffer
	streams := genericclioptions.IOStreams{Out: &buf, ErrOut: &buf}
	configFlags := genericclioptions.NewConfigFlags(true)
	cmd := NewDescribeExternalSecretCmd(streams, configFlags)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for no args, got nil")
	}
	if !strings.Contains(err.Error(), "ExternalSecret") {
		t.Errorf("error should contain resource type, got: %v", err)
	}
}

func TestPrintExternalSecretDetail(t *testing.T) {
	tests := []struct {
		name         string
		es           esv1.ExternalSecret
		wantContains []string
	}{
		{
			name: "basic detail",
			es: esv1.ExternalSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-es",
					Namespace: "default",
				},
				Spec: esv1.ExternalSecretSpec{
					SecretStoreRef: esv1.SecretStoreRef{
						Name: "my-store",
						Kind: "SecretStore",
					},
					RefreshInterval: &metav1.Duration{Duration: 3600000000000},
					Target: esv1.ExternalSecretTarget{
						Name:           "target-secret",
						CreationPolicy: esv1.CreatePolicyMerge,
					},
				},
				Status: esv1.ExternalSecretStatus{
					Conditions: []esv1.ExternalSecretStatusCondition{
						{
							Type:    esv1.ExternalSecretReady,
							Status:  corev1.ConditionTrue,
							Reason:  "SecretSynced",
							Message: "Secret was synced",
						},
					},
					SyncedResourceVersion: "12345",
				},
			},
			wantContains: []string{
				"Name:", "my-es",
				"Namespace:", "default",
				"Store:", "my-store",
				"Store Kind:", "SecretStore",
				"Target Secret:", "target-secret",
				"Refresh Interval:",
				"Creation Policy:", "Merge",
				"Conditions:",
				"Ready", "True", "SecretSynced",
				"Synced Resource Version:", "12345",
			},
		},
		{
			name: "no target name uses es name",
			es: esv1.ExternalSecret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-es",
					Namespace: "default",
				},
				Spec: esv1.ExternalSecretSpec{
					SecretStoreRef: esv1.SecretStoreRef{Name: "store"},
				},
			},
			wantContains: []string{"Target Secret:", "my-es"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := printExternalSecretDetail(&buf, tt.es)
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
