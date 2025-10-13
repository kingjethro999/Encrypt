#!/usr/bin/env python3
"""
Example usage of the encrypt package in Python code.
"""

from encrypt import get_secret, set_secret, get_all_secrets, is_vault_unlocked, get_status

def main():
    """Demonstrate encrypt package usage."""
    try:
        # Check if vault is unlocked
        if not is_vault_unlocked():
            print("❌ Vault is locked. Run 'encrypt setup <password>' to unlock secrets.")
            return
        
        # Get secrets
        try:
            api_key = get_secret('API_KEY')
            db_url = get_secret('DB_URL')
            
            print("✅ Secrets retrieved successfully!")
            print(f"API Key: ***{api_key[-4:] if api_key else 'Not found'}")
            print(f"DB URL: ***{db_url[-10:] if db_url else 'Not found'}")
            
        except ValueError as e:
            print(f"⚠️ Some secrets not found: {e}")
        
        # Get all secrets
        all_secrets = get_all_secrets()
        print(f"Available keys: {list(all_secrets.keys())}")
        
        # Get status
        status = get_status()
        print(f"Vault status: {'Unlocked' if not status['is_locked'] else 'Locked'}")
        print(f"Number of keys: {len(status['keys'])}")
        
    except Exception as e:
        print(f"Error: {e}")

if __name__ == "__main__":
    main()
