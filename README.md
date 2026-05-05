# envdiff

> Tool to diff and reconcile `.env` files across environments with redaction support.

---

## Installation

```bash
go install github.com/yourusername/envdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envdiff.git && cd envdiff && go build ./...
```

---

## Usage

Compare two `.env` files and highlight differences:

```bash
envdiff .env.development .env.production
```

Reconcile a local env file against a reference, with sensitive values redacted:

```bash
envdiff --redact .env.local .env.example
```

Output missing keys only:

```bash
envdiff --missing .env .env.example
```

### Example Output

```
~ DB_HOST        dev.db.local  →  prod.db.example.com
+ SENTRY_DSN     [redacted]
- DEBUG          true
```

**Flags**

| Flag | Description |
|------|-------------|
| `--redact` | Mask secret values in output |
| `--missing` | Show only keys missing from the target file |
| `--json` | Output diff as JSON |

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

[MIT](LICENSE) © 2024 yourusername