use crate::error::{EncryptError, Result};
use crate::vault::Vault;
use clap::{Parser, Subcommand};
use colored::*;
use indicatif::{ProgressBar, ProgressStyle};
use std::process;

#[derive(Parser)]
#[command(name = "encrypt")]
#[command(about = "A top-level secrets orchestrator. Not just another .env tool — this one encrypts, locks, and sets you up for secure local and team dev.")]
pub struct Cli {
    #[command(subcommand)]
    pub command: Commands,
}

#[derive(Subcommand)]
pub enum Commands {
    /// Create .encrypt vault
    Init,
    /// Encrypt and secure secrets with password
    Lockup {
        password: String,
    },
    /// Set up secrets on a new machine
    Setup {
        password: String,
    },
    /// Add/update a key
    Set {
        key_value: String,
    },
    /// Fetch decrypted value
    Get {
        key: String,
    },
    /// Decrypt everything into .env
    Unlock,
    /// Check if vault is locked, list keys
    Status,
    /// Remove vault (careful!)
    Reset,
}

pub fn run() {
    let cli = Cli::parse();

    let result = match cli.command {
        Commands::Init => init_command(),
        Commands::Lockup { password } => lockup_command(&password),
        Commands::Setup { password } => setup_command(&password),
        Commands::Set { key_value } => set_command(&key_value),
        Commands::Get { key } => get_command(&key),
        Commands::Unlock => unlock_command(),
        Commands::Status => status_command(),
        Commands::Reset => reset_command(),
    };

    if let Err(e) = result {
        eprintln!("❌ {}", e.to_string().red());
        process::exit(1);
    }
}

fn init_command() -> Result<()> {
    let pb = ProgressBar::new_spinner();
    pb.set_style(
        ProgressStyle::default_spinner()
            .template("{spinner:.green} {msg}")
            .unwrap()
            .tick_strings(&["⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"]),
    );
    pb.set_message("Initializing vault...");
    pb.enable_steady_tick(std::time::Duration::from_millis(100));

    let vault = Vault::new();
    vault.init()?;
    
    pb.finish_with_message("✅ Vault initialized successfully!");
    println!("📁 {}", "Created .encrypt directory with secure configuration.".blue());
    
    Ok(())
}

fn lockup_command(password: &str) -> Result<()> {
    let pb = ProgressBar::new_spinner();
    pb.set_style(
        ProgressStyle::default_spinner()
            .template("{spinner:.green} {msg}")
            .unwrap()
            .tick_strings(&["⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"]),
    );
    pb.set_message("Locking up secrets...");
    pb.enable_steady_tick(std::time::Duration::from_millis(100));

    let mut vault = Vault::new();

    if !vault.exists() {
        pb.finish_with_message("❌ Vault not found. Run 'encrypt init' first.");
        return Err(EncryptError::VaultNotFound);
    }

    vault.lockup(password)?;
    
    pb.finish_with_message("✅ Secrets locked successfully!");
    println!("🔒 {}", "Your secrets are now encrypted and secure.".blue());
    
    Ok(())
}

fn setup_command(password: &str) -> Result<()> {
    let pb = ProgressBar::new_spinner();
    pb.set_style(
        ProgressStyle::default_spinner()
            .template("{spinner:.green} {msg}")
            .unwrap()
            .tick_strings(&["⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"]),
    );
    pb.set_message("Setting up vault...");
    pb.enable_steady_tick(std::time::Duration::from_millis(100));

    let mut vault = Vault::new();

    if !vault.exists() {
        pb.finish_with_message("❌ Vault not found. Run 'encrypt init' first.");
        return Err(EncryptError::VaultNotFound);
    }

    vault.setup(password)?;
    
    pb.finish_with_message("✅ Vault unlocked successfully!");
    println!("🔓 {}", "Your secrets are now available for use.".blue());
    
    Ok(())
}

fn set_command(key_value: &str) -> Result<()> {
    let parts: Vec<&str> = key_value.splitn(2, '=').collect();
    if parts.len() != 2 {
        eprintln!("❌ {}", "Invalid format. Use: encrypt set KEY=value".red());
        process::exit(1);
    }

    let key = parts[0];
    let value = parts[1];

    let pb = ProgressBar::new_spinner();
    pb.set_style(
        ProgressStyle::default_spinner()
            .template("{spinner:.green} {msg}")
            .unwrap()
            .tick_strings(&["⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"]),
    );
    pb.set_message(format!("Setting secret: {}", key));
    pb.enable_steady_tick(std::time::Duration::from_millis(100));

    let mut vault = Vault::new();

    if !vault.exists() {
        pb.finish_with_message("❌ Vault not found. Run 'encrypt init' first.");
        return Err(EncryptError::VaultNotFound);
    }

    if !vault.unlocked()? {
        pb.finish_with_message("❌ Vault is locked. Run 'encrypt setup <password>' to unlock secrets.");
        return Err(EncryptError::VaultLocked);
    }

    vault.set(key, value)?;
    
    pb.finish_with_message(format!("✅ Secret '{}' set successfully!", key));
    
    Ok(())
}

