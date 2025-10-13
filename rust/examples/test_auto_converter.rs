use encrypt::sdk::SDK;
use std::env;

fn main() {
    println!("🧪 Testing Auto-Converter Engine (Rust)");
    println!("{}", "=".repeat(40));
    println!();
    
    // Test 1: Try to get secret when vault is locked (should fail without password)
    println!("1. Testing locked vault without password...");
    match SDK::get("API_KEY", None) {
        Ok(api_key) => {
            println!("❌ Unexpected: Got secret from locked vault: {}", api_key);
        }
        Err(e) => {
            println!("✅ Expected: Vault is locked - {}", e);
        }
    }
    
    // Test 2: Try with environment variable
    println!("\n2. Testing with ENCRYPT_PASSWORD environment variable...");
    env::set_var("ENCRYPT_PASSWORD", "mypassword");
    
    match SDK::get("API_KEY", None) {
        Ok(api_key) => {
            println!("✅ Success: Got secret with environment password - {}", api_key);
        }
        Err(e) => {
            println!("❌ Failed: Could not get secret with environment password - {}", e);
        }
    }
    
    // Test 3: Try with explicit password parameter
    println!("\n3. Testing with explicit password parameter...");
    env::remove_var("ENCRYPT_PASSWORD");  // Clear env var
    
    match SDK::get("API_KEY", Some("mypassword")) {
        Ok(api_key) => {
            println!("✅ Success: Got secret with explicit password - {}", api_key);
        }
        Err(e) => {
            println!("❌ Failed: Could not get secret with explicit password - {}", e);
        }
    }
    
    // Test 4: Test production-ready functions
    println!("\n4. Testing production-ready functions...");
    env::set_var("ENCRYPT_PASSWORD", "mypassword");
    
    match SDK::get_secret("API_KEY") {
        Ok(api_key) => {
            println!("✅ Success: get_secret() works - {}", api_key);
        }
        Err(e) => {
            println!("❌ Failed: get_secret() failed - {}", e);
        }
    }
    
    match SDK::get_secret("DB_URL") {
        Ok(db_url) => {
            println!("✅ Success: get_secret() works - {}", db_url);
        }
        Err(e) => {
            println!("❌ Failed: get_secret() failed - {}", e);
        }
    }
    
    // Test 5: Test development mode (should try common passwords)
    println!("\n5. Testing development mode...");
    env::remove_var("ENCRYPT_PASSWORD");
    env::set_var("NODE_ENV", "development");
    
    match SDK::get("API_KEY", None) {
        Ok(api_key) => {
            println!("✅ Success: Development mode worked - {}", api_key);
        }
        Err(e) => {
            println!("❌ Failed: Development mode failed - {}", e);
        }
    }
    
    println!("\n🎉 Auto-converter tests completed!");
}
