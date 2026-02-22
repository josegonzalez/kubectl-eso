package cmd

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"sort"

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
		Aliases: []string{"secrets"},
		Args:    requireNameArg("Secret"),
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

	clients, err := getClientsFn(configFlags)
	if err != nil {
		return err
	}

	namespace, err := getNamespaceFn(configFlags)
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

	// Resolve store: annotation first, then ExternalSecret lookup
	var store storeRef
	if secret.Annotations != nil {
		if storeName, ok := secret.Annotations[eso.AnnotationStore]; ok {
			kind := "SecretStore"
			if k, ok := secret.Annotations[eso.AnnotationStoreKind]; ok {
				kind = k
			}
			store = storeRef{name: storeName, kind: kind}
		}
	}
	if store.name == "" && secret.Labels != nil && secret.Labels[eso.LabelManaged] == "true" {
		storeMap := buildSecretStoreMap(context.TODO(), clients.CRClient, namespace)
		if ref, ok := storeMap[namespace+"/"+secret.Name]; ok {
			store = ref
		}
	}

	switch output {
	case "json":
		return printJSON(streams.Out, secret)
	case "yaml":
		return printYAML(streams.Out, secret)
	default:
		return printSecretDetail(streams.Out, secret.Name, secret.Namespace,
			string(secret.Type), secret.Labels, secret.Data, decode, store)
	}
}

func printSecretDetail(out io.Writer, name, namespace, secretType string, labels map[string]string, data map[string][]byte, decode bool, store storeRef) error {
	tw := newTableWriter(out)

	tw.fprintf("Name:\t%s\n", name)
	tw.fprintf("Namespace:\t%s\n", namespace)
	tw.fprintf("Type:\t%s\n", secretType)

	managed := "No"
	if labels != nil && labels[eso.LabelManaged] == "true" {
		managed = "Yes"
	}
	tw.fprintf("ESO-Managed:\t%s\n", managed)

	if managed == "No" {
		tw.fprintln("\nWARNING: This secret is not managed by External Secrets Operator")
	}

	if store.name != "" {
		tw.fprintf("Store:\t%s\n", store.name)
		tw.fprintf("Store Kind:\t%s\n", store.kind)
	}

	tw.fprintln("\nData:")
	tw.fprintf("  KEY\tVALUE\n")

	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := data[k]
		if decode {
			tw.fprintf("  %s\t%s\n", k, string(v))
		} else {
			tw.fprintf("  %s\t%s\n", k, base64.StdEncoding.EncodeToString(v))
		}
	}

	return tw.flush()
}
