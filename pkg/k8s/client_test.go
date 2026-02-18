package k8s

import (
	"testing"

	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	corev1 "k8s.io/api/core/v1"
)

func TestNewScheme(t *testing.T) {
	scheme, err := NewScheme()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify core types are registered
	if !scheme.Recognizes(corev1.SchemeGroupVersion.WithKind("Secret")) {
		t.Error("scheme does not recognize Secret")
	}

	// Verify ESO types are registered
	if !scheme.Recognizes(esv1beta1.SchemeGroupVersion.WithKind("ExternalSecret")) {
		t.Error("scheme does not recognize ExternalSecret")
	}
	if !scheme.Recognizes(esv1beta1.SchemeGroupVersion.WithKind("SecretStore")) {
		t.Error("scheme does not recognize SecretStore")
	}
	if !scheme.Recognizes(esv1beta1.SchemeGroupVersion.WithKind("ClusterSecretStore")) {
		t.Error("scheme does not recognize ClusterSecretStore")
	}
}
