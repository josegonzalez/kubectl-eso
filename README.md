# kubectl-eso

A kubectl plugin for managing Kubernetes Secrets in the context of the [External Secrets Operator](https://external-secrets.io).

## Features

- Annotate existing Secrets for ESO adoption (`creationPolicy: Merge`)
- List and inspect ExternalSecrets, SecretStores, and ClusterSecretStores
- View Secret data with optional base64 decoding
- Force re-sync ExternalSecrets
- Helm release secrets are automatically filtered out
- Full kubectl flag compatibility (`--namespace`, `--context`, `--kubeconfig`, etc.)
- Shell completion for bash, zsh, fish, and powershell

## Installation

### From source

```bash
go install github.com/josegonzalez/kubectl-eso/cmd/kubectl-eso@latest
```

### From release

Download the binary from the [releases page](https://github.com/josegonzalez/kubectl-eso/releases) and place it in your `$PATH`.

### With make

```bash
git clone https://github.com/josegonzalez/kubectl-eso.git
cd kubectl-eso
make install
```

This installs both `kubectl-eso` and the `kubectl_complete-eso` completion wrapper to `$GOPATH/bin`.

## Usage

```bash
# Annotate a secret for ESO adoption
kubectl eso annotate my-secret --store my-store

# Annotate with dry-run (outputs YAML)
kubectl eso annotate my-secret --store my-store --dry-run

# List ExternalSecrets
kubectl eso get external-secret
kubectl eso get es                    # short alias
kubectl eso get external-secrets -A   # all namespaces

# List Secrets with ESO-managed indicator
kubectl eso get secret

# List SecretStores
kubectl eso get secret-store
kubectl eso get ss                    # short alias

# List ClusterSecretStores
kubectl eso get cluster-secret-store
kubectl eso get css                   # short alias

# Describe an ExternalSecret
kubectl eso describe external-secret my-es

# View Secret data (base64 encoded)
kubectl eso describe secret my-secret

# View Secret data (decoded)
kubectl eso describe secret my-secret --decode

# Describe a SecretStore
kubectl eso describe secret-store my-store

# Describe a ClusterSecretStore
kubectl eso describe cluster-secret-store my-css

# Force re-sync an ExternalSecret
kubectl eso sync my-es

# Output as JSON or YAML
kubectl eso get external-secret -o json
kubectl eso describe secret my-secret -o yaml

# Print version
kubectl eso version
```

## Shell Completion

Generate completion scripts:

```bash
# bash
kubectl eso completion bash > ~/.bash_completions/kubectl-eso
source ~/.bash_completions/kubectl-eso

# zsh
kubectl eso completion zsh > "${fpath[1]}/_kubectl-eso"

# fish
kubectl eso completion fish > ~/.config/fish/completions/kubectl-eso.fish
```

For kubectl plugin completion (kubectl v1.26+), ensure `kubectl_complete-eso` is in your `$PATH`. The `make install` target handles this automatically.

## Resource Aliases

| Command | Aliases |
|---------|---------|
| `external-secret` | `external-secrets`, `es`, `ExternalSecret`, `ExternalSecrets` |
| `secret` | `secrets`, `Secret`, `Secrets` |
| `secret-store` | `secret-stores`, `ss`, `SecretStore`, `SecretStores` |
| `cluster-secret-store` | `cluster-secret-stores`, `css`, `ClusterSecretStore`, `ClusterSecretStores` |

## License

MIT
