package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	"github.com/josegonzalez/kubectl-eso/pkg/k8s"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"sigs.k8s.io/yaml"
)

// getClients builds k8s clients from configFlags.
func getClients(configFlags *genericclioptions.ConfigFlags) (*k8s.Clients, error) {
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

// getNamespace resolves the active namespace from configFlags.
func getNamespace(configFlags *genericclioptions.ConfigFlags) (string, error) {
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
func printConditionsTable(w io.Writer, conditions []esv1beta1.SecretStoreStatusCondition) {
	fmt.Fprintln(w, "\nConditions:")
	fmt.Fprintf(w, "  TYPE\tSTATUS\tREASON\tMESSAGE\tLAST TRANSITION\n")

	for _, c := range conditions {
		fmt.Fprintf(w, "  %s\t%s\t%s\t%s\t%s\n",
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

	fmt.Fprintln(out, string(data))

	return nil
}

// printYAML writes the object as YAML.
func printYAML(out io.Writer, obj interface{}) error {
	data, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}

	fmt.Fprint(out, string(data))

	return nil
}
