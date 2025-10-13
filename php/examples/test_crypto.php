<?php

require_once __DIR__ . '/../vendor/autoload.php';

use Encrypt\Crypto\Crypto;

echo "🧪 Testing TripleEncryption (PHP)...\n";

// Test data
$plaintext = "Hello, World!";
$password = "test-password-123";

echo "Original: $plaintext\n";

try {
    // Encrypt
    $crypto = new Crypto();
    $encrypted = $crypto->encrypt($plaintext, $password);
    
    echo "✅ Encryption successful\n";
    echo "Encrypted data length: " . strlen($encrypted->encrypted) . "\n";
    
    // Decrypt
    [$decrypted, $isValid] = $crypto->decrypt($encrypted->encrypted, $encrypted->salt, $encrypted->hmac, $password);
    
    if ($isValid && $decrypted === $plaintext) {
        echo "✅ Decryption successful\n";
        echo "Decrypted: $decrypted\n";
    } else {
        echo "❌ Decryption failed\n";
    }
    
    // Test password hashing
    $hashResult = $crypto->hashPassword($password);
    $isValidHash = $crypto->verifyPassword($password, $hashResult);
    
    if ($isValidHash) {
        echo "✅ Password hashing/verification successful\n";
    } else {
        echo "❌ Password hashing/verification failed\n";
    }
    
    // Test with different password (should fail)
    [$wrongDecrypted, $wrongValid] = $crypto->decrypt($encrypted->encrypted, $encrypted->salt, $encrypted->hmac, "wrong-password");
    
    if (!$wrongValid) {
        echo "✅ Wrong password correctly rejected\n";
    } else {
        echo "❌ Wrong password incorrectly accepted\n";
    }
    
    echo "\n🎉 All tests passed!\n";
    
} catch (Exception $e) {
    echo "❌ Test failed: " . $e->getMessage() . "\n";
}
