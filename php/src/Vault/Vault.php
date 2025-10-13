<?php

namespace Encrypt\Vault;

use Encrypt\Crypto\Crypto;
use Encrypt\Crypto\EncryptionResult;
use Encrypt\Exceptions\EncryptException;

/**
 * Vault management for the Encrypt tool.
 * Handles file operations, state management, and secret storage/retrieval.
 */
class Vault
{
    private const VAULT_DIR = '.encrypt';
    private const CONFIG_FILE = 'vault.lock';
    private const SECRETS_FILE = 'secrets.enc.json';
    private const GITIGNORE_FILE = '.gitignore';
    private const LOCK_FILE = 'vault.unlocked';

    private string $projectRoot;
    private string $vaultPath;
    private string $configPath;
    private string $secretsPath;
    private string $gitignorePath;
    private string $lockFilePath;
    private array $memoryCache = [];
    private Crypto $crypto;

    public function __construct(?string $projectRoot = null)
    {
        $this->projectRoot = $projectRoot ?? getcwd();
        $this->vaultPath = $this->projectRoot . DIRECTORY_SEPARATOR . self::VAULT_DIR;
        $this->configPath = $this->vaultPath . DIRECTORY_SEPARATOR . self::CONFIG_FILE;
        $this->secretsPath = $this->vaultPath . DIRECTORY_SEPARATOR . self::SECRETS_FILE;
        $this->gitignorePath = $this->projectRoot . DIRECTORY_SEPARATOR . self::GITIGNORE_FILE;
        $this->lockFilePath = $this->vaultPath . DIRECTORY_SEPARATOR . self::LOCK_FILE;
        $this->crypto = new Crypto();
    }

    /**
     * Initialize the vault directory structure
     */
    public function init(): void
    {
        if (!is_dir($this->vaultPath)) {
            if (!mkdir($this->vaultPath, 0755, true)) {
                throw new EncryptException('Failed to create vault directory');
            }
        }

        // Create initial empty secrets file
        if (!file_exists($this->secretsPath)) {
            if (file_put_contents($this->secretsPath, '{}') === false) {
                throw new EncryptException('Failed to create secrets file');
            }
        }

        // Create initial config file
        if (!file_exists($this->configPath)) {
            $initialConfig = [
                'password_hash' => '',
                'salt' => '',
                'hmac' => '',
                'created_at' => date('c'),
                'version' => '1.0.0'
            ];

            if (file_put_contents($this->configPath, json_encode($initialConfig, JSON_PRETTY_PRINT)) === false) {
                throw new EncryptException('Failed to create config file');
            }
        }

        // Update .gitignore
        $this->updateGitignore();
    }

    /**
     * Update .gitignore to exclude .encrypt directory
     */
    private function updateGitignore(): void
    {
        $gitignoreContent = file_exists($this->gitignorePath) ? file_get_contents($this->gitignorePath) : '';

        if (!str_contains($gitignoreContent, '.encrypt/')) {
            $file = fopen($this->gitignorePath, 'a');
            if ($file === false) {
                throw new EncryptException('Failed to open .gitignore file');
            }

            fwrite($file, "\n# Encrypt vault\n.encrypt/\n");
            fclose($file);
        }
    }

    /**
     * Check if vault exists
     */
    public function exists(): bool
    {
        return is_dir($this->vaultPath) && file_exists($this->configPath);
    }

    /**
     * Lock up secrets with password
     */
    public function lockup(string $password): void
    {
        // Load secrets from file if not in memory
        if (empty($this->memoryCache)) {
            $this->loadSecretsFromFile();
        }

        if (empty($this->memoryCache)) {
            throw new EncryptException('No secrets to lock. Use \'encrypt set\' to add secrets first.');
        }

        // Encrypt all secrets
        $encryptedSecrets = [];

        foreach ($this->memoryCache as $key => $value) {
            $result = $this->crypto->encrypt($value, $password);
            $encryptedSecrets[$key] = json_encode($result->toArray());
        }

        // Create vault config
        $configSalt = $this->crypto->generateSalt();
        $config = [
            'password_hash' => $this->crypto->hashPassword($password),
            'salt' => $configSalt,
            'hmac' => $this->crypto->generateHMAC(
                json_encode($encryptedSecrets),
                $this->crypto->deriveKey($password, $configSalt)
            ),
            'created_at' => date('c'),
            'version' => '1.0.0'
        ];

        // Write encrypted secrets and config
        if (file_put_contents($this->secretsPath, json_encode($encryptedSecrets, JSON_PRETTY_PRINT)) === false) {
            throw new EncryptException('Failed to write secrets file');
        }

        if (file_put_contents($this->configPath, json_encode($config, JSON_PRETTY_PRINT)) === false) {
            throw new EncryptException('Failed to write config file');
        }

        // Clear memory cache and remove lock file
        $this->memoryCache = [];
        if (file_exists($this->lockFilePath)) {
            unlink($this->lockFilePath);
        }
    }

