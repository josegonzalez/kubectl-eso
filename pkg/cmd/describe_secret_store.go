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

// NewDescribeSecretStoreCmd creates the describe secret-store subcommand.
func NewDescribeSecretStoreCmd(streams genericclioptions.IOStreams, configFlags *genericclioptions.ConfigFlags) *cobra.Command {
	return &cobra.Command{
		Use:     "secret-store <name>",
		Short:   "Show details of a SecretStore",
		Aliases: []string{"secret-stores", "secretstores.external-secrets.io"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDescribeStoreVariant(cmd, streams, configFlags, args[0], false)
		},
	}
}

// NewDescribeClusterSecretStoreCmd creates the describe cluster-secret-store subcommand.
func NewDescribeClusterSecretStoreCmd(streams genericclioptions.IOStreams, configFlags *genericclioptions.ConfigFlags) *cobra.Command {
	return &cobra.Command{
		Use:     "cluster-secret-store <name>",
		Short:   "Show details of a ClusterSecretStore",
		Aliases: []string{"cluster-secret-stores", "clustersecretstores.external-secrets.io"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDescribeStoreVariant(cmd, streams, configFlags, args[0], true)
		},
	}
}

func runDescribeStoreVariant(cmd *cobra.Command, streams genericclioptions.IOStreams, configFlags *genericclioptions.ConfigFlags, name string, clusterScoped bool) error {
	output, _ := cmd.Flags().GetString("output")

	clients, err := getClientsFn(configFlags)
	if err != nil {
		return err
	}

	if clusterScoped {
		var css esv1beta1.ClusterSecretStore
		if err := clients.CRClient.Get(context.TODO(), client.ObjectKey{Name: name}, &css); err != nil {
			return fmt.Errorf("failed to get ClusterSecretStore %q: %w", name, err)
		}

		switch output {
		case "json":
			return printJSON(streams.Out, css)
		case "yaml":
			return printYAML(streams.Out, css)
		default:
			return printStoreDetail(streams.Out, css.Name, "", css.Spec.Provider, css.Status.Conditions)
		}
	}

	namespace, err := getNamespaceFn(configFlags)
	if err != nil {
		return err
	}

	var ss esv1beta1.SecretStore
	if err := clients.CRClient.Get(context.TODO(), client.ObjectKey{Namespace: namespace, Name: name}, &ss); err != nil {
		return fmt.Errorf("failed to get SecretStore %q: %w", name, err)
	}

	switch output {
	case "json":
		return printJSON(streams.Out, ss)
	case "yaml":
		return printYAML(streams.Out, ss)
	default:
		return printStoreDetail(streams.Out, ss.Name, ss.Namespace, ss.Spec.Provider, ss.Status.Conditions)
	}
}

func printStoreDetail(out io.Writer, name, namespace string, provider *esv1beta1.SecretStoreProvider, conditions []esv1beta1.SecretStoreStatusCondition) error {
	tw := newTableWriter(out)

	tw.fprintf("Name:\t%s\n", name)

	if namespace != "" {
		tw.fprintf("Namespace:\t%s\n", namespace)
	} else {
		tw.fprintf("Scope:\tCluster\n")
	}

	tw.fprintf("Provider:\t%s\n", getProviderName(provider))
	tw.fprintf("Ready:\t%s\n", getStoreReadyStatus(conditions))

	printConditionsTable(tw, conditions)

	return tw.flush()
}
