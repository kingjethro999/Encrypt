#!/usr/bin/env node

const { execSync } = require('child_process');
const chalk = require('chalk');

console.log(chalk.blue.bold('🔐 Encrypt Tool Demo'));
console.log(chalk.gray('This demo shows the complete workflow of the encrypt tool.\n'));

function runCommand(command, description) {
  console.log(chalk.yellow(`📝 ${description}`));
  console.log(chalk.gray(`$ ${command}`));
  
  try {
    const output = execSync(command, { encoding: 'utf8', cwd: process.cwd() });
    if (output.trim()) {
      console.log(chalk.green(output.trim()));
    }
  } catch (error) {
    console.log(chalk.red(error.stdout || error.message));
  }
  console.log('');
}

// Clean up any existing vault
try {
  execSync('rm -rf .encrypt .env', { stdio: 'ignore' });
} catch (e) {}

console.log(chalk.blue.bold('1. Initialize the vault'));
runCommand('node dist/cli.js init', 'Create a new encrypted vault');

console.log(chalk.blue.bold('2. Setup the vault (unlock for first time)'));
runCommand('node dist/cli.js setup demo-password', 'Unlock the vault with a password');

console.log(chalk.blue.bold('3. Add some secrets'));
runCommand('node dist/cli.js set API_KEY=sk-1234567890abcdef', 'Add an API key');
runCommand('node dist/cli.js set DB_URL=postgres://user:pass@localhost:5432/mydb', 'Add a database URL');
runCommand('node dist/cli.js set JWT_SECRET=super-secret-jwt-key', 'Add a JWT secret');

console.log(chalk.blue.bold('4. Check vault status'));
runCommand('node dist/cli.js status', 'View vault status and available keys');

console.log(chalk.blue.bold('5. Retrieve secrets'));
runCommand('node dist/cli.js get API_KEY', 'Get the API key');
runCommand('node dist/cli.js get DB_URL', 'Get the database URL');

console.log(chalk.blue.bold('6. Test runtime SDK'));
console.log(chalk.yellow('📝 Using the encrypt package in code'));
console.log(chalk.gray('$ node example.js'));
try {
  const output = execSync('node example.js', { encoding: 'utf8' });
  console.log(chalk.green(output.trim()));
} catch (error) {
  console.log(chalk.red(error.stdout || error.message));
}
console.log('');

console.log(chalk.blue.bold('7. Create .env file'));
runCommand('node dist/cli.js unlock', 'Export all secrets to .env file');

console.log(chalk.blue.bold('8. Lock the vault'));
runCommand('node dist/cli.js lockup demo-password', 'Encrypt and lock all secrets');

console.log(chalk.blue.bold('9. Try to access locked secrets'));
runCommand('node dist/cli.js get API_KEY', 'Attempt to get secret from locked vault');

console.log(chalk.blue.bold('10. Unlock the vault again'));
runCommand('node dist/cli.js setup demo-password', 'Unlock the vault with the password');

console.log(chalk.blue.bold('11. Verify secrets are accessible again'));
runCommand('node dist/cli.js get API_KEY', 'Get the API key after unlocking');

console.log(chalk.green.bold('✅ Demo completed successfully!'));
console.log(chalk.gray('\nThe encrypt tool provides:'));
console.log(chalk.gray('• Secure encryption with triple-layer protection'));
console.log(chalk.gray('• Easy team onboarding with single command'));
console.log(chalk.gray('• Runtime SDK for in-code secret access'));
console.log(chalk.gray('• Git-safe encrypted storage'));
console.log(chalk.gray('• CLI and programmatic interfaces'));
