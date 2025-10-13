#!/usr/bin/env python3
"""
Test script for the TripleEncryption class.
"""

from encrypt.crypto import TripleEncryption

def test_encryption():
    """Test encryption and decryption functionality."""
    print("🧪 Testing TripleEncryption...")
    
    # Test data
    plaintext = "Hello, World!"
    password = "test-password-123"
    
    print(f"Original: {plaintext}")
    
    try:
        # Encrypt
        encrypted = TripleEncryption.encrypt(plaintext, password)
        print("✅ Encryption successful")
        print(f"Encrypted data length: {len(encrypted['encrypted'])}")
        
        # Decrypt
        decrypted, is_valid = TripleEncryption.decrypt(
            encrypted['encrypted'], 
            encrypted['salt'], 
            encrypted['hmac'], 
            password
        )
        
        if is_valid and decrypted == plaintext:
            print("✅ Decryption successful")
            print(f"Decrypted: {decrypted}")
        else:
            print("❌ Decryption failed")
        
        # Test password hashing
        hash_result = TripleEncryption.hash_password(password)
        is_valid_hash = TripleEncryption.verify_password(password, hash_result)
        
        if is_valid_hash:
            print("✅ Password hashing/verification successful")
        else:
            print("❌ Password hashing/verification failed")
        
        # Test with different password (should fail)
        wrong_decrypted, wrong_valid = TripleEncryption.decrypt(
            encrypted['encrypted'], 
            encrypted['salt'], 
            encrypted['hmac'], 
            "wrong-password"
        )
        
        if not wrong_valid:
            print("✅ Wrong password correctly rejected")
        else:
            print("❌ Wrong password incorrectly accepted")
        
        print("\n🎉 All tests passed!")
        
    except Exception as e:
        print(f"❌ Test failed: {str(e)}")

if __name__ == "__main__":
    test_encryption()
