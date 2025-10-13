import { Vault } from './vault';

// Global vault instance
let globalVault: Vault | null = null;

/**
 * Initialize the global vault instance
 */
function initVault(): Vault {
  if (!globalVault) {
    globalVault = new Vault();
  }
  return globalVault;
}

/**
 * Auto-unlock vault if password is available
 */
function autoUnlock(password?: string): void {
  const vault = initVault();
  
  // If already unlocked, no need to do anything
  if (vault.isUnlockedStatus()) {
    return;
  }
  
  // Try provided password first
  if (password) {
    try {
      vault.setup(password);
      return;
    } catch (error) {
      // Password might be wrong, continue to other methods
    }
  }
  
  // Try environment variable
  const envPassword = process.env.ENCRYPT_PASSWORD;
  if (envPassword) {
    try {
      vault.setup(envPassword);
      return;
    } catch (error) {
      // Environment password might be wrong, continue to other methods
    }
  }
  
  // In development, we can be more lenient
  if (process.env.NODE_ENV !== 'production') {
    // Try common development passwords
    const devPasswords = ['dev', 'development', 'test', 'password', '123456'];
    for (const devPassword of devPasswords) {
      try {
        vault.setup(devPassword);
        return;
      } catch (error) {
        // Continue to next password
      }
    }
  }
  
  // If we get here, we couldn't unlock the vault
  throw new Error('Vault is locked and no valid password found. Set ENCRYPT_PASSWORD environment variable or provide password parameter.');
}

/**
 * Get a secret value with auto-unlock support
 */
export function get(key: string, password?: string): string {
  autoUnlock(password);
  const vault = initVault();
  return vault.get(key);
}

/**
 * Set a secret value with auto-unlock support
 */
export function set(key: string, value: string, password?: string): void {
  autoUnlock(password);
  const vault = initVault();
  vault.set(key, value);
}

/**
 * Get all secrets with auto-unlock support
 */
export function all(password?: string): Record<string, string> {
  autoUnlock(password);
  const vault = initVault();
  return vault.all();
}

/**
 * Get vault status
 */
export function status() {
  const vault = initVault();
  return vault.status();
}

/**
 * Check if vault is unlocked
 */
export function isUnlocked(): boolean {
  const vault = initVault();
  return vault.isUnlockedStatus();
}

/**
 * Auto setup helper - checks if vault is locked and prompts for password
 */
export async function autoSetup(): Promise<void> {
  const vault = initVault();
  
  if (!vault.exists()) {
    throw new Error('Vault not found. Run "encrypt init" first.');
  }
  
  if (!vault.isUnlockedStatus()) {
    throw new Error('Vault is locked. Run "encrypt setup <password>" to unlock secrets.');
  }
}

/**
 * Get a secret with automatic environment variable support
 * This is the recommended function for production use
 */
export function getSecret(key: string): string {
  return get(key);
}

/**
 * Set a secret with automatic environment variable support
 * This is the recommended function for production use
 */
export function setSecret(key: string, value: string): void {
  set(key, value);
}

/**
 * Get all secrets with automatic environment variable support
 * This is the recommended function for production use
 */
export function getAllSecrets(): Record<string, string> {
  return all();
}

// Export the Vault class for advanced usage
export { Vault };

// Default export for convenience
export default {
  get,
  set,
  all,
  status,
  isUnlocked,
  autoSetup,
  getSecret,
  setSecret,
  getAllSecrets
};
