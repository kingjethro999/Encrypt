use thiserror::Error;

#[derive(Error, Debug)]
pub enum EncryptError {
    #[error("Vault not found. Run 'encrypt init' first.")]
    VaultNotFound,

    #[error("Vault is locked. Run 'encrypt setup <password>' to unlock secrets.")]
    VaultLocked,

    #[error("Invalid password.")]
    InvalidPassword,

    #[error("Secret \"{0}\" not found.")]
    SecretNotFound(String),

    #[error("Failed to decrypt secret: {0}")]
    DecryptionFailed(String),

    #[error("Vault is locked and no valid password found. Set ENCRYPT_PASSWORD environment variable or provide password parameter.")]
    NoValidPassword,

    #[error("No secrets to lock. Use 'encrypt set' to add secrets first.")]
    NoSecretsToLock,

    #[error("IO error: {0}")]
    Io(#[from] std::io::Error),

    #[error("JSON error: {0}")]
    Json(#[from] serde_json::Error),

    #[error("Crypto error: {0}")]
    Crypto(String),

    #[error("Base64 error: {0}")]
    Base64(#[from] base64::DecodeError),

    #[error("Hex error: {0}")]
    Hex(#[from] hex::FromHexError),
}

pub type Result<T> = std::result::Result<T, EncryptError>;
