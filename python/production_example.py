#!/usr/bin/env python3
"""
Production Usage Example for Python

This demonstrates how to use the encrypt tool in production
with automatic environment variable support.
"""

import os
from encrypt import get_secret, get_all_secrets

def main():
    """Demonstrate production usage patterns."""
    print("🚀 Production Usage Example (Python)")
    print("====================================\n")
    
    # Method 1: Using environment variable (Recommended for production)
    print("Method 1: Environment Variable (Recommended)")
    print("Set ENCRYPT_PASSWORD=your-password in your environment\n")
    
    try:
        # This will automatically unlock the vault using ENCRYPT_PASSWORD
        api_key = get_secret('API_KEY')
        db_url = get_secret('DB_URL')
        
        print("✅ Successfully retrieved secrets:")
        print(f"API Key: {'***' + api_key[-4:] if api_key else 'Not found'}")
        print(f"DB URL: {'***' + db_url[-10:] if db_url else 'Not found'}")
        
        # Use in your application
        config = {
            'api_key': api_key,
            'database': db_url,
            'port': int(os.getenv('PORT', 3000))
        }
        
        print("\n📋 Application config ready:", {
            'api_key': '***' + config['api_key'][-4:] if config['api_key'] else 'Not found',
            'database': '***' + config['database'][-10:] if config['database'] else 'Not found',
            'port': config['port']
        })
        
    except ValueError as e:
        print(f"❌ Error: {e}")
        print("\n💡 To fix this:")
        print("1. Set ENCRYPT_PASSWORD environment variable")
        print("2. Or provide password as second parameter")
        print("3. Or run 'encrypt setup <password>' first")
    
    print("\n" + "="*50)
    print("Method 2: Explicit Password Parameter")
    print("="*50)
    
    try:
        # This will use the provided password to unlock the vault
        from encrypt import get
        api_key = get('API_KEY', 'mypassword')
        print(f"✅ Success with explicit password: {'***' + api_key[-4:] if api_key else 'Not found'}")
    except ValueError as e:
        print(f"❌ Error with explicit password: {e}")
    
    print("\n" + "="*50)
    print("Method 3: Development Mode")
    print("="*50)
    
    # In development, you can set NODE_ENV=development
    # and it will try common passwords automatically
    os.environ['NODE_ENV'] = 'development'
    
    try:
        from encrypt import get
        api_key = get('API_KEY')
        print(f"✅ Development mode success: {'***' + api_key[-4:] if api_key else 'Not found'}")
    except ValueError as e:
        print(f"❌ Development mode failed: {e}")
    
    print("\n🎯 Production Deployment Examples:")
    print("===================================")
    print("Docker:")
    print("  ENV ENCRYPT_PASSWORD=your-production-password")
    print("")
    print("Kubernetes:")
    print("  env:")
    print("  - name: ENCRYPT_PASSWORD")
    print("    valueFrom:")
    print("      secretKeyRef:")
    print("        name: encrypt-secrets")
    print("        key: password")
    print("")
    print("Heroku:")
    print("  heroku config:set ENCRYPT_PASSWORD=your-password")
    print("")
    print("AWS Lambda:")
    print("  Set ENCRYPT_PASSWORD in environment variables")
    print("")
    print("Django settings.py:")
    print("  import os")
    print("  from encrypt import get_secret")
    print("  SECRET_KEY = get_secret('SECRET_KEY')")
    print("  DATABASE_URL = get_secret('DATABASE_URL')")

if __name__ == "__main__":
    main()
