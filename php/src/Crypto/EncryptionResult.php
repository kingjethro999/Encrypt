<?php

namespace Encrypt\Crypto;

/**
 * Represents the result of encryption operation
 */
class EncryptionResult
{
    public function __construct(
        public readonly string $encrypted,
        public readonly string $salt,
        public readonly string $hmac
    ) {}

    /**
     * Convert to array for JSON serialization
     */
    public function toArray(): array
    {
        return [
            'encrypted' => $this->encrypted,
            'salt' => $this->salt,
            'hmac' => $this->hmac
        ];
    }

    /**
     * Create from array (for JSON deserialization)
     */
    public static function fromArray(array $data): self
    {
        return new self(
            $data['encrypted'],
            $data['salt'],
            $data['hmac']
        );
    }
}
