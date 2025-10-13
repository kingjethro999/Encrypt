package main

import (
	"fmt"
	"os"

	"encrypt/internal/sdk"
)

func main() {
	fmt.Println("🧪 Testing Auto-Converter Engine (Go)")
	fmt.Println("=====================================")
	fmt.Println()
	
	// Test 1: Try to get secret when vault is locked (should fail without password)
	fmt.Println("1. Testing locked vault without password...")
	_, err := sdk.GetSDK().Get("API_KEY", "")
	if err != nil {
		fmt.Printf("✅ Expected: Vault is locked - %v\n", err)
	} else {
		fmt.Println("❌ Unexpected: Got secret from locked vault")
	}
	
	// Test 2: Try with environment variable
	fmt.Println("\n2. Testing with ENCRYPT_PASSWORD environment variable...")
	os.Setenv("ENCRYPT_PASSWORD", "mypassword")
	
	apiKey, err := sdk.GetSDK().Get("API_KEY", "")
	if err != nil {
		fmt.Printf("❌ Failed: Could not get secret with environment password - %v\n", err)
	} else {
		fmt.Printf("✅ Success: Got secret with environment password - %s\n", apiKey)
	}
	
	// Test 3: Try with explicit password parameter
	fmt.Println("\n3. Testing with explicit password parameter...")
	os.Unsetenv("ENCRYPT_PASSWORD") // Clear env var
	
	apiKey, err = sdk.GetSDK().Get("API_KEY", "mypassword")
	if err != nil {
		fmt.Printf("❌ Failed: Could not get secret with explicit password - %v\n", err)
	} else {
		fmt.Printf("✅ Success: Got secret with explicit password - %s\n", apiKey)
	}
	
	// Test 4: Test production-ready functions
	fmt.Println("\n4. Testing production-ready functions...")
	os.Setenv("ENCRYPT_PASSWORD", "mypassword")
	
	apiKey, err = sdk.GetSecret("API_KEY")
	if err != nil {
		fmt.Printf("❌ Failed: get_secret() failed - %v\n", err)
	} else {
		fmt.Printf("✅ Success: get_secret() works - %s\n", apiKey)
	}
	
	dbURL, err := sdk.GetSecret("DB_URL")
	if err != nil {
		fmt.Printf("❌ Failed: get_secret() failed - %v\n", err)
	} else {
		fmt.Printf("✅ Success: get_secret() works - %s\n", dbURL)
	}
	
	// Test 5: Test development mode (should try common passwords)
	fmt.Println("\n5. Testing development mode...")
	os.Unsetenv("ENCRYPT_PASSWORD")
	os.Setenv("NODE_ENV", "development")
	
	apiKey, err = sdk.GetSDK().Get("API_KEY", "")
	if err != nil {
		fmt.Printf("❌ Failed: Development mode failed - %v\n", err)
	} else {
		fmt.Printf("✅ Success: Development mode worked - %s\n", apiKey)
	}
	
	fmt.Println("\n🎉 Auto-converter tests completed!")
}
