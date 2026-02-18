package eso

const (
	// LabelManaged is the ESO standard label for managed secrets.
	LabelManaged = "reconcile.external-secrets.io/managed"

	// AnnotationImported indicates the secret was imported via kubectl-eso.
	AnnotationImported = "kubectl-eso.io/imported"

	// AnnotationImportedAt is the timestamp when the secret was imported.
	AnnotationImportedAt = "kubectl-eso.io/imported-at"

	// AnnotationStore is the name of the secret store.
	AnnotationStore = "kubectl-eso.io/store"

	// AnnotationStoreKind is the kind of the secret store.
	AnnotationStoreKind = "kubectl-eso.io/store-kind"

	// AnnotationForceSync forces a re-sync of an ExternalSecret.
	AnnotationForceSync = "kubectl-eso.io/force-sync"

	// HelmSecretType is the type used by Helm-managed release secrets.
	HelmSecretType = "helm.sh/release.v1"
)
