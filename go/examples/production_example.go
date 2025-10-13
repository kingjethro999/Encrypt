package main

import (
	"fmt"
	"os"

	"encrypt/internal/sdk"
)

type Config struct {
	APIKey    string
	Database  string
	Port      string
}

func main() {
	fmt.Println("🚀 Production Usage Example (Go)")
	fmt.Println("=================================")
	fmt.Println()
	
	// Method 1: Using environment variable (Recommended for production)
	fmt.Println("Method 1: Environment Variable (Recommended)")
	fmt.Println("Set ENCRYPT_PASSWORD=your-password in your environment")
	fmt.Println()
	
	apiKey, err := sdk.GetSecret("API_KEY")
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		fmt.Println("\n💡 To fix this:")
		fmt.Println("1. Set ENCRYPT_PASSWORD environment variable")
		fmt.Println("2. Or provide password as second parameter")
		fmt.Println("3. Or run 'encrypt setup <password>' first")
		return
	}
	
	fmt.Println("✅ Successfully retrieved secrets:")
	if len(apiKey) > 4 {
		fmt.Printf("API Key: ***%s\n", apiKey[len(apiKey)-4:])
	} else {
		fmt.Printf("API Key: ***%s\n", apiKey)
	}
	
	dbURL, err := sdk.GetSecret("DB_URL")
	if err != nil {
		fmt.Printf("DB URL: ***\n")
	} else {
		if len(dbURL) > 10 {
			fmt.Printf("DB URL: ***%s\n", dbURL[len(dbURL)-10:])
		} else {
			fmt.Printf("DB URL: ***%s\n", dbURL)
		}
	}
	
	// Use in your application
	config := Config{
		APIKey:   apiKey,
		Database: dbURL,
		Port:     getEnvOrDefault("PORT", "3000"),
	}
	
	fmt.Println("\n📋 Application config ready:")
	if len(config.APIKey) > 4 {
		fmt.Printf("  api_key: ***%s\n", config.APIKey[len(config.APIKey)-4:])
	} else {
		fmt.Printf("  api_key: ***%s\n", config.APIKey)
	}
	if len(config.Database) > 10 {
		fmt.Printf("  database: ***%s\n", config.Database[len(config.Database)-10:])
	} else {
		fmt.Printf("  database: ***%s\n", config.Database)
	}
	fmt.Printf("  port: %s\n", config.Port)
	
	fmt.Println("\n==================================================")
	fmt.Println("Method 2: Explicit Password Parameter")
	fmt.Println("==================================================")
	
	apiKey, err = sdk.GetSDK().Get("API_KEY", "mypassword")
	if err != nil {
		fmt.Printf("❌ Error with explicit password: %v\n", err)
	} else {
		if len(apiKey) > 4 {
			fmt.Printf("✅ Success with explicit password: ***%s\n", apiKey[len(apiKey)-4:])
		} else {
			fmt.Printf("✅ Success with explicit password: ***%s\n", apiKey)
		}
	}
	
	fmt.Println("\n==================================================")
	fmt.Println("Method 3: Development Mode")
	fmt.Println("==================================================")
	
	// In development, you can set NODE_ENV=development
	// and it will try common passwords automatically
	os.Setenv("NODE_ENV", "development")
	
	apiKey, err = sdk.GetSDK().Get("API_KEY", "")
	if err != nil {
		fmt.Printf("❌ Development mode failed: %v\n", err)
	} else {
		if len(apiKey) > 4 {
			fmt.Printf("✅ Development mode success: ***%s\n", apiKey[len(apiKey)-4:])
		} else {
			fmt.Printf("✅ Development mode success: ***%s\n", apiKey)
		}
	}
	
	fmt.Println("\n🎯 Production Deployment Examples:")
	fmt.Println("==================================")
	fmt.Println("Docker:")
	fmt.Println("  ENV ENCRYPT_PASSWORD=your-production-password")
	fmt.Println()
	fmt.Println("Kubernetes:")
	fmt.Println("  env:")
	fmt.Println("  - name: ENCRYPT_PASSWORD")
	fmt.Println("    valueFrom:")
	fmt.Println("      secretKeyRef:")
	fmt.Println("        name: encrypt-secrets")
	fmt.Println("        key: password")
	fmt.Println()
	fmt.Println("Heroku:")
	fmt.Println("  heroku config:set ENCRYPT_PASSWORD=your-password")
	fmt.Println()
	fmt.Println("AWS Lambda:")
	fmt.Println("  Set ENCRYPT_PASSWORD in environment variables")
	fmt.Println()
	fmt.Println("Go application:")
	fmt.Println("  import \"encrypt/internal/sdk\"")
	fmt.Println("  apiKey, err := sdk.GetSecret(\"API_KEY\")")
	fmt.Println("  databaseURL, err := sdk.GetSecret(\"DATABASE_URL\")")
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
