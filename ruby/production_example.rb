#!/usr/bin/env ruby
# frozen_string_literal: true

require_relative "lib/encrypt"

def main
  puts "🚀 Production Usage Example (Ruby)"
  puts "=" * 35
  puts
  
  # Method 1: Using environment variable (Recommended for production)
  puts "Method 1: Environment Variable (Recommended)"
  puts "Set ENCRYPT_PASSWORD=your-password in your environment"
  puts
  
  begin
    # This will automatically unlock the vault using ENCRYPT_PASSWORD
    api_key = Encrypt.get_secret('API_KEY')
    db_url = Encrypt.get_secret('DB_URL')
    
    puts "✅ Successfully retrieved secrets:"
    puts "API Key: ***#{api_key[-4..-1] if api_key}"
    puts "DB URL: ***#{db_url[-10..-1] if db_url}"
    
    # Use in your application
    config = {
      api_key: api_key,
      database: db_url,
      port: ENV['PORT'] || 3000
    }
    
    puts "\n📋 Application config ready:"
    puts "  api_key: ***#{config[:api_key][-4..-1] if config[:api_key]}"
    puts "  database: ***#{config[:database][-10..-1] if config[:database]}"
    puts "  port: #{config[:port]}"
    
  rescue StandardError => e
    puts "❌ Error: #{e.message}"
    puts "\n💡 To fix this:"
    puts "1. Set ENCRYPT_PASSWORD environment variable"
    puts "2. Or provide password as second parameter"
    puts "3. Or run 'encrypt setup <password>' first"
  end
  
  puts "\n" + "=" * 50
  puts "Method 2: Explicit Password Parameter"
  puts "=" * 50
  
  begin
    # This will use the provided password to unlock the vault
    api_key = Encrypt::SDK.get('API_KEY', 'mypassword')
    puts "✅ Success with explicit password: ***#{api_key[-4..-1] if api_key}"
  rescue StandardError => e
    puts "❌ Error with explicit password: #{e.message}"
  end
  
  puts "\n" + "=" * 50
  puts "Method 3: Development Mode"
  puts "=" * 50
  
  # In development, you can set NODE_ENV=development
  # and it will try common passwords automatically
  ENV['NODE_ENV'] = 'development'
  
  begin
    api_key = Encrypt::SDK.get('API_KEY')
    puts "✅ Development mode success: ***#{api_key[-4..-1] if api_key}"
  rescue StandardError => e
    puts "❌ Development mode failed: #{e.message}"
  end
  
  puts "\n🎯 Production Deployment Examples:"
  puts "=" * 35
  puts "Docker:"
  puts "  ENV ENCRYPT_PASSWORD=your-production-password"
  puts ""
  puts "Kubernetes:"
  puts "  env:"
  puts "  - name: ENCRYPT_PASSWORD"
  puts "    valueFrom:"
  puts "      secretKeyRef:"
  puts "        name: encrypt-secrets"
  puts "        key: password"
  puts ""
  puts "Heroku:"
  puts "  heroku config:set ENCRYPT_PASSWORD=your-password"
  puts ""
  puts "AWS Lambda:"
  puts "  Set ENCRYPT_PASSWORD in environment variables"
  puts ""
  puts "Rails application.rb:"
  puts "  require 'encrypt'"
  puts "  config.api_key = Encrypt.get_secret('API_KEY')"
  puts "  config.database_url = Encrypt.get_secret('DATABASE_URL')"
end

main
