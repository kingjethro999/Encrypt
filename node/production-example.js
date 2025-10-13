#!/usr/bin/env node

/**
 * Production Usage Example for Node.js
 * 
 * This demonstrates how to use the encrypt tool in production
 * with automatic environment variable support.
 */

const encrypt = require('./dist/index.js');

console.log('🚀 Production Usage Example (Node.js)');
console.log('=====================================\n');

// Method 1: Using environment variable (Recommended for production)
console.log('Method 1: Environment Variable (Recommended)');
console.log('Set ENCRYPT_PASSWORD=your-password in your environment\n');

try {
  // This will automatically unlock the vault using ENCRYPT_PASSWORD
  const apiKey = encrypt.getSecret('API_KEY');
  const dbUrl = encrypt.getSecret('DB_URL');
  
  console.log('✅ Successfully retrieved secrets:');
  console.log(`API Key: ${apiKey ? '***' + apiKey.slice(-4) : 'Not found'}`);
  console.log(`DB URL: ${dbUrl ? '***' + dbUrl.slice(-10) : 'Not found'}`);
  
  // Use in your application
  const config = {
    apiKey: apiKey,
    database: dbUrl,
    port: process.env.PORT || 3000
  };
  
  console.log('\n📋 Application config ready:', {
    apiKey: config.apiKey ? '***' + config.apiKey.slice(-4) : 'Not found',
    database: config.database ? '***' + config.database.slice(-10) : 'Not found',
    port: config.port
  });
  
} catch (error) {
  console.log('❌ Error:', error.message);
  console.log('\n💡 To fix this:');
  console.log('1. Set ENCRYPT_PASSWORD environment variable');
  console.log('2. Or provide password as second parameter');
  console.log('3. Or run "encrypt setup <password>" first');
}

console.log('\n' + '='.repeat(50));
console.log('Method 2: Explicit Password Parameter');
console.log('='.repeat(50));

try {
  // This will use the provided password to unlock the vault
  const apiKey = encrypt.get('API_KEY', 'mypassword');
  console.log('✅ Success with explicit password:', apiKey ? '***' + apiKey.slice(-4) : 'Not found');
} catch (error) {
  console.log('❌ Error with explicit password:', error.message);
}

console.log('\n' + '='.repeat(50));
console.log('Method 3: Development Mode');
console.log('='.repeat(50));

// In development, you can set NODE_ENV=development
// and it will try common passwords automatically
process.env.NODE_ENV = 'development';

try {
  const apiKey = encrypt.get('API_KEY');
  console.log('✅ Development mode success:', apiKey ? '***' + apiKey.slice(-4) : 'Not found');
} catch (error) {
  console.log('❌ Development mode failed:', error.message);
}

console.log('\n🎯 Production Deployment Examples:');
console.log('===================================');
console.log('Docker:');
console.log('  ENV ENCRYPT_PASSWORD=your-production-password');
console.log('');
console.log('Kubernetes:');
console.log('  env:');
console.log('  - name: ENCRYPT_PASSWORD');
console.log('    valueFrom:');
console.log('      secretKeyRef:');
console.log('        name: encrypt-secrets');
console.log('        key: password');
console.log('');
console.log('Heroku:');
console.log('  heroku config:set ENCRYPT_PASSWORD=your-password');
console.log('');
console.log('AWS Lambda:');
console.log('  Set ENCRYPT_PASSWORD in environment variables');
