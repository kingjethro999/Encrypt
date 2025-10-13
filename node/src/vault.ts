import * as fs from 'fs';
import * as path from 'path';
import { TripleEncryption } from './crypto';
import { VaultConfig, EncryptedSecrets, VaultStatus, EncryptionResult } from './types';

export class Vault {
  private static readonly VAULT_DIR = '.encrypt';
  private static readonly CONFIG_FILE = 'vault.lock';
  private static readonly SECRETS_FILE = 'secrets.enc.json';
  private static readonly GITIGNORE_FILE = '.gitignore';
  private static readonly LOCK_FILE = 'vault.unlocked';
  
  private vaultPath: string;
  private configPath: string;
  private secretsPath: string;
  private gitignorePath: string;
  private lockFilePath: string;
  private memoryCache: Map<string, string> = new Map();
  private isUnlocked: boolean = false;

  constructor(projectRoot: string = process.cwd()) {
    this.vaultPath = path.join(projectRoot, Vault.VAULT_DIR);
    this.configPath = path.join(this.vaultPath, Vault.CONFIG_FILE);
    this.secretsPath = path.join(this.vaultPath, Vault.SECRETS_FILE);
    this.gitignorePath = path.join(projectRoot, Vault.GITIGNORE_FILE);
    this.lockFilePath = path.join(this.vaultPath, Vault.LOCK_FILE);
  }

  /**
   * Initialize the vault directory structure
   */
  init(): void {
    if (!fs.existsSync(this.vaultPath)) {
      fs.mkdirSync(this.vaultPath, { recursive: true });
    }

    // Create initial empty secrets file
    if (!fs.existsSync(this.secretsPath)) {
      fs.writeFileSync(this.secretsPath, JSON.stringify({}, null, 2));
    }

    // Create initial config file
    if (!fs.existsSync(this.configPath)) {
      const initialConfig: VaultConfig = {
        password_hash: '',
        salt: '',
        hmac: '',
        created_at: new Date().toISOString(),
        version: '1.0.0'
      };
      fs.writeFileSync(this.configPath, JSON.stringify(initialConfig, null, 2));
    }

    // Create .gitignore to prevent accidental commits
    this.updateGitignore();
  }

  /**
   * Update .gitignore to exclude .encrypt directory
   */
  private updateGitignore(): void {
    const gitignoreContent = fs.existsSync(this.gitignorePath) 
      ? fs.readFileSync(this.gitignorePath, 'utf8') 
      : '';

    if (!gitignoreContent.includes('.encrypt/')) {
      const newContent = gitignoreContent + '\n# Encrypt vault\n.encrypt/\n';
      fs.writeFileSync(this.gitignorePath, newContent);
    }
  }

  /**
   * Check if vault exists
   */
  exists(): boolean {
    return fs.existsSync(this.vaultPath) && fs.existsSync(this.configPath);
  }

  /**
   * Lock up secrets with password
   */
  lockup(password: string): void {
    // Load secrets from file if not in memory
    if (this.memoryCache.size === 0) {
      this.loadSecretsFromFile();
    }
    
    if (this.memoryCache.size === 0) {
      throw new Error('No secrets to lock. Use "encrypt set" to add secrets first.');
    }

    // Encrypt all secrets
    const encryptedSecrets: EncryptedSecrets = {};
    
    for (const [key, value] of this.memoryCache.entries()) {
      const result: EncryptionResult = TripleEncryption.encrypt(value, password);
      encryptedSecrets[key] = JSON.stringify(result);
    }

    // Create vault config
    const configSalt = TripleEncryption.generateSalt();
    const config: VaultConfig = {
      password_hash: TripleEncryption.hashPassword(password),
      salt: configSalt,
      hmac: TripleEncryption.generateHMAC(JSON.stringify(encryptedSecrets), TripleEncryption.deriveKey(password, configSalt)),
      created_at: new Date().toISOString(),
      version: '1.0.0'
    };

    // Write encrypted secrets and config
    fs.writeFileSync(this.secretsPath, JSON.stringify(encryptedSecrets, null, 2));
    fs.writeFileSync(this.configPath, JSON.stringify(config, null, 2));

    // Clear memory cache
    this.memoryCache.clear();
    this.isUnlocked = false;
    // Remove lock file to indicate vault is locked
    if (fs.existsSync(this.lockFilePath)) {
      fs.unlinkSync(this.lockFilePath);
    }
  }

