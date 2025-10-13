# 🔐 Encrypt (Rust)

> *"A top-level secrets orchestrator. Not just another .env tool — this one encrypts, locks, and sets you up for secure local and team dev."*

## 🚀 Quick Start

### Installation

```bash
# Build from source
cd rust
cargo build --release

# Install globally
cargo install --path .

# Or run directly
cargo run -- init
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

### Rust

```rust
use encrypt::{get_secret, set_secret, get_all_secrets};

// Method 1: Auto-unlock with environment variable (Recommended for production)
// Set ENCRYPT_PASSWORD=your-password in your environment
let api_key = get_secret("API_KEY")?;
let db_url = get_secret("DB_URL")?;

// Method 2: Explicit password parameter
use encrypt::sdk::SDK;
let api_key = SDK::get("API_KEY", Some("your-password"))?;

// Method 3: Works when vault is already unlocked
let api_key = SDK::get("API_KEY", None)?;

// Use in your app
let config = Config {
    api_key,
    database: db_url,
    port: std::env::var("PORT").unwrap_or_else(|_| "3000".to_string()),
};

// Set secrets
set_secret("NEW_KEY", "new_value")?;

// Get all secrets
let all_secrets = get_all_secrets()?;
```

### Production Usage

```rust
// Set ENCRYPT_PASSWORD environment variable
// Works automatically in any environment
use encrypt::get_secret;

let api_key = get_secret("API_KEY")?;
let db_url = get_secret("DB_URL")?;
```

### Actix Web Example

```rust
use actix_web::{web, App, HttpServer, Result};
use encrypt::get_secret;

async fn index() -> Result<&'static str> {
    let api_key = get_secret("API_KEY")?;
    // Use api_key in your handler
    Ok("Hello world!")
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    HttpServer::new(|| {
        App::new()
            .route("/", web::get().to(index))
    })
    .bind("127.0.0.1:8080")?
    .run()
    .await
}
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
   Each secret value is encrypted using AES-256-GCM with a randomly generated nonce.

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
# Build in development mode
cargo build

# Run tests
cargo test

# Run examples
cargo run --example test_crypto
cargo run --example example
cargo run --example test_auto_converter
cargo run --example production_example

# Run demo
cargo run -- init
cargo run -- setup mypassword
cargo run -- set API_KEY=sk-1234567890abcdef
cargo run -- status
```

## 📁 Project Structure

```
src/
├── lib.rs           # Main library and convenience functions
├── main.rs          # CLI entry point
├── cli.rs           # CLI interface using clap
├── crypto.rs        # Triple-layer encryption implementation
├── vault.rs         # Vault management and file operations
├── sdk.rs           # Runtime SDK for in-code usage
└── error.rs         # Error types and handling

examples/
├── test_crypto.rs           # Crypto tests
├── example.rs               # Basic usage example
├── test_auto_converter.rs   # Auto-converter tests
└── production_example.rs    # Production usage example

.encrypt/                # Encrypted vault directory (created at runtime)
├── vault.lock           # Vault configuration and password hash
├── secrets.enc.json     # Encrypted secrets storage
└── vault.unlocked       # Lock file indicating vault status
```

## 🧪 Testing

```bash
# Test encryption/decryption
cargo run --example test_crypto

# Test runtime SDK
cargo run --example example

# Test auto-converter
cargo run --example test_auto_converter

# Test production usage
cargo run --example production_example
```

## 🛡️ Security Features

- **Triple-layer encryption** for maximum security
- **Password-based key derivation** using PBKDF2
- **HMAC signatures** to prevent tampering
- **Memory-only decryption** (secrets never written to disk when unlocked)
- **Git-safe** (only encrypted files are committed)
- **Zero-copy operations** where possible for performance

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

```rust
// Works automatically with ENCRYPT_PASSWORD environment variable
use encrypt::get_secret;

let api_key = get_secret("API_KEY")?;
let db_url = get_secret("DB_URL")?;

// No need to manually unlock the vault!
```

### Rocket Example

```rust
use rocket::{get, launch, routes, State};
use encrypt::get_secret;

#[get("/")]
fn index() -> &'static str {
    let api_key = get_secret("API_KEY").unwrap();
    // Use api_key in your handler
    "Hello, world!"
}

#[launch]
fn rocket() -> _ {
    rocket::build().mount("/", routes![index])
}
```

### Security Benefits

- ✅ **Secrets remain encrypted** in `.encrypt/` folder
- ✅ **Only decrypted in memory** during runtime
- ✅ **No plaintext secrets** ever written to disk
- ✅ **Environment-specific passwords** for dev/staging/prod
- ✅ **Zero configuration** required in your app code
- ✅ **Memory-safe** with Rust's ownership system

## 🚀 Why Encrypt?

- `.env` files are static and hard to share securely
- GitHub secrets don't help in local development
- Vault tools like HashiCorp are overkill for small projects
- You want an easy way to **lock your dev secrets before pushing** and **onboard teammates easily**

This tool solves that problem in a slick, dev-friendly way with Rust's performance and safety guarantees.

## 📄 License

MIT
