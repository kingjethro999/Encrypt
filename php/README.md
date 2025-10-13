# 🔐 Encrypt (PHP)

> *"A top-level secrets orchestrator. Not just another .env tool — this one encrypts, locks, and sets you up for secure local and team dev."*

## 🚀 Quick Start

### Installation

```bash
# Install dependencies
cd php
composer install

# Make executable
chmod +x bin/encrypt

# Run commands
./bin/encrypt init
```

### Basic Usage

```bash
# Initialize vault
./bin/encrypt init

# Add secrets
./bin/encrypt set API_KEY=your-api-key-here
./bin/encrypt set DB_URL=postgres://localhost:5432/mydb

# Lock secrets before committing
./bin/encrypt lockup mySuperSecurePassword

# New developer setup
./bin/encrypt setup mySuperSecurePassword
```

## 💻 In-Code Usage

### PHP

```php
<?php

require_once 'vendor/autoload.php';

use Encrypt\SDK\SDK;

// Method 1: Auto-unlock with environment variable (Recommended for production)
// Set ENCRYPT_PASSWORD=your-password in your environment
$apiKey = SDK::getInstance()->getSecret('API_KEY');
$dbUrl = SDK::getInstance()->getSecret('DB_URL');

// Method 2: Explicit password parameter
$apiKey = SDK::getInstance()->get('API_KEY', 'your-password');

// Method 3: Works when vault is already unlocked
$apiKey = SDK::getInstance()->get('API_KEY', '');

// Use in your app
$config = [
    'api_key' => $apiKey,
    'database' => $dbUrl,
    'port' => $_ENV['PORT'] ?? '3000'
];

// Set secrets
SDK::getInstance()->setSecret('NEW_KEY', 'new_value');

// Get all secrets
$allSecrets = SDK::getInstance()->getAllSecrets();
```

### Production Usage

```php
<?php

// Set ENCRYPT_PASSWORD environment variable
// Works automatically in any environment
use Encrypt\SDK\SDK;

$apiKey = SDK::getInstance()->getSecret('API_KEY');
$dbUrl = SDK::getInstance()->getSecret('DB_URL');
```

### Laravel Example

```php
<?php

// config/app.php
use Encrypt\SDK\SDK;

return [
    'api_key' => SDK::getInstance()->getSecret('API_KEY'),
    'database_url' => SDK::getInstance()->getSecret('DATABASE_URL'),
    // ... other config
];
```

### Symfony Example

```php
<?php

// config/services.yaml
use Encrypt\SDK\SDK;

parameters:
    app.api_key: '%env(ENCRYPT_PASSWORD)%'
    
services:
    app.api_client:
        class: App\Service\ApiClient
        arguments:
            $apiKey: !php/const Encrypt\SDK\SDK::getInstance()->getSecret('API_KEY')
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
./bin/encrypt init
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
./bin/encrypt lockup mySuperSecurePassword
```

This:
- Encrypts all secret values in `.encrypt/secrets.enc.json`
- Stores an encrypted hash of your password
- Prevents accidental push of plaintext secrets

### 👤 New Developer Setup

```bash
git clone your-repo
cd your-repo
./bin/encrypt setup mySuperSecurePassword
```

This:
- Prompts for password
- Decrypts secrets into memory
- Your app works 🎉

## 🔧 Development

```bash
# Install dependencies
composer install

# Run tests
composer test

# Run examples
php examples/test_crypto.php
php examples/example.php
php examples/test_auto_converter.php
php examples/production_example.php

# Run demo
./bin/encrypt init
./bin/encrypt setup mypassword
./bin/encrypt set API_KEY=sk-1234567890abcdef
./bin/encrypt status
```

## 📁 Project Structure

```
src/
├── Application.php           # Main application class
├── Commands/
│   ├── BaseCommand.php       # Base command class
│   ├── InitCommand.php       # Init command
│   ├── LockupCommand.php     # Lockup command
│   ├── SetupCommand.php      # Setup command
│   ├── SetCommand.php        # Set command
│   ├── GetCommand.php        # Get command
│   ├── UnlockCommand.php     # Unlock command
│   ├── StatusCommand.php     # Status command
│   └── ResetCommand.php      # Reset command
├── Crypto/
│   ├── Crypto.php            # Triple-layer encryption implementation
│   └── EncryptionResult.php  # Encryption result class
├── Vault/
│   ├── Vault.php             # Vault management and file operations
│   └── VaultStatus.php       # Vault status class
├── SDK/
│   └── SDK.php               # Runtime SDK for in-code usage
└── Exceptions/
    └── EncryptException.php  # Base exception class

examples/
├── test_crypto.php           # Crypto tests
├── example.php               # Basic usage example
├── test_auto_converter.php   # Auto-converter tests
└── production_example.php    # Production usage example

bin/
└── encrypt                   # CLI executable

.encrypt/                # Encrypted vault directory (created at runtime)
├── vault.lock           # Vault configuration and password hash
├── secrets.enc.json     # Encrypted secrets storage
└── vault.unlocked       # Lock file indicating vault status
```

## 🧪 Testing

```bash
# Test encryption/decryption
php examples/test_crypto.php

# Test runtime SDK
php examples/example.php

# Test auto-converter
php examples/test_auto_converter.php

# Test production usage
php examples/production_example.php
```

## 🛡️ Security Features

- **Triple-layer encryption** for maximum security
- **Password-based key derivation** using PBKDF2
- **HMAC signatures** to prevent tampering
- **Memory-only decryption** (secrets never written to disk when unlocked)
- **Git-safe** (only encrypted files are committed)
- **Object-oriented design** with proper encapsulation
- **Type safety** with PHP 8+ features

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

```php
<?php
// Works automatically with ENCRYPT_PASSWORD environment variable
use Encrypt\SDK\SDK;

$apiKey = SDK::getInstance()->getSecret('API_KEY');
$dbUrl = SDK::getInstance()->getSecret('DB_URL');

// No need to manually unlock the vault!
```

### Laravel Service Provider Example

```php
<?php

namespace App\Providers;

use Encrypt\SDK\SDK;
use Illuminate\Support\ServiceProvider;

class EncryptServiceProvider extends ServiceProvider
{
    public function register(): void
    {
        $this->app->singleton('encrypt', function () {
            return SDK::getInstance();
        });
    }

    public function boot(): void
    {
        // Set secrets in config
        config([
            'services.api.key' => SDK::getInstance()->getSecret('API_KEY'),
            'database.connections.mysql.password' => SDK::getInstance()->getSecret('DB_PASSWORD'),
        ]);
    }
}
```

### Security Benefits

- ✅ **Secrets remain encrypted** in `.encrypt/` folder
- ✅ **Only decrypted in memory** during runtime
- ✅ **No plaintext secrets** ever written to disk
- ✅ **Environment-specific passwords** for dev/staging/prod
- ✅ **Zero configuration** required in your app code
- ✅ **Object-oriented design** with proper encapsulation
- ✅ **Type safety** with PHP 8+ features

## 🚀 Why Encrypt?

- `.env` files are static and hard to share securely
- GitHub secrets don't help in local development
- Vault tools like HashiCorp are overkill for small projects
- You want an easy way to **lock your dev secrets before pushing** and **onboard teammates easily**

This tool solves that problem in a slick, dev-friendly way with PHP's object-oriented design and modern features.

## 📄 License

MIT
