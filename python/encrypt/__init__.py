"""
Encrypt - A top-level secrets orchestrator for Python.

This package provides secure secret management with triple-layer encryption,
replacing .env files with an encrypted local secrets vault.
"""

from .vault import Vault
from .crypto import TripleEncryption
from .sdk import (
    get, set, all_secrets, status, is_unlocked, auto_setup,
    get_secret, set_secret, get_all_secrets
)

__version__ = "1.0.0"
__author__ = "Encrypt Team"

# Convenience aliases for backward compatibility
get_status = status
is_vault_unlocked = is_unlocked
auto_setup_vault = auto_setup

# Export main classes and functions
__all__ = [
    'Vault',
    'TripleEncryption', 
    'get_secret',
    'set_secret',
    'get_all_secrets',
    'get_status',
    'is_vault_unlocked',
    'auto_setup_vault',
    'get',
    'set',
    'all_secrets',
    'status',
    'is_unlocked',
    'auto_setup'
]
