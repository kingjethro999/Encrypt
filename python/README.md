# 🔐 Encrypt (Python)

> *"A top-level secrets orchestrator. Not just another .env tool — this one encrypts, locks, and sets you up for secure local and team dev."*

## 🚀 Quick Start

### Installation

```bash
# Install from source
cd python
pip install -e .

# Or install dependencies manually
pip install -r requirements.txt
```

### Basic Usage

```bash
# Initialize vault
encrypt init

# Add secrets
encrypt set API_KEY=your-api-key-here
encrypt set DB_URL=postgres://localhost:5432/mydb

# Lock secrets before committing
encrypt lockup mySuperSecurePassword

# New developer setup
encrypt setup mySuperSecurePassword
```

## 💻 In-Code Usage

### Python

```python
from encrypt import get_secret, set_secret, get_all_secrets

# Method 1: Auto-unlock with environment variable (Recommended for production)
# Set ENCRYPT_PASSWORD=your-password in your environment
api_key = get_secret('API_KEY')
db_url = get_secret('DB_URL')

# Method 2: Explicit password parameter
from encrypt import get
api_key = get('API_KEY', 'your-password')

# Method 3: Works when vault is already unlocked
api_key = get('API_KEY')

# Use in your app
config = {
    'api_key': api_key,
    'database': db_url
}

# Set secrets
set_secret('NEW_KEY', 'new_value')

# Get all secrets
all_secrets = get_all_secrets()
```

### Production Usage

```python
# Set ENCRYPT_PASSWORD environment variable
# Works automatically in any environment
api_key = get_secret('API_KEY')
db_url = get_secret('DB_URL')
```

### Alternative Import Style

```python
import encrypt

# Get secrets with auto-unlock
api_key = encrypt.get('API_KEY')
db_url = encrypt.get('DB_URL')

# Set secrets
encrypt.set('NEW_KEY', 'new_value')

# Check status
status = encrypt.status()
print(f"Vault unlocked: {encrypt.is_unlocked()}")
```

## 🧪 CLI Commands

| Command                     | Description                              |
| --------------------------- | ---------------------------------------- |
| `encrypt init`              | Create `.encrypt` vault                  |
| `encrypt lockup <password>` | Encrypt and secure secrets with password |
| `encrypt setup <password>`  | Set up secrets on a new machine          |
| `encrypt set KEY=value`     | Add/update a key                         |
| `encrypt get KEY`           | Fetch decrypted value                    |
| `encrypt unlock`            | Decrypt everything into `.env`           |
| `encrypt status`            | Check if vault is locked, list keys      |
| `encrypt reset`             | Remove vault (careful!)                  |

## 🔒 Triple Encryption Phases

1. **Phase 1: AES-256-CBC Encryption**
   Each secret value is encrypted using AES-256-CBC with a randomly generated IV.

2. **Phase 2: Password Hashing (PBKDF2)**
   The user's master password is used to derive an encryption key securely.

3. **Phase 3: HMAC Signatures**
   Encrypted secrets are signed with HMAC to prevent tampering.

## 🧾 Example Workflow

### 🔐 Initial Setup

```bash
encrypt init
```

Creates:
```
/.encrypt/
  ├── vault.lock (encrypted storage)
  ├── secrets.enc.json
  └── .gitignore (ensures raw secrets never get committed)
```

### 🔒 Lock Secrets Before Commit

```bash
encrypt lockup mySuperSecurePassword
```

This:
- Encrypts all secret values in `.encrypt/secrets.enc.json`
- Stores an encrypted hash of your password
- Prevents accidental push of plaintext secrets

### 👤 New Developer Setup

```bash
git clone your-repo
cd your-repo
encrypt setup mySuperSecurePassword
```

This:
- Prompts for password
- Decrypts secrets into memory
- Your app works 🎉

## 🔧 Development

```bash
# Install in development mode
pip install -e .

# Run tests
python -m pytest tests/

# Run demo
python demo.py
```

## 📁 Project Structure

```
encrypt/
├── __init__.py      # Package initialization and SDK exports
├── cli.py           # CLI interface using Click and Rich
├── crypto.py        # Triple-layer encryption implementation
├── vault.py         # Vault management and file operations
└── sdk.py           # Runtime SDK for in-code usage

.encrypt/            # Encrypted vault directory (created at runtime)
├── vault.lock       # Vault configuration and password hash
├── secrets.enc.json # Encrypted secrets storage
└── vault.unlocked   # Lock file indicating vault status
```

## 🧪 Testing

```bash
# Test encryption/decryption
python test_crypto.py

# Run complete demo
python demo.py

# Test runtime SDK
python example.py
```

## 🛡️ Security Features

- **Triple-layer encryption** for maximum security
- **Password-based key derivation** using PBKDF2
- **HMAC signatures** to prevent tampering
- **Memory-only decryption** (secrets never written to disk when unlocked)
- **Git-safe** (only encrypted files are committed)

## 🚀 Production Deployment

### Environment Variable Method (Recommended)

Set the `ENCRYPT_PASSWORD` environment variable in your production environment:

```bash
# Docker
ENV ENCRYPT_PASSWORD=your-production-password

# Kubernetes
env:
- name: ENCRYPT_PASSWORD
  valueFrom:
    secretKeyRef:
      name: encrypt-secrets
      key: password

# Heroku
heroku config:set ENCRYPT_PASSWORD=your-password

# AWS Lambda
# Set ENCRYPT_PASSWORD in environment variables
```

### Your Application Code

```python
# Works automatically with ENCRYPT_PASSWORD environment variable
from encrypt import get_secret

api_key = get_secret('API_KEY')
db_url = get_secret('DB_URL')

# No need to manually unlock the vault!
```

### Django Example

```python
# settings.py
import os
from encrypt import get_secret

SECRET_KEY = get_secret('SECRET_KEY')
DATABASE_URL = get_secret('DATABASE_URL')
DEBUG = os.getenv('DEBUG', 'False').lower() == 'true'
```

### Security Benefits

- ✅ **Secrets remain encrypted** in `.encrypt/` folder
- ✅ **Only decrypted in memory** during runtime
- ✅ **No plaintext secrets** ever written to disk
- ✅ **Environment-specific passwords** for dev/staging/prod
- ✅ **Zero configuration** required in your app code

## 🚀 Why Encrypt?

- `.env` files are static and hard to share securely
- GitHub secrets don't help in local development
- Vault tools like HashiCorp are overkill for small projects
- You want an easy way to **lock your dev secrets before pushing** and **onboard teammates easily**

This tool solves that problem in a slick, dev-friendly way.

## 📄 License

MIT
