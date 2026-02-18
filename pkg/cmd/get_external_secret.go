package cmd

import (
	"context"
	"fmt"
	"io"

	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NewGetExternalSecretCmd creates the get external-secret subcommand.
func NewGetExternalSecretCmd(streams genericclioptions.IOStreams, configFlags *genericclioptions.ConfigFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "external-secret",
		Short:   "List ExternalSecrets",
		Aliases: []string{"external-secrets", "es", "ExternalSecret", "ExternalSecrets"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGetExternalSecret(cmd, streams, configFlags)
		},
	}

	cmd.Flags().BoolP("all-namespaces", "A", false, "List across all namespaces")

	return cmd
}

func runGetExternalSecret(cmd *cobra.Command, streams genericclioptions.IOStreams, configFlags *genericclioptions.ConfigFlags) error {
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

	var listOpts []client.ListOption
	if !allNamespaces {
		listOpts = append(listOpts, client.InNamespace(namespace))
	}

	var esList esv1beta1.ExternalSecretList
	if err := clients.CRClient.List(context.TODO(), &esList, listOpts...); err != nil {
		return fmt.Errorf("failed to list ExternalSecrets: %w", err)
	}

	switch output {
	case "json":
		return printJSON(streams.Out, esList)
	case "yaml":
		return printYAML(streams.Out, esList)
	default:
		return printExternalSecretTable(streams.Out, esList.Items, allNamespaces, noHeaders)
	}
}

func printExternalSecretTable(out io.Writer, items []esv1beta1.ExternalSecret, allNamespaces bool, noHeaders bool) error {
	tw := newTableWriter(out)

	if !noHeaders {
		if allNamespaces {
			tw.fprint("NAMESPACE\t")
		}
		tw.fprintln("NAME\tSTORE\tREFRESH INTERVAL\tREADY\tAGE")
	}

	for _, es := range items {
		if allNamespaces {
			tw.fprintf("%s\t", es.Namespace)
		}

		store := ""
		if es.Spec.SecretStoreRef.Name != "" {
			store = es.Spec.SecretStoreRef.Name
		}

		refreshInterval := ""
		if es.Spec.RefreshInterval != nil {
			refreshInterval = es.Spec.RefreshInterval.Duration.String()
		}

		ready := getESReadyStatus(es.Status.Conditions)
		age := formatAge(es.CreationTimestamp.Time)

		tw.fprintf("%s\t%s\t%s\t%s\t%s\n", es.Name, store, refreshInterval, ready, age)
	}

	return tw.flush()
}
