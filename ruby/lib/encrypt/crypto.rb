# frozen_string_literal: true

require "openssl"
require "securerandom"

module Encrypt
  # Triple-layer encryption implementation for the Encrypt tool.
  # Provides AES-256-CBC encryption with PBKDF2 key derivation and HMAC signatures.
  class Crypto
    ALGORITHM = "AES-256-CBC"
    KEY_LENGTH = 32
    IV_LENGTH = 16
    SALT_LENGTH = 32
    ITERATIONS = 100_000

    class << self
      # Generate a random salt for key derivation
      def generate_salt
        SecureRandom.hex(SALT_LENGTH)
      end

      # Derive encryption key from password using PBKDF2
      def derive_key(password, salt)
        OpenSSL::PKCS5.pbkdf2_hmac(
          password,
          [salt].pack("H*"),
          ITERATIONS,
          KEY_LENGTH,
          "SHA256"
        )
      end

      # Generate HMAC signature for data integrity
      def generate_hmac(data, key)
        OpenSSL::HMAC.hexdigest("SHA256", key, data)
      end

      # Triple-layer encryption:
      # 1. Generate random salt
      # 2. Derive key using PBKDF2
      # 3. AES-256-CBC encryption with HMAC
      def encrypt(plaintext, password)
        # Phase 1: Generate salt
        salt = generate_salt

        # Phase 2: Derive key from password
        key = derive_key(password, salt)

        # Phase 3: AES-256-CBC encryption
        cipher = OpenSSL::Cipher.new(ALGORITHM)
        cipher.encrypt
        cipher.key = key
        iv = cipher.random_iv

        encrypted = cipher.update(plaintext) + cipher.final

        # Combine IV and encrypted data
        encrypted_data = iv.unpack1("H*") + ":" + encrypted.unpack1("H*")

        # Phase 3: Generate HMAC signature
        hmac = generate_hmac(encrypted_data, key)

        {
          encrypted: encrypted_data,
          salt: salt,
          hmac: hmac
        }
      end

      # Triple-layer decryption:
      # 1. Verify HMAC signature
      # 2. Derive key using PBKDF2
      # 3. AES-256-CBC decryption
      def decrypt(encrypted_data, salt, hmac, password)
        # Phase 2: Derive key from password
        key = derive_key(password, salt)

        # Phase 3: Verify HMAC signature
        expected_hmac = generate_hmac(encrypted_data, key)
        return ["", false] unless expected_hmac == hmac

        # Phase 3: AES-256-CBC decryption
        iv_hex, encrypted_hex = encrypted_data.split(":", 2)
        iv = [iv_hex].pack("H*")
        encrypted = [encrypted_hex].pack("H*")

        decipher = OpenSSL::Cipher.new(ALGORITHM)
        decipher.decrypt
        decipher.key = key
        decipher.iv = iv

        decrypted = decipher.update(encrypted) + decipher.final

        [decrypted, true]
      rescue StandardError
        ["", false]
      end

      # Hash password for storage using PBKDF2
      def hash_password(password)
        salt = generate_salt
        salt_bytes = [salt].pack("H*")
        hash_bytes = OpenSSL::PKCS5.pbkdf2_hmac(
          password,
          salt_bytes,
          ITERATIONS,
          64,
          "SHA256"
        )
        salt + ":" + hash_bytes.unpack1("H*")
      end

      # Verify password against stored hash
      def verify_password(password, stored_hash)
        salt, stored_hash_hex = stored_hash.split(":", 2)
        salt_bytes = [salt].pack("H*")
        computed_hash = OpenSSL::PKCS5.pbkdf2_hmac(
          password,
          salt_bytes,
          ITERATIONS,
          64,
          "SHA256"
        )
        computed_hash.unpack1("H*") == stored_hash_hex
      rescue StandardError
        false
      end
    end
  end
end
