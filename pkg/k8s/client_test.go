package k8s

import (
	"testing"

	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
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

func TestNewClients(t *testing.T) {
	// Use a minimal REST config pointing to a dummy server.
	// Client creation succeeds without actually connecting.
	config := &rest.Config{Host: "https://127.0.0.1:0"}

	clients, err := NewClients(config)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if clients.CRClient == nil {
		t.Error("expected non-nil CRClient")
	}
	if clients.Clientset == nil {
		t.Error("expected non-nil Clientset")
	}
}
