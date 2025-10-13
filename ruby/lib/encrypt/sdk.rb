# frozen_string_literal: true

module Encrypt
  # Runtime SDK for the Encrypt tool.
  # Provides easy-to-use functions for accessing secrets in Ruby code.
  module SDK
    class << self
      # Global vault instance
      @global_vault = nil

      # Get or create the global vault instance
      def global_vault
        @global_vault ||= Vault.new
      end

      # Auto-unlock vault if password is available
      def auto_unlock(password = nil)
        vault = global_vault

        # If already unlocked, no need to do anything
        return if vault.unlocked?

        # Try provided password first
        if password
          begin
            vault.setup(password)
            return
          rescue StandardError
            # Password might be wrong, continue to other methods
          end
        end

        # Try environment variable
        env_password = ENV["ENCRYPT_PASSWORD"]
        if env_password
          begin
            vault.setup(env_password)
            return
          rescue StandardError
            # Environment password might be wrong, continue to other methods
          end
        end

        # In development, we can be more lenient
        if ENV["NODE_ENV"] != "production"
          # Try common development passwords
          dev_passwords = %w[dev development test password 123456]
          dev_passwords.each do |dev_password|
            begin
              vault.setup(dev_password)
              return
            rescue StandardError
              # Continue to next password
            end
          end
        end

        # If we get here, we couldn't unlock the vault
        raise "Vault is locked and no valid password found. Set ENCRYPT_PASSWORD environment variable or provide password parameter."
      end

      # Get a secret value from the vault with auto-unlock support
      def get(key, password = nil)
        auto_unlock(password)
        vault = global_vault
        vault.get(key)
      end

      # Set a secret value in the vault with auto-unlock support
      def set(key, value, password = nil)
        auto_unlock(password)
        vault = global_vault
        vault.set(key, value)
      end

      # Get all secrets from the vault with auto-unlock support
      def all_secrets(password = nil)
        auto_unlock(password)
        vault = global_vault
        vault.all
      end

      # Get vault status
      def status
        vault = global_vault
        vault.status
      end

      # Check if vault is unlocked
      def is_unlocked?
        vault = global_vault
        vault.unlocked?
      end

      # Auto setup helper - checks if vault is locked and provides helpful error
      def auto_setup
        vault = global_vault

        raise "Vault not found. Run 'encrypt init' first." unless vault.exists?

        raise "Vault is locked. Run 'encrypt setup <password>' to unlock secrets." unless vault.unlocked?
      end

      # Get a secret with automatic environment variable support
      # This is the recommended function for production use
      def get_secret(key)
        get(key)
      end

      # Set a secret with automatic environment variable support
      # This is the recommended function for production use
      def set_secret(key, value)
        set(key, value)
      end

      # Get all secrets with automatic environment variable support
      # This is the recommended function for production use
      def get_all_secrets
        all_secrets
      end
    end
  end
end
