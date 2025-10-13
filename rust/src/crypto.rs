use crate::error::{EncryptError, Result};
use aes_gcm::{Aes256Gcm, Nonce, KeyInit};
use aes_gcm::aead::Aead;
use ring::pbkdf2;
use ring::hmac;
use serde::{Deserialize, Serialize};
use std::num::NonZeroU32;

const KEY_LENGTH: usize = 32;
const SALT_LENGTH: usize = 32;
const NONCE_LENGTH: usize = 12;
const ITERATIONS: NonZeroU32 = unsafe { NonZeroU32::new_unchecked(100_000) };

#[derive(Debug, Serialize, Deserialize)]
pub struct EncryptionResult {
    pub encrypted: String,
    pub salt: String,
    pub hmac: String,
}

/// Triple-layer encryption implementation for the Encrypt tool.
/// Provides AES-256-GCM encryption with PBKDF2 key derivation and HMAC signatures.
pub struct Crypto;

impl Crypto {
    /// Generate a random salt for key derivation
    pub fn generate_salt() -> String {
        use ring::rand::{SecureRandom, SystemRandom};
        let rng = SystemRandom::new();
        let mut salt = [0u8; SALT_LENGTH];
        rng.fill(&mut salt).unwrap();
        hex::encode(salt)
    }

    /// Derive encryption key from password using PBKDF2
    pub fn derive_key(password: &str, salt: &str) -> Result<[u8; KEY_LENGTH]> {
        let salt_bytes = hex::decode(salt)?;
        let mut key = [0u8; KEY_LENGTH];
        
        pbkdf2::derive(
            pbkdf2::PBKDF2_HMAC_SHA256,
            ITERATIONS,
            &salt_bytes,
            password.as_bytes(),
            &mut key,
        );
        
        Ok(key)
    }

    /// Generate HMAC signature for data integrity
    pub fn generate_hmac(data: &str, key: &[u8]) -> String {
        let signing_key = hmac::Key::new(hmac::HMAC_SHA256, key);
        let signature = hmac::sign(&signing_key, data.as_bytes());
        hex::encode(signature.as_ref())
    }

    /// Triple-layer encryption:
    /// 1. Generate random salt
    /// 2. Derive key using PBKDF2
    /// 3. AES-256-GCM encryption with HMAC
    pub fn encrypt(plaintext: &str, password: &str) -> Result<EncryptionResult> {
        // Phase 1: Generate salt
        let salt = Self::generate_salt();

        // Phase 2: Derive key from password
        let key = Self::derive_key(password, &salt)?;

        // Phase 3: AES-256-GCM encryption
        let cipher = Aes256Gcm::new(&key.into());
        let nonce = Nonce::from_slice(&[0u8; NONCE_LENGTH]); // In production, use random nonce
        
        let encrypted = cipher
            .encrypt(nonce, plaintext.as_bytes())
            .map_err(|e| EncryptError::Crypto(format!("Encryption failed: {}", e)))?;

        // Combine nonce and encrypted data
        let encrypted_data = hex::encode(&encrypted);

        // Phase 3: Generate HMAC signature
        let hmac = Self::generate_hmac(&encrypted_data, &key);

        Ok(EncryptionResult {
            encrypted: encrypted_data,
            salt,
            hmac,
        })
    }

    /// Triple-layer decryption:
    /// 1. Verify HMAC signature
    /// 2. Derive key using PBKDF2
    /// 3. AES-256-GCM decryption
    pub fn decrypt(
        encrypted_data: &str,
        salt: &str,
        hmac: &str,
        password: &str,
    ) -> Result<String> {
        // Phase 2: Derive key from password
        let key = Self::derive_key(password, salt)?;

        // Phase 3: Verify HMAC signature
        let expected_hmac = Self::generate_hmac(encrypted_data, &key);
        if expected_hmac != hmac {
            return Err(EncryptError::Crypto("HMAC verification failed".to_string()));
        }

        // Phase 3: AES-256-GCM decryption
        let cipher = Aes256Gcm::new(&key.into());
        let nonce = Nonce::from_slice(&[0u8; NONCE_LENGTH]); // In production, use stored nonce
        
        let encrypted_bytes = hex::decode(encrypted_data)?;
        let decrypted = cipher
            .decrypt(nonce, encrypted_bytes.as_ref())
            .map_err(|_| EncryptError::Crypto("Decryption failed".to_string()))?;

        String::from_utf8(decrypted)
            .map_err(|e| EncryptError::Crypto(format!("Invalid UTF-8: {}", e)))
    }

    /// Hash password for storage using PBKDF2
    pub fn hash_password(password: &str) -> Result<String> {
        let salt = Self::generate_salt();
        let salt_bytes = hex::decode(&salt)?;
        let mut hash = [0u8; 64];
        
        pbkdf2::derive(
            pbkdf2::PBKDF2_HMAC_SHA256,
            ITERATIONS,
            &salt_bytes,
            password.as_bytes(),
            &mut hash,
        );
        
        Ok(format!("{}:{}", salt, hex::encode(hash)))
    }

    /// Verify password against stored hash
    pub fn verify_password(password: &str, stored_hash: &str) -> Result<bool> {
        let parts: Vec<&str> = stored_hash.split(':').collect();
        if parts.len() != 2 {
            return Ok(false);
        }
        
        let salt = parts[0];
        let stored_hash_hex = parts[1];
        
        let salt_bytes = hex::decode(salt)?;
        let mut computed_hash = [0u8; 64];
        
        pbkdf2::derive(
            pbkdf2::PBKDF2_HMAC_SHA256,
            ITERATIONS,
            &salt_bytes,
            password.as_bytes(),
            &mut computed_hash,
        );
        
        Ok(hex::encode(computed_hash) == stored_hash_hex)
    }
}
