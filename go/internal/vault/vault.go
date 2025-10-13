package vault

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"encrypt/internal/crypto"
)

const (
	VaultDir     = ".encrypt"
	ConfigFile   = "vault.lock"
	SecretsFile  = "secrets.enc.json"
	GitignoreFile = ".gitignore"
	LockFile     = "vault.unlocked"
)

// VaultConfig represents the vault configuration
type VaultConfig struct {
	PasswordHash string    `json:"password_hash"`
	Salt         string    `json:"salt"`
	HMAC         string    `json:"hmac"`
	CreatedAt    time.Time `json:"created_at"`
	Version      string    `json:"version"`
}

// VaultStatus represents the vault status
type VaultStatus struct {
	IsLocked     bool      `json:"is_locked"`
	Keys         []string  `json:"keys"`
	LastModified time.Time `json:"last_modified"`
}

// VaultLockFile represents the lock file structure
type VaultLockFile struct {
	Unlocked  bool      `json:"unlocked"`
	Timestamp time.Time `json:"timestamp"`
}

// Vault manages encrypted secrets storage
type Vault struct {
	projectRoot  string
	vaultPath    string
	configPath   string
	secretsPath  string
	gitignorePath string
	lockFilePath string
	memoryCache  map[string]string
	crypto       *crypto.Crypto
}

// NewVault creates a new Vault instance
func NewVault() *Vault {
	projectRoot, _ := os.Getwd()
	vaultPath := filepath.Join(projectRoot, VaultDir)
	
	return &Vault{
		projectRoot:  projectRoot,
		vaultPath:    vaultPath,
		configPath:   filepath.Join(vaultPath, ConfigFile),
		secretsPath:  filepath.Join(vaultPath, SecretsFile),
		gitignorePath: filepath.Join(projectRoot, GitignoreFile),
		lockFilePath: filepath.Join(vaultPath, LockFile),
		memoryCache:  make(map[string]string),
		crypto:       crypto.NewCrypto(),
	}
}

// Init initializes the vault directory structure
func (v *Vault) Init() error {
	// Create vault directory
	if err := os.MkdirAll(v.vaultPath, 0755); err != nil {
		return fmt.Errorf("failed to create vault directory: %w", err)
	}

	// Create initial empty secrets file
	if _, err := os.Stat(v.secretsPath); os.IsNotExist(err) {
		if err := os.WriteFile(v.secretsPath, []byte("{}"), 0644); err != nil {
			return fmt.Errorf("failed to create secrets file: %w", err)
		}
	}

	// Create initial config file
	if _, err := os.Stat(v.configPath); os.IsNotExist(err) {
		initialConfig := VaultConfig{
			PasswordHash: "",
			Salt:         "",
			HMAC:         "",
			CreatedAt:    time.Now(),
			Version:      "1.0.0",
		}
		
		configData, err := json.MarshalIndent(initialConfig, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal config: %w", err)
		}
		
		if err := os.WriteFile(v.configPath, configData, 0644); err != nil {
			return fmt.Errorf("failed to create config file: %w", err)
		}
	}

	// Update .gitignore
	return v.updateGitignore()
}

// updateGitignore updates .gitignore to exclude .encrypt directory
func (v *Vault) updateGitignore() error {
	var gitignoreContent string
	
	if _, err := os.Stat(v.gitignorePath); err == nil {
		content, err := os.ReadFile(v.gitignorePath)
		if err != nil {
			return fmt.Errorf("failed to read .gitignore: %w", err)
		}
		gitignoreContent = string(content)
	}

	if !strings.Contains(gitignoreContent, ".encrypt/") {
		gitignoreContent += "\n# Encrypt vault\n.encrypt/\n"
		
		file, err := os.OpenFile(v.gitignorePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return fmt.Errorf("failed to open .gitignore: %w", err)
		}
		defer file.Close()
		
		if _, err := file.WriteString("\n# Encrypt vault\n.encrypt/\n"); err != nil {
			return fmt.Errorf("failed to write to .gitignore: %w", err)
		}
	}

	return nil
}

// Exists checks if vault exists
func (v *Vault) Exists() bool {
	_, err := os.Stat(v.vaultPath)
	if os.IsNotExist(err) {
		return false
	}
	
	_, err = os.Stat(v.configPath)
	return err == nil
}

