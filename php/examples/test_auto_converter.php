<?php

require_once __DIR__ . '/../vendor/autoload.php';

use Encrypt\SDK\SDK;

echo "🧪 Testing Auto-Converter Engine (PHP)\n";
echo "=====================================\n\n";

// Test 1: Try to get secret when vault is locked (should fail without password)
echo "1. Testing locked vault without password...\n";
try {
    $apiKey = SDK::getInstance()->get('API_KEY', '');
    echo "❌ Unexpected: Got secret from locked vault: $apiKey\n";
} catch (Exception $e) {
    echo "✅ Expected: Vault is locked - " . $e->getMessage() . "\n";
}

// Test 2: Try with environment variable
echo "\n2. Testing with ENCRYPT_PASSWORD environment variable...\n";
putenv('ENCRYPT_PASSWORD=mypassword');

try {
    $apiKey = SDK::getInstance()->get('API_KEY', '');
    echo "✅ Success: Got secret with environment password - $apiKey\n";
} catch (Exception $e) {
    echo "❌ Failed: Could not get secret with environment password - " . $e->getMessage() . "\n";
}

// Test 3: Try with explicit password parameter
echo "\n3. Testing with explicit password parameter...\n";
putenv('ENCRYPT_PASSWORD'); // Clear env var

try {
    $apiKey = SDK::getInstance()->get('API_KEY', 'mypassword');
    echo "✅ Success: Got secret with explicit password - $apiKey\n";
} catch (Exception $e) {
    echo "❌ Failed: Could not get secret with explicit password - " . $e->getMessage() . "\n";
}

// Test 4: Test production-ready functions
echo "\n4. Testing production-ready functions...\n";
putenv('ENCRYPT_PASSWORD=mypassword');

try {
    $apiKey = SDK::getInstance()->getSecret('API_KEY');
    echo "✅ Success: get_secret() works - $apiKey\n";
} catch (Exception $e) {
    echo "❌ Failed: get_secret() failed - " . $e->getMessage() . "\n";
}

try {
    $dbUrl = SDK::getInstance()->getSecret('DB_URL');
    echo "✅ Success: get_secret() works - $dbUrl\n";
} catch (Exception $e) {
    echo "❌ Failed: get_secret() failed - " . $e->getMessage() . "\n";
}

// Test 5: Test development mode (should try common passwords)
echo "\n5. Testing development mode...\n";
putenv('ENCRYPT_PASSWORD'); // Clear env var
putenv('NODE_ENV=development');

try {
    $apiKey = SDK::getInstance()->get('API_KEY', '');
    echo "✅ Success: Development mode worked - $apiKey\n";
} catch (Exception $e) {
    echo "❌ Failed: Development mode failed - " . $e->getMessage() . "\n";
}

echo "\n🎉 Auto-converter tests completed!\n";
