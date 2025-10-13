package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"encrypt/internal/vault"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	// Colors
	red    = color.New(color.FgRed).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
	blue   = color.New(color.FgBlue).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "A top-level secrets orchestrator. Not just another .env tool — this one encrypts, locks, and sets you up for secure local and team dev.",
	Long: `A top-level secrets orchestrator. Not just another .env tool — this one encrypts, locks, and sets you up for secure local and team dev.

This tool provides:
• Secure encryption with triple-layer protection
• Easy team onboarding with single command
• Runtime SDK for in-code secret access
• Git-safe encrypted storage
• Beautiful CLI with colorized output`,
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create .encrypt vault",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("⠋ Initializing vault...")
		
		v := vault.NewVault()
		if err := v.Init(); err != nil {
			fmt.Printf("\r❌ Failed to initialize vault\n")
			fmt.Printf("%s\n", red(err.Error()))
			os.Exit(1)
		}
		
		fmt.Printf("\r✅ Vault initialized successfully!\n")
		fmt.Printf("📁 %s\n", blue("Created .encrypt directory with secure configuration."))
	},
}

// lockupCmd represents the lockup command
var lockupCmd = &cobra.Command{
	Use:   "lockup [password]",
	Short: "Encrypt and secure secrets with password",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		password := args[0]
		
		fmt.Print("⠋ Locking up secrets...")
		
		v := vault.NewVault()
		if !v.Exists() {
			fmt.Printf("\r❌ Vault not found. Run 'encrypt init' first.\n")
			os.Exit(1)
		}
		
		if err := v.Lockup(password); err != nil {
			fmt.Printf("\r❌ Failed to lock secrets\n")
			fmt.Printf("%s\n", red(err.Error()))
			os.Exit(1)
		}
		
		fmt.Printf("\r✅ Secrets locked successfully!\n")
		fmt.Printf("🔒 %s\n", blue("Your secrets are now encrypted and secure."))
	},
}

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup [password]",
	Short: "Set up secrets on a new machine",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		password := args[0]
		
		fmt.Print("⠋ Setting up vault...")
		
		v := vault.NewVault()
		if !v.Exists() {
			fmt.Printf("\r❌ Vault not found. Run 'encrypt init' first.\n")
			os.Exit(1)
		}
		
		if err := v.Setup(password); err != nil {
			fmt.Printf("\r❌ Failed to setup vault\n")
			fmt.Printf("%s\n", red(err.Error()))
			os.Exit(1)
		}
		
		fmt.Printf("\r✅ Vault unlocked successfully!\n")
		fmt.Printf("🔓 %s\n", blue("Your secrets are now available for use."))
	},
}

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set [key=value]",
	Short: "Add/update a key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		keyValue := args[0]
		
		parts := strings.SplitN(keyValue, "=", 2)
		if len(parts) != 2 {
			fmt.Printf("❌ %s\n", red("Invalid format. Use: encrypt set KEY=value"))
			os.Exit(1)
		}
		
		key, value := parts[0], parts[1]
		
		fmt.Printf("⠋ Setting secret: %s", key)
		
		v := vault.NewVault()
		if !v.Exists() {
			fmt.Printf("\r❌ Vault not found. Run 'encrypt init' first.\n")
			os.Exit(1)
		}
		
		if !v.Unlocked() {
			fmt.Printf("\r❌ Vault is locked. Run 'encrypt setup <password>' to unlock secrets.\n")
			os.Exit(1)
		}
		
		if err := v.Set(key, value); err != nil {
			fmt.Printf("\r❌ Failed to set secret\n")
			fmt.Printf("%s\n", red(err.Error()))
			os.Exit(1)
		}
		
		fmt.Printf("\r✅ Secret '%s' set successfully!\n", key)
	},
}

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Fetch decrypted value",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		
		v := vault.NewVault()
		if !v.Exists() {
			fmt.Printf("❌ %s\n", red("Vault not found. Run 'encrypt init' first."))
			os.Exit(1)
		}
		
		if !v.Unlocked() {
			fmt.Printf("❌ %s\n", red("Vault is locked. Run 'encrypt setup <password>' to unlock secrets."))
			os.Exit(1)
		}
		
		value, err := v.Get(key)
		if err != nil {
			fmt.Printf("❌ %s\n", red(err.Error()))
			os.Exit(1)
		}
		
		fmt.Println(value)
	},
}

