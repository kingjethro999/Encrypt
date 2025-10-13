<?php

require_once __DIR__ . '/../vendor/autoload.php';

use Encrypt\SDK\SDK;

echo "💻 Example usage of the encrypt package in PHP code\n";
echo "==================================================\n";

try {
    // Check if vault is unlocked
    if (!SDK::getInstance()->isUnlocked()) {
        echo "❌ Vault is locked. Run 'encrypt setup <password>' to unlock secrets.\n";
        exit(1);
    }
    
    // Get secrets
    try {
        $apiKey = SDK::getInstance()->getSecret('API_KEY');
        echo "✅ Secrets retrieved successfully!\n";
        if (strlen($apiKey) > 4) {
            echo "API Key: ***" . substr($apiKey, -4) . "\n";
        } else {
            echo "API Key: ***$apiKey\n";
        }
    } catch (Exception $e) {
        echo "⚠️ API_KEY not found: " . $e->getMessage() . "\n";
    }
    
    try {
        $dbUrl = SDK::getInstance()->getSecret('DB_URL');
        if (strlen($dbUrl) > 10) {
            echo "DB URL: ***" . substr($dbUrl, -10) . "\n";
        } else {
            echo "DB URL: ***$dbUrl\n";
        }
    } catch (Exception $e) {
        echo "⚠️ DB_URL not found: " . $e->getMessage() . "\n";
    }
    
    // Get all secrets
    try {
        $allSecrets = SDK::getInstance()->getAllSecrets();
        echo "Available keys: " . implode(', ', array_keys($allSecrets)) . "\n";
    } catch (Exception $e) {
        echo "⚠️ Error getting all secrets: " . $e->getMessage() . "\n";
    }
    
    // Get status
    try {
        $status = SDK::getInstance()->status();
        $statusText = $status->isLocked ? 'Locked' : 'Unlocked';
        echo "Vault status: $statusText\n";
        echo "Number of keys: " . count($status->keys) . "\n";
    } catch (Exception $e) {
        echo "⚠️ Error getting status: " . $e->getMessage() . "\n";
    }
    
} catch (Exception $e) {
    echo "Error: " . $e->getMessage() . "\n";
}
