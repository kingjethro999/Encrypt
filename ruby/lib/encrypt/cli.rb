# frozen_string_literal: true

require "thor"
require "colorize"
require "tty-spinner"

module Encrypt
  # CLI interface for the Encrypt tool.
  # Provides command-line access to all vault operations with beautiful formatting.
  class CLI < Thor
    desc "init", "Create .encrypt vault"
    def init
      spinner = TTY::Spinner.new("Initializing vault...", format: :dots)
      spinner.auto_spin

      begin
        vault = Vault.new
        vault.init
        spinner.success("Vault initialized successfully!")
        puts "📁 Created .encrypt directory with secure configuration.".blue
      rescue StandardError => e
        spinner.error("Failed to initialize vault")
        puts "❌ #{e.message}".red
        exit 1
      end
    end

    desc "lockup PASSWORD", "Encrypt and secure secrets with password"
    def lockup(password)
      spinner = TTY::Spinner.new("Locking up secrets...", format: :dots)
      spinner.auto_spin

      begin
        vault = Vault.new

        unless vault.exists?
          spinner.error("Vault not found. Run 'encrypt init' first.")
          exit 1
        end

        vault.lockup(password)
        spinner.success("Secrets locked successfully!")
        puts "🔒 Your secrets are now encrypted and secure.".blue
      rescue StandardError => e
        spinner.error("Failed to lock secrets")
        puts "❌ #{e.message}".red
        exit 1
      end
    end

    desc "setup PASSWORD", "Set up secrets on a new machine"
    def setup(password)
      spinner = TTY::Spinner.new("Setting up vault...", format: :dots)
      spinner.auto_spin

      begin
        vault = Vault.new

        unless vault.exists?
          spinner.error("Vault not found. Run 'encrypt init' first.")
          exit 1
        end

        vault.setup(password)
        spinner.success("Vault unlocked successfully!")
        puts "🔓 Your secrets are now available for use.".blue
      rescue StandardError => e
        spinner.error("Failed to setup vault")
        puts "❌ #{e.message}".red
        exit 1
      end
    end

    desc "set KEY=VALUE", "Add/update a key"
    def set(key_value)
      unless key_value.include?("=")
        puts "❌ Invalid format. Use: encrypt set KEY=value".red
        exit 1
      end

      key, value = key_value.split("=", 2)
      spinner = TTY::Spinner.new("Setting secret: #{key}", format: :dots)
      spinner.auto_spin

      begin
        vault = Vault.new

        unless vault.exists?
          spinner.error("Vault not found. Run 'encrypt init' first.")
          exit 1
        end

        unless vault.unlocked?
          spinner.error("Vault is locked. Run 'encrypt setup <password>' to unlock secrets.")
          exit 1
        end

        vault.set(key, value)
        spinner.success("Secret '#{key}' set successfully!")
      rescue StandardError => e
        spinner.error("Failed to set secret")
        puts "❌ #{e.message}".red
        exit 1
      end
    end

    desc "get KEY", "Fetch decrypted value"
    def get(key)
      begin
        vault = Vault.new

        unless vault.exists?
          puts "❌ Vault not found. Run 'encrypt init' first.".red
          exit 1
        end

        unless vault.unlocked?
          puts "❌ Vault is locked. Run 'encrypt setup <password>' to unlock secrets.".red
          exit 1
        end

        value = vault.get(key)
        puts value
      rescue StandardError => e
        puts "❌ #{e.message}".red
        exit 1
      end
    end

    desc "unlock", "Decrypt everything into .env"
    def unlock
      spinner = TTY::Spinner.new("Unlocking secrets...", format: :dots)
      spinner.auto_spin

      begin
        vault = Vault.new

        unless vault.exists?
          spinner.error("Vault not found. Run 'encrypt init' first.")
          exit 1
        end

        unless vault.unlocked?
          spinner.error("Vault is locked. Run 'encrypt setup <password>' to unlock secrets.")
          exit 1
        end

        secrets = vault.all
        env_content = secrets.map { |k, v| "#{k}=#{v}" }.join("\n")

        File.write(".env", env_content)
        spinner.success("Secrets unlocked and written to .env")
      rescue StandardError => e
        spinner.error("Failed to unlock secrets")
        puts "❌ #{e.message}".red
        exit 1
      end
    end

    desc "status", "Check if vault is locked, list keys"
    def status
      begin
        vault = Vault.new

        unless vault.exists?
          puts "⚠️ Vault not found. Run 'encrypt init' first.".yellow
          return
        end

        status_info = vault.status

        puts "📊 Vault Status:".blue
        puts "┌─────────────────┬──────────────────────────────┐"
        puts "│ Property        │ Value                        │"
        puts "├─────────────────┼──────────────────────────────┤"

        status_icon = status_info[:is_locked] ? "🔒" : "🔓"
        status_text = status_info[:is_locked] ? "Locked" : "Unlocked"
        status_color = status_info[:is_locked] ? "red" : "green"

        puts "│ Status          │ #{status_icon} #{status_text.send(status_color.to_sym).ljust(25)} │"
        puts "│ Keys            │ #{status_info[:keys].length.to_s.ljust(25)} │"

        if status_info[:last_modified]
          puts "│ Last Modified   │ #{status_info[:last_modified].ljust(25)} │"
        end

        puts "└─────────────────┴──────────────────────────────┘"

        if status_info[:keys].any?
          puts "\n🔑 Available keys:".blue
          status_info[:keys].each do |key|
            puts "  • #{key}"
          end
        end
      rescue StandardError => e
        puts "❌ #{e.message}".red
        exit 1
      end
    end

    desc "reset", "Remove vault (careful!)"
    def reset
      print "Are you sure you want to reset the vault? This will delete all encrypted secrets. (y/N): "
      response = STDIN.gets.chomp.downcase

      unless response == "y" || response == "yes"
        puts "⚠️ Reset cancelled.".yellow
        return
      end

      spinner = TTY::Spinner.new("Resetting vault...", format: :dots)
      spinner.auto_spin

      begin
        vault = Vault.new
        vault.reset
        spinner.success("Vault reset successfully!")
      rescue StandardError => e
        spinner.error("Failed to reset vault")
        puts "❌ #{e.message}".red
        exit 1
      end
    end

    def self.exit_on_failure?
      true
    end
  end
end
