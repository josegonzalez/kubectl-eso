package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/josegonzalez/kubectl-eso/pkg/eso"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// NewAnnotateCmd creates the annotate subcommand.
func NewAnnotateCmd(streams genericclioptions.IOStreams, configFlags *genericclioptions.ConfigFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "annotate <secret-name>",
		Short: "Annotate a Secret for ESO adoption",
		Long:  "Annotate an existing Kubernetes Secret so it can be adopted by a future ExternalSecret with creationPolicy: Merge.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAnnotate(cmd, streams, configFlags, args[0])
		},
	}

	cmd.Flags().Bool("dry-run", false, "Output annotated Secret as YAML without applying")
	cmd.Flags().String("store", "", "Name of the SecretStore or ClusterSecretStore")
	cmd.Flags().String("store-kind", "SecretStore", "Kind of the store (SecretStore or ClusterSecretStore)")

	return cmd
}

func runAnnotate(cmd *cobra.Command, streams genericclioptions.IOStreams, configFlags *genericclioptions.ConfigFlags, name string) error {
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	store, _ := cmd.Flags().GetString("store")
	storeKind, _ := cmd.Flags().GetString("store-kind")

	clients, err := getClients(configFlags)
	if err != nil {
		return err
	}

	namespace, err := getNamespace(configFlags)
	if err != nil {
		return err
	}

	secret, err := clients.Clientset.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get Secret %q: %w", name, err)
	}

	// Reject Helm secrets
	if string(secret.Type) == eso.HelmSecretType {
		return fmt.Errorf("secret %q is a Helm-managed release secret and cannot be annotated by this plugin", name)
	}

	// Set labels
	if secret.Labels == nil {
		secret.Labels = make(map[string]string)
	}
	secret.Labels[eso.LabelManaged] = "true"

	// Set annotations
	if secret.Annotations == nil {
		secret.Annotations = make(map[string]string)
	}
	secret.Annotations[eso.AnnotationImported] = "true"
	secret.Annotations[eso.AnnotationImportedAt] = time.Now().UTC().Format(time.RFC3339)

	if store != "" {
		secret.Annotations[eso.AnnotationStore] = store
		secret.Annotations[eso.AnnotationStoreKind] = storeKind
	}

	if dryRun {
		// Clear managed fields for cleaner output
		secret.ManagedFields = nil
		return printYAML(streams.Out, secret)
	}

	_, err = clients.Clientset.CoreV1().Secrets(namespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update Secret %q: %w", name, err)
	}

	fmt.Fprintf(streams.Out, "secret/%s annotated for ESO adoption\n", name)
	return nil
}
