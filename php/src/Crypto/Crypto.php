<?php

namespace Encrypt\Crypto;

/**
 * Triple-layer encryption implementation for the Encrypt tool.
 * Provides AES-256-CBC encryption with PBKDF2 key derivation and HMAC signatures.
 */
class Crypto
{
    private const ALGORITHM = 'AES-256-CBC';
    private const KEY_LENGTH = 32;
    private const IV_LENGTH = 16;
    private const SALT_LENGTH = 32;
    private const ITERATIONS = 100000;

    /**
     * Generate a random salt for key derivation
     */
    public function generateSalt(): string
    {
        return bin2hex(random_bytes(self::SALT_LENGTH));
    }

    /**
     * Derive encryption key from password using PBKDF2
     */
    public function deriveKey(string $password, string $salt): string
    {
        return hash_pbkdf2('sha256', $password, hex2bin($salt), self::ITERATIONS, self::KEY_LENGTH, true);
    }

    /**
     * Generate HMAC signature for data integrity
     */
    public function generateHMAC(string $data, string $key): string
    {
        return hash_hmac('sha256', $data, $key);
    }

    /**
     * Triple-layer encryption:
     * 1. Generate random salt
     * 2. Derive key using PBKDF2
     * 3. AES-256-CBC encryption with HMAC
     */
    public function encrypt(string $plaintext, string $password): EncryptionResult
    {
        // Phase 1: Generate salt
        $salt = $this->generateSalt();

        // Phase 2: Derive key from password
        $key = $this->deriveKey($password, $salt);

        // Phase 3: AES-256-CBC encryption
        $iv = random_bytes(self::IV_LENGTH);
        $encrypted = openssl_encrypt($plaintext, self::ALGORITHM, $key, OPENSSL_RAW_DATA, $iv);

        if ($encrypted === false) {
            throw new \RuntimeException('Encryption failed');
        }

        // Combine IV and encrypted data
        $encryptedData = bin2hex($iv) . ':' . bin2hex($encrypted);

        // Phase 3: Generate HMAC signature
        $hmac = $this->generateHMAC($encryptedData, $key);

        return new EncryptionResult($encryptedData, $salt, $hmac);
    }

    /**
     * Triple-layer decryption:
     * 1. Verify HMAC signature
     * 2. Derive key using PBKDF2
     * 3. AES-256-CBC decryption
     */
    public function decrypt(string $encryptedData, string $salt, string $hmac, string $password): array
    {
        // Phase 2: Derive key from password
        $key = $this->deriveKey($password, $salt);

        // Phase 3: Verify HMAC signature
        $expectedHMAC = $this->generateHMAC($encryptedData, $key);
        if (!hash_equals($expectedHMAC, $hmac)) {
            return ['', false];
        }

        // Phase 3: AES-256-CBC decryption
        $parts = explode(':', $encryptedData, 2);
        if (count($parts) !== 2) {
            return ['', false];
        }

        $iv = hex2bin($parts[0]);
        $encrypted = hex2bin($parts[1]);

        $decrypted = openssl_decrypt($encrypted, self::ALGORITHM, $key, OPENSSL_RAW_DATA, $iv);

        if ($decrypted === false) {
            return ['', false];
        }

        return [$decrypted, true];
    }

    /**
     * Hash password for storage using PBKDF2
     */
    public function hashPassword(string $password): string
    {
        $salt = $this->generateSalt();
        $saltBytes = hex2bin($salt);
        $hash = hash_pbkdf2('sha256', $password, $saltBytes, self::ITERATIONS, 64, true);
        return $salt . ':' . bin2hex($hash);
    }

    /**
     * Verify password against stored hash
     */
    public function verifyPassword(string $password, string $storedHash): bool
    {
        $parts = explode(':', $storedHash, 2);
        if (count($parts) !== 2) {
            return false;
        }

        $salt = $parts[0];
        $storedHashHex = $parts[1];

        $saltBytes = hex2bin($salt);
        $computedHash = hash_pbkdf2('sha256', $password, $saltBytes, self::ITERATIONS, 64, true);

        return hash_equals(bin2hex($computedHash), $storedHashHex);
    }
}
