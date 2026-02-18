package eso

import "testing"

func TestConstants(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{"LabelManaged", LabelManaged, "reconcile.external-secrets.io/managed"},
		{"AnnotationImported", AnnotationImported, "kubectl-eso.io/imported"},
		{"AnnotationImportedAt", AnnotationImportedAt, "kubectl-eso.io/imported-at"},
		{"AnnotationStore", AnnotationStore, "kubectl-eso.io/store"},
		{"AnnotationStoreKind", AnnotationStoreKind, "kubectl-eso.io/store-kind"},
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
