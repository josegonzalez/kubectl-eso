package eso

const (
	// LabelManaged is the ESO standard label for managed secrets.
	LabelManaged = "reconcile.external-secrets.io/managed"

	// AnnotationForceSync forces a re-sync of an ExternalSecret.
	AnnotationForceSync = "kubectl-eso.io/force-sync"

	// HelmSecretType is the type used by Helm-managed release secrets.
	HelmSecretType = "helm.sh/release.v1"
)
