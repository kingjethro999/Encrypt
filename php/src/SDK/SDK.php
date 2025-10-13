<?php

namespace Encrypt\SDK;

use Encrypt\Exceptions\EncryptException;
use Encrypt\Vault\Vault;
use Encrypt\Vault\VaultStatus;

/**
 * Runtime SDK for the Encrypt tool.
 * Provides easy-to-use functions for accessing secrets in PHP code.
 */
class SDK
{
    private static ?SDK $instance = null;
    private Vault $globalVault;

    private function __construct()
    {
        $this->globalVault = new Vault();
    }

    /**
     * Get the singleton SDK instance
     */
    public static function getInstance(): self
    {
        if (self::$instance === null) {
            self::$instance = new self();
        }
        return self::$instance;
    }

    /**
     * Auto-unlock vault if password is available
     */
    private function autoUnlock(?string $password = null): void
    {
        $vault = $this->globalVault;

        // If already unlocked, no need to do anything
        if ($vault->unlocked()) {
            return;
        }

        // Try provided password first
        if ($password !== null) {
            try {
                $vault->setup($password);
                return;
            } catch (EncryptException) {
                // Password might be wrong, continue to other methods
            }
        }

        // Try environment variable
        $envPassword = getenv('ENCRYPT_PASSWORD');
        if ($envPassword !== false) {
            try {
                $vault->setup($envPassword);
                return;
            } catch (EncryptException) {
                // Environment password might be wrong, continue to other methods
            }
        }

        // In development, we can be more lenient
        if (getenv('NODE_ENV') !== 'production') {
            // Try common development passwords
            $devPasswords = ['dev', 'development', 'test', 'password', '123456'];
            foreach ($devPasswords as $devPassword) {
                try {
                    $vault->setup($devPassword);
                    return;
                } catch (EncryptException) {
                    // Continue to next password
                }
            }
        }

        // If we get here, we couldn't unlock the vault
        throw new EncryptException('Vault is locked and no valid password found. Set ENCRYPT_PASSWORD environment variable or provide password parameter.');
    }

    /**
     * Get a secret value from the vault with auto-unlock support
     */
    public function get(string $key, ?string $password = null): string
    {
        $this->autoUnlock($password);
        return $this->globalVault->get($key);
    }

    /**
     * Set a secret value in the vault with auto-unlock support
     */
    public function set(string $key, string $value, ?string $password = null): void
    {
        $this->autoUnlock($password);
        $this->globalVault->set($key, $value);
    }

    /**
     * Get all secrets from the vault with auto-unlock support
     */
    public function allSecrets(?string $password = null): array
    {
        $this->autoUnlock($password);
        return $this->globalVault->all();
    }

    /**
     * Get vault status
     */
    public function status(): VaultStatus
    {
        return $this->globalVault->status();
    }

    /**
     * Check if vault is unlocked
     */
    public function isUnlocked(): bool
    {
        return $this->globalVault->unlocked();
    }

    /**
     * Auto setup helper - checks if vault is locked and provides helpful error
     */
    public function autoSetup(): void
    {
        $vault = $this->globalVault;

        if (!$vault->exists()) {
            throw new EncryptException('Vault not found. Run \'encrypt init\' first.');
        }

        if (!$vault->unlocked()) {
            throw new EncryptException('Vault is locked. Run \'encrypt setup <password>\' to unlock secrets.');
        }
    }

    /**
     * Get a secret with automatic environment variable support
     * This is the recommended function for production use
     */
    public function getSecret(string $key): string
    {
        return $this->get($key);
    }

    /**
     * Set a secret with automatic environment variable support
     * This is the recommended function for production use
     */
    public function setSecret(string $key, string $value): void
    {
        $this->set($key, $value);
    }

    /**
     * Get all secrets with automatic environment variable support
     * This is the recommended function for production use
     */
    public function getAllSecrets(): array
    {
        return $this->allSecrets();
    }
}

/**
 * Convenience functions for backward compatibility
 */

/**
 * Get a secret with automatic environment variable support
 */
function getSecret(string $key): string
{
    return SDK::getInstance()->getSecret($key);
}

/**
 * Set a secret with automatic environment variable support
 */
function setSecret(string $key, string $value): void
{
    SDK::getInstance()->setSecret($key, $value);
}

/**
 * Get all secrets with automatic environment variable support
 */
function getAllSecrets(): array
{
    return SDK::getInstance()->getAllSecrets();
}

/**
 * Get vault status
 */
function getStatus(): VaultStatus
{
    return SDK::getInstance()->status();
}

/**
 * Check if vault is unlocked
 */
function isVaultUnlocked(): bool
{
    return SDK::getInstance()->isUnlocked();
}

/**
 * Auto setup helper
 */
function autoSetupVault(): void
{
    SDK::getInstance()->autoSetup();
}
