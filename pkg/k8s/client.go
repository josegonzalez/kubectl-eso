package k8s

import (
	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NewScheme creates a runtime.Scheme with core and ESO types registered.
func NewScheme() (*runtime.Scheme, error) {
	s := runtime.NewScheme()
	if err := corev1.AddToScheme(s); err != nil {
		return nil, err
	}
	if err := esv1beta1.AddToScheme(s); err != nil {
		return nil, err
	}
	return s, nil
}

// Clients holds both the controller-runtime client and typed kubernetes clientset.
type Clients struct {
	// CRClient is a controller-runtime client for ESO CRDs and core resources.
	CRClient client.Client
	// Clientset is the typed Kubernetes clientset for core resources.
	Clientset kubernetes.Interface
}

// NewClients creates both clients from a rest.Config.
func NewClients(config *rest.Config) (*Clients, error) {
	scheme, err := NewScheme()
	if err != nil {
		return nil, err
	}

	crClient, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Clients{
		CRClient:  crClient,
		Clientset: clientset,
	}, nil
}
