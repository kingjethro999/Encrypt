<?php

namespace Encrypt\Vault;

/**
 * Represents the vault status
 */
class VaultStatus
{
    public function __construct(
        public readonly bool $isLocked,
        public readonly array $keys,
        public readonly ?int $lastModified
    ) {}

    /**
     * Convert to array for JSON serialization
     */
    public function toArray(): array
    {
        return [
            'is_locked' => $this->isLocked,
            'keys' => $this->keys,
            'last_modified' => $this->lastModified ? date('c', $this->lastModified) : null
        ];
    }
}
