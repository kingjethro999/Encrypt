#!/usr/bin/env python3
"""
Complete demo of the Encrypt Python tool.
"""

import subprocess
import sys
import os
from rich.console import Console
from rich import print as rprint

console = Console()

def run_command(command, description):
    """Run a command and display the result."""
    console.print(f"[yellow]📝 {description}[/yellow]")
    console.print(f"[gray]$ {command}[/gray]")
    
    try:
        result = subprocess.run(command, shell=True, capture_output=True, text=True, cwd=os.getcwd())
        if result.stdout.strip():
            console.print(f"[green]{result.stdout.strip()}[/green]")
        if result.stderr.strip():
            console.print(f"[red]{result.stderr.strip()}[/red]")
    except Exception as e:
        console.print(f"[red]Error: {e}[/red]")
    
    console.print()

def main():
    """Run the complete demo."""
    console.print("[bold blue]🔐 Encrypt Tool Demo (Python)[/bold blue]")
    console.print("[gray]This demo shows the complete workflow of the encrypt Python tool.[/gray]\n")
    
    # Clean up any existing vault
    try:
        subprocess.run("rm -rf .encrypt .env", shell=True, capture_output=True)
    except:
        pass
    
    console.print("[bold blue]1. Initialize the vault[/bold blue]")
    run_command("encrypt init", "Create a new encrypted vault")
    
    console.print("[bold blue]2. Setup the vault (unlock for first time)[/bold blue]")
    run_command("encrypt setup demo-password", "Unlock the vault with a password")
    
    console.print("[bold blue]3. Add some secrets[/bold blue]")
    run_command("encrypt set API_KEY=sk-1234567890abcdef", "Add an API key")
    run_command("encrypt set DB_URL=postgres://user:pass@localhost:5432/mydb", "Add a database URL")
    run_command("encrypt set JWT_SECRET=super-secret-jwt-key", "Add a JWT secret")
    
    console.print("[bold blue]4. Check vault status[/bold blue]")
    run_command("encrypt status", "View vault status and available keys")
    
    console.print("[bold blue]5. Retrieve secrets[/bold blue]")
    run_command("encrypt get API_KEY", "Get the API key")
    run_command("encrypt get DB_URL", "Get the database URL")
    
    console.print("[bold blue]6. Test runtime SDK[/bold blue]")
    console.print("[yellow]📝 Using the encrypt package in Python code[/yellow]")
    console.print("[gray]$ python example.py[/gray]")
    try:
        result = subprocess.run([sys.executable, "example.py"], capture_output=True, text=True, cwd=os.getcwd())
        if result.stdout.strip():
            console.print(f"[green]{result.stdout.strip()}[/green]")
        if result.stderr.strip():
            console.print(f"[red]{result.stderr.strip()}[/red]")
    except Exception as e:
        console.print(f"[red]Error: {e}[/red]")
    console.print()
    
    console.print("[bold blue]7. Create .env file[/bold blue]")
    run_command("encrypt unlock", "Export all secrets to .env file")
    
    console.print("[bold blue]8. Lock the vault[/bold blue]")
    run_command("encrypt lockup demo-password", "Encrypt and lock all secrets")
    
    console.print("[bold blue]9. Try to access locked secrets[/bold blue]")
    run_command("encrypt get API_KEY", "Attempt to get secret from locked vault")
    
    console.print("[bold blue]10. Unlock the vault again[/bold blue]")
    run_command("encrypt setup demo-password", "Unlock the vault with the password")
    
    console.print("[bold blue]11. Verify secrets are accessible again[/bold blue]")
    run_command("encrypt get API_KEY", "Get the API key after unlocking")
    
    console.print("[bold green]✅ Demo completed successfully![/bold green]")
    console.print("[gray]\nThe encrypt Python tool provides:[/gray]")
    console.print("[gray]• Secure encryption with triple-layer protection[/gray]")
    console.print("[gray]• Easy team onboarding with single command[/gray]")
    console.print("[gray]• Runtime SDK for in-code secret access[/gray]")
    console.print("[gray]• Git-safe encrypted storage[/gray]")
    console.print("[gray]• Beautiful CLI with Rich formatting[/gray]")
    console.print("[gray]• Full Python package with setup.py[/gray]")

if __name__ == "__main__":
    main()
