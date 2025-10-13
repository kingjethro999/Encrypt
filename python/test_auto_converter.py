#!/usr/bin/env python3
"""
Test the auto-converter functionality for Python.
"""

import os
from encrypt import get, get_secret, set_secret, get_all_secrets

def test_auto_converter():
    """Test the auto-converter engine functionality."""
    print("🧪 Testing Auto-Converter Engine (Python)")
    print("========================================\n")
    
    try:
        # Test 1: Try to get secret when vault is locked (should fail without password)
        print("1. Testing locked vault without password...")
        try:
            api_key = get('API_KEY')
            print("❌ Unexpected: Got secret from locked vault")
        except ValueError as e:
            print(f"✅ Expected: Vault is locked - {e}")
        
        # Test 2: Try with environment variable
        print("\n2. Testing with ENCRYPT_PASSWORD environment variable...")
        os.environ['ENCRYPT_PASSWORD'] = 'mypassword'
        
        try:
            api_key = get('API_KEY')
            print(f"✅ Success: Got secret with environment password - {api_key}")
        except ValueError as e:
            print(f"❌ Failed: Could not get secret with environment password - {e}")
        
        # Test 3: Try with explicit password parameter
        print("\n3. Testing with explicit password parameter...")
        del os.environ['ENCRYPT_PASSWORD']  # Clear env var
        
        try:
            api_key = get('API_KEY', 'mypassword')
            print(f"✅ Success: Got secret with explicit password - {api_key}")
        except ValueError as e:
            print(f"❌ Failed: Could not get secret with explicit password - {e}")
        
        # Test 4: Test production-ready functions
        print("\n4. Testing production-ready functions...")
        os.environ['ENCRYPT_PASSWORD'] = 'mypassword'
        
        try:
            api_key = get_secret('API_KEY')
            db_url = get_secret('DB_URL')
            print(f"✅ Success: get_secret() works - {api_key}")
            print(f"✅ Success: get_secret() works - {db_url}")
        except ValueError as e:
            print(f"❌ Failed: get_secret() failed - {e}")
        
        # Test 5: Test development mode (should try common passwords)
        print("\n5. Testing development mode...")
        del os.environ['ENCRYPT_PASSWORD']
        os.environ['NODE_ENV'] = 'development'
        
        try:
            api_key = get('API_KEY')
            print(f"✅ Success: Development mode worked - {api_key}")
        except ValueError as e:
            print(f"❌ Failed: Development mode failed - {e}")
        
        print("\n🎉 Auto-converter tests completed!")
        
    except Exception as e:
        print(f"❌ Test failed: {e}")

if __name__ == "__main__":
    test_auto_converter()
