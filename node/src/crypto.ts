import * as crypto from 'crypto';
import { EncryptionResult, DecryptionResult } from './types';

export class TripleEncryption {
  private static readonly ALGORITHM = 'aes-256-cbc';
  private static readonly KEY_LENGTH = 32;
  private static readonly IV_LENGTH = 16;
  private static readonly SALT_LENGTH = 32;
  private static readonly ITERATIONS = 100000;

  /**
   * Phase 1: Generate a random salt
   */
  static generateSalt(): string {
    return crypto.randomBytes(this.SALT_LENGTH).toString('hex');
  }

  /**
   * Phase 2: Derive key from password using PBKDF2
   */
  static deriveKey(password: string, salt: string): Buffer {
    return crypto.pbkdf2Sync(password, salt, this.ITERATIONS, this.KEY_LENGTH, 'sha512');
  }

  /**
   * Phase 3: Generate HMAC signature
   */
  static generateHMAC(data: string, key: Buffer): string {
    return crypto.createHmac('sha256', key).update(data).digest('hex');
  }

  /**
   * Triple-layer encryption
   */
  static encrypt(plaintext: string, password: string): EncryptionResult {
    // Phase 1: Generate salt
    const salt = this.generateSalt();
    
    // Phase 2: Derive key from password
    const key = this.deriveKey(password, salt);
    
    // Phase 3: AES-256-CBC encryption
    const iv = crypto.randomBytes(this.IV_LENGTH);
    const cipher = crypto.createCipheriv(this.ALGORITHM, key, iv);
    
    let encrypted = cipher.update(plaintext, 'utf8', 'hex');
    encrypted += cipher.final('hex');
    
    // Combine encrypted data with IV
    const encryptedData = iv.toString('hex') + ':' + encrypted;
    
    // Phase 3: Generate HMAC signature
    const hmac = this.generateHMAC(encryptedData, key);
    
    return {
      encrypted: encryptedData,
      salt,
      hmac
    };
  }

  /**
   * Triple-layer decryption
   */
  static decrypt(encryptedData: string, salt: string, hmac: string, password: string): DecryptionResult {
    try {
      // Phase 2: Derive key from password
      const key = this.deriveKey(password, salt);
      
      // Phase 3: Verify HMAC signature
      const expectedHMAC = this.generateHMAC(encryptedData, key);
      if (expectedHMAC !== hmac) {
        return { decrypted: '', isValid: false };
      }
      
      // Phase 3: AES-256-CBC decryption
      const [ivHex, encrypted] = encryptedData.split(':');
      const iv = Buffer.from(ivHex, 'hex');
      
      const decipher = crypto.createDecipheriv(this.ALGORITHM, key, iv);
      
      let decrypted = decipher.update(encrypted, 'hex', 'utf8');
      decrypted += decipher.final('utf8');
      
      return { decrypted, isValid: true };
    } catch (error) {
      return { decrypted: '', isValid: false };
    }
  }

  /**
   * Hash password for storage
   */
  static hashPassword(password: string): string {
    const salt = this.generateSalt();
    const hash = crypto.pbkdf2Sync(password, salt, this.ITERATIONS, 64, 'sha512');
    return salt + ':' + hash.toString('hex');
  }

  /**
   * Verify password against hash
   */
  static verifyPassword(password: string, hash: string): boolean {
    const [salt, storedHash] = hash.split(':');
    const computedHash = crypto.pbkdf2Sync(password, salt, this.ITERATIONS, 64, 'sha512');
    return computedHash.toString('hex') === storedHash;
  }
}
