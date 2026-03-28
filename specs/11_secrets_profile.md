# Spec 11 — Secrets & Profile

## Overview

The profile service provides persistent storage for app metadata, properties, and secrets. The secrets service builds on top of it to manage encrypted credentials and SSH key pairs.

## Profile Service

**Key:** `"profile"`

### Profile Structure

Stored as JSON at `{runDir}/etc/profile.json`:

```go
type Profile struct {
    Name       string            `json:"name"`
    Version    string            `json:"version"`
    Properties map[string]string `json:"properties"`
    Secrets    map[string]string `json:"secrets"`
}
```

### Methods

| Method                    | Description                        |
| ------------------------- | ---------------------------------- |
| `GetProfile()`            | Load full profile from file        |
| `GetPublicProfile()`      | Return profile without secrets     |
| `HasProperty(key)`        | Check if property exists           |
| `GetProperty(key)`        | Get property value                 |
| `SetProperty(key, value)` | Write property                     |
| `HasSecret(key)`          | Check if secret exists             |
| `GetSecret(key)`          | Get secret value                   |
| `SetSecret(key, value)`   | Write secret                       |
| `GetUserDirectory()`      | Create and return `{runDir}/home/` |

### Init Behavior

On init, reads or creates the profile file, updating name and version from the kernel.

## Secrets Service

**Key:** `"secrets"`

### Methods

| Method                  | Description                         |
| ----------------------- | ----------------------------------- |
| `GetSecret(key)`        | Retrieve and base64-decode a secret |
| `SetSecret(key, value)` | Base64-encode and store a secret    |

### SSH Key Generation

On first init, generates an ECDSA P-256 key pair:

```go
type KeyPair struct {
    PublicKey  []byte `json:"publicKey"`   // DER-encoded PKIX
    PrivateKey []byte `json:"privateKey"`  // DER-encoded ECPrivateKey
}
```

Keys are stored as base64-encoded strings under:

- `"ssh.private"` — private key
- `"ssh.public"` — public key

## Data Flow

```
Secrets Service
    └── Profile Service
        └── profile.json (filesystem)
```
