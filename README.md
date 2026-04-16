# vaultline

A lightweight CLI for syncing secrets from Vault to local `.env` files safely.

---

## Installation

```bash
go install github.com/yourusername/vaultline@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/vaultline.git
cd vaultline
go build -o vaultline .
```

---

## Usage

Configure your Vault address and token, then sync secrets to a local `.env` file:

```bash
# Set required environment variables
export VAULT_ADDR="https://vault.example.com"
export VAULT_TOKEN="s.your-vault-token"

# Sync secrets from a Vault path to a .env file
vaultline sync --path secret/myapp/prod --out .env

# Preview secrets without writing to disk
vaultline sync --path secret/myapp/prod --dry-run
```

Your `.env` file will be populated with the key-value pairs stored at the given Vault path:

```
DB_HOST=prod-db.internal
DB_PASSWORD=supersecret
API_KEY=abc123
```

> **Note:** vaultline will never overwrite existing keys unless the `--force` flag is provided.

---

## Configuration

| Flag | Description | Default |
|------|-------------|---------|
| `--path` | Vault secret path | _(required)_ |
| `--out` | Output `.env` file path | `.env` |
| `--dry-run` | Print secrets without writing | `false` |
| `--force` | Overwrite existing keys | `false` |

---

## License

MIT © [yourusername](https://github.com/yourusername)