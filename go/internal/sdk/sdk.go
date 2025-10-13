package sdk

import (
	"fmt"
	"os"

	"encrypt/internal/vault"
)

// SDK provides runtime SDK functionality for the Encrypt tool
type SDK struct {
	globalVault *vault.Vault
}

var globalSDK *SDK

// GetSDK returns the global SDK instance
func GetSDK() *SDK {
	if globalSDK == nil {
		globalSDK = &SDK{
			globalVault: vault.NewVault(),
		}
	}
	return globalSDK
}

// autoUnlock auto-unlocks vault if password is available
func (s *SDK) autoUnlock(password string) error {
	v := s.globalVault

	// If already unlocked, no need to do anything
	if v.Unlocked() {
		return nil
	}

	// Try provided password first
	if password != "" {
		if err := v.Setup(password); err == nil {
			return nil
		}
	}

	// Try environment variable
	if envPassword := os.Getenv("ENCRYPT_PASSWORD"); envPassword != "" {
		if err := v.Setup(envPassword); err == nil {
			return nil
		}
	}

	// In development, we can be more lenient
	if os.Getenv("NODE_ENV") != "production" {
		// Try common development passwords
		devPasswords := []string{"dev", "development", "test", "password", "123456"}
		for _, devPassword := range devPasswords {
			if err := v.Setup(devPassword); err == nil {
				return nil
			}
		}
	}

	// If we get here, we couldn't unlock the vault
	return fmt.Errorf("vault is locked and no valid password found. Set ENCRYPT_PASSWORD environment variable or provide password parameter")
}

// Get gets a secret value from the vault with auto-unlock support
func (s *SDK) Get(key, password string) (string, error) {
	if err := s.autoUnlock(password); err != nil {
		return "", err
	}
	return s.globalVault.Get(key)
}

// Set sets a secret value in the vault with auto-unlock support
func (s *SDK) Set(key, value, password string) error {
	if err := s.autoUnlock(password); err != nil {
		return err
	}
	return s.globalVault.Set(key, value)
}

// AllSecrets gets all secrets from the vault with auto-unlock support
func (s *SDK) AllSecrets(password string) (map[string]string, error) {
	if err := s.autoUnlock(password); err != nil {
		return nil, err
	}
	return s.globalVault.All()
}

// Status gets vault status
func (s *SDK) Status() (*vault.VaultStatus, error) {
	return s.globalVault.Status()
}

// IsUnlocked checks if vault is unlocked
func (s *SDK) IsUnlocked() bool {
	return s.globalVault.Unlocked()
}

// AutoSetup auto setup helper - checks if vault is locked and provides helpful error
func (s *SDK) AutoSetup() error {
	v := s.globalVault

	if !v.Exists() {
		return fmt.Errorf("vault not found. Run 'encrypt init' first")
	}

	if !v.Unlocked() {
		return fmt.Errorf("vault is locked. Run 'encrypt setup <password>' to unlock secrets")
	}

	return nil
}

// GetSecret gets a secret with automatic environment variable support
// This is the recommended function for production use
func (s *SDK) GetSecret(key string) (string, error) {
	return s.Get(key, "")
}

// SetSecret sets a secret with automatic environment variable support
// This is the recommended function for production use
func (s *SDK) SetSecret(key, value string) error {
	return s.Set(key, value, "")
}

// GetAllSecrets gets all secrets with automatic environment variable support
// This is the recommended function for production use
func (s *SDK) GetAllSecrets() (map[string]string, error) {
	return s.AllSecrets("")
}

// Convenience functions for backward compatibility
func GetSecret(key string) (string, error) {
	return GetSDK().GetSecret(key)
}

func SetSecret(key, value string) error {
	return GetSDK().SetSecret(key, value)
}

func GetAllSecrets() (map[string]string, error) {
	return GetSDK().GetAllSecrets()
}

func GetStatus() (*vault.VaultStatus, error) {
	return GetSDK().Status()
}

func IsVaultUnlocked() bool {
	return GetSDK().IsUnlocked()
}

func AutoSetupVault() error {
	return GetSDK().AutoSetup()
}
