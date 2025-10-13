#!/usr/bin/env ruby
# frozen_string_literal: true

require_relative "lib/encrypt"

def main
  puts "💻 Example usage of the encrypt package in Ruby code"
  puts "=" * 50
  
  begin
    # Check if vault is unlocked
    unless Encrypt.is_vault_unlocked?
      puts "❌ Vault is locked. Run 'encrypt setup <password>' to unlock secrets."
      return
    end
    
    # Get secrets
    begin
      api_key = Encrypt.get_secret('API_KEY')
      db_url = Encrypt.get_secret('DB_URL')
      
      puts "✅ Secrets retrieved successfully!"
      puts "API Key: ***#{api_key[-4..-1] if api_key}"
      puts "DB URL: ***#{db_url[-10..-1] if db_url}"
      
    rescue StandardError => e
      puts "⚠️ Some secrets not found: #{e.message}"
    end
    
    # Get all secrets
    all_secrets = Encrypt.get_all_secrets
    puts "Available keys: #{all_secrets.keys}"
    
    # Get status
    status = Encrypt.get_status
    puts "Vault status: #{status[:is_locked] ? 'Locked' : 'Unlocked'}"
    puts "Number of keys: #{status[:keys].length}"
    
  rescue StandardError => e
    puts "Error: #{e.message}"
  end
end

main