  /**
   * Setup/unlock vault with password
   */
  setup(password: string): void {
    if (!this.exists()) {
      throw new Error('Vault not found. Run "encrypt init" first.');
    }

    const config: VaultConfig = JSON.parse(fs.readFileSync(this.configPath, 'utf8'));
    
    // If this is a fresh vault (no password set), just unlock it
    if (!config.password_hash) {
      this.memoryCache.clear();
      this.isUnlocked = true;
      // Create lock file to indicate vault is unlocked
      fs.writeFileSync(this.lockFilePath, JSON.stringify({ unlocked: true, timestamp: new Date().toISOString() }));
      return;
    }
    
    // Verify password
    if (!TripleEncryption.verifyPassword(password, config.password_hash)) {
      throw new Error('Invalid password.');
    }

    // Load and decrypt secrets
    const encryptedSecrets: EncryptedSecrets = JSON.parse(fs.readFileSync(this.secretsPath, 'utf8'));
    
    this.memoryCache.clear();
    
    for (const [key, encryptedData] of Object.entries(encryptedSecrets)) {
      // Check if the data is already decrypted (plain text) or encrypted
      if (typeof encryptedData === 'string' && encryptedData.startsWith('{')) {
        // This is encrypted data stored as JSON string, decrypt it
        const result = JSON.parse(encryptedData) as EncryptionResult;
        const decrypted = TripleEncryption.decrypt(result.encrypted, result.salt, result.hmac, password);
        
        if (!decrypted.isValid) {
          throw new Error(`Failed to decrypt secret: ${key}`);
        }
        
        this.memoryCache.set(key, decrypted.decrypted);
      } else {
        // This is plain text data
        this.memoryCache.set(key, encryptedData as string);
      }
    }

    this.isUnlocked = true;
    // Save decrypted secrets to file for easy access
    this.saveSecretsToFile();
    // Create lock file to indicate vault is unlocked
    fs.writeFileSync(this.lockFilePath, JSON.stringify({ unlocked: true, timestamp: new Date().toISOString() }));
  }

  /**
   * Set a secret (only works when unlocked)
   */
  set(key: string, value: string): void {
    if (!this.isUnlockedStatus()) {
      throw new Error('Vault is locked. Run "encrypt setup <password>" to unlock secrets.');
    }
    
    // Load secrets from file if not in memory
    if (this.memoryCache.size === 0) {
      this.loadSecretsFromFile();
    }
    
    this.memoryCache.set(key, value);
    this.saveSecretsToFile();
  }

  /**
   * Get a secret (only works when unlocked)
   */
  get(key: string): string {
    if (!this.isUnlockedStatus()) {
      throw new Error('Vault is locked. Run "encrypt setup <password>" to unlock secrets.');
    }
    
    // Load secrets from file if not in memory
    if (this.memoryCache.size === 0) {
      this.loadSecretsFromFile();
    }
    
    if (!this.memoryCache.has(key)) {
      throw new Error(`Secret "${key}" not found.`);
    }
    
    return this.memoryCache.get(key)!;
  }

  /**
   * Get all secrets (only works when unlocked)
   */
  all(): Record<string, string> {
    if (!this.isUnlockedStatus()) {
      throw new Error('Vault is locked. Run "encrypt setup <password>" to unlock secrets.');
    }
    
    // Load secrets from file if not in memory
    if (this.memoryCache.size === 0) {
      this.loadSecretsFromFile();
    }
    
    const result: Record<string, string> = {};
    for (const [key, value] of this.memoryCache.entries()) {
      result[key] = value;
    }
    return result;
  }

  /**
   * Get vault status
   */
  status(): VaultStatus {
    // Check if vault is unlocked by looking for lock file
    const isUnlocked = fs.existsSync(this.lockFilePath);
    
    // Load secrets from file if unlocked and not in memory
    if (isUnlocked && this.memoryCache.size === 0) {
      this.loadSecretsFromFile();
    }
    
    const keys = Array.from(this.memoryCache.keys());
    const lastModified = this.exists() ? fs.statSync(this.configPath).mtime.toISOString() : undefined;
    
    return {
      isLocked: !isUnlocked,
      keys,
      lastModified
    };
  }

  /**
   * Reset/remove vault
   */
  reset(): void {
    if (fs.existsSync(this.vaultPath)) {
      fs.rmSync(this.vaultPath, { recursive: true, force: true });
    }
    this.memoryCache.clear();
    this.isUnlocked = false;
  }

  /**
   * Check if vault is unlocked
   */
  isUnlockedStatus(): boolean {
    return fs.existsSync(this.lockFilePath);
  }

  /**
   * Load secrets from file (for unlocked vault)
   */
  private loadSecretsFromFile(): void {
    if (fs.existsSync(this.secretsPath)) {
      const secretsData = fs.readFileSync(this.secretsPath, 'utf8');
      const secrets = JSON.parse(secretsData);
      
      this.memoryCache.clear();
      for (const [key, value] of Object.entries(secrets)) {
        // Check if the value is encrypted (starts with {) or plain text
        if (typeof value === 'string' && value.startsWith('{')) {
          // This is encrypted data, we need to decrypt it
          // But we don't have the password here, so we can't decrypt
          // This should not happen in an unlocked vault
          throw new Error(`Secret "${key}" is encrypted but vault is unlocked. This should not happen.`);
        } else {
          // This is plain text data
          this.memoryCache.set(key, value as string);
        }
      }
    }
  }

  /**
   * Save secrets to file (for unlocked vault)
   */
  private saveSecretsToFile(): void {
    const secrets: Record<string, string> = {};
    for (const [key, value] of this.memoryCache.entries()) {
      secrets[key] = value;
    }
    fs.writeFileSync(this.secretsPath, JSON.stringify(secrets, null, 2));
  }
}