// unlockCmd represents the unlock command
var unlockCmd = &cobra.Command{
	Use:   "unlock",
	Short: "Decrypt everything into .env",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("⠋ Unlocking secrets...")
		
		v := vault.NewVault()
		if !v.Exists() {
			fmt.Printf("\r❌ Vault not found. Run 'encrypt init' first.\n")
			os.Exit(1)
		}
		
		if !v.Unlocked() {
			fmt.Printf("\r❌ Vault is locked. Run 'encrypt setup <password>' to unlock secrets.\n")
			os.Exit(1)
		}
		
		secrets, err := v.All()
		if err != nil {
			fmt.Printf("\r❌ Failed to unlock secrets\n")
			fmt.Printf("%s\n", red(err.Error()))
			os.Exit(1)
		}
		
		var envContent strings.Builder
		for k, v := range secrets {
			envContent.WriteString(fmt.Sprintf("%s=%s\n", k, v))
		}
		
		if err := os.WriteFile(".env", []byte(envContent.String()), 0644); err != nil {
			fmt.Printf("\r❌ Failed to write .env file\n")
			fmt.Printf("%s\n", red(err.Error()))
			os.Exit(1)
		}
		
		fmt.Printf("\r✅ Secrets unlocked and written to .env\n")
	},
}

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check if vault is locked, list keys",
	Run: func(cmd *cobra.Command, args []string) {
		v := vault.NewVault()
		if !v.Exists() {
			fmt.Printf("⚠️ %s\n", yellow("Vault not found. Run 'encrypt init' first."))
			return
		}
		
		status, err := v.Status()
		if err != nil {
			fmt.Printf("❌ %s\n", red(err.Error()))
			os.Exit(1)
		}
		
		fmt.Printf("📊 %s\n", blue("Vault Status:"))
		fmt.Println("┌─────────────────┬──────────────────────────────┐")
		fmt.Println("│ Property        │ Value                        │")
		fmt.Println("├─────────────────┼──────────────────────────────┤")
		
		statusIcon := "🔒"
		statusText := "Locked"
		statusColor := red
		if !status.IsLocked {
			statusIcon = "🔓"
			statusText = "Unlocked"
			statusColor = green
		}
		
		fmt.Printf("│ Status          │ %s %s │\n", statusIcon, statusColor(statusText))
		fmt.Printf("│ Keys            │ %-28d │\n", len(status.Keys))
		
		if !status.LastModified.IsZero() {
			fmt.Printf("│ Last Modified   │ %-28s │\n", status.LastModified.Format("2006-01-02T15:04:05.000Z"))
		}
		
		fmt.Println("└─────────────────┴──────────────────────────────┘")
		
		if len(status.Keys) > 0 {
			fmt.Printf("\n🔑 %s\n", blue("Available keys:"))
			for _, key := range status.Keys {
				fmt.Printf("  • %s\n", key)
			}
		}
	},
}

// resetCmd represents the reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Remove vault (careful!)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Are you sure you want to reset the vault? This will delete all encrypted secrets. (y/N): ")
		
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))
		
		if response != "y" && response != "yes" {
			fmt.Printf("⚠️ %s\n", yellow("Reset cancelled."))
			return
		}
		
		fmt.Print("⠋ Resetting vault...")
		
		v := vault.NewVault()
		if err := v.Reset(); err != nil {
			fmt.Printf("\r❌ Failed to reset vault\n")
			fmt.Printf("%s\n", red(err.Error()))
			os.Exit(1)
		}
		
		fmt.Printf("\r✅ Vault reset successfully!\n")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	// Add commands
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(lockupCmd)
	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(unlockCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(resetCmd)
	
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("❌ %s\n", red(err.Error()))
		os.Exit(1)
	}
}
