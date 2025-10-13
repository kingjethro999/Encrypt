# frozen_string_literal: true

require_relative "encrypt/version"
require_relative "encrypt/crypto"
require_relative "encrypt/vault"
require_relative "encrypt/sdk"

module Encrypt
  class Error < StandardError; end
  
  # Global vault instance for SDK usage
  @global_vault = nil

  def self.global_vault
    @global_vault ||= Vault.new
  end

  def self.global_vault=(vault)
    @global_vault = vault
  end

  # Convenience methods for backward compatibility
  def self.get_secret(key)
    SDK.get_secret(key)
  end

  def self.set_secret(key, value)
    SDK.set_secret(key, value)
  end

  def self.get_all_secrets
    SDK.get_all_secrets
  end

  def self.get_status
    SDK.status
  end

  def self.is_vault_unlocked?
    SDK.is_unlocked?
  end

  def self.auto_setup_vault
    SDK.auto_setup
  end
end
