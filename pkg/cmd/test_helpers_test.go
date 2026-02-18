package cmd

import (
	"fmt"

	"github.com/josegonzalez/kubectl-eso/pkg/k8s"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"sigs.k8s.io/controller-runtime/pkg/client"
	fakecr "sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
)

// setupFakeClients overrides getClientsFn and getNamespaceFn for testing.
// It returns a cleanup function that restores the originals.
func setupFakeClients(namespace string, crObjects []client.Object, k8sObjects ...runtime.Object) func() {
	origClients := getClientsFn
	origNamespace := getNamespaceFn

	scheme, _ := k8s.NewScheme()
	crClient := fakecr.NewClientBuilder().WithScheme(scheme).WithObjects(crObjects...).Build()
	clientset := fake.NewClientset(k8sObjects...)

	getClientsFn = func(_ *genericclioptions.ConfigFlags) (*k8s.Clients, error) {
		return &k8s.Clients{
			CRClient:  crClient,
			Clientset: clientset,
		}, nil
	}

	getNamespaceFn = func(_ *genericclioptions.ConfigFlags) (string, error) {
		return namespace, nil
	}

	return func() {
		getClientsFn = origClients
		getNamespaceFn = origNamespace
	}
}

// setupFailingClients overrides getClientsFn to return an error.
func setupFailingClients() func() {
	origClients := getClientsFn
	origNamespace := getNamespaceFn

	getClientsFn = func(_ *genericclioptions.ConfigFlags) (*k8s.Clients, error) {
		return nil, fmt.Errorf("client error")
	}

	getNamespaceFn = func(_ *genericclioptions.ConfigFlags) (string, error) {
		return "default", nil
	}

	return func() {
		getClientsFn = origClients
		getNamespaceFn = origNamespace
	}
}

// setupFakeClientsWithInterceptors is like setupFakeClients but allows injecting
// interceptor functions for the CR client (e.g., to simulate List/Update errors).
func setupFakeClientsWithInterceptors(namespace string, crObjects []client.Object, funcs interceptor.Funcs, k8sObjects ...runtime.Object) func() {
	origClients := getClientsFn
	origNamespace := getNamespaceFn

	scheme, _ := k8s.NewScheme()
	crClient := fakecr.NewClientBuilder().WithScheme(scheme).WithObjects(crObjects...).WithInterceptorFuncs(funcs).Build()
	clientset := fake.NewClientset(k8sObjects...)

	getClientsFn = func(_ *genericclioptions.ConfigFlags) (*k8s.Clients, error) {
		return &k8s.Clients{
			CRClient:  crClient,
			Clientset: clientset,
		}, nil
	}

	getNamespaceFn = func(_ *genericclioptions.ConfigFlags) (string, error) {
		return namespace, nil
	}

	return func() {
		getClientsFn = origClients
		getNamespaceFn = origNamespace
	}
}

// setupFakeClientsWithClientset is like setupFakeClients but accepts a pre-built
// clientset, allowing callers to inject reactors for simulating k8s API errors.
func setupFakeClientsWithClientset(namespace string, crObjects []client.Object, clientset kubernetes.Interface) func() {
	origClients := getClientsFn
	origNamespace := getNamespaceFn

	scheme, _ := k8s.NewScheme()
	crClient := fakecr.NewClientBuilder().WithScheme(scheme).WithObjects(crObjects...).Build()

	getClientsFn = func(_ *genericclioptions.ConfigFlags) (*k8s.Clients, error) {
		return &k8s.Clients{
			CRClient:  crClient,
			Clientset: clientset,
		}, nil
	}

	getNamespaceFn = func(_ *genericclioptions.ConfigFlags) (string, error) {
		return namespace, nil
	}

	return func() {
		getClientsFn = origClients
		getNamespaceFn = origNamespace
	}
}

// setupFailingNamespace overrides getNamespaceFn to return an error
// while getClientsFn returns valid fake clients.
func setupFailingNamespace() func() {
	origClients := getClientsFn
	origNamespace := getNamespaceFn

	scheme, _ := k8s.NewScheme()
	crClient := fakecr.NewClientBuilder().WithScheme(scheme).Build()
	clientset := fake.NewClientset()

	getClientsFn = func(_ *genericclioptions.ConfigFlags) (*k8s.Clients, error) {
		return &k8s.Clients{
			CRClient:  crClient,
			Clientset: clientset,
		}, nil
	}

	getNamespaceFn = func(_ *genericclioptions.ConfigFlags) (string, error) {
		return "", fmt.Errorf("namespace error")
	}

	return func() {
		getClientsFn = origClients
		getNamespaceFn = origNamespace
	}
}
