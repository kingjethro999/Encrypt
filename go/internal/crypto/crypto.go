package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
)

const (
	KeyLength   = 32
	SaltLength  = 32
	IVLength    = 16
	Iterations  = 100000
)

// EncryptionResult represents the result of encryption
type EncryptionResult struct {
	Encrypted string `json:"encrypted"`
	Salt      string `json:"salt"`
	HMAC      string `json:"hmac"`
}

// Crypto provides triple-layer encryption functionality
type Crypto struct{}

// NewCrypto creates a new Crypto instance
func NewCrypto() *Crypto {
	return &Crypto{}
}

// GenerateSalt generates a random salt for key derivation
func (c *Crypto) GenerateSalt() (string, error) {
	salt := make([]byte, SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}
	return hex.EncodeToString(salt), nil
}

// DeriveKey derives encryption key from password using PBKDF2
func (c *Crypto) DeriveKey(password, salt string) ([]byte, error) {
	saltBytes, err := hex.DecodeString(salt)
	if err != nil {
		return nil, fmt.Errorf("failed to decode salt: %w", err)
	}
	
	key := pbkdf2.Key([]byte(password), saltBytes, Iterations, KeyLength, sha256.New)
	return key, nil
}

// GenerateHMAC generates HMAC signature for data integrity
func (c *Crypto) GenerateHMAC(data string, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// Encrypt performs triple-layer encryption:
// 1. Generate random salt
// 2. Derive key using PBKDF2
// 3. AES-256-CBC encryption with HMAC
func (c *Crypto) Encrypt(plaintext, password string) (*EncryptionResult, error) {
	// Phase 1: Generate salt
	salt, err := c.GenerateSalt()
	if err != nil {
		return nil, err
	}

	// Phase 2: Derive key from password
	key, err := c.DeriveKey(password, salt)
	if err != nil {
		return nil, err
	}

	// Phase 3: AES-256-CBC encryption
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// Generate random IV
	iv := make([]byte, IVLength)
	if _, err := rand.Read(iv); err != nil {
		return nil, fmt.Errorf("failed to generate IV: %w", err)
	}

	// Encrypt the plaintext
	mode := cipher.NewCBCEncrypter(block, iv)
	paddedPlaintext := c.pkcs7Pad([]byte(plaintext), aes.BlockSize)
	encrypted := make([]byte, len(paddedPlaintext))
	mode.CryptBlocks(encrypted, paddedPlaintext)

	// Combine IV and encrypted data
	encryptedData := hex.EncodeToString(iv) + ":" + hex.EncodeToString(encrypted)

	// Phase 3: Generate HMAC signature
	hmac := c.GenerateHMAC(encryptedData, key)

	return &EncryptionResult{
		Encrypted: encryptedData,
		Salt:      salt,
		HMAC:      hmac,
	}, nil
}

// Decrypt performs triple-layer decryption:
// 1. Verify HMAC signature
// 2. Derive key using PBKDF2
// 3. AES-256-CBC decryption
func (c *Crypto) Decrypt(encryptedData, salt, hmac, password string) (string, bool) {
	// Phase 2: Derive key from password
	key, err := c.DeriveKey(password, salt)
	if err != nil {
		return "", false
	}

	// Phase 3: Verify HMAC signature
	expectedHMAC := c.GenerateHMAC(encryptedData, key)
	if expectedHMAC != hmac {
		return "", false
	}

	// Phase 3: AES-256-CBC decryption
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", false
	}

	// Split IV and encrypted data
	parts := splitString(encryptedData, ":", 2)
	if len(parts) != 2 {
		return "", false
	}

	iv, err := hex.DecodeString(parts[0])
	if err != nil {
		return "", false
	}

	encrypted, err := hex.DecodeString(parts[1])
	if err != nil {
		return "", false
	}

	// Decrypt
	mode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(encrypted))
	mode.CryptBlocks(decrypted, encrypted)

	// Remove padding
	unpadded, err := c.pkcs7Unpad(decrypted, aes.BlockSize)
	if err != nil {
		return "", false
	}

	return string(unpadded), true
}

// HashPassword hashes password for storage using PBKDF2
func (c *Crypto) HashPassword(password string) (string, error) {
	salt, err := c.GenerateSalt()
	if err != nil {
		return "", err
	}

	saltBytes, err := hex.DecodeString(salt)
	if err != nil {
		return "", err
	}

	hash := pbkdf2.Key([]byte(password), saltBytes, Iterations, 64, sha256.New)
	return salt + ":" + hex.EncodeToString(hash), nil
}

// VerifyPassword verifies password against stored hash
func (c *Crypto) VerifyPassword(password, storedHash string) bool {
	parts := splitString(storedHash, ":", 2)
	if len(parts) != 2 {
		return false
	}

	salt := parts[0]
	storedHashHex := parts[1]

	saltBytes, err := hex.DecodeString(salt)
	if err != nil {
		return false
	}

	computedHash := pbkdf2.Key([]byte(password), saltBytes, Iterations, 64, sha256.New)
	return hex.EncodeToString(computedHash) == storedHashHex
}

// pkcs7Pad adds PKCS7 padding to data
func (c *Crypto) pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := make([]byte, padding)
	for i := range padtext {
		padtext[i] = byte(padding)
	}
	return append(data, padtext...)
}

// pkcs7Unpad removes PKCS7 padding from data
func (c *Crypto) pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, fmt.Errorf("invalid padding")
	}

	unpadding := int(data[length-1])
	if unpadding > blockSize || unpadding == 0 {
		return nil, fmt.Errorf("invalid padding")
	}

	if length < unpadding {
		return nil, fmt.Errorf("invalid padding")
	}

	return data[:(length - unpadding)], nil
}

// splitString splits a string by delimiter with limit
func splitString(s, delimiter string, limit int) []string {
	parts := make([]string, 0, limit)
	start := 0
	
	for i := 0; i < limit-1; i++ {
		pos := findString(s, delimiter, start)
		if pos == -1 {
			break
		}
		parts = append(parts, s[start:pos])
		start = pos + len(delimiter)
	}
	
	parts = append(parts, s[start:])
	return parts
}

// findString finds the first occurrence of substr in s starting from start
func findString(s, substr string, start int) int {
	if start >= len(s) {
		return -1
	}
	
	for i := start; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	
	return -1
}
