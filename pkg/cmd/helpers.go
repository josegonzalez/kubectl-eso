package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"
	"time"

	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	"github.com/josegonzalez/kubectl-eso/pkg/k8s"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"sigs.k8s.io/yaml"
)

// tableWriter wraps a tabwriter.Writer and tracks the first write error.
type tableWriter struct {
	w   *tabwriter.Writer
	err error
}

func newTableWriter(out io.Writer) *tableWriter {
	return &tableWriter{w: tabwriter.NewWriter(out, 0, 4, 2, ' ', 0)}
}

func (tw *tableWriter) fprintf(format string, a ...any) {
	if tw.err != nil {
		return
	}
	_, tw.err = fmt.Fprintf(tw.w, format, a...)
}

func (tw *tableWriter) fprintln(a ...any) {
	if tw.err != nil {
		return
	}
	_, tw.err = fmt.Fprintln(tw.w, a...)
}

func (tw *tableWriter) fprint(a ...any) {
	if tw.err != nil {
		return
	}
	_, tw.err = fmt.Fprint(tw.w, a...)
}

func (tw *tableWriter) flush() error {
	if tw.err != nil {
		return tw.err
	}
	return tw.w.Flush()
}

// getClientsFn builds k8s clients from configFlags.
var getClientsFn = func(configFlags *genericclioptions.ConfigFlags) (*k8s.Clients, error) {
	config, err := configFlags.ToRESTConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get REST config: %w", err)
	}

	clients, err := k8s.NewClients(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clients: %w", err)
	}

	return clients, nil
}

// getNamespaceFn resolves the active namespace from configFlags.
var getNamespaceFn = func(configFlags *genericclioptions.ConfigFlags) (string, error) {
	namespace, _, err := configFlags.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return "", fmt.Errorf("failed to get namespace: %w", err)
	}

	return namespace, nil
}

// getStoreReadyStatus extracts Ready status from SecretStore conditions.
func getStoreReadyStatus(conditions []esv1beta1.SecretStoreStatusCondition) string {
	for _, c := range conditions {
		if c.Type == "Ready" {
			return string(c.Status)
		}
	}

	return "Unknown"
}

// getESReadyStatus extracts Ready status from ExternalSecret conditions.
func getESReadyStatus(conditions []esv1beta1.ExternalSecretStatusCondition) string {
	for _, c := range conditions {
		if c.Type == esv1beta1.ExternalSecretReady {
			if c.Status == "True" {
				return "True"
			}

			return "False"
		}
	}

	return "Unknown"
}

// printConditionsTable prints a SecretStore conditions table.
func printConditionsTable(tw *tableWriter, conditions []esv1beta1.SecretStoreStatusCondition) {
	tw.fprintln("\nConditions:")
	tw.fprintf("  TYPE\tSTATUS\tREASON\tMESSAGE\tLAST TRANSITION\n")

	for _, c := range conditions {
		tw.fprintf("  %s\t%s\t%s\t%s\t%s\n",
			c.Type, c.Status, c.Reason, c.Message, c.LastTransitionTime.Time)
	}
}

// formatAge formats a time as a human-readable age string.
func formatAge(t time.Time) string {
	if t.IsZero() {
		return "<unknown>"
	}

	d := time.Since(t)

	switch {
	case d < time.Minute:
		return fmt.Sprintf("%ds", int(d.Seconds()))
	case d < time.Hour:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd", int(d.Hours()/24))
	}
}

// printJSON writes the object as indented JSON.
func printJSON(out io.Writer, obj interface{}) error {
	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(out, string(data))

	return err
}

// printYAML writes the object as YAML.
func printYAML(out io.Writer, obj interface{}) error {
	data, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(out, string(data))

	return err
}
