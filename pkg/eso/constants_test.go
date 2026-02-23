package eso

import "testing"

func TestConstants(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{"LabelManaged", LabelManaged, "reconcile.external-secrets.io/managed"},
		{"AnnotationForceSync", AnnotationForceSync, "kubectl-eso.io/force-sync"},
		{"HelmSecretType", HelmSecretType, "helm.sh/release.v1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("%s = %q, want %q", tt.name, tt.got, tt.want)
			}
		})
	}
}
