#!/usr/bin/env node

// Test the auto-converter functionality
const encrypt = require('./dist/index.js');

console.log('🧪 Testing Auto-Converter Engine (Node.js)');
console.log('==========================================\n');

async function testAutoConverter() {
  try {
    // Test 1: Try to get secret when vault is locked (should fail without password)
    console.log('1. Testing locked vault without password...');
    try {
      const apiKey = encrypt.get('API_KEY');
      console.log('❌ Unexpected: Got secret from locked vault');
    } catch (error) {
      console.log('✅ Expected: Vault is locked -', error.message);
    }
    
    // Test 2: Try with environment variable
    console.log('\n2. Testing with ENCRYPT_PASSWORD environment variable...');
    process.env.ENCRYPT_PASSWORD = 'mypassword';
    
    try {
      const apiKey = encrypt.get('API_KEY');
      console.log('✅ Success: Got secret with environment password -', apiKey);
    } catch (error) {
      console.log('❌ Failed: Could not get secret with environment password -', error.message);
    }
    
    // Test 3: Try with explicit password parameter
    console.log('\n3. Testing with explicit password parameter...');
    delete process.env.ENCRYPT_PASSWORD; // Clear env var
    
    try {
      const apiKey = encrypt.get('API_KEY', 'mypassword');
      console.log('✅ Success: Got secret with explicit password -', apiKey);
    } catch (error) {
      console.log('❌ Failed: Could not get secret with explicit password -', error.message);
    }
    
    // Test 4: Test production-ready functions
    console.log('\n4. Testing production-ready functions...');
    process.env.ENCRYPT_PASSWORD = 'mypassword';
    
    try {
      const apiKey = encrypt.getSecret('API_KEY');
      const dbUrl = encrypt.getSecret('DB_URL');
      console.log('✅ Success: getSecret() works -', apiKey);
      console.log('✅ Success: getSecret() works -', dbUrl);
    } catch (error) {
      console.log('❌ Failed: getSecret() failed -', error.message);
    }
    
    // Test 5: Test development mode (should try common passwords)
    console.log('\n5. Testing development mode...');
    delete process.env.ENCRYPT_PASSWORD;
    process.env.NODE_ENV = 'development';
    
    try {
      const apiKey = encrypt.get('API_KEY');
      console.log('✅ Success: Development mode worked -', apiKey);
    } catch (error) {
      console.log('❌ Failed: Development mode failed -', error.message);
    }
    
    console.log('\n🎉 Auto-converter tests completed!');
    
  } catch (error) {
    console.error('❌ Test failed:', error.message);
  }
}

testAutoConverter();
