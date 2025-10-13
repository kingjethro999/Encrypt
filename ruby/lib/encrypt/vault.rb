# frozen_string_literal: true

require "json"
require "fileutils"
require "time"

module Encrypt
  # Vault management for the Encrypt tool.
  # Handles file operations, state management, and secret storage/retrieval.
  class Vault
    VAULT_DIR = ".encrypt"
    CONFIG_FILE = "vault.lock"
    SECRETS_FILE = "secrets.enc.json"
    GITIGNORE_FILE = ".gitignore"
    LOCK_FILE = "vault.unlocked"

    def initialize(project_root = nil)
      @project_root = project_root || Dir.pwd
      @vault_path = File.join(@project_root, VAULT_DIR)
      @config_path = File.join(@vault_path, CONFIG_FILE)
      @secrets_path = File.join(@vault_path, SECRETS_FILE)
      @gitignore_path = File.join(@project_root, GITIGNORE_FILE)
      @lock_file_path = File.join(@vault_path, LOCK_FILE)
      @memory_cache = {}
    end

    # Initialize the vault directory structure
    def init
      FileUtils.mkdir_p(@vault_path) unless Dir.exist?(@vault_path)

      # Create initial empty secrets file
      unless File.exist?(@secrets_path)
        File.write(@secrets_path, JSON.generate({}))
      end

      # Create initial config file
      unless File.exist?(@config_path)
        initial_config = {
          password_hash: "",
          salt: "",
          hmac: "",
          created_at: Time.now.iso8601,
          version: "1.0.0"
        }
        File.write(@config_path, JSON.generate(initial_config))
      end

      # Update .gitignore
      update_gitignore
    end

    # Update .gitignore to exclude .encrypt directory
    def update_gitignore
      gitignore_content = File.exist?(@gitignore_path) ? File.read(@gitignore_path) : ""

      unless gitignore_content.include?(".encrypt/")
        File.open(@gitignore_path, "a") do |f|
          f.puts "\n# Encrypt vault\n.encrypt/"
        end
      end
    end

    # Check if vault exists
    def exists?
      Dir.exist?(@vault_path) && File.exist?(@config_path)
    end

    # Lock up secrets with password
    def lockup(password)
      # Load secrets from file if not in memory
      load_secrets_from_file if @memory_cache.empty?

      raise "No secrets to lock. Use 'encrypt set' to add secrets first." if @memory_cache.empty?

      # Encrypt all secrets
      encrypted_secrets = {}

      @memory_cache.each do |key, value|
        result = Crypto.encrypt(value, password)
        encrypted_secrets[key] = JSON.generate(result)
      end

      # Create vault config
      config_salt = Crypto.generate_salt
      config = {
        password_hash: Crypto.hash_password(password),
        salt: config_salt,
        hmac: Crypto.generate_hmac(
          JSON.generate(encrypted_secrets),
          Crypto.derive_key(password, config_salt)
        ),
        created_at: Time.now.iso8601,
        version: "1.0.0"
      }

      # Write encrypted secrets and config
      File.write(@secrets_path, JSON.generate(encrypted_secrets))
      File.write(@config_path, JSON.generate(config))

      # Clear memory cache and remove lock file
      @memory_cache.clear
      File.delete(@lock_file_path) if File.exist?(@lock_file_path)
    end

    # Setup/unlock vault with password
    def setup(password)
      raise "Vault not found. Run 'encrypt init' first." unless exists?

      config = JSON.parse(File.read(@config_path))

      # If this is a fresh vault (no password set), just unlock it
      if config["password_hash"].empty?
        @memory_cache.clear
        create_lock_file
        return
      end

      # Verify password
      raise "Invalid password." unless Crypto.verify_password(password, config["password_hash"])

      # Load and decrypt secrets
      encrypted_secrets = JSON.parse(File.read(@secrets_path))

      @memory_cache.clear

      encrypted_secrets.each do |key, encrypted_data|
        # Check if the data is already decrypted (plain text) or encrypted
        if encrypted_data.is_a?(String) && encrypted_data.start_with?("{")
          # This is encrypted data stored as JSON string, decrypt it
          result = JSON.parse(encrypted_data)
          decrypted, is_valid = Crypto.decrypt(
            result["encrypted"], result["salt"], result["hmac"], password
          )

          raise "Failed to decrypt secret: #{key}" unless is_valid

          @memory_cache[key] = decrypted
        else
          # This is plain text data
          @memory_cache[key] = encrypted_data
        end
      end

      # Save decrypted secrets to file for easy access
      save_secrets_to_file
      create_lock_file
    end

    # Set a secret (only works when unlocked)
    def set(key, value)
      raise "Vault is locked. Run 'encrypt setup <password>' to unlock secrets." unless unlocked?

      # Load secrets from file if not in memory
      load_secrets_from_file if @memory_cache.empty?

      @memory_cache[key] = value
      save_secrets_to_file
    end

    # Get a secret (only works when unlocked)
    def get(key)
      raise "Vault is locked. Run 'encrypt setup <password>' to unlock secrets." unless unlocked?

      # Load secrets from file if not in memory
      load_secrets_from_file if @memory_cache.empty?

      raise "Secret \"#{key}\" not found." unless @memory_cache.key?(key)

      @memory_cache[key]
    end

    # Get all secrets (only works when unlocked)
    def all
      raise "Vault is locked. Run 'encrypt setup <password>' to unlock secrets." unless unlocked?

      # Load secrets from file if not in memory
      load_secrets_from_file if @memory_cache.empty?

      @memory_cache.dup
    end

    # Get vault status
    def status
      # Check if vault is unlocked by looking for lock file
      is_unlocked = File.exist?(@lock_file_path)

      # Load secrets from file if unlocked and not in memory
      if is_unlocked && @memory_cache.empty?
        load_secrets_from_file
      end

      keys = @memory_cache.keys
      last_modified = nil
      if exists?
        last_modified = File.mtime(@config_path).iso8601
      end

      {
        is_locked: !is_unlocked,
        keys: keys,
        last_modified: last_modified
      }
    end

    # Reset/remove vault
    def reset
      FileUtils.rm_rf(@vault_path) if Dir.exist?(@vault_path)
      @memory_cache.clear
    end

    # Check if vault is unlocked
    def unlocked?
      File.exist?(@lock_file_path)
    end

    private

    # Create lock file to indicate vault is unlocked
    def create_lock_file
      lock_data = {
        unlocked: true,
        timestamp: Time.now.iso8601
      }
      File.write(@lock_file_path, JSON.generate(lock_data))
    end

    # Load secrets from file (for unlocked vault)
    def load_secrets_from_file
      return unless File.exist?(@secrets_path)

      secrets = JSON.parse(File.read(@secrets_path))

      @memory_cache.clear
      secrets.each do |key, value|
        # Check if the value is encrypted (starts with {) or plain text
        if value.is_a?(String) && value.start_with?("{")
          # This is encrypted data, we need to decrypt it
          # But we don't have the password here, so we can't decrypt
          # This should not happen in an unlocked vault
          raise "Secret \"#{key}\" is encrypted but vault is unlocked. This should not happen."
        else
          # This is plain text data
          @memory_cache[key] = value
        end
      end
    end

    # Save secrets to file (for unlocked vault)
    def save_secrets_to_file
      File.write(@secrets_path, JSON.generate(@memory_cache))
    end
  end
end
