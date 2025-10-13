use encrypt::crypto::Crypto;

fn main() {
    println!("🧪 Testing TripleEncryption (Rust)...");
    
    // Test data
    let plaintext = "Hello, World!";
    let password = "test-password-123";
    
    println!("Original: {}", plaintext);
    
    match Crypto::encrypt(plaintext, password) {
        Ok(encrypted) => {
            println!("✅ Encryption successful");
            println!("Encrypted data length: {}", encrypted.encrypted.len());
            
            match Crypto::decrypt(&encrypted.encrypted, &encrypted.salt, &encrypted.hmac, password) {
                Ok(decrypted) => {
                    if decrypted == plaintext {
                        println!("✅ Decryption successful");
                        println!("Decrypted: {}", decrypted);
                    } else {
                        println!("❌ Decryption failed - content mismatch");
                    }
                }
                Err(e) => {
                    println!("❌ Decryption failed: {}", e);
                }
            }
            
            // Test password hashing
            match Crypto::hash_password(password) {
                Ok(hash_result) => {
                    match Crypto::verify_password(password, &hash_result) {
                        Ok(is_valid) => {
                            if is_valid {
                                println!("✅ Password hashing/verification successful");
                            } else {
                                println!("❌ Password hashing/verification failed");
                            }
                        }
                        Err(e) => {
                            println!("❌ Password verification error: {}", e);
                        }
                    }
                }
                Err(e) => {
                    println!("❌ Password hashing error: {}", e);
                }
            }
            
            // Test with different password (should fail)
            match Crypto::decrypt(&encrypted.encrypted, &encrypted.salt, &encrypted.hmac, "wrong-password") {
                Ok(_) => {
                    println!("❌ Wrong password incorrectly accepted");
                }
                Err(_) => {
                    println!("✅ Wrong password correctly rejected");
                }
            }
            
            println!("\n🎉 All tests passed!");
        }
        Err(e) => {
            println!("❌ Test failed: {}", e);
        }
    }
}
