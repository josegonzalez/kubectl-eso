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

// storeItem is an intermediate representation for both SecretStore and ClusterSecretStore table rows.
type storeItem struct {
	Name, Namespace, Ready, Provider, Age string
}

// NewGetSecretStoreCmd creates the get secret-store subcommand.
func NewGetSecretStoreCmd(streams genericclioptions.IOStreams, configFlags *genericclioptions.ConfigFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "secret-store",
		Short:   "List SecretStores",
		Aliases: []string{"secret-stores", "ss", "SecretStore", "SecretStores"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGetStoreVariant(cmd, streams, configFlags, false)
		},
	}

	cmd.Flags().BoolP("all-namespaces", "A", false, "List across all namespaces")

	return cmd
}

// NewGetClusterSecretStoreCmd creates the get cluster-secret-store subcommand.
func NewGetClusterSecretStoreCmd(streams genericclioptions.IOStreams, configFlags *genericclioptions.ConfigFlags) *cobra.Command {
	return &cobra.Command{
		Use:     "cluster-secret-store",
		Short:   "List ClusterSecretStores",
		Aliases: []string{"cluster-secret-stores", "css", "ClusterSecretStore", "ClusterSecretStores"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGetStoreVariant(cmd, streams, configFlags, true)
		},
	}
}

func runGetStoreVariant(cmd *cobra.Command, streams genericclioptions.IOStreams, configFlags *genericclioptions.ConfigFlags, clusterScoped bool) error {
	output, _ := cmd.Flags().GetString("output")
	noHeaders, _ := cmd.Flags().GetBool("no-headers")

	clients, err := getClientsFn(configFlags)
	if err != nil {
		return err
	}

	if clusterScoped {
		var cssList esv1beta1.ClusterSecretStoreList
		if err := clients.CRClient.List(context.TODO(), &cssList, &client.ListOptions{}); err != nil {
			return fmt.Errorf("failed to list ClusterSecretStores: %w", err)
		}

		switch output {
		case "json":
			return printJSON(streams.Out, cssList)
		case "yaml":
			return printYAML(streams.Out, cssList)
		default:
			items := make([]storeItem, len(cssList.Items))
			for i, css := range cssList.Items {
				items[i] = storeItem{
					Name:     css.Name,
					Ready:    getStoreReadyStatus(css.Status.Conditions),
					Provider: getProviderName(css.Spec.Provider),
					Age:      formatAge(css.CreationTimestamp.Time),
				}
			}

			return printStoreTable(streams.Out, items, false, noHeaders)
		}
	}

	allNamespaces, _ := cmd.Flags().GetBool("all-namespaces")

	namespace, err := getNamespaceFn(configFlags)
	if err != nil {
		return err
	}

	var listOpts []client.ListOption
	if !allNamespaces {
		listOpts = append(listOpts, client.InNamespace(namespace))
	}

	var ssList esv1beta1.SecretStoreList
	if err := clients.CRClient.List(context.TODO(), &ssList, listOpts...); err != nil {
		return fmt.Errorf("failed to list SecretStores: %w", err)
	}

	switch output {
	case "json":
		return printJSON(streams.Out, ssList)
	case "yaml":
		return printYAML(streams.Out, ssList)
	default:
		items := make([]storeItem, len(ssList.Items))
		for i, ss := range ssList.Items {
			items[i] = storeItem{
				Name:      ss.Name,
				Namespace: ss.Namespace,
				Ready:     getStoreReadyStatus(ss.Status.Conditions),
				Provider:  getProviderName(ss.Spec.Provider),
				Age:       formatAge(ss.CreationTimestamp.Time),
			}
		}

		return printStoreTable(streams.Out, items, allNamespaces, noHeaders)
	}
}

func printStoreTable(out io.Writer, items []storeItem, allNamespaces bool, noHeaders bool) error {
	tw := newTableWriter(out)

	if !noHeaders {
		if allNamespaces {
			tw.fprint("NAMESPACE\t")
		}
		tw.fprintln("NAME\tREADY\tPROVIDER\tAGE")
	}

	for _, item := range items {
		if allNamespaces {
			tw.fprintf("%s\t", item.Namespace)
		}

		tw.fprintf("%s\t%s\t%s\t%s\n", item.Name, item.Ready, item.Provider, item.Age)
	}

	return tw.flush()
}

func getProviderName(provider *esv1beta1.SecretStoreProvider) string {
	if provider == nil {
		return "<none>"
	}

	if provider.AWS != nil {
		return "AWS"
	}

	if provider.AzureKV != nil {
		return "AzureKV"
	}

	if provider.GCPSM != nil {
		return "GCPSM"
	}

	if provider.Vault != nil {
		return "Vault"
	}

	if provider.Kubernetes != nil {
		return "Kubernetes"
	}

	if provider.Fake != nil {
		return "Fake"
	}

	return "<unknown>"
}
