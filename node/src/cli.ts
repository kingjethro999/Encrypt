#!/usr/bin/env node

import { Command } from 'commander';
import inquirer from 'inquirer';
import chalk from 'chalk';
import ora from 'ora';
import { Vault } from './vault';
import { TripleEncryption } from './crypto';

// Global vault instance
let globalVault: Vault | null = null;

function getVault(): Vault {
  if (!globalVault) {
    globalVault = new Vault();
  }
  return globalVault;
}

const program = new Command();

program
  .name('encrypt')
  .description('A top-level secrets orchestrator. Not just another .env tool — this one encrypts, locks, and sets you up for secure local and team dev.')
  .version('1.0.0');

// Initialize vault
program
  .command('init')
  .description('Create .encrypt vault')
  .action(() => {
    const spinner = ora('Initializing vault...').start();
    
    try {
      const vault = getVault();
      vault.init();
      spinner.succeed(chalk.green('Vault initialized successfully!'));
      console.log(chalk.blue('Created .encrypt directory with secure configuration.'));
    } catch (error) {
      spinner.fail(chalk.red('Failed to initialize vault'));
      console.error(chalk.red(error instanceof Error ? error.message : 'Unknown error'));
      process.exit(1);
    }
  });

// Lock up secrets
program
  .command('lockup <password>')
  .description('Encrypt and secure secrets with password')
  .action(async (password: string) => {
    const spinner = ora('Locking up secrets...').start();
    
    try {
      const vault = getVault();
      
      if (!vault.exists()) {
        spinner.fail(chalk.red('Vault not found. Run "encrypt init" first.'));
        process.exit(1);
      }
      
      vault.lockup(password);
      spinner.succeed(chalk.green('Secrets locked successfully!'));
      console.log(chalk.blue('Your secrets are now encrypted and secure.'));
    } catch (error) {
      spinner.fail(chalk.red('Failed to lock secrets'));
      console.error(chalk.red(error instanceof Error ? error.message : 'Unknown error'));
      process.exit(1);
    }
  });

// Setup/unlock vault
program
  .command('setup <password>')
  .description('Set up secrets on a new machine')
  .action(async (password: string) => {
    const spinner = ora('Setting up vault...').start();
    
    try {
      const vault = getVault();
      
      if (!vault.exists()) {
        spinner.fail(chalk.red('Vault not found. Run "encrypt init" first.'));
        process.exit(1);
      }
      
      vault.setup(password);
      spinner.succeed(chalk.green('Vault unlocked successfully!'));
      console.log(chalk.blue('Your secrets are now available for use.'));
    } catch (error) {
      spinner.fail(chalk.red('Failed to setup vault'));
      console.error(chalk.red(error instanceof Error ? error.message : 'Unknown error'));
      process.exit(1);
    }
  });

// Set a secret
program
  .command('set <key=value>')
  .description('Add/update a key')
  .action(async (keyValue: string) => {
    const [key, ...valueParts] = keyValue.split('=');
    const value = valueParts.join('=');
    
    if (!key || !value) {
      console.error(chalk.red('Invalid format. Use: encrypt set KEY=value'));
      process.exit(1);
    }
    
    const spinner = ora(`Setting secret: ${key}`).start();
    
    try {
      const vault = getVault();
      
      if (!vault.exists()) {
        spinner.fail(chalk.red('Vault not found. Run "encrypt init" first.'));
        process.exit(1);
      }
      
      if (!vault.isUnlockedStatus()) {
        spinner.fail(chalk.red('Vault is locked. Run "encrypt setup <password>" to unlock secrets.'));
        process.exit(1);
      }
      
      vault.set(key, value);
      spinner.succeed(chalk.green(`Secret "${key}" set successfully!`));
    } catch (error) {
      spinner.fail(chalk.red('Failed to set secret'));
      console.error(chalk.red(error instanceof Error ? error.message : 'Unknown error'));
      process.exit(1);
    }
  });

// Get a secret
program
  .command('get <key>')
  .description('Fetch decrypted value')
  .action(async (key: string) => {
    try {
      const vault = getVault();
      
      if (!vault.exists()) {
        console.error(chalk.red('Vault not found. Run "encrypt init" first.'));
        process.exit(1);
      }
      
      if (!vault.isUnlockedStatus()) {
        console.error(chalk.red('Vault is locked. Run "encrypt setup <password>" to unlock secrets.'));
        process.exit(1);
      }
      
      const value = vault.get(key);
      console.log(value);
    } catch (error) {
      console.error(chalk.red(error instanceof Error ? error.message : 'Unknown error'));
      process.exit(1);
    }
  });

// Unlock all secrets
program
  .command('unlock')
  .description('Decrypt everything into .env')
  .action(async () => {
    const spinner = ora('Unlocking secrets...').start();
    
    try {
      const vault = getVault();
      
      if (!vault.exists()) {
        spinner.fail(chalk.red('Vault not found. Run "encrypt init" first.'));
        process.exit(1);
      }
      
      if (!vault.isUnlockedStatus()) {
        spinner.fail(chalk.red('Vault is locked. Run "encrypt setup <password>" to unlock secrets.'));
        process.exit(1);
      }
      
      const secrets = vault.all();
      const envContent = Object.entries(secrets)
        .map(([key, value]) => `${key}=${value}`)
        .join('\n');
      
      require('fs').writeFileSync('.env', envContent);
      spinner.succeed(chalk.green('Secrets unlocked and written to .env'));
    } catch (error) {
      spinner.fail(chalk.red('Failed to unlock secrets'));
      console.error(chalk.red(error instanceof Error ? error.message : 'Unknown error'));
      process.exit(1);
    }
  });

// Check status
program
  .command('status')
  .description('Check if vault is locked, list keys')
  .action(async () => {
    try {
      const vault = getVault();
      
      if (!vault.exists()) {
        console.log(chalk.yellow('Vault not found. Run "encrypt init" first.'));
        return;
      }
      
      const status = vault.status();
      
      console.log(chalk.blue('Vault Status:'));
      console.log(`  Status: ${status.isLocked ? chalk.red('🔒 Locked') : chalk.green('🔓 Unlocked')}`);
      console.log(`  Keys: ${status.keys.length}`);
      
      if (status.keys.length > 0) {
        console.log(chalk.blue('  Available keys:'));
        status.keys.forEach(key => {
          console.log(`    - ${key}`);
        });
      }
      
      if (status.lastModified) {
        console.log(`  Last modified: ${status.lastModified}`);
      }
    } catch (error) {
      console.error(chalk.red(error instanceof Error ? error.message : 'Unknown error'));
      process.exit(1);
    }
  });

// Reset vault
program
  .command('reset')
  .description('Remove vault (careful!)')
  .action(async () => {
    const { confirm } = await inquirer.prompt([
      {
        type: 'confirm',
        name: 'confirm',
        message: 'Are you sure you want to reset the vault? This will delete all encrypted secrets.',
        default: false
      }
    ]);
    
    if (!confirm) {
      console.log(chalk.yellow('Reset cancelled.'));
      return;
    }
    
    const spinner = ora('Resetting vault...').start();
    
    try {
      const vault = getVault();
      vault.reset();
      spinner.succeed(chalk.green('Vault reset successfully!'));
    } catch (error) {
      spinner.fail(chalk.red('Failed to reset vault'));
      console.error(chalk.red(error instanceof Error ? error.message : 'Unknown error'));
      process.exit(1);
    }
  });

// Parse command line arguments
program.parse();
