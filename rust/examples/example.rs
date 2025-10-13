use encrypt::{get_secret, get_all_secrets, get_status, is_vault_unlocked};

fn main() {
    println!("💻 Example usage of the encrypt package in Rust code");
    println!("{}", "=".repeat(50));
    
    // Check if vault is unlocked
    match is_vault_unlocked() {
        Ok(unlocked) => {
            if !unlocked {
                println!("❌ Vault is locked. Run 'encrypt setup <password>' to unlock secrets.");
                return;
            }
        }
        Err(e) => {
            println!("❌ Error checking vault status: {}", e);
            return;
        }
    }
    
    // Get secrets
    match get_secret("API_KEY") {
        Ok(api_key) => {
            println!("✅ Secrets retrieved successfully!");
            if api_key.len() > 4 {
                println!("API Key: ***{}", &api_key[api_key.len()-4..]);
            } else {
                println!("API Key: ***{}", api_key);
            }
        }
        Err(e) => {
            println!("⚠️ API_KEY not found: {}", e);
        }
    }
    
    match get_secret("DB_URL") {
        Ok(db_url) => {
            if db_url.len() > 10 {
                println!("DB URL: ***{}", &db_url[db_url.len()-10..]);
            } else {
                println!("DB URL: ***{}", db_url);
            }
        }
        Err(e) => {
            println!("⚠️ DB_URL not found: {}", e);
        }
    }
    
    // Get all secrets
    match get_all_secrets() {
        Ok(all_secrets) => {
            println!("Available keys: {:?}", all_secrets.keys().collect::<Vec<_>>());
        }
        Err(e) => {
            println!("⚠️ Error getting all secrets: {}", e);
        }
    }
    
    // Get status
    match get_status() {
        Ok(status) => {
            println!("Vault status: {}", if status.is_locked { "Locked" } else { "Unlocked" });
            println!("Number of keys: {}", status.keys.len());
        }
        Err(e) => {
            println!("⚠️ Error getting status: {}", e);
        }
    }
}
