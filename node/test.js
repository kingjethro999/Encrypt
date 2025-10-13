// Simple test to verify encrypt functionality
const { TripleEncryption } = require('./dist/crypto.js');

console.log('🧪 Testing TripleEncryption...');

// Test encryption/decryption
const plaintext = 'Hello, World!';
const password = 'test-password-123';

console.log('Original:', plaintext);

try {
  // Encrypt
  const encrypted = TripleEncryption.encrypt(plaintext, password);
  console.log('✅ Encryption successful');
  console.log('Encrypted data length:', encrypted.encrypted.length);
  
  // Decrypt
  const decrypted = TripleEncryption.decrypt(encrypted.encrypted, encrypted.salt, encrypted.hmac, password);
  
  if (decrypted.isValid && decrypted.decrypted === plaintext) {
    console.log('✅ Decryption successful');
    console.log('Decrypted:', decrypted.decrypted);
  } else {
    console.log('❌ Decryption failed');
  }
  
  // Test password hashing
  const hash = TripleEncryption.hashPassword(password);
  const isValid = TripleEncryption.verifyPassword(password, hash);
  
  if (isValid) {
    console.log('✅ Password hashing/verification successful');
  } else {
    console.log('❌ Password hashing/verification failed');
  }
  
} catch (error) {
  console.error('❌ Test failed:', error.message);
}
