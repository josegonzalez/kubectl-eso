package cmd

import (
	"context"
	"fmt"
	"io"

	"github.com/josegonzalez/kubectl-eso/pkg/eso"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// NewGetSecretCmd creates the get secret subcommand.
func NewGetSecretCmd(streams genericclioptions.IOStreams, configFlags *genericclioptions.ConfigFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "secret",
		Short:   "List Secrets with ESO-managed indicator",
		Aliases: []string{"secrets"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGetSecret(cmd, streams, configFlags)
		},
	}

	cmd.Flags().BoolP("all-namespaces", "A", false, "List across all namespaces")

	return cmd
}

func runGetSecret(cmd *cobra.Command, streams genericclioptions.IOStreams, configFlags *genericclioptions.ConfigFlags) error {
	allNamespaces, _ := cmd.Flags().GetBool("all-namespaces")
	output, _ := cmd.Flags().GetString("output")
	noHeaders, _ := cmd.Flags().GetBool("no-headers")

	clients, err := getClientsFn(configFlags)
	if err != nil {
		return err
	}

	namespace, err := getNamespaceFn(configFlags)
	if err != nil {
		return err
	}

	listNS := namespace
	if allNamespaces {
		listNS = ""
	}

	secretList, err := clients.Clientset.CoreV1().Secrets(listNS).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list Secrets: %w", err)
	}

	// Filter out Helm secrets
	var filtered []corev1.Secret
	for _, s := range secretList.Items {
		if string(s.Type) != eso.HelmSecretType {
			filtered = append(filtered, s)
		}
	}

	switch output {
	case "json":
		return printJSON(streams.Out, filtered)
	case "yaml":
		return printYAML(streams.Out, filtered)
	default:
		return printSecretTable(streams.Out, filtered, allNamespaces, noHeaders)
	}
}

func printSecretTable(out io.Writer, secrets []corev1.Secret, allNamespaces bool, noHeaders bool) error {
	tw := newTableWriter(out)

	if !noHeaders {
		if allNamespaces {
			tw.fprint("NAMESPACE\t")
		}
		tw.fprintln("NAME\tTYPE\tESO-MANAGED\tSTORE\tAGE")
	}

	for _, s := range secrets {
		if allNamespaces {
			tw.fprintf("%s\t", s.Namespace)
		}

		managed := "No"
		if s.Labels != nil && s.Labels[eso.LabelManaged] == "true" {
			managed = "Yes"
		}

		store := ""
		if s.Annotations != nil {
			store = s.Annotations[eso.AnnotationStore]
		}

		age := formatAge(s.CreationTimestamp.Time)

		tw.fprintf("%s\t%s\t%s\t%s\t%s\n", s.Name, s.Type, managed, store, age)
	}

	return tw.flush()
}
