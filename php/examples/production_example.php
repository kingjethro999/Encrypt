<?php

require_once __DIR__ . '/../vendor/autoload.php';

use Encrypt\SDK\SDK;

class Config
{
    public function __construct(
        public string $apiKey,
        public string $database,
        public string $port
    ) {}
}

echo "🚀 Production Usage Example (PHP)\n";
echo "==================================\n\n";

// Method 1: Using environment variable (Recommended for production)
echo "Method 1: Environment Variable (Recommended)\n";
echo "Set ENCRYPT_PASSWORD=your-password in your environment\n\n";

try {
    $apiKey = SDK::getInstance()->getSecret('API_KEY');
    echo "✅ Successfully retrieved secrets:\n";
    if (strlen($apiKey) > 4) {
        echo "API Key: ***" . substr($apiKey, -4) . "\n";
    } else {
        echo "API Key: ***$apiKey\n";
    }
    
    $dbUrl = SDK::getInstance()->getSecret('DB_URL');
    if (strlen($dbUrl) > 10) {
        echo "DB URL: ***" . substr($dbUrl, -10) . "\n";
    } else {
        echo "DB URL: ***$dbUrl\n";
    }
    
    // Use in your application
    $config = new Config(
        $apiKey,
        $dbUrl,
        $_ENV['PORT'] ?? '3000'
    );
    
    echo "\n📋 Application config ready:\n";
    if (strlen($config->apiKey) > 4) {
        echo "  api_key: ***" . substr($config->apiKey, -4) . "\n";
    } else {
        echo "  api_key: ***{$config->apiKey}\n";
    }
    if (strlen($config->database) > 10) {
        echo "  database: ***" . substr($config->database, -10) . "\n";
    } else {
        echo "  database: ***{$config->database}\n";
    }
    echo "  port: {$config->port}\n";
    
} catch (Exception $e) {
    echo "❌ Error: " . $e->getMessage() . "\n";
    echo "\n💡 To fix this:\n";
    echo "1. Set ENCRYPT_PASSWORD environment variable\n";
    echo "2. Or provide password as second parameter\n";
    echo "3. Or run 'encrypt setup <password>' first\n";
}

echo "\n==================================================\n";
echo "Method 2: Explicit Password Parameter\n";
echo "==================================================\n";

try {
    $apiKey = SDK::getInstance()->get('API_KEY', 'mypassword');
    if (strlen($apiKey) > 4) {
        echo "✅ Success with explicit password: ***" . substr($apiKey, -4) . "\n";
    } else {
        echo "✅ Success with explicit password: ***$apiKey\n";
    }
} catch (Exception $e) {
    echo "❌ Error with explicit password: " . $e->getMessage() . "\n";
}

echo "\n==================================================\n";
echo "Method 3: Development Mode\n";
echo "==================================================\n";

// In development, you can set NODE_ENV=development
// and it will try common passwords automatically
putenv('NODE_ENV=development');

try {
    $apiKey = SDK::getInstance()->get('API_KEY', '');
    if (strlen($apiKey) > 4) {
        echo "✅ Development mode success: ***" . substr($apiKey, -4) . "\n";
    } else {
        echo "✅ Development mode success: ***$apiKey\n";
    }
} catch (Exception $e) {
    echo "❌ Development mode failed: " . $e->getMessage() . "\n";
}

echo "\n🎯 Production Deployment Examples:\n";
echo "==================================\n";
echo "Docker:\n";
echo "  ENV ENCRYPT_PASSWORD=your-production-password\n\n";
echo "Kubernetes:\n";
echo "  env:\n";
echo "  - name: ENCRYPT_PASSWORD\n";
echo "    valueFrom:\n";
echo "      secretKeyRef:\n";
echo "        name: encrypt-secrets\n";
echo "        key: password\n\n";
echo "Heroku:\n";
echo "  heroku config:set ENCRYPT_PASSWORD=your-password\n\n";
echo "AWS Lambda:\n";
echo "  Set ENCRYPT_PASSWORD in environment variables\n\n";
echo "PHP application:\n";
echo "  use Encrypt\\SDK\\SDK;\n";
echo "  \$apiKey = SDK::getInstance()->getSecret('API_KEY');\n";
echo "  \$databaseURL = SDK::getInstance()->getSecret('DATABASE_URL');\n";
