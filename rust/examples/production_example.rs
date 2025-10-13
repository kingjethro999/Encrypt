use encrypt::{get_secret, set_secret, SDK};

fn main() {
    println!("🚀 Production Usage Example (Rust)");
    println!("{}", "=".repeat(35));
    println!();
    
    // Method 1: Using environment variable (Recommended for production)
    println!("Method 1: Environment Variable (Recommended)");
    println!("Set ENCRYPT_PASSWORD=your-password in your environment");
    println!();
    
    match get_secret("API_KEY") {
        Ok(api_key) => {
            println!("✅ Successfully retrieved secrets:");
            if api_key.len() > 4 {
                println!("API Key: ***{}", &api_key[api_key.len()-4..]);
            } else {
                println!("API Key: ***{}", api_key);
            }
            
            // Use in your application
            let config = Config {
                api_key: api_key.clone(),
                database: get_secret("DB_URL").unwrap_or_default(),
                port: std::env::var("PORT").unwrap_or_else(|_| "3000".to_string()),
            };
            
            println!("\n📋 Application config ready:");
            if config.api_key.len() > 4 {
                println!("  api_key: ***{}", &config.api_key[config.api_key.len()-4..]);
            } else {
                println!("  api_key: ***{}", config.api_key);
            }
            if config.database.len() > 10 {
                println!("  database: ***{}", &config.database[config.database.len()-10..]);
            } else {
                println!("  database: ***{}", config.database);
            }
            println!("  port: {}", config.port);
            
        }
        Err(e) => {
            println!("❌ Error: {}", e);
            println!("\n💡 To fix this:");
            println!("1. Set ENCRYPT_PASSWORD environment variable");
            println!("2. Or provide password as second parameter");
            println!("3. Or run 'encrypt setup <password>' first");
        }
    }
    
    println!("\n{}", "=".repeat(50));
    println!("Method 2: Explicit Password Parameter");
    println!("{}", "=".repeat(50));
    
    match SDK::get("API_KEY", Some("mypassword")) {
        Ok(api_key) => {
            if api_key.len() > 4 {
                println!("✅ Success with explicit password: ***{}", &api_key[api_key.len()-4..]);
            } else {
                println!("✅ Success with explicit password: ***{}", api_key);
            }
        }
        Err(e) => {
            println!("❌ Error with explicit password: {}", e);
        }
    }
    
    println!("\n{}", "=".repeat(50));
    println!("Method 3: Development Mode");
    println!("{}", "=".repeat(50));
    
    // In development, you can set NODE_ENV=development
    // and it will try common passwords automatically
    std::env::set_var("NODE_ENV", "development");
    
    match SDK::get("API_KEY", None) {
        Ok(api_key) => {
            if api_key.len() > 4 {
                println!("✅ Development mode success: ***{}", &api_key[api_key.len()-4..]);
            } else {
                println!("✅ Development mode success: ***{}", api_key);
            }
        }
        Err(e) => {
            println!("❌ Development mode failed: {}", e);
        }
    }
    
    println!("\n🎯 Production Deployment Examples:");
    println!("{}", "=".repeat(35));
    println!("Docker:");
    println!("  ENV ENCRYPT_PASSWORD=your-production-password");
    println!();
    println!("Kubernetes:");
    println!("  env:");
    println!("  - name: ENCRYPT_PASSWORD");
    println!("    valueFrom:");
    println!("      secretKeyRef:");
    println!("        name: encrypt-secrets");
    println!("        key: password");
    println!();
    println!("Heroku:");
    println!("  heroku config:set ENCRYPT_PASSWORD=your-password");
    println!();
    println!("AWS Lambda:");
    println!("  Set ENCRYPT_PASSWORD in environment variables");
    println!();
    println!("Rust application:");
    println!("  use encrypt::get_secret;");
    println!("  let api_key = get_secret(\"API_KEY\")?;");
    println!("  let database_url = get_secret(\"DATABASE_URL\")?;");
}

#[derive(Debug)]
struct Config {
    api_key: String,
    database: String,
    port: String,
}
