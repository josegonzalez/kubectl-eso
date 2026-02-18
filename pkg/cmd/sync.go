package cmd

import (
	"context"
	"fmt"
	"strconv"
	"time"

	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	"github.com/josegonzalez/kubectl-eso/pkg/eso"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NewSyncCmd creates the sync subcommand.
func NewSyncCmd(streams genericclioptions.IOStreams, configFlags *genericclioptions.ConfigFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "sync <name>",
		Short: "Force re-sync an ExternalSecret",
		Long:  "Force re-sync an ExternalSecret by setting a force-sync annotation.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSync(cmd, streams, configFlags, args[0])
		},
	}
}

func runSync(cmd *cobra.Command, streams genericclioptions.IOStreams, configFlags *genericclioptions.ConfigFlags, name string) error {
	clients, err := getClientsFn(configFlags)
	if err != nil {
		return err
	}

	namespace, err := getNamespaceFn(configFlags)
	if err != nil {
		return err
	}

	var es esv1beta1.ExternalSecret
	key := client.ObjectKey{Namespace: namespace, Name: name}
	if err := clients.CRClient.Get(context.TODO(), key, &es); err != nil {
		return fmt.Errorf("failed to get ExternalSecret %q: %w", name, err)
	}

	if es.Annotations == nil {
		es.Annotations = make(map[string]string)
	}
	es.Annotations[eso.AnnotationForceSync] = strconv.FormatInt(time.Now().Unix(), 10)

	if err := clients.CRClient.Update(context.TODO(), &es); err != nil {
		return fmt.Errorf("failed to update ExternalSecret %q: %w", name, err)
	}

	_, err = fmt.Fprintf(streams.Out, "externalsecret/%s sync triggered\n", name)
	return err
}
