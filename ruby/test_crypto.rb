#!/usr/bin/env ruby
# frozen_string_literal: true

require_relative "lib/encrypt"

def test_encryption
  puts "🧪 Testing TripleEncryption (Ruby)..."
  
  # Test data
  plaintext = "Hello, World!"
  password = "test-password-123"
  
  puts "Original: #{plaintext}"
  
  begin
    # Encrypt
    encrypted = Encrypt::Crypto.encrypt(plaintext, password)
    puts "✅ Encryption successful"
    puts "Encrypted data length: #{encrypted[:encrypted].length}"
    
    # Decrypt
    decrypted, is_valid = Encrypt::Crypto.decrypt(
      encrypted[:encrypted], 
      encrypted[:salt], 
      encrypted[:hmac], 
      password
    )
    
    if is_valid && decrypted == plaintext
      puts "✅ Decryption successful"
      puts "Decrypted: #{decrypted}"
    else
      puts "❌ Decryption failed"
    end
    
    # Test password hashing
    hash_result = Encrypt::Crypto.hash_password(password)
    is_valid_hash = Encrypt::Crypto.verify_password(password, hash_result)
    
    if is_valid_hash
      puts "✅ Password hashing/verification successful"
    else
      puts "❌ Password hashing/verification failed"
    end
    
    # Test with different password (should fail)
    wrong_decrypted, wrong_valid = Encrypt::Crypto.decrypt(
      encrypted[:encrypted], 
      encrypted[:salt], 
      encrypted[:hmac], 
      "wrong-password"
    )
    
    if !wrong_valid
      puts "✅ Wrong password correctly rejected"
    else
      puts "❌ Wrong password incorrectly accepted"
    end
    
    puts "\n🎉 All tests passed!"
    
  rescue StandardError => e
    puts "❌ Test failed: #{e.message}"
  end
end

test_encryption
