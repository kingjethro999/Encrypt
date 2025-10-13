use crate::error::{EncryptError, Result};
use crate::vault::{Vault, VaultStatus};
use std::collections::HashMap;
use std::env;

/// Runtime SDK for the Encrypt tool.
/// Provides easy-to-use functions for accessing secrets in Rust code.
pub struct SDK;

impl SDK {
    /// Get or create the global vault instance
    fn global_vault() -> &'static mut Vault {
        crate::global_vault()
    }

    /// Auto-unlock vault if password is available
    fn auto_unlock(password: Option<&str>) -> Result<()> {
        let vault = Self::global_vault();

        // If already unlocked, no need to do anything
        if vault.unlocked()? {
            return Ok(());
        }

        // Try provided password first
        if let Some(pwd) = password {
            if vault.setup(pwd).is_ok() {
                return Ok(());
            }
        }

        // Try environment variable
        if let Ok(env_password) = env::var("ENCRYPT_PASSWORD") {
            if vault.setup(&env_password).is_ok() {
                return Ok(());
            }
        }

        // In development, we can be more lenient
        if env::var("NODE_ENV").unwrap_or_default() != "production" {
            // Try common development passwords
            let dev_passwords = ["dev", "development", "test", "password", "123456"];
            for dev_password in &dev_passwords {
                if vault.setup(dev_password).is_ok() {
                    return Ok(());
                }
            }
        }

        // If we get here, we couldn't unlock the vault
        Err(EncryptError::NoValidPassword)
    }

    /// Get a secret value from the vault with auto-unlock support
    pub fn get(key: &str, password: Option<&str>) -> Result<String> {
        Self::auto_unlock(password)?;
        let vault = Self::global_vault();
        vault.get(key)
    }

    /// Set a secret value in the vault with auto-unlock support
    pub fn set(key: &str, value: &str, password: Option<&str>) -> Result<()> {
        Self::auto_unlock(password)?;
        let vault = Self::global_vault();
        vault.set(key, value)
    }

    /// Get all secrets from the vault with auto-unlock support
    pub fn all_secrets(password: Option<&str>) -> Result<HashMap<String, String>> {
        Self::auto_unlock(password)?;
        let vault = Self::global_vault();
        vault.all()
    }

    /// Get vault status
    pub fn status() -> Result<VaultStatus> {
        let vault = Self::global_vault();
        vault.status()
    }

    /// Check if vault is unlocked
    pub fn is_unlocked() -> Result<bool> {
        let vault = Self::global_vault();
        vault.unlocked()
    }

    /// Auto setup helper - checks if vault is locked and provides helpful error
    pub fn auto_setup() -> Result<()> {
        let vault = Self::global_vault();

        if !vault.exists() {
            return Err(EncryptError::VaultNotFound);
        }

        if !vault.unlocked()? {
            return Err(EncryptError::VaultLocked);
        }

        Ok(())
    }

    /// Get a secret with automatic environment variable support
    /// This is the recommended function for production use
    pub fn get_secret(key: &str) -> Result<String> {
        Self::get(key, None)
    }

    /// Set a secret with automatic environment variable support
    /// This is the recommended function for production use
    pub fn set_secret(key: &str, value: &str) -> Result<()> {
        Self::set(key, value, None)
    }

    /// Get all secrets with automatic environment variable support
    /// This is the recommended function for production use
    pub fn get_all_secrets() -> Result<HashMap<String, String>> {
        Self::all_secrets(None)
    }
}
