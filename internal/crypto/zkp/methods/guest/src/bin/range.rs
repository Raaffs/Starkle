#![no_main]
#![no_std]

extern crate alloc;
use alloc::vec::Vec;
use alloc::string::String;
use core::fmt::Write; 

use risc0_zkvm::guest::env;
use risc0_zkvm::sha::{Impl, Sha256, Digest};

risc0_zkvm::guest::entry!(main);

// Helper to convert hex string back to raw bytes
fn decode_hex(s: &str) -> Vec<u8> {
    let mut res = Vec::with_capacity(s.len() / 2);
    let bytes = s.as_bytes();
    for i in (0..s.len()).step_by(2) {
        let high = char::from(bytes[i]).to_digit(16).expect("Invalid hex high") as u8;
        let low = char::from(bytes[i+1]).to_digit(16).expect("Invalid hex low") as u8;
        res.push((high << 4) | low);
    }
    res
}

fn hash_data_go_style(h1: &str, h2: &str) -> Digest {
    let mut combined = String::with_capacity(h1.len() + h2.len());
    combined.push_str(h1);
    combined.push_str(h2);
    *Impl::hash_bytes(combined.as_bytes())
}


fn to_hex_string(bytes: &[u8]) -> String {
    let mut s = String::with_capacity(bytes.len() * 2);
    for &b in bytes {
        // This now works because core::fmt::Write is in scope
        write!(s, "{:02x}", b).unwrap();
    }
    s
}

fn digest_to_hex(digest: &Digest) -> String {
    let mut s = String::with_capacity(64);
    for &b in digest.as_bytes() {
        write!(s, "{:02x}", b).unwrap();
    }
    s
}


fn main() {
    let value: u32 = env::read();
    let salt_hex: String = env::read();
    let siblings: Vec<String> = env::read(); 
    let lower_bound: u32 = env::read();
    let upper_bound: u32 = env::read();
    let expected_root_hex: String = env::read();

    assert!(value >= lower_bound && value <= upper_bound, "Value not in range");

    // 1. Calculate Leaf Hash (Raw Bytes)
    let mut hasher_input = Vec::new();
    hasher_input.extend_from_slice(&value.to_le_bytes()); 
    hasher_input.extend_from_slice(salt_hex.as_bytes());


    let leaf_digest = *Impl::hash_bytes(&hasher_input);

    let mut current_hash_hex = to_hex_string(leaf_digest.as_bytes());

    env::log(&current_hash_hex);
    for sibling in siblings {
        let mut combined = String::with_capacity(128);
        if current_hash_hex < sibling {
            combined.push_str(&current_hash_hex);
            combined.push_str(&sibling);
        } else {
            combined.push_str(&sibling);
            combined.push_str(&current_hash_hex);
        }
        current_hash_hex = digest_to_hex(&Impl::hash_bytes(combined.as_bytes()));
    }

    // 3. Final Check
    assert_eq!(current_hash_hex, expected_root_hex, "Merkle Root Mismatch!");

    env::commit(&lower_bound);
    env::commit(&upper_bound);
    env::commit(&current_hash_hex);
}