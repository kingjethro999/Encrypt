# 🔐 Encrypt

> *"A top-level secrets orchestrator. Not just another .env tool — this one encrypts, locks, and sets you up for secure local and team dev."*

## 🚀 Quick Start

### Installation

```bash
# Clone the repository
git clone <your-repo-url>
cd best-encrypt

# Install dependencies
npm install

# Build the project
npm run build

# Make CLI globally available (optional)
npm link
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

### Demo

Run the complete demo to see all features:

```bash
node demo.js
```

## 🧠 Core Concept

Encrypt replaces `.env` files with an **encrypted local secrets vault** that:

- ✅ **Encrypts secrets** with triple-layer encryption (AES-256 + PBKDF2 + HMAC)
- ✅ **Locks secrets** before committing to Git
- ✅ **Onboards teammates** with a single command
- ✅ **Runtime access** in your code without environment variables
- ✅ **Cross-platform** support (Node.js, Python, Rust coming soon)

## 💻 In-Code Usage

### JavaScript/TypeScript

```js
import encrypt from 'encrypt'

// Method 1: Auto-unlock with environment variable (Recommended for production)
// Set ENCRYPT_PASSWORD=your-password in your environment
const openaiKey = encrypt.getSecret('OPENAI_API_KEY')
const dbUrl = encrypt.getSecret('DB_URL')

// Method 2: Explicit password parameter
const openaiKey = encrypt.get('OPENAI_API_KEY', 'your-password')

// Method 3: Works when vault is already unlocked
const openaiKey = encrypt.get('OPENAI_API_KEY')

// Use in your app
const config = {
  apiKey: openaiKey,
  database: dbUrl
}
```

### Production Usage

```js
// Set ENCRYPT_PASSWORD environment variable
// Works automatically in any environment
const apiKey = encrypt.getSecret('API_KEY')
const dbUrl = encrypt.getSecret('DB_URL')
```

### Python

```python
from encrypt import get_secret

# Method 1: Auto-unlock with environment variable (Recommended for production)
# Set ENCRYPT_PASSWORD=your-password in your environment
openai_key = get_secret("OPENAI_API_KEY")
db_url = get_secret("DB_URL")

# Method 2: Explicit password parameter
from encrypt import get
openai_key = get("OPENAI_API_KEY", "your-password")
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

1. **Phase 1: AES-256-GCM Encryption**
   Each secret value is encrypted using AES-256-GCM with a randomly generated IV.

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
  ├── .vault (encrypted storage)
  ├── config.json
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
# Install dependencies
npm install

# Build
npm run build

# Run in development
npm run dev

# Test
npm test
```

## 📁 Project Structure

```
src/
├── cli.ts          # CLI interface using Commander.js
├── crypto.ts       # Triple-layer encryption implementation
├── vault.ts        # Vault management and file operations
├── index.ts        # Runtime SDK for in-code usage
└── types.ts        # TypeScript type definitions

dist/               # Compiled JavaScript output
.encrypt/           # Encrypted vault directory (created at runtime)
├── vault.lock      # Vault configuration and password hash
├── secrets.enc.json # Encrypted secrets storage
└── vault.unlocked  # Lock file indicating vault status
```

## 🧪 Testing

```bash
# Test encryption/decryption
node test.js

# Run complete demo
node demo.js

# Test runtime SDK
node example.js
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

```js
// Works automatically with ENCRYPT_PASSWORD environment variable
const apiKey = encrypt.getSecret('API_KEY')
const dbUrl = encrypt.getSecret('DB_URL')

// No need to manually unlock the vault!
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