    /**
     * Setup/unlock vault with password
     */
    public function setup(string $password): void
    {
        if (!$this->exists()) {
            throw new EncryptException('Vault not found. Run \'encrypt init\' first.');
        }

        $configData = file_get_contents($this->configPath);
        if ($configData === false) {
            throw new EncryptException('Failed to read config file');
        }

        $config = json_decode($configData, true);
        if ($config === null) {
            throw new EncryptException('Invalid config file');
        }

        // If this is a fresh vault (no password set), just unlock it
        if (empty($config['password_hash'])) {
            $this->memoryCache = [];
            $this->createLockFile();
            return;
        }

        // Verify password
        if (!$this->crypto->verifyPassword($password, $config['password_hash'])) {
            throw new EncryptException('Invalid password.');
        }

        // Load and decrypt secrets
        $secretsData = file_get_contents($this->secretsPath);
        if ($secretsData === false) {
            throw new EncryptException('Failed to read secrets file');
        }

        $encryptedSecrets = json_decode($secretsData, true);
        if ($encryptedSecrets === null) {
            throw new EncryptException('Invalid secrets file');
        }

        $this->memoryCache = [];

        foreach ($encryptedSecrets as $key => $encryptedData) {
            // Check if the data is already decrypted (plain text) or encrypted
            if (str_starts_with($encryptedData, '{')) {
                // This is encrypted data stored as JSON string, decrypt it
                $resultData = json_decode($encryptedData, true);
                if ($resultData === null) {
                    throw new EncryptException("Failed to parse encryption result for key: $key");
                }

                $result = EncryptionResult::fromArray($resultData);
                [$decrypted, $isValid] = $this->crypto->decrypt($result->encrypted, $result->salt, $result->hmac, $password);

                if (!$isValid) {
                    throw new EncryptException("Failed to decrypt secret: $key");
                }

                $this->memoryCache[$key] = $decrypted;
            } else {
                // This is plain text data
                $this->memoryCache[$key] = $encryptedData;
            }
        }

        // Save decrypted secrets to file for easy access
        $this->saveSecretsToFile();
        $this->createLockFile();
    }

    /**
     * Set a secret (only works when unlocked)
     */
    public function set(string $key, string $value): void
    {
        if (!$this->unlocked()) {
            throw new EncryptException('Vault is locked. Run \'encrypt setup <password>\' to unlock secrets.');
        }

        // Load secrets from file if not in memory
        if (empty($this->memoryCache)) {
            $this->loadSecretsFromFile();
        }

        $this->memoryCache[$key] = $value;
        $this->saveSecretsToFile();
    }

    /**
     * Get a secret (only works when unlocked)
     */
    public function get(string $key): string
    {
        if (!$this->unlocked()) {
            throw new EncryptException('Vault is locked. Run \'encrypt setup <password>\' to unlock secrets.');
        }

        // Load secrets from file if not in memory
        if (empty($this->memoryCache)) {
            $this->loadSecretsFromFile();
        }

        if (!array_key_exists($key, $this->memoryCache)) {
            throw new EncryptException("Secret \"$key\" not found.");
        }

        return $this->memoryCache[$key];
    }

    /**
     * Get all secrets (only works when unlocked)
     */
    public function all(): array
    {
        if (!$this->unlocked()) {
            throw new EncryptException('Vault is locked. Run \'encrypt setup <password>\' to unlock secrets.');
        }

        // Load secrets from file if not in memory
        if (empty($this->memoryCache)) {
            $this->loadSecretsFromFile();
        }

        return $this->memoryCache;
    }

    /**
     * Get vault status
     */
    public function status(): VaultStatus
    {
        // Check if vault is unlocked by looking for lock file
        $isUnlocked = $this->unlocked();

        // Load secrets from file if unlocked and not in memory
        if ($isUnlocked && empty($this->memoryCache)) {
            $this->loadSecretsFromFile();
        }

        $keys = array_keys($this->memoryCache);
        $lastModified = null;

        if ($this->exists()) {
            $lastModified = filemtime($this->configPath);
        }

        return new VaultStatus(!$isUnlocked, $keys, $lastModified);
    }

    /**
     * Reset/remove vault
     */
    public function reset(): void
    {
        if (is_dir($this->vaultPath)) {
            $this->removeDirectory($this->vaultPath);
        }
        $this->memoryCache = [];
    }

    /**
     * Check if vault is unlocked
     */
    public function unlocked(): bool
    {
        return file_exists($this->lockFilePath);
    }

    /**
     * Create lock file to indicate vault is unlocked
     */
    private function createLockFile(): void
    {
        $lockData = [
            'unlocked' => true,
            'timestamp' => date('c')
        ];

        if (file_put_contents($this->lockFilePath, json_encode($lockData, JSON_PRETTY_PRINT)) === false) {
            throw new EncryptException('Failed to create lock file');
        }
    }

    /**
     * Load secrets from file (for unlocked vault)
     */
    private function loadSecretsFromFile(): void
    {
        if (!file_exists($this->secretsPath)) {
            return;
        }

        $data = file_get_contents($this->secretsPath);
        if ($data === false) {
            throw new EncryptException('Failed to read secrets file');
        }

        $secrets = json_decode($data, true);
        if ($secrets === null) {
            throw new EncryptException('Invalid secrets file');
        }

        $this->memoryCache = [];
        foreach ($secrets as $key => $value) {
            // Check if the value is encrypted (starts with {) or plain text
            if (str_starts_with($value, '{')) {
                // This is encrypted data, we need to decrypt it
                // But we don't have the password here, so we can't decrypt
                // This should not happen in an unlocked vault
                throw new EncryptException("Secret \"$key\" is encrypted but vault is unlocked. This should not happen.");
            } else {
                // This is plain text data
                $this->memoryCache[$key] = $value;
            }
        }
    }

    /**
     * Save secrets to file (for unlocked vault)
     */
    private function saveSecretsToFile(): void
    {
        if (file_put_contents($this->secretsPath, json_encode($this->memoryCache, JSON_PRETTY_PRINT)) === false) {
            throw new EncryptException('Failed to save secrets file');
        }
    }

    /**
     * Recursively remove directory
     */
    private function removeDirectory(string $dir): void
    {
        if (!is_dir($dir)) {
            return;
        }

        $files = array_diff(scandir($dir), ['.', '..']);
        foreach ($files as $file) {
            $path = $dir . DIRECTORY_SEPARATOR . $file;
            is_dir($path) ? $this->removeDirectory($path) : unlink($path);
        }
        rmdir($dir);
    }
}
