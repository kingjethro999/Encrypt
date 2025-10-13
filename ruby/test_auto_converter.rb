#!/usr/bin/env ruby
# frozen_string_literal: true

require_relative "lib/encrypt"

def test_auto_converter
  puts "🧪 Testing Auto-Converter Engine (Ruby)"
  puts "=" * 40
  puts
  
  begin
    # Test 1: Try to get secret when vault is locked (should fail without password)
    puts "1. Testing locked vault without password..."
    begin
      api_key = Encrypt::SDK.get('API_KEY')
      puts "❌ Unexpected: Got secret from locked vault"
    rescue StandardError => e
      puts "✅ Expected: Vault is locked - #{e.message}"
    end
    
    # Test 2: Try with environment variable
    puts "\n2. Testing with ENCRYPT_PASSWORD environment variable..."
    ENV['ENCRYPT_PASSWORD'] = 'mypassword'
    
    begin
      api_key = Encrypt::SDK.get('API_KEY')
      puts "✅ Success: Got secret with environment password - #{api_key}"
    rescue StandardError => e
      puts "❌ Failed: Could not get secret with environment password - #{e.message}"
    end
    
    # Test 3: Try with explicit password parameter
    puts "\n3. Testing with explicit password parameter..."
    ENV.delete('ENCRYPT_PASSWORD')  # Clear env var
    
    begin
      api_key = Encrypt::SDK.get('API_KEY', 'mypassword')
      puts "✅ Success: Got secret with explicit password - #{api_key}"
    rescue StandardError => e
      puts "❌ Failed: Could not get secret with explicit password - #{e.message}"
    end
    
    # Test 4: Test production-ready functions
    puts "\n4. Testing production-ready functions..."
    ENV['ENCRYPT_PASSWORD'] = 'mypassword'
    
    begin
      api_key = Encrypt::SDK.get_secret('API_KEY')
      db_url = Encrypt::SDK.get_secret('DB_URL')
      puts "✅ Success: get_secret() works - #{api_key}"
      puts "✅ Success: get_secret() works - #{db_url}"
    rescue StandardError => e
      puts "❌ Failed: get_secret() failed - #{e.message}"
    end
    
    # Test 5: Test development mode (should try common passwords)
    puts "\n5. Testing development mode..."
    ENV.delete('ENCRYPT_PASSWORD')
    ENV['NODE_ENV'] = 'development'
    
    begin
      api_key = Encrypt::SDK.get('API_KEY')
      puts "✅ Success: Development mode worked - #{api_key}"
    rescue StandardError => e
      puts "❌ Failed: Development mode failed - #{e.message}"
    end
    
    puts "\n🎉 Auto-converter tests completed!"
    
  rescue StandardError => e
    puts "❌ Test failed: #{e.message}"
  end
end

test_auto_converter
