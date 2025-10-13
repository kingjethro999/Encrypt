//! # Encrypt - A top-level secrets orchestrator
//! 
//! Not just another .env tool — this one encrypts, locks, and sets you up for secure local and team dev.

pub mod crypto;
pub mod vault;
pub mod sdk;
pub mod cli;
pub mod error;

pub use crypto::Crypto;
pub use vault::Vault;
pub use sdk::SDK;
pub use error::{EncryptError, Result};

// Global vault instance for SDK usage
static mut GLOBAL_VAULT: Option<Vault> = None;
static INIT: std::sync::Once = std::sync::Once::new();

/// Get or create the global vault instance
pub fn global_vault() -> &'static mut Vault {
    unsafe {
        INIT.call_once(|| {
            GLOBAL_VAULT = Some(Vault::new());
        });
        GLOBAL_VAULT.as_mut().unwrap()
    }
}

/// Convenience functions for backward compatibility
pub fn get_secret(key: &str) -> Result<String> {
    SDK::get_secret(key)
}

pub fn set_secret(key: &str, value: &str) -> Result<()> {
    SDK::set_secret(key, value)
}

pub fn get_all_secrets() -> Result<std::collections::HashMap<String, String>> {
    SDK::get_all_secrets()
}

pub fn get_status() -> Result<vault::VaultStatus> {
    SDK::status()
}

pub fn is_vault_unlocked() -> Result<bool> {
    SDK::is_unlocked()
}

pub fn auto_setup_vault() -> Result<()> {
    SDK::auto_setup()
}
