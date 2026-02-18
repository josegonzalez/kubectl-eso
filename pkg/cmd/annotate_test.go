package cmd

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/josegonzalez/kubectl-eso/pkg/eso"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestAnnotateSecret(t *testing.T) {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-secret",
			Namespace: "default",
		},
		Type: corev1.SecretTypeOpaque,
		Data: map[string][]byte{
			"key": []byte("value"),
		},
	}

	clientset := fake.NewSimpleClientset(secret)

	// Simulate the annotation logic
	s, err := clientset.CoreV1().Secrets("default").Get(context.TODO(), "my-secret", metav1.GetOptions{})
	if err != nil {
		t.Fatalf("failed to get secret: %v", err)
	}

	if s.Labels == nil {
		s.Labels = make(map[string]string)
	}
	s.Labels[eso.LabelManaged] = "true"

	if s.Annotations == nil {
		s.Annotations = make(map[string]string)
	}
	s.Annotations[eso.AnnotationImported] = "true"
	s.Annotations[eso.AnnotationStore] = "my-store"
	s.Annotations[eso.AnnotationStoreKind] = "SecretStore"

	_, err = clientset.CoreV1().Secrets("default").Update(context.TODO(), s, metav1.UpdateOptions{})
	if err != nil {
		t.Fatalf("failed to update secret: %v", err)
	}

	// Verify annotations
	updated, err := clientset.CoreV1().Secrets("default").Get(context.TODO(), "my-secret", metav1.GetOptions{})
	if err != nil {
		t.Fatalf("failed to get updated secret: %v", err)
	}

	if updated.Labels[eso.LabelManaged] != "true" {
		t.Errorf("expected label %s=true, got %s", eso.LabelManaged, updated.Labels[eso.LabelManaged])
	}
	if updated.Annotations[eso.AnnotationImported] != "true" {
		t.Errorf("expected annotation %s=true", eso.AnnotationImported)
	}
	if updated.Annotations[eso.AnnotationStore] != "my-store" {
		t.Errorf("expected store annotation my-store, got %s", updated.Annotations[eso.AnnotationStore])
	}
	if updated.Annotations[eso.AnnotationStoreKind] != "SecretStore" {
		t.Errorf("expected store-kind annotation SecretStore, got %s", updated.Annotations[eso.AnnotationStoreKind])
	}
}

func TestAnnotateDryRun(t *testing.T) {
	// Test that dry-run output contains expected annotations
	var buf bytes.Buffer
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-secret",
			Namespace: "default",
			Labels: map[string]string{
				eso.LabelManaged: "true",
			},
			Annotations: map[string]string{
				eso.AnnotationImported:  "true",
				eso.AnnotationStore:     "my-store",
				eso.AnnotationStoreKind: "SecretStore",
			},
		},
		Type: corev1.SecretTypeOpaque,
	}

	err := printYAML(&buf, secret)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	for _, want := range []string{
		eso.LabelManaged,
		eso.AnnotationImported,
		eso.AnnotationStore,
		eso.AnnotationStoreKind,
	} {
		if !strings.Contains(output, want) {
			t.Errorf("dry-run output missing %q:\n%s", want, output)
		}
	}
}
