package cmd

import (
	"bytes"
	"context"
	"strings"
	"testing"

	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	"github.com/josegonzalez/kubectl-eso/pkg/eso"
	"github.com/josegonzalez/kubectl-eso/pkg/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestSyncAnnotation(t *testing.T) {
	scheme, err := k8s.NewScheme()
	if err != nil {
		t.Fatalf("failed to create scheme: %v", err)
	}

	es := &esv1beta1.ExternalSecret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-es",
			Namespace: "default",
		},
		Spec: esv1beta1.ExternalSecretSpec{
			SecretStoreRef: esv1beta1.SecretStoreRef{Name: "store"},
		},
	}

	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(es).Build()

	// Simulate sync logic
	var fetched esv1beta1.ExternalSecret
	err = fakeClient.Get(context.TODO(), client.ObjectKey{Namespace: "default", Name: "my-es"}, &fetched)
	if err != nil {
		t.Fatalf("failed to get ExternalSecret: %v", err)
	}

	if fetched.Annotations == nil {
		fetched.Annotations = make(map[string]string)
	}
	fetched.Annotations[eso.AnnotationForceSync] = "1234567890"

	err = fakeClient.Update(context.TODO(), &fetched)
	if err != nil {
		t.Fatalf("failed to update ExternalSecret: %v", err)
	}

	// Verify
	var updated esv1beta1.ExternalSecret
	err = fakeClient.Get(context.TODO(), client.ObjectKey{Namespace: "default", Name: "my-es"}, &updated)
	if err != nil {
		t.Fatalf("failed to get updated ExternalSecret: %v", err)
	}

	if updated.Annotations[eso.AnnotationForceSync] != "1234567890" {
		t.Errorf("expected force-sync annotation 1234567890, got %s", updated.Annotations[eso.AnnotationForceSync])
	}
}

func TestRunSync(t *testing.T) {
	es := &esv1beta1.ExternalSecret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-es",
			Namespace: "default",
		},
		Spec: esv1beta1.ExternalSecretSpec{
			SecretStoreRef: esv1beta1.SecretStoreRef{Name: "store"},
		},
	}

	t.Run("normal flow", func(t *testing.T) {
		cleanup := setupFakeClients("default", []client.Object{es})
		defer cleanup()

		var buf bytes.Buffer
		streams := genericclioptions.IOStreams{Out: &buf}
		root := NewRootCmd(streams)
		root.SetArgs([]string{"sync", "my-es"})

		if err := root.Execute(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !strings.Contains(buf.String(), "sync triggered") {
			t.Errorf("expected sync triggered message, got %q", buf.String())
		}
	})

	t.Run("missing ExternalSecret", func(t *testing.T) {
		cleanup := setupFakeClients("default", nil)
		defer cleanup()

		var buf bytes.Buffer
		streams := genericclioptions.IOStreams{Out: &buf, ErrOut: &buf}
		root := NewRootCmd(streams)
		root.SetArgs([]string{"sync", "nonexistent"})

		err := root.Execute()
		if err == nil {
			t.Fatal("expected error for missing ExternalSecret, got nil")
		}
	})
}
