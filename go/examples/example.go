package main

import (
	"fmt"
	"log"

	"encrypt/internal/sdk"
)

func main() {
	fmt.Println("💻 Example usage of the encrypt package in Go code")
	fmt.Println("==================================================")
	
	// Check if vault is unlocked
	if !sdk.IsVaultUnlocked() {
		fmt.Println("❌ Vault is locked. Run 'encrypt setup <password>' to unlock secrets.")
		return
	}
	
	// Get secrets
	apiKey, err := sdk.GetSecret("API_KEY")
	if err != nil {
		fmt.Printf("⚠️ API_KEY not found: %v\n", err)
	} else {
		fmt.Println("✅ Secrets retrieved successfully!")
		if len(apiKey) > 4 {
			fmt.Printf("API Key: ***%s\n", apiKey[len(apiKey)-4:])
		} else {
			fmt.Printf("API Key: ***%s\n", apiKey)
		}
	}
	
	dbURL, err := sdk.GetSecret("DB_URL")
	if err != nil {
		fmt.Printf("⚠️ DB_URL not found: %v\n", err)
	} else {
		if len(dbURL) > 10 {
			fmt.Printf("DB URL: ***%s\n", dbURL[len(dbURL)-10:])
		} else {
			fmt.Printf("DB URL: ***%s\n", dbURL)
		}
	}
	
	// Get all secrets
	allSecrets, err := sdk.GetAllSecrets()
	if err != nil {
		fmt.Printf("⚠️ Error getting all secrets: %v\n", err)
	} else {
		keys := make([]string, 0, len(allSecrets))
		for key := range allSecrets {
			keys = append(keys, key)
		}
		fmt.Printf("Available keys: %v\n", keys)
	}
	
	// Get status
	status, err := sdk.GetStatus()
	if err != nil {
		fmt.Printf("⚠️ Error getting status: %v\n", err)
	} else {
		statusText := "Locked"
		if !status.IsLocked {
			statusText = "Unlocked"
		}
		fmt.Printf("Vault status: %s\n", statusText)
		fmt.Printf("Number of keys: %d\n", len(status.Keys))
	}
}
