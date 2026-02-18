package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	corev1 "k8s.io/api/core/v1"
)

// errWriter is an io.Writer that always returns an error.
type errWriter struct{}

func (errWriter) Write([]byte) (int, error) {
	return 0, errors.New("write error")
}

func TestGetStoreReadyStatus(t *testing.T) {
	tests := []struct {
		name       string
		conditions []esv1beta1.SecretStoreStatusCondition
		want       string
	}{
		{"nil conditions", nil, "Unknown"},
		{"empty conditions", []esv1beta1.SecretStoreStatusCondition{}, "Unknown"},
		{
			"ready true",
			[]esv1beta1.SecretStoreStatusCondition{
				{Type: esv1beta1.SecretStoreReady, Status: corev1.ConditionTrue},
			},
			"True",
		},
		{
			"ready false",
			[]esv1beta1.SecretStoreStatusCondition{
				{Type: esv1beta1.SecretStoreReady, Status: corev1.ConditionFalse},
			},
			"False",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getStoreReadyStatus(tt.conditions)
			if got != tt.want {
				t.Errorf("getStoreReadyStatus() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGetESReadyStatus(t *testing.T) {
	tests := []struct {
		name       string
		conditions []esv1beta1.ExternalSecretStatusCondition
		want       string
	}{
		{"nil conditions", nil, "Unknown"},
		{"empty conditions", []esv1beta1.ExternalSecretStatusCondition{}, "Unknown"},
		{
			"ready true",
			[]esv1beta1.ExternalSecretStatusCondition{
				{Type: esv1beta1.ExternalSecretReady, Status: corev1.ConditionTrue},
			},
			"True",
		},
		{
			"ready false",
			[]esv1beta1.ExternalSecretStatusCondition{
				{Type: esv1beta1.ExternalSecretReady, Status: corev1.ConditionFalse},
			},
			"False",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getESReadyStatus(tt.conditions)
			if got != tt.want {
				t.Errorf("getESReadyStatus() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPrintConditionsTable(t *testing.T) {
	conditions := []esv1beta1.SecretStoreStatusCondition{
		{
			Type:    esv1beta1.SecretStoreReady,
			Status:  corev1.ConditionTrue,
			Reason:  "Valid",
			Message: "store is valid",
		},
	}

	var buf bytes.Buffer
	tw := newTableWriter(&buf)
	printConditionsTable(tw, conditions)
	if err := tw.flush(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	for _, want := range []string{"Conditions:", "TYPE", "STATUS", "REASON", "MESSAGE", "Ready", "True", "Valid", "store is valid"} {
		if !strings.Contains(output, want) {
			t.Errorf("output missing %q:\n%s", want, output)
		}
	}
}

func TestFormatAge(t *testing.T) {
	tests := []struct {
		name string
		t    time.Time
		want string
	}{
		{"zero time", time.Time{}, "<unknown>"},
		{"seconds ago", time.Now().Add(-30 * time.Second), "30s"},
		{"minutes ago", time.Now().Add(-5 * time.Minute), "5m"},
		{"hours ago", time.Now().Add(-3 * time.Hour), "3h"},
		{"days ago", time.Now().Add(-48 * time.Hour), "2d"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatAge(tt.t)
			if got != tt.want {
				t.Errorf("formatAge() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPrintJSON(t *testing.T) {
	var buf bytes.Buffer
	obj := map[string]string{"key": "value"}
	if err := printJSON(&buf, obj); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, `"key": "value"`) {
		t.Errorf("output missing expected JSON:\n%s", output)
	}
}

func TestPrintYAML(t *testing.T) {
	var buf bytes.Buffer
	obj := map[string]string{"key": "value"}
	if err := printYAML(&buf, obj); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "key: value") {
		t.Errorf("output missing expected YAML:\n%s", output)
	}
}

func TestTableWriterErrorPropagation(t *testing.T) {
	var buf bytes.Buffer
	tw := newTableWriter(&buf)

	// Simulate a prior write error
	tw.err = errors.New("prior error")

	// All subsequent writes should be no-ops due to the error guard
	tw.fprintf("should not write")
	tw.fprintln("should not write")
	tw.fprint("should not write")

	// flush should return the prior error without writing
	if err := tw.flush(); err == nil {
		t.Fatal("expected error from flush")
	}

	// Buffer should be empty since all writes were skipped
	if buf.Len() != 0 {
		t.Errorf("expected empty buffer, got %q", buf.String())
	}
}

func TestTableWriterFlushError(t *testing.T) {
	ew := errWriter{}
	tw := newTableWriter(ew)

	// Write to the tabwriter buffer (succeeds because tabwriter buffers)
	tw.fprintf("hello")

	// flush writes to the underlying errWriter, which fails
	err := tw.flush()
	if err == nil {
		t.Fatal("expected error from flush with errWriter")
	}
}

func TestPrintJSONWriteError(t *testing.T) {
	ew := errWriter{}
	err := printJSON(ew, map[string]string{"key": "value"})
	if err == nil {
		t.Fatal("expected error writing JSON to errWriter")
	}
}

func TestPrintJSONMarshalError(t *testing.T) {
	var buf bytes.Buffer
	err := printJSON(&buf, make(chan int))
	if err == nil {
		t.Fatal("expected marshal error for channel type")
	}
}

func TestPrintYAMLWriteError(t *testing.T) {
	ew := errWriter{}
	err := printYAML(ew, map[string]string{"key": "value"})
	if err == nil {
		t.Fatal("expected error writing YAML to errWriter")
	}
}

func TestPrintYAMLMarshalError(t *testing.T) {
	var buf bytes.Buffer
	// yaml.Marshal uses JSON under the hood, so a channel triggers a marshal error
	err := printYAML(&buf, make(chan int))
	if err == nil {
		t.Fatal("expected marshal error for channel type")
	}
}

func TestGetProviderNameVariants(t *testing.T) {
	tests := []struct {
		name     string
		provider *esv1beta1.SecretStoreProvider
		want     string
	}{
		{"nil provider", nil, "<none>"},
		{"AzureKV", &esv1beta1.SecretStoreProvider{AzureKV: &esv1beta1.AzureKVProvider{}}, "AzureKV"},
		{"GCPSM", &esv1beta1.SecretStoreProvider{GCPSM: &esv1beta1.GCPSMProvider{}}, "GCPSM"},
		{"Kubernetes", &esv1beta1.SecretStoreProvider{Kubernetes: &esv1beta1.KubernetesProvider{}}, "Kubernetes"},
		{"Fake", &esv1beta1.SecretStoreProvider{Fake: &esv1beta1.FakeProvider{}}, "Fake"},
		{"unknown", &esv1beta1.SecretStoreProvider{}, "<unknown>"},
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

func TestGetESReadyStatusFalse(t *testing.T) {
	// Test the specific code path where Status != "True" returns "False"
	conditions := []esv1beta1.ExternalSecretStatusCondition{
		{Type: esv1beta1.ExternalSecretReady, Status: "Unknown"},
	}
	got := getESReadyStatus(conditions)
	if got != "False" {
		t.Errorf("getESReadyStatus() = %q, want %q", got, "False")
	}
}
