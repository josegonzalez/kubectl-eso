package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/josegonzalez/kubectl-eso/pkg/eso"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

func TestDescribeSecretNoArgs(t *testing.T) {
	var buf bytes.Buffer
	streams := genericclioptions.IOStreams{Out: &buf, ErrOut: &buf}
	configFlags := genericclioptions.NewConfigFlags(true)
	cmd := NewDescribeSecretCmd(streams, configFlags)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for no args, got nil")
	}
	if !strings.Contains(err.Error(), "Secret") {
		t.Errorf("error should contain resource type, got: %v", err)
	}
}

func TestPrintSecretDetail(t *testing.T) {
	tests := []struct {
		name         string
		secretName   string
		namespace    string
		secretType   string
		labels       map[string]string
		data         map[string][]byte
		decode       bool
		store        storeRef
		wantContains []string
		wantMissing  []string
	}{
		{
			name:       "managed secret with decoded values",
			secretName: "my-secret",
			namespace:  "default",
			secretType: "Opaque",
			labels: map[string]string{
				eso.LabelManaged: "true",
			},
			store: storeRef{name: "my-store", kind: "SecretStore"},
			data: map[string][]byte{
				"username": []byte("admin"),
				"password": []byte("secret123"),
			},
			decode: true,
			wantContains: []string{
				"Name:", "my-secret",
				"Namespace:", "default",
				"Type:", "Opaque",
				"ESO-Managed:", "Yes",
				"Store:", "my-store",
				"Store Kind:", "SecretStore",
				"username", "admin",
				"password", "secret123",
			},
			wantMissing: []string{"WARNING"},
		},
		{
			name:       "unmanaged secret shows warning",
			secretName: "plain-secret",
			namespace:  "default",
			secretType: "Opaque",
			labels:     nil,
			data: map[string][]byte{
				"key": []byte("value"),
			},
			decode:       false,
			wantContains: []string{"ESO-Managed:", "No", "WARNING"},
		},
		{
			name:       "encoded values",
			secretName: "encoded",
			namespace:  "default",
			secretType: "Opaque",
			data: map[string][]byte{
				"token": []byte("mysecrettoken"),
			},
			decode:       false,
			wantContains: []string{"bXlzZWNyZXR0b2tlbg=="},
			wantMissing:  []string{"mysecrettoken"},
		},
		{
			name:       "store kind shown for ClusterSecretStore",
			secretName: "css-secret",
			namespace:  "default",
			secretType: "Opaque",
			labels:     map[string]string{eso.LabelManaged: "true"},
			store:      storeRef{name: "global-vault", kind: "ClusterSecretStore"},
			data:       map[string][]byte{"key": []byte("val")},
			decode:     false,
			wantContains: []string{
				"Store:", "global-vault",
				"Store Kind:", "ClusterSecretStore",
			},
		},
		{
			name:        "no store lines when store is empty",
			secretName:  "no-store",
			namespace:   "default",
			secretType:  "Opaque",
			labels:      map[string]string{eso.LabelManaged: "true"},
			data:        map[string][]byte{"key": []byte("val")},
			decode:      false,
			wantMissing: []string{"Store:"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := printSecretDetail(&buf, tt.secretName, tt.namespace, tt.secretType, tt.labels, tt.data, tt.decode, tt.store)
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
