export interface VaultConfig {
  password_hash: string;
  salt: string;
  hmac: string;
  created_at: string;
  version: string;
}

export interface EncryptedSecrets {
  [key: string]: string;
}

export interface VaultStatus {
  isLocked: boolean;
  keys: string[];
  lastModified?: string;
}

export interface EncryptionResult {
  encrypted: string;
  salt: string;
  hmac: string;
}

export interface DecryptionResult {
  decrypted: string;
  isValid: boolean;
}