fn get_command(key: &str) -> Result<()> {
    let mut vault = Vault::new();

    if !vault.exists() {
        eprintln!("❌ {}", "Vault not found. Run 'encrypt init' first.".red());
        process::exit(1);
    }

    if !vault.unlocked()? {
        eprintln!("❌ {}", "Vault is locked. Run 'encrypt setup <password>' to unlock secrets.".red());
        process::exit(1);
    }

    let value = vault.get(key)?;
    println!("{}", value);
    
    Ok(())
}

fn unlock_command() -> Result<()> {
    let pb = ProgressBar::new_spinner();
    pb.set_style(
        ProgressStyle::default_spinner()
            .template("{spinner:.green} {msg}")
            .unwrap()
            .tick_strings(&["⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"]),
    );
    pb.set_message("Unlocking secrets...");
    pb.enable_steady_tick(std::time::Duration::from_millis(100));

    let mut vault = Vault::new();

    if !vault.exists() {
        pb.finish_with_message("❌ Vault not found. Run 'encrypt init' first.");
        return Err(EncryptError::VaultNotFound);
    }

    if !vault.unlocked()? {
        pb.finish_with_message("❌ Vault is locked. Run 'encrypt setup <password>' to unlock secrets.");
        return Err(EncryptError::VaultLocked);
    }

    let secrets = vault.all()?;
    let env_content = secrets
        .iter()
        .map(|(k, v)| format!("{}={}", k, v))
        .collect::<Vec<_>>()
        .join("\n");

    std::fs::write(".env", env_content)?;
    
    pb.finish_with_message("✅ Secrets unlocked and written to .env");
    
    Ok(())
}

fn status_command() -> Result<()> {
    let mut vault = Vault::new();

    if !vault.exists() {
        println!("⚠️ {}", "Vault not found. Run 'encrypt init' first.".yellow());
        return Ok(());
    }

    let status = vault.status()?;

    println!("📊 {}", "Vault Status:".blue());
    println!("┌─────────────────┬──────────────────────────────┐");
    println!("│ Property        │ Value                        │");
    println!("├─────────────────┼──────────────────────────────┤");

    let status_icon = if status.is_locked { "🔒" } else { "🔓" };
    let status_text = if status.is_locked { "Locked" } else { "Unlocked" };
    let status_color = if status.is_locked { "red" } else { "green" };

    println!(
        "│ Status          │ {} {} │",
        status_icon,
        if status_color == "red" {
            status_text.red()
        } else {
            status_text.green()
        }
    );
    println!("│ Keys            │ {} │", status.keys.len());

    if let Some(last_modified) = status.last_modified {
        println!("│ Last Modified   │ {} │", last_modified.format("%Y-%m-%dT%H:%M:%S%.3fZ"));
    }

    println!("└─────────────────┴──────────────────────────────┘");

    if !status.keys.is_empty() {
        println!("\n🔑 {}", "Available keys:".blue());
        for key in &status.keys {
            println!("  • {}", key);
        }
    }

    Ok(())
}

fn reset_command() -> Result<()> {
    print!("Are you sure you want to reset the vault? This will delete all encrypted secrets. (y/N): ");
    use std::io::{self, Write};
    io::stdout().flush()?;

    let mut input = String::new();
    io::stdin().read_line(&mut input)?;
    let response = input.trim().to_lowercase();

    if response != "y" && response != "yes" {
        println!("⚠️ {}", "Reset cancelled.".yellow());
        return Ok(());
    }

    let pb = ProgressBar::new_spinner();
    pb.set_style(
        ProgressStyle::default_spinner()
            .template("{spinner:.green} {msg}")
            .unwrap()
            .tick_strings(&["⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"]),
    );
    pb.set_message("Resetting vault...");
    pb.enable_steady_tick(std::time::Duration::from_millis(100));

    let vault = Vault::new();
    vault.reset()?;
    
    pb.finish_with_message("✅ Vault reset successfully!");
    
    Ok(())
}
