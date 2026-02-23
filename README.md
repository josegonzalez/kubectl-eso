# kubectl-eso

A kubectl plugin for managing Kubernetes Secrets in the context of
the [External Secrets Operator](https://external-secrets.io) (v1 API).
It supports listing and inspecting ExternalSecrets, SecretStores, and
ClusterSecretStores, viewing Secret data with optional base64
decoding, forcing re-sync of ExternalSecrets, and shell completion
for bash, zsh, fish, and powershell. Helm release secrets are
automatically filtered out, and all standard kubectl flags
are supported.

## Installation

### Via krew

```bash
kubectl krew install eso
```

### Via go install

```bash
go install github.com/josegonzalez/kubectl-eso/cmd/kubectl-eso@latest
```

### Pre-built binaries

Download from the [releases page](https://github.com/josegonzalez/kubectl-eso/releases).

## Building from Source

```bash
# Build
make build

# Run tests
make test

# Lint
make lint

# Install kubectl-eso and kubectl_complete-eso to ~/.krew/bin
make install
```

The `make install` target also installs the
`kubectl_complete-eso` completion wrapper for kubectl plugin
completion (kubectl v1.26+).

## Usage

```text
kubectl eso COMMAND [SUBCOMMAND] [ARGS...] [FLAGS...]
```

| Command | Description |
| ------- | ----------- |
| `get` | List resources (ExternalSecrets, Secrets, stores) |
| `describe` | Show detailed info for a resource |
| `sync` | Force re-sync an ExternalSecret |
| `completion` | Generate shell completion scripts |
| `version` | Print version information |

### Global Flags

All standard kubectl flags are supported:

| Flag | Short | Description |
| ---- | ----- | ----------- |
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

| Flag | Short | Description | Default |
| ---- | ----- | ----------- | ------- |
| `--output` | `-o` | Output format (`table`, `json`, `yaml`) | `table` |
| `--no-headers` | | Omit table header row | `false` |

## Commands

### `kubectl eso get external-secret [flags]`

Lists ExternalSecrets with target Secret, store ref,
refresh interval, and ready status.

**Requires:** `list` on `externalsecrets`

**Examples:**

```bash
# List ExternalSecrets
kubectl eso get external-secret
kubectl eso get external-secrets -A   # all namespaces

# Output as JSON
kubectl eso get external-secret -o json
```

### `kubectl eso get secret [flags]`

Lists Secrets (excluding Helm release secrets) with an
ESO-managed indicator.

**Requires:** `list` on `secrets`

**Table columns:** NAME, NAMESPACE (if `-A`), TYPE, ESO-MANAGED, STORE, AGE

**Examples:**

```bash
# List Secrets with ESO-managed indicator
kubectl eso get secret

# All namespaces
kubectl eso get secret -A
```

### `kubectl eso get secret-store [flags]`

Lists SecretStores with health status, provider type,
and age.

**Requires:** `list` on `secretstores`

**Examples:**

```bash
# List SecretStores
kubectl eso get secret-store
kubectl eso get secret-stores    # plural alias
```

### `kubectl eso get cluster-secret-store [flags]`

Lists ClusterSecretStores with health status, provider
type, and age.

**Requires:** `list` on `clustersecretstores`

**Examples:**

```bash
# List ClusterSecretStores
kubectl eso get cluster-secret-store
kubectl eso get cluster-secret-stores    # plural alias
```

### `kubectl eso describe external-secret NAME [flags]`

Shows detailed sync status including conditions, target
secret, store ref, refresh interval, and last sync time.

**Requires:** `get` on `externalsecrets`

**Examples:**

```bash
# Describe an ExternalSecret
kubectl eso describe external-secret my-es

# Output as YAML
kubectl eso describe external-secret my-es -o yaml
```

### `kubectl eso describe secret NAME [flags]`

Views Secret data. Base64-encoded by default, decoded with
`--decode`/`-d`. Shows a warning if the Secret is not
ESO-managed. Errors if the Secret is a Helm release secret.

**Requires:** `get` on `secrets`

**Flags:**

| Flag | Short | Description | Default |
| ---- | ----- | ----------- | ------- |
| `--decode` | `-d` | Decode base64 Secret data | `false` |

**Examples:**

```bash
# View Secret data (base64 encoded)
kubectl eso describe secret my-secret

# View Secret data (decoded)
kubectl eso describe secret my-secret --decode

# Output as YAML
kubectl eso describe secret my-secret -o yaml
```

### `kubectl eso describe secret-store NAME [flags]`

Shows SecretStore health details including conditions,
provider type, and last transition time.

**Requires:** `get` on `secretstores`

**Examples:**

```bash
# Describe a SecretStore
kubectl eso describe secret-store my-store
```

### `kubectl eso describe cluster-secret-store NAME [flags]`

Shows ClusterSecretStore health details including conditions,
provider type, and last transition time.

**Requires:** `get` on `clustersecretstores`

**Examples:**

```bash
# Describe a ClusterSecretStore
kubectl eso describe cluster-secret-store my-css
```

### `kubectl eso sync NAME`

Forces re-sync of an ExternalSecret by setting a
`kubectl-eso.io/force-sync` annotation with the current
unix timestamp.

**Requires:** `get` and `update` on `externalsecrets`

**Examples:**

```bash
# Force re-sync an ExternalSecret
kubectl eso sync my-es
```

### `kubectl eso completion SHELL`

Generates shell completion scripts.

**Examples:**

```bash
# bash
kubectl eso completion bash > ~/.bash_completions/kubectl-eso
source ~/.bash_completions/kubectl-eso

# zsh
kubectl eso completion zsh > "${fpath[1]}/_kubectl-eso"

# fish
kubectl eso completion fish > ~/.config/fish/completions/kubectl-eso.fish
```

For kubectl plugin completion (kubectl v1.26+), ensure
`kubectl_complete-eso` is in your `$PATH`. The `make install`
target handles this automatically.

### `kubectl eso version`

Prints version, commit, and build date.

**Examples:**

```bash
kubectl eso version
```

## Resource Aliases

| Command | Aliases |
| ------- | ------- |
| `external-secret` | `external-secrets`, `externalsecrets.external-secrets.io` |
| `secret` | `secrets` |
| `secret-store` | `secret-stores`, `secretstores.external-secrets.io` |
| `cluster-secret-store` | `cluster-secret-stores`, `clustersecretstores.external-secrets.io` |

## Helm Secret Filtering

Secrets with type `helm.sh/release.v1` are automatically excluded:

- `get secret` silently filters them from the listing
- `describe secret` returns an error

## RBAC Requirements

The following Role and ClusterRole provide the minimum
permissions needed for all plugin commands:

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: kubectl-eso
rules:
  - apiGroups: [""]
    resources: ["secrets"]
    verbs: ["get", "list"]
  - apiGroups: ["external-secrets.io"]
    resources: ["externalsecrets"]
    verbs: ["get", "list", "update"]
  - apiGroups: ["external-secrets.io"]
    resources: ["secretstores"]
    verbs: ["get", "list"]
```

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: kubectl-eso
rules:
  - apiGroups: ["external-secrets.io"]
    resources: ["clustersecretstores"]
    verbs: ["get", "list"]
```

## License

MIT
