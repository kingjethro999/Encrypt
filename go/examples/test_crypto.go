package main

import (
	"fmt"
	"log"

	"encrypt/internal/crypto"
)

func main() {
	fmt.Println("🧪 Testing TripleEncryption (Go)...")
	
	// Test data
	plaintext := "Hello, World!"
	password := "test-password-123"
	
	fmt.Printf("Original: %s\n", plaintext)
	
	// Encrypt
	encrypted, err := crypto.NewCrypto().Encrypt(plaintext, password)
	if err != nil {
		log.Fatalf("❌ Test failed: %v", err)
	}
	
	fmt.Println("✅ Encryption successful")
	fmt.Printf("Encrypted data length: %d\n", len(encrypted.Encrypted))
	
	// Decrypt
	decrypted, isValid := crypto.NewCrypto().Decrypt(encrypted.Encrypted, encrypted.Salt, encrypted.HMAC, password)
	if !isValid || decrypted != plaintext {
		log.Fatal("❌ Decryption failed")
	}
	
	fmt.Println("✅ Decryption successful")
	fmt.Printf("Decrypted: %s\n", decrypted)
	
	// Test password hashing
	hashResult, err := crypto.NewCrypto().HashPassword(password)
	if err != nil {
		log.Fatalf("❌ Password hashing failed: %v", err)
	}
	
	isValidHash := crypto.NewCrypto().VerifyPassword(password, hashResult)
	if !isValidHash {
		log.Fatal("❌ Password hashing/verification failed")
	}
	
	fmt.Println("✅ Password hashing/verification successful")
	
	// Test with different password (should fail)
	_, wrongValid := crypto.NewCrypto().Decrypt(encrypted.Encrypted, encrypted.Salt, encrypted.HMAC, "wrong-password")
	if wrongValid {
		log.Fatal("❌ Wrong password incorrectly accepted")
	}
	
	fmt.Println("✅ Wrong password correctly rejected")
	fmt.Println("\n🎉 All tests passed!")
}
