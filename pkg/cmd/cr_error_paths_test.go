package cmd

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
)

func TestRunGetExternalSecretListError(t *testing.T) {
	cleanup := setupFakeClientsWithInterceptors("default", nil, interceptor.Funcs{
		List: func(ctx context.Context, c client.WithWatch, list client.ObjectList, opts ...client.ListOption) error {
			return fmt.Errorf("list error")
		},
	})
	defer cleanup()

	var buf bytes.Buffer
	streams := genericclioptions.IOStreams{Out: &buf, ErrOut: &buf}
	configFlags := genericclioptions.NewConfigFlags(true)
	cmd := NewGetCmd(streams, configFlags)
	cmd.SetArgs([]string{"external-secret"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error from list, got nil")
	}
	if !strings.Contains(err.Error(), "failed to list ExternalSecrets") {
		t.Errorf("expected 'failed to list ExternalSecrets', got %v", err)
	}
}

func TestRunGetSecretStoreListError(t *testing.T) {
	cleanup := setupFakeClientsWithInterceptors("default", nil, interceptor.Funcs{
		List: func(ctx context.Context, c client.WithWatch, list client.ObjectList, opts ...client.ListOption) error {
			return fmt.Errorf("list error")
		},
	})
	defer cleanup()

	var buf bytes.Buffer
	streams := genericclioptions.IOStreams{Out: &buf, ErrOut: &buf}
	configFlags := genericclioptions.NewConfigFlags(true)
	cmd := NewGetCmd(streams, configFlags)
	cmd.SetArgs([]string{"secret-store"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error from list, got nil")
	}
	if !strings.Contains(err.Error(), "failed to list SecretStores") {
		t.Errorf("expected 'failed to list SecretStores', got %v", err)
	}
}

func TestRunGetClusterSecretStoreListError(t *testing.T) {
	cleanup := setupFakeClientsWithInterceptors("default", nil, interceptor.Funcs{
		List: func(ctx context.Context, c client.WithWatch, list client.ObjectList, opts ...client.ListOption) error {
			return fmt.Errorf("list error")
		},
	})
	defer cleanup()

	var buf bytes.Buffer
	streams := genericclioptions.IOStreams{Out: &buf, ErrOut: &buf}
	configFlags := genericclioptions.NewConfigFlags(true)
	cmd := NewGetCmd(streams, configFlags)
	cmd.SetArgs([]string{"cluster-secret-store"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error from list, got nil")
	}
	if !strings.Contains(err.Error(), "failed to list ClusterSecretStores") {
		t.Errorf("expected 'failed to list ClusterSecretStores', got %v", err)
	}
}

func TestRunSyncUpdateError(t *testing.T) {
	es := &esv1beta1.ExternalSecret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-es",
			Namespace: "default",
		},
		Spec: esv1beta1.ExternalSecretSpec{
			SecretStoreRef: esv1beta1.SecretStoreRef{Name: "store"},
		},
	}

	cleanup := setupFakeClientsWithInterceptors("default", []client.Object{es}, interceptor.Funcs{
		Update: func(ctx context.Context, c client.WithWatch, obj client.Object, opts ...client.UpdateOption) error {
			return fmt.Errorf("update error")
		},
	})
	defer cleanup()

	var buf bytes.Buffer
	streams := genericclioptions.IOStreams{Out: &buf, ErrOut: &buf}
	root := NewRootCmd(streams)
	root.SetArgs([]string{"sync", "my-es"})

	err := root.Execute()
	if err == nil {
		t.Fatal("expected error from update, got nil")
	}
	if !strings.Contains(err.Error(), "failed to update ExternalSecret") {
		t.Errorf("expected 'failed to update ExternalSecret', got %v", err)
	}
}

func TestRunGetExternalSecretNoHeaders(t *testing.T) {
	es := &esv1beta1.ExternalSecret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-es",
			Namespace: "default",
		},
		Spec: esv1beta1.ExternalSecretSpec{
			SecretStoreRef: esv1beta1.SecretStoreRef{Name: "store"},
		},
	}

	cleanup := setupFakeClients("default", []client.Object{es})
	defer cleanup()

	var buf bytes.Buffer
	streams := genericclioptions.IOStreams{Out: &buf}
	configFlags := genericclioptions.NewConfigFlags(true)
	cmd := NewGetCmd(streams, configFlags)
	cmd.SetArgs([]string{"external-secret", "--no-headers"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if strings.Contains(output, "NAME") {
		t.Errorf("expected no headers in output, but found NAME:\n%s", output)
	}
	if !strings.Contains(output, "my-es") {
		t.Errorf("expected data in output:\n%s", output)
	}
}

func TestRunGetSecretNoHeaders(t *testing.T) {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-secret",
			Namespace: "default",
		},
		Type: corev1.SecretTypeOpaque,
	}

	cleanup := setupFakeClients("default", nil, secret)
	defer cleanup()

	var buf bytes.Buffer
	streams := genericclioptions.IOStreams{Out: &buf}
	configFlags := genericclioptions.NewConfigFlags(true)
	cmd := NewGetCmd(streams, configFlags)
	cmd.SetArgs([]string{"secret", "--no-headers"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if strings.Contains(output, "NAME") {
		t.Errorf("expected no headers in output, but found NAME:\n%s", output)
	}
	if !strings.Contains(output, "my-secret") {
		t.Errorf("expected data in output:\n%s", output)
	}
}

func TestRunGetSecretListError(t *testing.T) {
	clientset := fake.NewClientset()
	clientset.PrependReactor("list", "secrets", func(action k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, fmt.Errorf("list error")
	})

	cleanup := setupFakeClientsWithClientset("default", nil, clientset)
	defer cleanup()

	var buf bytes.Buffer
	streams := genericclioptions.IOStreams{Out: &buf, ErrOut: &buf}
	configFlags := genericclioptions.NewConfigFlags(true)
	cmd := NewGetCmd(streams, configFlags)
	cmd.SetArgs([]string{"secret"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error from list, got nil")
	}
	if !strings.Contains(err.Error(), "failed to list Secrets") {
		t.Errorf("expected 'failed to list Secrets', got %v", err)
	}
}

func TestRunAnnotateUpdateError(t *testing.T) {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-secret",
			Namespace: "default",
		},
		Type: corev1.SecretTypeOpaque,
	}

	clientset := fake.NewClientset(secret)
	clientset.PrependReactor("update", "secrets", func(action k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, fmt.Errorf("update error")
	})

	cleanup := setupFakeClientsWithClientset("default", nil, clientset)
	defer cleanup()

	var buf bytes.Buffer
	streams := genericclioptions.IOStreams{Out: &buf, ErrOut: &buf}
	configFlags := genericclioptions.NewConfigFlags(true)
	cmd := NewAnnotateCmd(streams, configFlags)
	cmd.SetArgs([]string{"my-secret"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error from update, got nil")
	}
	if !strings.Contains(err.Error(), "failed to update Secret") {
		t.Errorf("expected 'failed to update Secret', got %v", err)
	}
}

func TestRunGetSecretStoreNoHeaders(t *testing.T) {
	ss := &esv1beta1.SecretStore{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-store",
			Namespace: "default",
		},
		Spec: esv1beta1.SecretStoreSpec{
			Provider: &esv1beta1.SecretStoreProvider{
				AWS: &esv1beta1.AWSProvider{},
			},
		},
	}

	cleanup := setupFakeClients("default", []client.Object{ss})
	defer cleanup()

	var buf bytes.Buffer
	streams := genericclioptions.IOStreams{Out: &buf}
	configFlags := genericclioptions.NewConfigFlags(true)
	cmd := NewGetCmd(streams, configFlags)
	cmd.SetArgs([]string{"secret-store", "--no-headers"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if strings.Contains(output, "NAME") {
		t.Errorf("expected no headers in output, but found NAME:\n%s", output)
	}
	if !strings.Contains(output, "my-store") {
		t.Errorf("expected data in output:\n%s", output)
	}
}
