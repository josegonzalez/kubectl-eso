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

// NewDescribeExternalSecretCmd creates the describe external-secret subcommand.
func NewDescribeExternalSecretCmd(streams genericclioptions.IOStreams, configFlags *genericclioptions.ConfigFlags) *cobra.Command {
	return &cobra.Command{
		Use:     "external-secret <name>",
		Short:   "Show details of an ExternalSecret",
		Aliases: []string{"external-secrets", "externalsecrets.external-secrets.io"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDescribeExternalSecret(cmd, streams, configFlags, args[0])
		},
	}
}

func runDescribeExternalSecret(cmd *cobra.Command, streams genericclioptions.IOStreams, configFlags *genericclioptions.ConfigFlags, name string) error {
	output, _ := cmd.Flags().GetString("output")

	clients, err := getClientsFn(configFlags)
	if err != nil {
		return err
	}

	namespace, err := getNamespaceFn(configFlags)
	if err != nil {
		return err
	}

	var es esv1beta1.ExternalSecret
	if err := clients.CRClient.Get(context.TODO(), client.ObjectKey{Namespace: namespace, Name: name}, &es); err != nil {
		return fmt.Errorf("failed to get ExternalSecret %q: %w", name, err)
	}

	switch output {
	case "json":
		return printJSON(streams.Out, es)
	case "yaml":
		return printYAML(streams.Out, es)
	default:
		return printExternalSecretDetail(streams.Out, es)
	}
}

func printExternalSecretDetail(out io.Writer, es esv1beta1.ExternalSecret) error {
	tw := newTableWriter(out)

	tw.fprintf("Name:\t%s\n", es.Name)
	tw.fprintf("Namespace:\t%s\n", es.Namespace)

	// Store ref
	tw.fprintf("Store:\t%s\n", es.Spec.SecretStoreRef.Name)
	if es.Spec.SecretStoreRef.Kind != "" {
		tw.fprintf("Store Kind:\t%s\n", es.Spec.SecretStoreRef.Kind)
	}

	// Target
	if es.Spec.Target.Name != "" {
		tw.fprintf("Target Secret:\t%s\n", es.Spec.Target.Name)
	} else {
		tw.fprintf("Target Secret:\t%s\n", es.Name)
	}

	// Refresh interval
	if es.Spec.RefreshInterval != nil {
		tw.fprintf("Refresh Interval:\t%s\n", es.Spec.RefreshInterval.Duration)
	}

	// Creation policy
	if es.Spec.Target.CreationPolicy != "" {
		tw.fprintf("Creation Policy:\t%s\n", es.Spec.Target.CreationPolicy)
	}

	// Status
	tw.fprintln("\nConditions:")
	tw.fprintf("  TYPE\tSTATUS\tREASON\tMESSAGE\tLAST TRANSITION\n")
	for _, c := range es.Status.Conditions {
		tw.fprintf("  %s\t%s\t%s\t%s\t%s\n",
			c.Type, c.Status, c.Reason, c.Message, c.LastTransitionTime.Time)
	}

	if es.Status.SyncedResourceVersion != "" {
		tw.fprintf("\nSynced Resource Version:\t%s\n", es.Status.SyncedResourceVersion)
	}

	return tw.flush()
}
