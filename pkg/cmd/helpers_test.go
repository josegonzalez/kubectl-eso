package cmd

import (
	"bytes"
	"strings"
	"testing"
	"text/tabwriter"
	"time"

	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	corev1 "k8s.io/api/core/v1"
)

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
	w := tabwriter.NewWriter(&buf, 0, 4, 2, ' ', 0)
	printConditionsTable(w, conditions)
	w.Flush()

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
