# 🔐 Encrypt (Go)

> *"A top-level secrets orchestrator. Not just another .env tool — this one encrypts, locks, and sets you up for secure local and team dev."*

## 🚀 Quick Start

### Installation

```bash
# Build from source
cd go
go build -o encrypt main.go

# Or install dependencies and run
go mod tidy
go run main.go init
```

### Basic Usage

```bash
# Initialize vault
./encrypt init

# Add secrets
./encrypt set API_KEY=your-api-key-here
./encrypt set DB_URL=postgres://localhost:5432/mydb

# Lock secrets before committing
./encrypt lockup mySuperSecurePassword

# New developer setup
./encrypt setup mySuperSecurePassword
```

## 💻 In-Code Usage

### Go

```go
package main

import (
    "fmt"
    "log"
    
    "encrypt/internal/sdk"
)

func main() {
    // Method 1: Auto-unlock with environment variable (Recommended for production)
    // Set ENCRYPT_PASSWORD=your-password in your environment
    apiKey, err := sdk.GetSecret("API_KEY")
    if err != nil {
        log.Fatal(err)
    }
    
    dbURL, err := sdk.GetSecret("DB_URL")
    if err != nil {
        log.Fatal(err)
    }
    
    // Method 2: Explicit password parameter
    apiKey, err = sdk.GetSDK().Get("API_KEY", "your-password")
    if err != nil {
        log.Fatal(err)
    }
    
    // Method 3: Works when vault is already unlocked
    apiKey, err = sdk.GetSDK().Get("API_KEY", "")
    if err != nil {
        log.Fatal(err)
    }
    
    // Use in your app
    config := Config{
        APIKey:   apiKey,
        Database: dbURL,
        Port:     "3000",
    }
    
    // Set secrets
    err = sdk.SetSecret("NEW_KEY", "new_value")
    if err != nil {
        log.Fatal(err)
    }
    
    // Get all secrets
    allSecrets, err := sdk.GetAllSecrets()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Available keys: %v\n", allSecrets)
}

type Config struct {
    APIKey   string
    Database string
    Port     string
}
```

### Production Usage

```go
// Set ENCRYPT_PASSWORD environment variable
// Works automatically in any environment
import "encrypt/internal/sdk"

apiKey, err := sdk.GetSecret("API_KEY")
if err != nil {
    log.Fatal(err)
}

dbURL, err := sdk.GetSecret("DB_URL")
if err != nil {
    log.Fatal(err)
}
```

### Gin Web Framework Example

```go
package main

import (
    "net/http"
    
    "encrypt/internal/sdk"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    
    r.GET("/", func(c *gin.Context) {
        apiKey, err := sdk.GetSecret("API_KEY")
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        
        c.JSON(http.StatusOK, gin.H{
            "message": "Hello World",
            "api_key": apiKey[:4] + "***", // Mask the key
        })
    })
    
    r.Run(":8080")
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
# Build in development mode
go build -o encrypt main.go

# Run tests
go test ./...

# Run examples
go run examples/test_crypto.go
go run examples/example.go
go run examples/test_auto_converter.go
go run examples/production_example.go

# Run demo
./encrypt init
./encrypt setup mypassword
./encrypt set API_KEY=sk-1234567890abcdef
./encrypt status
```

## 📁 Project Structure

```
internal/
├── cli/
│   └── cli.go           # CLI interface using cobra
├── crypto/
│   └── crypto.go        # Triple-layer encryption implementation
├── vault/
│   └── vault.go         # Vault management and file operations
└── sdk/
    └── sdk.go           # Runtime SDK for in-code usage

examples/
├── test_crypto.go           # Crypto tests
├── example.go               # Basic usage example
├── test_auto_converter.go   # Auto-converter tests
└── production_example.go    # Production usage example

main.go              # CLI entry point
go.mod               # Go module definition

.encrypt/                # Encrypted vault directory (created at runtime)
├── vault.lock           # Vault configuration and password hash
├── secrets.enc.json     # Encrypted secrets storage
└── vault.unlocked       # Lock file indicating vault status
```

## 🧪 Testing

```bash
# Test encryption/decryption
go run examples/test_crypto.go

# Test runtime SDK
go run examples/example.go

# Test auto-converter
go run examples/test_auto_converter.go

# Test production usage
go run examples/production_example.go
```

## 🛡️ Security Features

- **Triple-layer encryption** for maximum security
- **Password-based key derivation** using PBKDF2
- **HMAC signatures** to prevent tampering
- **Memory-only decryption** (secrets never written to disk when unlocked)
- **Git-safe** (only encrypted files are committed)
- **Memory-safe** with Go's garbage collection
- **Type-safe** with Go's static typing

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

```go
// Works automatically with ENCRYPT_PASSWORD environment variable
import "encrypt/internal/sdk"

apiKey, err := sdk.GetSecret("API_KEY")
if err != nil {
    log.Fatal(err)
}

dbURL, err := sdk.GetSecret("DB_URL")
if err != nil {
    log.Fatal(err)
}

// No need to manually unlock the vault!
```

### Echo Web Framework Example

```go
package main

import (
    "net/http"
    
    "encrypt/internal/sdk"
    "github.com/labstack/echo/v4"
)

func main() {
    e := echo.New()
    
    e.GET("/", func(c echo.Context) error {
        apiKey, err := sdk.GetSecret("API_KEY")
        if err != nil {
            return c.JSON(http.StatusInternalServerError, map[string]string{
                "error": err.Error(),
            })
        }
        
        return c.JSON(http.StatusOK, map[string]string{
            "message": "Hello World",
            "api_key": apiKey[:4] + "***", // Mask the key
        })
    })
    
    e.Logger.Fatal(e.Start(":8080"))
}
```

### Security Benefits

- ✅ **Secrets remain encrypted** in `.encrypt/` folder
- ✅ **Only decrypted in memory** during runtime
- ✅ **No plaintext secrets** ever written to disk
- ✅ **Environment-specific passwords** for dev/staging/prod
- ✅ **Zero configuration** required in your app code
- ✅ **Memory-safe** with Go's garbage collection
- ✅ **Type-safe** with Go's static typing

## 🚀 Why Encrypt?

- `.env` files are static and hard to share securely
- GitHub secrets don't help in local development
- Vault tools like HashiCorp are overkill for small projects
- You want an easy way to **lock your dev secrets before pushing** and **onboard teammates easily**

This tool solves that problem in a slick, dev-friendly way with Go's performance and safety guarantees.

## 📄 License

MIT