// Lockup locks up secrets with password
func (v *Vault) Lockup(password string) error {
	// Load secrets from file if not in memory
	if len(v.memoryCache) == 0 {
		if err := v.loadSecretsFromFile(); err != nil {
			return err
		}
	}

	if len(v.memoryCache) == 0 {
		return fmt.Errorf("no secrets to lock. Use 'encrypt set' to add secrets first")
	}

	// Encrypt all secrets
	encryptedSecrets := make(map[string]string)

	for key, value := range v.memoryCache {
		result, err := v.crypto.Encrypt(value, password)
		if err != nil {
			return fmt.Errorf("failed to encrypt secret %s: %w", key, err)
		}
		
		resultJSON, err := json.Marshal(result)
		if err != nil {
			return fmt.Errorf("failed to marshal encryption result: %w", err)
		}
		
		encryptedSecrets[key] = string(resultJSON)
	}

	// Create vault config
	configSalt, err := v.crypto.GenerateSalt()
	if err != nil {
		return fmt.Errorf("failed to generate config salt: %w", err)
	}

	configKey, err := v.crypto.DeriveKey(password, configSalt)
	if err != nil {
		return fmt.Errorf("failed to derive config key: %w", err)
	}

	secretsJSON, err := json.Marshal(encryptedSecrets)
	if err != nil {
		return fmt.Errorf("failed to marshal encrypted secrets: %w", err)
	}

	config := VaultConfig{
		PasswordHash: func() string {
			hash, _ := v.crypto.HashPassword(password)
			return hash
		}(),
		Salt:      configSalt,
		HMAC:      v.crypto.GenerateHMAC(string(secretsJSON), configKey),
		CreatedAt: time.Now(),
		Version:   "1.0.0",
	}

	// Write encrypted secrets and config
	secretsData, err := json.MarshalIndent(encryptedSecrets, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal secrets: %w", err)
	}

	if err := os.WriteFile(v.secretsPath, secretsData, 0644); err != nil {
		return fmt.Errorf("failed to write secrets file: %w", err)
	}

	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(v.configPath, configData, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	// Clear memory cache and remove lock file
	v.memoryCache = make(map[string]string)
	if err := os.Remove(v.lockFilePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove lock file: %w", err)
	}

	return nil
}

// Setup unlocks vault with password
func (v *Vault) Setup(password string) error {
	if !v.Exists() {
		return fmt.Errorf("vault not found. Run 'encrypt init' first")
	}

	// Read config
	configData, err := os.ReadFile(v.configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var config VaultConfig
	if err := json.Unmarshal(configData, &config); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// If this is a fresh vault (no password set), just unlock it
	if config.PasswordHash == "" {
		v.memoryCache = make(map[string]string)
		return v.createLockFile()
	}

	// Verify password
	if !v.crypto.VerifyPassword(password, config.PasswordHash) {
		return fmt.Errorf("invalid password")
	}

	// Load and decrypt secrets
	secretsData, err := os.ReadFile(v.secretsPath)
	if err != nil {
		return fmt.Errorf("failed to read secrets file: %w", err)
	}

	var encryptedSecrets map[string]string
	if err := json.Unmarshal(secretsData, &encryptedSecrets); err != nil {
		return fmt.Errorf("failed to unmarshal secrets: %w", err)
	}

	v.memoryCache = make(map[string]string)

	for key, encryptedData := range encryptedSecrets {
		// Check if the data is already decrypted (plain text) or encrypted
		if strings.HasPrefix(encryptedData, "{") {
			// This is encrypted data stored as JSON string, decrypt it
			var result crypto.EncryptionResult
			if err := json.Unmarshal([]byte(encryptedData), &result); err != nil {
				return fmt.Errorf("failed to unmarshal encryption result: %w", err)
			}

			decrypted, isValid := v.crypto.Decrypt(result.Encrypted, result.Salt, result.HMAC, password)
			if !isValid {
				return fmt.Errorf("failed to decrypt secret: %s", key)
			}

			v.memoryCache[key] = decrypted
		} else {
			// This is plain text data
			v.memoryCache[key] = encryptedData
		}
	}

	// Save decrypted secrets to file for easy access
	if err := v.saveSecretsToFile(); err != nil {
		return err
	}

	return v.createLockFile()
}

// Set sets a secret (only works when unlocked)
func (v *Vault) Set(key, value string) error {
	if !v.Unlocked() {
		return fmt.Errorf("vault is locked. Run 'encrypt setup <password>' to unlock secrets")
	}

	// Load secrets from file if not in memory
	if len(v.memoryCache) == 0 {
		if err := v.loadSecretsFromFile(); err != nil {
			return err
		}
	}

	v.memoryCache[key] = value
	return v.saveSecretsToFile()
}

// Get gets a secret (only works when unlocked)
func (v *Vault) Get(key string) (string, error) {
	if !v.Unlocked() {
		return "", fmt.Errorf("vault is locked. Run 'encrypt setup <password>' to unlock secrets")
	}

	// Load secrets from file if not in memory
	if len(v.memoryCache) == 0 {
		if err := v.loadSecretsFromFile(); err != nil {
			return "", err
		}
	}

	value, exists := v.memoryCache[key]
	if !exists {
		return "", fmt.Errorf("secret \"%s\" not found", key)
	}

	return value, nil
}

// All gets all secrets (only works when unlocked)
func (v *Vault) All() (map[string]string, error) {
	if !v.Unlocked() {
		return nil, fmt.Errorf("vault is locked. Run 'encrypt setup <password>' to unlock secrets")
	}

	// Load secrets from file if not in memory
	if len(v.memoryCache) == 0 {
		if err := v.loadSecretsFromFile(); err != nil {
			return nil, err
		}
	}

	// Return a copy
	result := make(map[string]string)
	for k, v := range v.memoryCache {
		result[k] = v
	}

	return result, nil
}

// Status gets vault status
func (v *Vault) Status() (*VaultStatus, error) {
	// Check if vault is unlocked by looking for lock file
	isUnlocked := v.Unlocked()

	// Load secrets from file if unlocked and not in memory
	if isUnlocked && len(v.memoryCache) == 0 {
		if err := v.loadSecretsFromFile(); err != nil {
			return nil, err
		}
	}

	keys := make([]string, 0, len(v.memoryCache))
	for key := range v.memoryCache {
		keys = append(keys, key)
	}

	var lastModified time.Time
	if v.Exists() {
		if info, err := os.Stat(v.configPath); err == nil {
			lastModified = info.ModTime()
		}
	}

	return &VaultStatus{
		IsLocked:     !isUnlocked,
		Keys:         keys,
		LastModified: lastModified,
	}, nil
}

// Reset removes vault
func (v *Vault) Reset() error {
	if err := os.RemoveAll(v.vaultPath); err != nil {
		return fmt.Errorf("failed to remove vault directory: %w", err)
	}
	v.memoryCache = make(map[string]string)
	return nil
}

// Unlocked checks if vault is unlocked
func (v *Vault) Unlocked() bool {
	_, err := os.Stat(v.lockFilePath)
	return err == nil
}

// createLockFile creates lock file to indicate vault is unlocked
func (v *Vault) createLockFile() error {
	lockData := VaultLockFile{
		Unlocked:  true,
		Timestamp: time.Now(),
	}

	data, err := json.MarshalIndent(lockData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal lock data: %w", err)
	}

	return os.WriteFile(v.lockFilePath, data, 0644)
}

// loadSecretsFromFile loads secrets from file (for unlocked vault)
func (v *Vault) loadSecretsFromFile() error {
	if _, err := os.Stat(v.secretsPath); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(v.secretsPath)
	if err != nil {
		return fmt.Errorf("failed to read secrets file: %w", err)
	}

	var secrets map[string]string
	if err := json.Unmarshal(data, &secrets); err != nil {
		return fmt.Errorf("failed to unmarshal secrets: %w", err)
	}

	v.memoryCache = make(map[string]string)
	for key, value := range secrets {
		// Check if the value is encrypted (starts with {) or plain text
		if strings.HasPrefix(value, "{") {
			// This is encrypted data, we need to decrypt it
			// But we don't have the password here, so we can't decrypt
			// This should not happen in an unlocked vault
			return fmt.Errorf("secret \"%s\" is encrypted but vault is unlocked. This should not happen", key)
		} else {
			// This is plain text data
			v.memoryCache[key] = value
		}
	}

	return nil
}

// saveSecretsToFile saves secrets to file (for unlocked vault)
func (v *Vault) saveSecretsToFile() error {
	data, err := json.MarshalIndent(v.memoryCache, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal secrets: %w", err)
	}

	return os.WriteFile(v.secretsPath, data, 0644)
}
