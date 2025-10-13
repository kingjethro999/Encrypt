"""
Vault management for the Encrypt tool.

Handles file operations, state management, and secret storage/retrieval.
"""

import os
import json
from datetime import datetime
from typing import Dict, List, Optional
from .crypto import TripleEncryption


class Vault:
    """Vault class for managing encrypted secrets."""
    
    VAULT_DIR = '.encrypt'
    CONFIG_FILE = 'vault.lock'
    SECRETS_FILE = 'secrets.enc.json'
    GITIGNORE_FILE = '.gitignore'
    LOCK_FILE = 'vault.unlocked'
    
    def __init__(self, project_root: str = None):
        """Initialize vault with project root directory."""
        if project_root is None:
            project_root = os.getcwd()
        
        self.vault_path = os.path.join(project_root, Vault.VAULT_DIR)
        self.config_path = os.path.join(self.vault_path, Vault.CONFIG_FILE)
        self.secrets_path = os.path.join(self.vault_path, Vault.SECRETS_FILE)
        self.gitignore_path = os.path.join(project_root, Vault.GITIGNORE_FILE)
        self.lock_file_path = os.path.join(self.vault_path, Vault.LOCK_FILE)
        self.memory_cache: Dict[str, str] = {}
    
    def init(self) -> None:
        """Initialize the vault directory structure."""
        if not os.path.exists(self.vault_path):
            os.makedirs(self.vault_path, exist_ok=True)
        
        # Create initial empty secrets file
        if not os.path.exists(self.secrets_path):
            with open(self.secrets_path, 'w') as f:
                json.dump({}, f, indent=2)
        
        # Create initial config file
        if not os.path.exists(self.config_path):
            initial_config = {
                'password_hash': '',
                'salt': '',
                'hmac': '',
                'created_at': datetime.now().isoformat(),
                'version': '1.0.0'
            }
            with open(self.config_path, 'w') as f:
                json.dump(initial_config, f, indent=2)
        
        # Update .gitignore
        self._update_gitignore()
    
    def _update_gitignore(self) -> None:
        """Update .gitignore to exclude .encrypt directory."""
        gitignore_content = ''
        if os.path.exists(self.gitignore_path):
            with open(self.gitignore_path, 'r') as f:
                gitignore_content = f.read()
        
        if '.encrypt/' not in gitignore_content:
            with open(self.gitignore_path, 'a') as f:
                f.write('\n# Encrypt vault\n.encrypt/\n')
    
    def exists(self) -> bool:
        """Check if vault exists."""
        return os.path.exists(self.vault_path) and os.path.exists(self.config_path)
    
    def lockup(self, password: str) -> None:
        """Lock up secrets with password."""
        # Load secrets from file if not in memory
        if not self.memory_cache:
            self._load_secrets_from_file()
        
        if not self.memory_cache:
            raise ValueError('No secrets to lock. Use "encrypt set" to add secrets first.')
        
        # Encrypt all secrets
        encrypted_secrets = {}
        
        for key, value in self.memory_cache.items():
            result = TripleEncryption.encrypt(value, password)
            encrypted_secrets[key] = json.dumps(result)
        
        # Create vault config
        config_salt = TripleEncryption.generate_salt()
        config = {
            'password_hash': TripleEncryption.hash_password(password),
            'salt': config_salt,
            'hmac': TripleEncryption.generate_hmac(
                json.dumps(encrypted_secrets),
                TripleEncryption.derive_key(password, config_salt)
            ),
            'created_at': datetime.now().isoformat(),
            'version': '1.0.0'
        }
        
        # Write encrypted secrets and config
        with open(self.secrets_path, 'w') as f:
            json.dump(encrypted_secrets, f, indent=2)
        
        with open(self.config_path, 'w') as f:
            json.dump(config, f, indent=2)
        
        # Clear memory cache and remove lock file
        self.memory_cache.clear()
        if os.path.exists(self.lock_file_path):
            os.remove(self.lock_file_path)
    
    def setup(self, password: str) -> None:
        """Setup/unlock vault with password."""
        if not self.exists():
            raise ValueError('Vault not found. Run "encrypt init" first.')
        
        with open(self.config_path, 'r') as f:
            config = json.load(f)
        
        # If this is a fresh vault (no password set), just unlock it
        if not config['password_hash']:
            self.memory_cache.clear()
            self._create_lock_file()
            return
        
        # Verify password
        if not TripleEncryption.verify_password(password, config['password_hash']):
            raise ValueError('Invalid password.')
        
        # Load and decrypt secrets
        with open(self.secrets_path, 'r') as f:
            encrypted_secrets = json.load(f)
        
        self.memory_cache.clear()
        
        for key, encrypted_data in encrypted_secrets.items():
            # Check if the data is already decrypted (plain text) or encrypted
            if isinstance(encrypted_data, str) and encrypted_data.startswith('{'):
                # This is encrypted data stored as JSON string, decrypt it
                result = json.loads(encrypted_data)
                decrypted, is_valid = TripleEncryption.decrypt(
                    result['encrypted'], result['salt'], result['hmac'], password
                )
                
                if not is_valid:
                    raise ValueError(f'Failed to decrypt secret: {key}')
                
                self.memory_cache[key] = decrypted
            else:
                # This is plain text data
                self.memory_cache[key] = encrypted_data
        
        # Save decrypted secrets to file for easy access
        self._save_secrets_to_file()
        self._create_lock_file()
    
    def set(self, key: str, value: str) -> None:
        """Set a secret (only works when unlocked)."""
        if not self.is_unlocked_status():
            raise ValueError('Vault is locked. Run "encrypt setup <password>" to unlock secrets.')
        
        # Load secrets from file if not in memory
        if not self.memory_cache:
            self._load_secrets_from_file()
        
        self.memory_cache[key] = value
        self._save_secrets_to_file()
    
    def get(self, key: str) -> str:
        """Get a secret (only works when unlocked)."""
        if not self.is_unlocked_status():
            raise ValueError('Vault is locked. Run "encrypt setup <password>" to unlock secrets.')
        
        # Load secrets from file if not in memory
        if not self.memory_cache:
            self._load_secrets_from_file()
        
        if key not in self.memory_cache:
            raise ValueError(f'Secret "{key}" not found.')
        
        return self.memory_cache[key]
    
    def all(self) -> Dict[str, str]:
        """Get all secrets (only works when unlocked)."""
        if not self.is_unlocked_status():
            raise ValueError('Vault is locked. Run "encrypt setup <password>" to unlock secrets.')
        
        # Load secrets from file if not in memory
        if not self.memory_cache:
            self._load_secrets_from_file()
        
        return self.memory_cache.copy()
    
    def status(self) -> Dict:
        """Get vault status."""
        # Check if vault is unlocked by looking for lock file
        is_unlocked = os.path.exists(self.lock_file_path)
        
        # Load secrets from file if unlocked and not in memory
        if is_unlocked and not self.memory_cache:
            self._load_secrets_from_file()
        
        keys = list(self.memory_cache.keys())
        last_modified = None
        if self.exists():
            last_modified = datetime.fromtimestamp(
                os.path.getmtime(self.config_path)
            ).isoformat()
        
        return {
            'is_locked': not is_unlocked,
            'keys': keys,
            'last_modified': last_modified
        }
    
    def reset(self) -> None:
        """Reset/remove vault."""
        if os.path.exists(self.vault_path):
            import shutil
            shutil.rmtree(self.vault_path)
        self.memory_cache.clear()
    
    def is_unlocked_status(self) -> bool:
        """Check if vault is unlocked."""
        return os.path.exists(self.lock_file_path)
    
    def _create_lock_file(self) -> None:
        """Create lock file to indicate vault is unlocked."""
        lock_data = {
            'unlocked': True,
            'timestamp': datetime.now().isoformat()
        }
        with open(self.lock_file_path, 'w') as f:
            json.dump(lock_data, f, indent=2)
    
    def _load_secrets_from_file(self) -> None:
        """Load secrets from file (for unlocked vault)."""
        if os.path.exists(self.secrets_path):
            with open(self.secrets_path, 'r') as f:
                secrets = json.load(f)
            
            self.memory_cache.clear()
            for key, value in secrets.items():
                # Check if the value is encrypted (starts with {) or plain text
                if isinstance(value, str) and value.startswith('{'):
                    # This is encrypted data, we need to decrypt it
                    # But we don't have the password here, so we can't decrypt
                    # This should not happen in an unlocked vault
                    raise ValueError(f'Secret "{key}" is encrypted but vault is unlocked. This should not happen.')
                else:
                    # This is plain text data
                    self.memory_cache[key] = value
    
    def _save_secrets_to_file(self) -> None:
        """Save secrets to file (for unlocked vault)."""
        with open(self.secrets_path, 'w') as f:
            json.dump(self.memory_cache, f, indent=2)
