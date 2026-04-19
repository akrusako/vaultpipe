# vaultpipe

> Stream secrets from Vault into process environments without writing to disk.

---

## Installation

```bash
go install github.com/yourname/vaultpipe@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourname/vaultpipe/releases).

---

## Usage

```bash
vaultpipe [flags] -- <command>
```

**Example:** Run a server with database credentials injected from Vault:

```bash
vaultpipe \
  --addr https://vault.example.com \
  --path secret/data/myapp/db \
  -- ./server
```

Secrets at the specified path are resolved and injected as environment variables into the child process. Nothing is written to disk.

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--addr` | Vault server address | `$VAULT_ADDR` |
| `--token` | Vault token | `$VAULT_TOKEN` |
| `--path` | Secret path to read | *(required)* |
| `--prefix` | Prefix for injected env vars | *(none)* |

**Authentication** falls back to standard Vault environment variables (`VAULT_ADDR`, `VAULT_TOKEN`, `VAULT_ROLE_ID`, etc.) if flags are not provided.

---

## How It Works

1. Authenticates with Vault using the provided credentials
2. Fetches secrets from the specified path
3. Merges secrets into the current environment
4. Executes the given command via `execve` — no subprocess, no temp files

---

## License

MIT © [yourname](https://github.com/yourname)