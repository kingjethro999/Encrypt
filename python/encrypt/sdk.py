"""
Runtime SDK for the Encrypt tool.

Provides easy-to-use functions for accessing secrets in Python code.
"""

import os
from typing import Dict, Any, Optional
from .vault import Vault

# Global vault instance
_global_vault = None


def _get_vault() -> Vault:
    """Get or create the global vault instance."""
    global _global_vault
    if _global_vault is None:
        _global_vault = Vault()
    return _global_vault


def _auto_unlock(password: Optional[str] = None) -> None:
    """
    Auto-unlock vault if password is available.
    
    Args:
        password: Optional password to try first
    """
    vault = _get_vault()
    
    # If already unlocked, no need to do anything
    if vault.is_unlocked_status():
        return
    
    # Try provided password first
    if password:
        try:
            vault.setup(password)
            return
        except ValueError:
            # Password might be wrong, continue to other methods
            pass
    
    # Try environment variable
    env_password = os.getenv('ENCRYPT_PASSWORD')
    if env_password:
        try:
            vault.setup(env_password)
            return
        except ValueError:
            # Environment password might be wrong, continue to other methods
            pass
    
    # In development, we can be more lenient
    if os.getenv('NODE_ENV') != 'production':
        # Try common development passwords
        dev_passwords = ['dev', 'development', 'test', 'password', '123456']
        for dev_password in dev_passwords:
            try:
                vault.setup(dev_password)
                return
            except ValueError:
                # Continue to next password
                pass
    
    # If we get here, we couldn't unlock the vault
    raise ValueError('Vault is locked and no valid password found. Set ENCRYPT_PASSWORD environment variable or provide password parameter.')


def get(key: str, password: Optional[str] = None) -> str:
    """
    Get a secret value from the vault with auto-unlock support.
    
    Args:
        key: The secret key to retrieve
        password: Optional password to unlock vault
        
    Returns:
        The secret value
        
    Raises:
        ValueError: If vault is locked or key not found
    """
    _auto_unlock(password)
    vault = _get_vault()
    return vault.get(key)


def set(key: str, value: str, password: Optional[str] = None) -> None:
    """
    Set a secret value in the vault with auto-unlock support.
    
    Args:
        key: The secret key
        value: The secret value
        password: Optional password to unlock vault
        
    Raises:
        ValueError: If vault is locked
    """
    _auto_unlock(password)
    vault = _get_vault()
    vault.set(key, value)


def all_secrets(password: Optional[str] = None) -> Dict[str, str]:
    """
    Get all secrets from the vault with auto-unlock support.
    
    Args:
        password: Optional password to unlock vault
        
    Returns:
        Dictionary of all secrets
        
    Raises:
        ValueError: If vault is locked
    """
    _auto_unlock(password)
    vault = _get_vault()
    return vault.all()


def status() -> Dict[str, Any]:
    """
    Get vault status.
    
    Returns:
        Dictionary with vault status information
    """
    vault = _get_vault()
    return vault.status()


def is_unlocked() -> bool:
    """
    Check if vault is unlocked.
    
    Returns:
        True if vault is unlocked, False otherwise
    """
    vault = _get_vault()
    return vault.is_unlocked_status()


def auto_setup() -> None:
    """
    Auto setup helper - checks if vault is locked and provides helpful error.
    
    Raises:
        ValueError: If vault is not found or locked
    """
    vault = _get_vault()
    
    if not vault.exists():
        raise ValueError('Vault not found. Run "encrypt init" first.')
    
    if not vault.is_unlocked_status():
        raise ValueError('Vault is locked. Run "encrypt setup <password>" to unlock secrets.')


def get_secret(key: str) -> str:
    """
    Get a secret with automatic environment variable support.
    This is the recommended function for production use.
    
    Args:
        key: The secret key to retrieve
        
    Returns:
        The secret value
    """
    return get(key)


def set_secret(key: str, value: str) -> None:
    """
    Set a secret with automatic environment variable support.
    This is the recommended function for production use.
    
    Args:
        key: The secret key
        value: The secret value
    """
    set(key, value)


def get_all_secrets() -> Dict[str, str]:
    """
    Get all secrets with automatic environment variable support.
    This is the recommended function for production use.
    
    Returns:
        Dictionary of all secrets
    """
    return all_secrets()
