# kubectl-eso Documentation

## Commands

### annotate

Annotates an existing Kubernetes Secret for ESO adoption. Adds labels and annotations so a future ExternalSecret with `creationPolicy: Merge` can adopt the Secret.

**Labels/annotations applied:**

- Label: `reconcile.external-secrets.io/managed: "true"`
- Annotation: `kubectl-eso.io/imported: "true"`
- Annotation: `kubectl-eso.io/imported-at: <RFC3339 timestamp>`
- Annotation: `kubectl-eso.io/store: <name>` (if `--store` provided)
- Annotation: `kubectl-eso.io/store-kind: <kind>` (if `--store` provided)

```bash
kubectl eso annotate <secret-name> [--store <name>] [--store-kind <kind>] [--dry-run]
```

| Flag | Description | Default |
|------|-------------|---------|
| `--dry-run` | Output annotated Secret as YAML without applying | `false` |
| `--store` | Name of the SecretStore or ClusterSecretStore | |
| `--store-kind` | Kind of the store | `SecretStore` |

### get external-secret

Lists ExternalSecrets with target Secret, store ref, refresh interval, and ready status.

```bash
kubectl eso get external-secret [-A] [-o table|json|yaml] [--no-headers]
```

### get secret

Lists Secrets (excluding Helm release secrets) with an ESO-managed indicator.

```bash
kubectl eso get secret [-A] [-o table|json|yaml] [--no-headers]
```

**Table columns:** NAME, NAMESPACE (if `-A`), TYPE, ESO-MANAGED, STORE, AGE

### get secret-store

Lists SecretStores with health status, provider type, and age.

```bash
kubectl eso get secret-store [-A] [-o table|json|yaml] [--no-headers]
```

### get cluster-secret-store

Lists ClusterSecretStores with health status, provider type, and age.

```bash
kubectl eso get cluster-secret-store [-o table|json|yaml] [--no-headers]
```

### describe external-secret

Shows detailed sync status including conditions, target secret, store ref, refresh interval, and last sync time.

```bash
kubectl eso describe external-secret <name> [-o table|json|yaml]
```

### describe secret

Views Secret data. Base64-encoded by default, decoded with `--decode`/`-d`. Shows a warning if the Secret is not ESO-managed. Errors if the Secret is a Helm release secret.

```bash
kubectl eso describe secret <name> [--decode|-d] [-o table|json|yaml]
```

### describe secret-store

Shows SecretStore health details including conditions, provider type, and last transition time.

```bash
kubectl eso describe secret-store <name> [-o table|json|yaml]
```

### describe cluster-secret-store

Shows ClusterSecretStore health details including conditions, provider type, and last transition time.

```bash
kubectl eso describe cluster-secret-store <name> [-o table|json|yaml]
```

### sync

Forces re-sync of an ExternalSecret by setting a `kubectl-eso.io/force-sync` annotation with the current unix timestamp.

```bash
kubectl eso sync <name>
```

### completion

Generates shell completion scripts.

```bash
kubectl eso completion <bash|zsh|fish|powershell>
```

### version

Prints version, commit, and build date.

```bash
kubectl eso version
```

## Global Flags

All standard kubectl flags are supported:

| Flag | Short | Description |
|------|-------|-------------|
| `--kubeconfig` | | Path to kubeconfig file |
| `--context` | | Kubeconfig context to use |
| `--namespace` | `-n` | Target namespace |
| `--cluster` | | Kubeconfig cluster to use |
| `--user` | | Kubeconfig user to use |
| `--server` | `-s` | Kubernetes API server address |
| `--token` | | Bearer token for authentication |
| `--as` | | Username to impersonate |
| `--as-group` | | Group to impersonate |
| `--as-uid` | | UID to impersonate |
| `--certificate-authority` | | CA cert file |
| `--client-certificate` | | TLS client cert file |
| `--client-key` | | TLS client key file |
| `--insecure-skip-tls-verify` | | Skip server cert verification |
| `--tls-server-name` | | Server name for TLS cert validation |
| `--request-timeout` | | Request timeout duration |
| `--cache-dir` | | HTTP cache directory |

Additional flags:

| Flag | Short | Description | Default | Applies To |
|------|-------|-------------|---------|------------|
| `--output` | `-o` | Output format: `table`, `json`, `yaml` | `table` | get, describe |
| `--no-headers` | | Omit table header row | `false` | get |

## Helm Secret Filtering

Secrets with type `helm.sh/release.v1` are automatically excluded:

- `get secret` silently filters them from the listing
- `describe secret` returns an error

## RBAC Requirements

| Subcommand | Resource | Verbs |
|-----------|----------|-------|
| `annotate` | `secrets` | `get`, `update` |
| `get secret` | `secrets` | `list` |
| `describe secret` | `secrets` | `get` |
| `get external-secret` | `externalsecrets` | `list` |
| `describe external-secret` | `externalsecrets` | `get` |
| `sync` | `externalsecrets` | `get`, `update` |
| `get secret-store` | `secretstores` | `list` |
| `describe secret-store` | `secretstores` | `get` |
| `get cluster-secret-store` | `clustersecretstores` | `list` |
| `describe cluster-secret-store` | `clustersecretstores` | `get` |
