package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	esv1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1"
	"github.com/josegonzalez/kubectl-eso/pkg/k8s"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	fakecr "sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
)

// errWriter is an io.Writer that always returns an error.
type errWriter struct{}

func (errWriter) Write([]byte) (int, error) {
	return 0, errors.New("write error")
}

func TestRequireNameArg(t *testing.T) {
	validate := requireNameArg("ExternalSecret")

	t.Run("no args returns error", func(t *testing.T) {
		err := validate(nil, []string{})
		if err == nil {
			t.Fatal("expected error for no args")
		}
		if !strings.Contains(err.Error(), "ExternalSecret") {
			t.Errorf("error should contain resource type, got: %v", err)
		}
	})

	t.Run("one arg returns nil", func(t *testing.T) {
		err := validate(nil, []string{"my-secret"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("two args returns error", func(t *testing.T) {
		err := validate(nil, []string{"a", "b"})
		if err == nil {
			t.Fatal("expected error for two args")
		}
		if !strings.Contains(err.Error(), "ExternalSecret") {
			t.Errorf("error should contain resource type, got: %v", err)
		}
	})
}

func TestGetStoreReadyStatus(t *testing.T) {
	tests := []struct {
		name       string
		conditions []esv1.SecretStoreStatusCondition
		want       string
	}{
		{"nil conditions", nil, "Unknown"},
		{"empty conditions", []esv1.SecretStoreStatusCondition{}, "Unknown"},
		{
			"ready true",
			[]esv1.SecretStoreStatusCondition{
				{Type: esv1.SecretStoreReady, Status: corev1.ConditionTrue},
			},
			"True",
		},
		{
			"ready false",
			[]esv1.SecretStoreStatusCondition{
				{Type: esv1.SecretStoreReady, Status: corev1.ConditionFalse},
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
		conditions []esv1.ExternalSecretStatusCondition
		want       string
	}{
		{"nil conditions", nil, "Unknown"},
		{"empty conditions", []esv1.ExternalSecretStatusCondition{}, "Unknown"},
		{
			"ready true",
			[]esv1.ExternalSecretStatusCondition{
				{Type: esv1.ExternalSecretReady, Status: corev1.ConditionTrue},
			},
			"True",
		},
		{
			"ready false",
			[]esv1.ExternalSecretStatusCondition{
				{Type: esv1.ExternalSecretReady, Status: corev1.ConditionFalse},
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
	conditions := []esv1.SecretStoreStatusCondition{
		{
			Type:    esv1.SecretStoreReady,
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
		provider *esv1.SecretStoreProvider
		want     string
	}{
		{"nil provider", nil, "<none>"},
		{"AzureKV", &esv1.SecretStoreProvider{AzureKV: &esv1.AzureKVProvider{}}, "AzureKV"},
		{"GCPSM", &esv1.SecretStoreProvider{GCPSM: &esv1.GCPSMProvider{}}, "GCPSM"},
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

func TestGetESReadyStatusFalse(t *testing.T) {
	// Test the specific code path where Status != "True" returns "False"
	conditions := []esv1.ExternalSecretStatusCondition{
		{Type: esv1.ExternalSecretReady, Status: "Unknown"},
	}
	got := getESReadyStatus(conditions)
	if got != "False" {
		t.Errorf("getESReadyStatus() = %q, want %q", got, "False")
	}
}

func newFakeCRClient(objects ...client.Object) client.Client {
	scheme, _ := k8s.NewScheme()
	return fakecr.NewClientBuilder().WithScheme(scheme).WithObjects(objects...).Build()
}

func newFakeCRClientWithInterceptors(funcs interceptor.Funcs, objects ...client.Object) client.Client {
	scheme, _ := k8s.NewScheme()
	return fakecr.NewClientBuilder().WithScheme(scheme).WithObjects(objects...).WithInterceptorFuncs(funcs).Build()
}

func TestBuildSecretStoreMap(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		crClient := newFakeCRClient()
		m := buildSecretStoreMap(context.Background(), crClient, "default")
		if len(m) != 0 {
			t.Errorf("expected empty map, got %v", m)
		}
	})

	t.Run("explicit target name", func(t *testing.T) {
		es := &esv1.ExternalSecret{
			ObjectMeta: metav1.ObjectMeta{Name: "my-es", Namespace: "default"},
			Spec: esv1.ExternalSecretSpec{
				SecretStoreRef: esv1.SecretStoreRef{Name: "vault", Kind: "SecretStore"},
				Target:         esv1.ExternalSecretTarget{Name: "my-secret"},
			},
		}
		crClient := newFakeCRClient(es)
		m := buildSecretStoreMap(context.Background(), crClient, "default")
		ref, ok := m["default/my-secret"]
		if !ok {
			t.Fatal("expected key default/my-secret in map")
		}
		if ref.name != "vault" || ref.kind != "SecretStore" {
			t.Errorf("got ref %+v, want {name:vault kind:SecretStore}", ref)
		}
	})

	t.Run("default target name from ES name", func(t *testing.T) {
		es := &esv1.ExternalSecret{
			ObjectMeta: metav1.ObjectMeta{Name: "my-es", Namespace: "default"},
			Spec: esv1.ExternalSecretSpec{
				SecretStoreRef: esv1.SecretStoreRef{Name: "vault"},
			},
		}
		crClient := newFakeCRClient(es)
		m := buildSecretStoreMap(context.Background(), crClient, "default")
		ref, ok := m["default/my-es"]
		if !ok {
			t.Fatal("expected key default/my-es in map")
		}
		if ref.name != "vault" {
			t.Errorf("got name %q, want %q", ref.name, "vault")
		}
	})

	t.Run("default kind to SecretStore", func(t *testing.T) {
		es := &esv1.ExternalSecret{
			ObjectMeta: metav1.ObjectMeta{Name: "my-es", Namespace: "default"},
			Spec: esv1.ExternalSecretSpec{
				SecretStoreRef: esv1.SecretStoreRef{Name: "vault"},
				Target:         esv1.ExternalSecretTarget{Name: "my-secret"},
			},
		}
		crClient := newFakeCRClient(es)
		m := buildSecretStoreMap(context.Background(), crClient, "default")
		if m["default/my-secret"].kind != "SecretStore" {
			t.Errorf("expected kind SecretStore, got %q", m["default/my-secret"].kind)
		}
	})

	t.Run("ClusterSecretStore kind", func(t *testing.T) {
		es := &esv1.ExternalSecret{
			ObjectMeta: metav1.ObjectMeta{Name: "my-es", Namespace: "default"},
			Spec: esv1.ExternalSecretSpec{
				SecretStoreRef: esv1.SecretStoreRef{Name: "global-vault", Kind: "ClusterSecretStore"},
				Target:         esv1.ExternalSecretTarget{Name: "my-secret"},
			},
		}
		crClient := newFakeCRClient(es)
		m := buildSecretStoreMap(context.Background(), crClient, "default")
		if m["default/my-secret"].kind != "ClusterSecretStore" {
			t.Errorf("expected kind ClusterSecretStore, got %q", m["default/my-secret"].kind)
		}
	})

	t.Run("first match wins for duplicate targets", func(t *testing.T) {
		es1 := &esv1.ExternalSecret{
			ObjectMeta: metav1.ObjectMeta{Name: "es-a", Namespace: "default"},
			Spec: esv1.ExternalSecretSpec{
				SecretStoreRef: esv1.SecretStoreRef{Name: "store-a"},
				Target:         esv1.ExternalSecretTarget{Name: "my-secret"},
			},
		}
		es2 := &esv1.ExternalSecret{
			ObjectMeta: metav1.ObjectMeta{Name: "es-b", Namespace: "default"},
			Spec: esv1.ExternalSecretSpec{
				SecretStoreRef: esv1.SecretStoreRef{Name: "store-b"},
				Target:         esv1.ExternalSecretTarget{Name: "my-secret"},
			},
		}
		crClient := newFakeCRClient(es1, es2)
		m := buildSecretStoreMap(context.Background(), crClient, "default")
		ref := m["default/my-secret"]
		// first-match-wins: should be store-a since es-a sorts before es-b
		if ref.name != "store-a" {
			t.Errorf("expected first match store-a, got %q", ref.name)
		}
	})

	t.Run("list error returns empty map", func(t *testing.T) {
		crClient := newFakeCRClientWithInterceptors(interceptor.Funcs{
			List: func(_ context.Context, _ client.WithWatch, _ client.ObjectList, _ ...client.ListOption) error {
				return fmt.Errorf("list error")
			},
		})
		m := buildSecretStoreMap(context.Background(), crClient, "default")
		if len(m) != 0 {
			t.Errorf("expected empty map on list error, got %v", m)
		}
	})

	t.Run("empty namespace lists all", func(t *testing.T) {
		es := &esv1.ExternalSecret{
			ObjectMeta: metav1.ObjectMeta{Name: "my-es", Namespace: "production"},
			Spec: esv1.ExternalSecretSpec{
				SecretStoreRef: esv1.SecretStoreRef{Name: "vault"},
				Target:         esv1.ExternalSecretTarget{Name: "my-secret"},
			},
		}
		crClient := newFakeCRClient(es)
		m := buildSecretStoreMap(context.Background(), crClient, "")
		if _, ok := m["production/my-secret"]; !ok {
			t.Error("expected key production/my-secret in map when namespace is empty")
		}
	})
}
