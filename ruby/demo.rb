#!/usr/bin/env ruby
# frozen_string_literal: true

require_relative "lib/encrypt"

def run_command(command, description)
  puts "📝 #{description}".yellow
  puts "$ #{command}".gray
  puts
  
  begin
    result = `#{command} 2>&1`
    if result.strip.length > 0
      puts result.strip.green
    end
  rescue StandardError => e
    puts "Error: #{e.message}".red
  end
  
  puts
end

def main
  puts "🔐 Encrypt Tool Demo (Ruby)".bold.blue
  puts "This demo shows the complete workflow of the encrypt Ruby tool.".gray
  puts
  
  # Clean up any existing vault
  begin
    `rm -rf .encrypt .env`
  rescue StandardError
    # Ignore errors
  end
  
  puts "1. Initialize the vault".bold.blue
  run_command("ruby exe/encrypt init", "Create a new encrypted vault")
  
  puts "2. Setup the vault (unlock for first time)".bold.blue
  run_command("ruby exe/encrypt setup demo-password", "Unlock the vault with a password")
  
  puts "3. Add some secrets".bold.blue
  run_command("ruby exe/encrypt set API_KEY=sk-1234567890abcdef", "Add an API key")
  run_command("ruby exe/encrypt set DB_URL=postgres://user:pass@localhost:5432/mydb", "Add a database URL")
  run_command("ruby exe/encrypt set JWT_SECRET=super-secret-jwt-key", "Add a JWT secret")
  
  puts "4. Check vault status".bold.blue
  run_command("ruby exe/encrypt status", "View vault status and available keys")
  
  puts "5. Retrieve secrets".bold.blue
  run_command("ruby exe/encrypt get API_KEY", "Get the API key")
  run_command("ruby exe/encrypt get DB_URL", "Get the database URL")
  
  puts "6. Test runtime SDK".bold.blue
  puts "📝 Using the encrypt package in Ruby code".yellow
  puts "$ ruby example.rb".gray
  begin
    result = `ruby example.rb 2>&1`
    if result.strip.length > 0
      puts result.strip.green
    end
  rescue StandardError => e
    puts "Error: #{e.message}".red
  end
  puts
  
  puts "7. Create .env file".bold.blue
  run_command("ruby exe/encrypt unlock", "Export all secrets to .env file")
  
  puts "8. Lock the vault".bold.blue
  run_command("ruby exe/encrypt lockup demo-password", "Encrypt and lock all secrets")
  
  puts "9. Try to access locked secrets".bold.blue
  run_command("ruby exe/encrypt get API_KEY", "Attempt to get secret from locked vault")
  
  puts "10. Unlock the vault again".bold.blue
  run_command("ruby exe/encrypt setup demo-password", "Unlock the vault with the password")
  
  puts "11. Verify secrets are accessible again".bold.blue
  run_command("ruby exe/encrypt get API_KEY", "Get the API key after unlocking")
  
  puts "✅ Demo completed successfully!".bold.green
  puts "\nThe encrypt Ruby tool provides:".gray
  puts "• Secure encryption with triple-layer protection".gray
  puts "• Easy team onboarding with single command".gray
  puts "• Runtime SDK for in-code secret access".gray
  puts "• Git-safe encrypted storage".gray
  puts "• Beautiful CLI with colorized output".gray
  puts "• Full Ruby gem with proper structure".gray
end

main
