#!/usr/bin/env python3
"""
CLI interface for the Encrypt tool.

Provides command-line access to all vault operations with rich formatting.
"""

import sys
import os
import json
from typing import Optional
import click
from rich.console import Console
from rich.prompt import Prompt, Confirm
from rich.table import Table
from rich import print as rprint

from .vault import Vault
from .crypto import TripleEncryption

# Global vault instance
_global_vault = None

console = Console()


def get_vault() -> Vault:
    """Get or create the global vault instance."""
    global _global_vault
    if _global_vault is None:
        _global_vault = Vault()
    return _global_vault


@click.group()
@click.version_option(version="1.0.0")
def cli():
    """🔐 Encrypt - A top-level secrets orchestrator.
    
    Not just another .env tool — this one encrypts, locks, and sets you up for secure local and team dev.
    """
    pass


@cli.command()
def init():
    """Create .encrypt vault."""
    with console.status("Initializing vault..."):
        try:
            vault = get_vault()
            vault.init()
            console.print("✅ [green]Vault initialized successfully![/green]")
            console.print("📁 [blue]Created .encrypt directory with secure configuration.[/blue]")
        except Exception as e:
            console.print(f"❌ [red]Failed to initialize vault[/red]")
            console.print(f"[red]{str(e)}[/red]")
            sys.exit(1)


@cli.command()
@click.argument('password')
def lockup(password: str):
    """Encrypt and secure secrets with password."""
    with console.status("Locking up secrets..."):
        try:
            vault = get_vault()
            
            if not vault.exists():
                console.print("❌ [red]Vault not found. Run 'encrypt init' first.[/red]")
                sys.exit(1)
            
            vault.lockup(password)
            console.print("✅ [green]Secrets locked successfully![/green]")
            console.print("🔒 [blue]Your secrets are now encrypted and secure.[/blue]")
        except Exception as e:
            console.print("❌ [red]Failed to lock secrets[/red]")
            console.print(f"[red]{str(e)}[/red]")
            sys.exit(1)


@cli.command()
@click.argument('password')
def setup(password: str):
    """Set up secrets on a new machine."""
    with console.status("Setting up vault..."):
        try:
            vault = get_vault()
            
            if not vault.exists():
                console.print("❌ [red]Vault not found. Run 'encrypt init' first.[/red]")
                sys.exit(1)
            
            vault.setup(password)
            console.print("✅ [green]Vault unlocked successfully![/green]")
            console.print("🔓 [blue]Your secrets are now available for use.[/blue]")
        except Exception as e:
            console.print("❌ [red]Failed to setup vault[/red]")
            console.print(f"[red]{str(e)}[/red]")
            sys.exit(1)


@cli.command()
@click.argument('key_value')
def set(key_value: str):
    """Add/update a key."""
    if '=' not in key_value:
        console.print("❌ [red]Invalid format. Use: encrypt set KEY=value[/red]")
        sys.exit(1)
    
    key, value = key_value.split('=', 1)
    
    with console.status(f"Setting secret: {key}"):
        try:
            vault = get_vault()
            
            if not vault.exists():
                console.print("❌ [red]Vault not found. Run 'encrypt init' first.[/red]")
                sys.exit(1)
            
            if not vault.is_unlocked_status():
                console.print("❌ [red]Vault is locked. Run 'encrypt setup <password>' to unlock secrets.[/red]")
                sys.exit(1)
            
            vault.set(key, value)
            console.print(f"✅ [green]Secret '{key}' set successfully![/green]")
        except Exception as e:
            console.print("❌ [red]Failed to set secret[/red]")
            console.print(f"[red]{str(e)}[/red]")
            sys.exit(1)


@cli.command()
@click.argument('key')
def get(key: str):
    """Fetch decrypted value."""
    try:
        vault = get_vault()
        
        if not vault.exists():
            console.print("❌ [red]Vault not found. Run 'encrypt init' first.[/red]")
            sys.exit(1)
        
        if not vault.is_unlocked_status():
            console.print("❌ [red]Vault is locked. Run 'encrypt setup <password>' to unlock secrets.[/red]")
            sys.exit(1)
        
        value = vault.get(key)
        console.print(value)
    except Exception as e:
        console.print(f"[red]{str(e)}[/red]")
        sys.exit(1)


@cli.command()
def unlock():
    """Decrypt everything into .env."""
    with console.status("Unlocking secrets..."):
        try:
            vault = get_vault()
            
            if not vault.exists():
                console.print("❌ [red]Vault not found. Run 'encrypt init' first.[/red]")
                sys.exit(1)
            
            if not vault.is_unlocked_status():
                console.print("❌ [red]Vault is locked. Run 'encrypt setup <password>' to unlock secrets.[/red]")
                sys.exit(1)
            
            secrets = vault.all()
            env_content = '\n'.join([f"{k}={v}" for k, v in secrets.items()])
            
            with open('.env', 'w') as f:
                f.write(env_content)
            
            console.print("✅ [green]Secrets unlocked and written to .env[/green]")
        except Exception as e:
            console.print("❌ [red]Failed to unlock secrets[/red]")
            console.print(f"[red]{str(e)}[/red]")
            sys.exit(1)


@cli.command()
def status():
    """Check if vault is locked, list keys."""
    try:
        vault = get_vault()
        
        if not vault.exists():
            console.print("⚠️ [yellow]Vault not found. Run 'encrypt init' first.[/yellow]")
            return
        
        status_info = vault.status()
        
        console.print("📊 [blue]Vault Status:[/blue]")
        
        # Status table
        table = Table(show_header=True, header_style="bold blue")
        table.add_column("Property", style="cyan")
        table.add_column("Value", style="white")
        
        status_icon = "🔓" if not status_info['is_locked'] else "🔒"
        status_text = "Unlocked" if not status_info['is_locked'] else "Locked"
        status_color = "green" if not status_info['is_locked'] else "red"
        
        table.add_row("Status", f"{status_icon} [{status_color}]{status_text}[/{status_color}]")
        table.add_row("Keys", str(len(status_info['keys'])))
        
        if status_info['last_modified']:
            table.add_row("Last Modified", status_info['last_modified'])
        
        console.print(table)
        
        if status_info['keys']:
            console.print("\n🔑 [blue]Available keys:[/blue]")
            for key in status_info['keys']:
                console.print(f"  • {key}")
        
    except Exception as e:
        console.print(f"[red]{str(e)}[/red]")
        sys.exit(1)


@cli.command()
def reset():
    """Remove vault (careful!)."""
    if not Confirm.ask("Are you sure you want to reset the vault? This will delete all encrypted secrets."):
        console.print("⚠️ [yellow]Reset cancelled.[/yellow]")
        return
    
    with console.status("Resetting vault..."):
        try:
            vault = get_vault()
            vault.reset()
            console.print("✅ [green]Vault reset successfully![/green]")
        except Exception as e:
            console.print("❌ [red]Failed to reset vault[/red]")
            console.print(f"[red]{str(e)}[/red]")
            sys.exit(1)


if __name__ == '__main__':
    cli()
