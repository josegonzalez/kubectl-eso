package cmd

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/josegonzalez/kubectl-eso/pkg/eso"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// NewDescribeSecretCmd creates the describe secret subcommand.
func NewDescribeSecretCmd(streams genericclioptions.IOStreams, configFlags *genericclioptions.ConfigFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "secret <name>",
		Short:   "Show details of a Secret",
		Aliases: []string{"secrets", "Secret", "Secrets"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDescribeSecret(cmd, streams, configFlags, args[0])
		},
	}

	cmd.Flags().BoolP("decode", "d", false, "Decode base64 secret values")

	return cmd
}

func runDescribeSecret(cmd *cobra.Command, streams genericclioptions.IOStreams, configFlags *genericclioptions.ConfigFlags, name string) error {
	decode, _ := cmd.Flags().GetBool("decode")
	output, _ := cmd.Flags().GetString("output")

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
		return fmt.Errorf("secret %q is a Helm-managed release secret and cannot be described by this plugin", name)
	}

	switch output {
	case "json":
		return printJSON(streams.Out, secret)
	case "yaml":
		return printYAML(streams.Out, secret)
	default:
		return printSecretDetail(streams.Out, secret.Name, secret.Namespace,
			string(secret.Type), secret.Labels, secret.Annotations, secret.Data, decode)
	}
}

func printSecretDetail(out io.Writer, name, namespace, secretType string, labels, annotations map[string]string, data map[string][]byte, decode bool) error {
	w := tabwriter.NewWriter(out, 0, 4, 2, ' ', 0)
	defer w.Flush()

	fmt.Fprintf(w, "Name:\t%s\n", name)
	fmt.Fprintf(w, "Namespace:\t%s\n", namespace)
	fmt.Fprintf(w, "Type:\t%s\n", secretType)

	managed := "No"
	if labels != nil && labels[eso.LabelManaged] == "true" {
		managed = "Yes"
	}
	fmt.Fprintf(w, "ESO-Managed:\t%s\n", managed)

	if managed == "No" {
		fmt.Fprintln(w, "\nWARNING: This secret is not managed by External Secrets Operator")
	}

	if annotations != nil {
		if store, ok := annotations[eso.AnnotationStore]; ok {
			fmt.Fprintf(w, "Store:\t%s\n", store)
		}
	}

	fmt.Fprintln(w, "\nData:")
	fmt.Fprintf(w, "  KEY\tVALUE\n")

	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := data[k]
		if decode {
			fmt.Fprintf(w, "  %s\t%s\n", k, string(v))
		} else {
			fmt.Fprintf(w, "  %s\t%s\n", k, base64.StdEncoding.EncodeToString(v))
		}
	}

	return nil
}
