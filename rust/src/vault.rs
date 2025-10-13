use crate::crypto::{Crypto, EncryptionResult};
use crate::error::{EncryptError, Result};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::fs;
use std::path::PathBuf;
use chrono::{DateTime, Utc};

const VAULT_DIR: &str = ".encrypt";
const CONFIG_FILE: &str = "vault.lock";
const SECRETS_FILE: &str = "secrets.enc.json";
const GITIGNORE_FILE: &str = ".gitignore";
const LOCK_FILE: &str = "vault.unlocked";

#[derive(Debug, Serialize, Deserialize)]
pub struct VaultConfig {
    pub password_hash: String,
    pub salt: String,
    pub hmac: String,
    pub created_at: DateTime<Utc>,
    pub version: String,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct VaultStatus {
    pub is_locked: bool,
    pub keys: Vec<String>,
    pub last_modified: Option<DateTime<Utc>>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct LockFile {
    pub unlocked: bool,
    pub timestamp: DateTime<Utc>,
}

/// Vault management for the Encrypt tool.
/// Handles file operations, state management, and secret storage/retrieval.
pub struct Vault {
    project_root: PathBuf,
    vault_path: PathBuf,
    config_path: PathBuf,
    secrets_path: PathBuf,
    gitignore_path: PathBuf,
    lock_file_path: PathBuf,
    memory_cache: HashMap<String, String>,
}

impl Vault {
    pub fn new() -> Self {
        let project_root = std::env::current_dir().unwrap_or_else(|_| PathBuf::from("."));
        let vault_path = project_root.join(VAULT_DIR);
        let config_path = vault_path.join(CONFIG_FILE);
        let secrets_path = vault_path.join(SECRETS_FILE);
        let gitignore_path = project_root.join(GITIGNORE_FILE);
        let lock_file_path = vault_path.join(LOCK_FILE);

        Self {
            project_root,
            vault_path,
            config_path,
            secrets_path,
            gitignore_path,
            lock_file_path,
            memory_cache: HashMap::new(),
        }
    }

    /// Initialize the vault directory structure
    pub fn init(&self) -> Result<()> {
        fs::create_dir_all(&self.vault_path)?;

        // Create initial empty secrets file
        if !self.secrets_path.exists() {
            fs::write(&self.secrets_path, "{}")?;
        }

        // Create initial config file
        if !self.config_path.exists() {
            let initial_config = VaultConfig {
                password_hash: String::new(),
                salt: String::new(),
                hmac: String::new(),
                created_at: Utc::now(),
                version: "1.0.0".to_string(),
            };
            fs::write(&self.config_path, serde_json::to_string_pretty(&initial_config)?)?;
        }

        // Update .gitignore
        self.update_gitignore()?;

        Ok(())
    }

    /// Update .gitignore to exclude .encrypt directory
    fn update_gitignore(&self) -> Result<()> {
        let gitignore_content = if self.gitignore_path.exists() {
            fs::read_to_string(&self.gitignore_path)?
        } else {
            String::new()
        };

        if !gitignore_content.contains(".encrypt/") {
            let mut file = fs::OpenOptions::new()
                .create(true)
                .append(true)
                .open(&self.gitignore_path)?;
            use std::io::Write;
            writeln!(file, "\n# Encrypt vault\n.encrypt/")?;
        }

        Ok(())
    }

    /// Check if vault exists
    pub fn exists(&self) -> bool {
        self.vault_path.exists() && self.config_path.exists()
    }

    /// Lock up secrets with password
    pub fn lockup(&mut self, password: &str) -> Result<()> {
        // Load secrets from file if not in memory
        if self.memory_cache.is_empty() {
            self.load_secrets_from_file()?;
        }

        if self.memory_cache.is_empty() {
            return Err(EncryptError::NoSecretsToLock);
        }

        // Encrypt all secrets
        let mut encrypted_secrets = HashMap::new();

        for (key, value) in &self.memory_cache {
            let result = Crypto::encrypt(value, password)?;
            encrypted_secrets.insert(key.clone(), serde_json::to_string(&result)?);
        }

        // Create vault config
        let config_salt = Crypto::generate_salt();
        let config = VaultConfig {
            password_hash: Crypto::hash_password(password)?,
            salt: config_salt.clone(),
            hmac: Crypto::generate_hmac(
                &serde_json::to_string(&encrypted_secrets)?,
                &Crypto::derive_key(password, &config_salt)?,
            ),
            created_at: Utc::now(),
            version: "1.0.0".to_string(),
        };

        // Write encrypted secrets and config
        fs::write(&self.secrets_path, serde_json::to_string_pretty(&encrypted_secrets)?)?;
        fs::write(&self.config_path, serde_json::to_string_pretty(&config)?)?;

        // Clear memory cache and remove lock file
        self.memory_cache.clear();
        if self.lock_file_path.exists() {
            fs::remove_file(&self.lock_file_path)?;
        }

        Ok(())
    }

    /// Setup/unlock vault with password
    pub fn setup(&mut self, password: &str) -> Result<()> {
        if !self.exists() {
            return Err(EncryptError::VaultNotFound);
        }

        let config: VaultConfig = serde_json::from_str(&fs::read_to_string(&self.config_path)?)?;

        // If this is a fresh vault (no password set), just unlock it
        if config.password_hash.is_empty() {
            self.memory_cache.clear();
            self.create_lock_file()?;
            return Ok(());
        }

        // Verify password
        if !Crypto::verify_password(password, &config.password_hash)? {
            return Err(EncryptError::InvalidPassword);
        }

        // Load and decrypt secrets
        let encrypted_secrets: HashMap<String, String> =
            serde_json::from_str(&fs::read_to_string(&self.secrets_path)?)?;

        self.memory_cache.clear();

        for (key, encrypted_data) in encrypted_secrets {
            // Check if the data is already decrypted (plain text) or encrypted
            if encrypted_data.starts_with('{') {
                // This is encrypted data stored as JSON string, decrypt it
                let result: EncryptionResult = serde_json::from_str(&encrypted_data)?;
                let decrypted = Crypto::decrypt(&result.encrypted, &result.salt, &result.hmac, password)?;
                self.memory_cache.insert(key, decrypted);
            } else {
                // This is plain text data
                self.memory_cache.insert(key, encrypted_data);
            }
        }

        // Save decrypted secrets to file for easy access
        self.save_secrets_to_file()?;
        self.create_lock_file()?;

        Ok(())
    }

    /// Set a secret (only works when unlocked)
    pub fn set(&mut self, key: &str, value: &str) -> Result<()> {
        if !self.unlocked()? {
            return Err(EncryptError::VaultLocked);
        }

        // Load secrets from file if not in memory
        if self.memory_cache.is_empty() {
            self.load_secrets_from_file()?;
        }

        self.memory_cache.insert(key.to_string(), value.to_string());
        self.save_secrets_to_file()?;

        Ok(())
    }

    /// Get a secret (only works when unlocked)
    pub fn get(&mut self, key: &str) -> Result<String> {
        if !self.unlocked()? {
            return Err(EncryptError::VaultLocked);
        }

        // Load secrets from file if not in memory
        if self.memory_cache.is_empty() {
            self.load_secrets_from_file()?;
        }

        self.memory_cache
            .get(key)
            .ok_or_else(|| EncryptError::SecretNotFound(key.to_string()))
            .map(|s| s.clone())
    }

    /// Get all secrets (only works when unlocked)
    pub fn all(&mut self) -> Result<HashMap<String, String>> {
        if !self.unlocked()? {
            return Err(EncryptError::VaultLocked);
        }

        // Load secrets from file if not in memory
        if self.memory_cache.is_empty() {
            self.load_secrets_from_file()?;
        }

        Ok(self.memory_cache.clone())
    }

    /// Get vault status
    pub fn status(&mut self) -> Result<VaultStatus> {
        // Check if vault is unlocked by looking for lock file
        let is_unlocked = self.unlocked()?;

        // Load secrets from file if unlocked and not in memory
        if is_unlocked && self.memory_cache.is_empty() {
            self.load_secrets_from_file()?;
        }

        let keys: Vec<String> = self.memory_cache.keys().cloned().collect();
        let last_modified = if self.exists() {
            Some(DateTime::from(fs::metadata(&self.config_path)?.modified()?))
        } else {
            None
        };

        Ok(VaultStatus {
            is_locked: !is_unlocked,
            keys,
            last_modified,
        })
    }

    /// Reset/remove vault
    pub fn reset(&self) -> Result<()> {
        if self.vault_path.exists() {
            fs::remove_dir_all(&self.vault_path)?;
        }
        Ok(())
    }

    /// Check if vault is unlocked
    pub fn unlocked(&self) -> Result<bool> {
        Ok(self.lock_file_path.exists())
    }

    /// Create lock file to indicate vault is unlocked
    fn create_lock_file(&self) -> Result<()> {
        let lock_data = LockFile {
            unlocked: true,
            timestamp: Utc::now(),
        };
        fs::write(&self.lock_file_path, serde_json::to_string_pretty(&lock_data)?)?;
        Ok(())
    }

    /// Load secrets from file (for unlocked vault)
    fn load_secrets_from_file(&mut self) -> Result<()> {
        if !self.secrets_path.exists() {
            return Ok(());
        }

        let secrets: HashMap<String, String> =
            serde_json::from_str(&fs::read_to_string(&self.secrets_path)?)?;

        self.memory_cache.clear();
        for (key, value) in secrets {
            // Check if the value is encrypted (starts with {) or plain text
            if value.starts_with('{') {
                // This is encrypted data, we need to decrypt it
                // But we don't have the password here, so we can't decrypt
                // This should not happen in an unlocked vault
                return Err(EncryptError::Crypto(format!(
                    "Secret \"{}\" is encrypted but vault is unlocked. This should not happen.",
                    key
                )));
            } else {
                // This is plain text data
                self.memory_cache.insert(key, value);
            }
        }

        Ok(())
    }

    /// Save secrets to file (for unlocked vault)
    fn save_secrets_to_file(&self) -> Result<()> {
        fs::write(&self.secrets_path, serde_json::to_string_pretty(&self.memory_cache)?)?;
        Ok(())
    }
}
