"""
Triple-layer encryption implementation for the Encrypt tool.

Provides AES-256-CBC encryption with PBKDF2 key derivation and HMAC signatures.
"""

import os
import json
import hashlib
from typing import Dict, Tuple
from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes
from cryptography.hazmat.primitives.kdf.pbkdf2 import PBKDF2HMAC
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.backends import default_backend


class TripleEncryption:
    """Triple-layer encryption class with AES-256-CBC, PBKDF2, and HMAC."""
    
    ALGORITHM = "AES-256-CBC"
    KEY_LENGTH = 32
    IV_LENGTH = 16
    SALT_LENGTH = 32
    ITERATIONS = 100000
    
    @staticmethod
    def generate_salt() -> str:
        """Generate a random salt for key derivation."""
        return os.urandom(TripleEncryption.SALT_LENGTH).hex()
    
    @staticmethod
    def derive_key(password: str, salt: str) -> bytes:
        """Derive encryption key from password using PBKDF2."""
        salt_bytes = bytes.fromhex(salt)
        kdf = PBKDF2HMAC(
            algorithm=hashes.SHA512(),
            length=TripleEncryption.KEY_LENGTH,
            salt=salt_bytes,
            iterations=TripleEncryption.ITERATIONS,
            backend=default_backend()
        )
        return kdf.derive(password.encode('utf-8'))
    
    @staticmethod
    def generate_hmac(data: str, key: bytes) -> str:
        """Generate HMAC signature for data integrity."""
        h = hashlib.sha256()
        h.update(key)
        h.update(data.encode('utf-8'))
        return h.hexdigest()
    
    @staticmethod
    def encrypt(plaintext: str, password: str) -> Dict[str, str]:
        """
        Triple-layer encryption:
        1. Generate random salt
        2. Derive key using PBKDF2
        3. AES-256-CBC encryption with HMAC
        """
        # Phase 1: Generate salt
        salt = TripleEncryption.generate_salt()
        
        # Phase 2: Derive key from password
        key = TripleEncryption.derive_key(password, salt)
        
        # Phase 3: AES-256-CBC encryption
        iv = os.urandom(TripleEncryption.IV_LENGTH)
        cipher = Cipher(
            algorithms.AES(key),
            modes.CBC(iv),
            backend=default_backend()
        )
        encryptor = cipher.encryptor()
        
        # Pad the plaintext to block size
        plaintext_bytes = plaintext.encode('utf-8')
        padding_length = TripleEncryption.IV_LENGTH - (len(plaintext_bytes) % TripleEncryption.IV_LENGTH)
        padded_plaintext = plaintext_bytes + bytes([padding_length] * padding_length)
        
        encrypted = encryptor.update(padded_plaintext) + encryptor.finalize()
        
        # Combine IV and encrypted data
        encrypted_data = iv.hex() + ':' + encrypted.hex()
        
        # Phase 3: Generate HMAC signature
        hmac = TripleEncryption.generate_hmac(encrypted_data, key)
        
        return {
            'encrypted': encrypted_data,
            'salt': salt,
            'hmac': hmac
        }
    
    @staticmethod
    def decrypt(encrypted_data: str, salt: str, hmac: str, password: str) -> Tuple[str, bool]:
        """
        Triple-layer decryption:
        1. Verify HMAC signature
        2. Derive key using PBKDF2
        3. AES-256-CBC decryption
        """
        try:
            # Phase 2: Derive key from password
            key = TripleEncryption.derive_key(password, salt)
            
            # Phase 3: Verify HMAC signature
            expected_hmac = TripleEncryption.generate_hmac(encrypted_data, key)
            if expected_hmac != hmac:
                return '', False
            
            # Phase 3: AES-256-CBC decryption
            iv_hex, encrypted_hex = encrypted_data.split(':', 1)
            iv = bytes.fromhex(iv_hex)
            encrypted = bytes.fromhex(encrypted_hex)
            
            cipher = Cipher(
                algorithms.AES(key),
                modes.CBC(iv),
                backend=default_backend()
            )
            decryptor = cipher.decryptor()
            
            decrypted_padded = decryptor.update(encrypted) + decryptor.finalize()
            
            # Remove padding
            padding_length = decrypted_padded[-1]
            decrypted = decrypted_padded[:-padding_length]
            
            return decrypted.decode('utf-8'), True
            
        except Exception:
            return '', False
    
    @staticmethod
    def hash_password(password: str) -> str:
        """Hash password for storage using PBKDF2."""
        salt = TripleEncryption.generate_salt()
        salt_bytes = bytes.fromhex(salt)
        kdf = PBKDF2HMAC(
            algorithm=hashes.SHA512(),
            length=64,
            salt=salt_bytes,
            iterations=TripleEncryption.ITERATIONS,
            backend=default_backend()
        )
        hash_bytes = kdf.derive(password.encode('utf-8'))
        return salt + ':' + hash_bytes.hex()
    
    @staticmethod
    def verify_password(password: str, stored_hash: str) -> bool:
        """Verify password against stored hash."""
        try:
            salt, stored_hash_hex = stored_hash.split(':', 1)
            salt_bytes = bytes.fromhex(salt)
            kdf = PBKDF2HMAC(
                algorithm=hashes.SHA512(),
                length=64,
                salt=salt_bytes,
                iterations=TripleEncryption.ITERATIONS,
                backend=default_backend()
            )
            computed_hash = kdf.derive(password.encode('utf-8'))
            return computed_hash.hex() == stored_hash_hex
        except Exception:
            return False
