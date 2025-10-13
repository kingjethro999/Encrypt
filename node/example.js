// Example usage of the encrypt package
const encrypt = require('./dist/index.js');

async function example() {
  try {
    // Check if vault is unlocked
    if (!encrypt.isUnlocked()) {
      console.log('❌ Vault is locked. Run "encrypt setup <password>" to unlock secrets.');
      return;
    }

    // Get secrets
    const apiKey = encrypt.get('API_KEY');
    const dbUrl = encrypt.get('DB_URL');
    
    console.log('✅ Secrets retrieved successfully!');
    console.log('API Key:', apiKey ? '***' + apiKey.slice(-4) : 'Not found');
    console.log('DB URL:', dbUrl ? '***' + dbUrl.slice(-10) : 'Not found');
    
    // Get all secrets
    const allSecrets = encrypt.all();
    console.log('Available keys:', Object.keys(allSecrets));
    
  } catch (error) {
    console.error('Error:', error.message);
  }
}

example();
