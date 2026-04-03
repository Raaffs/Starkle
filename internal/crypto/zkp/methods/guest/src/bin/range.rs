#![no_main]
#![no_std]

extern crate alloc;
use alloc::vec::Vec;
use alloc::string::String;
// 1. IMPORT THIS: Required for write! to work on Strings
use core::fmt::Write; 

use risc0_zkvm::guest::env;
use risc0_zkvm::sha::{Impl, Sha256, Digest};

risc0_zkvm::guest::entry!(main);

fn to_hex_string(bytes: &[u8]) -> String {
    let mut s = String::with_capacity(bytes.len() * 2);
    for &b in bytes {
        // This now works because core::fmt::Write is in scope
        write!(s, "{:02x}", b).unwrap();
    }
    s
}

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

fn main() {
    let value: u32 = env::read();
    let salt_hex: String = env::read();
    let all_leaves_hex: Vec<String> = env::read();
    
    let lower_bound: u32 = env::read();
    let upper_bound: u32 = env::read();

    assert!(value >= lower_bound && value <= upper_bound, "Value not in range");

    let mut hasher_input = Vec::new();
    // 2. ADD AMPERSAND: Borrow the array as a slice
    hasher_input.extend_from_slice(&value.to_le_bytes()); 
    hasher_input.extend_from_slice(salt_hex.as_bytes());

    let raw_leaf_hash = *Impl::hash_bytes(&hasher_input);
    let calculated_leaf_hex = to_hex_string(raw_leaf_hash.as_bytes());

    assert!(
        all_leaves_hex.iter().any(|l| l == &calculated_leaf_hex),
        "Leaf mismatch!"
    );

    let mut current_level = all_leaves_hex;
    while current_level.len() > 1 {
        if current_level.len() % 2 != 0 {
            let last = current_level.last().unwrap().clone();
            current_level.push(last);
        }

        let mut next_level = Vec::with_capacity(current_level.len() / 2);
        for i in (0..current_level.len()).step_by(2) {
            let h1 = &current_level[i];
            let h2 = &current_level[i + 1];
            let (first, second) = if h1 < h2 { (h1, h2) } else { (h2, h1) };
            let parent_digest = hash_data_go_style(first, second);
            next_level.push(to_hex_string(parent_digest.as_bytes()));
        }
        current_level = next_level;
    }

    let final_root_bytes = decode_hex(&current_level[0]);
    let final_root_digest = Digest::try_from(final_root_bytes.as_slice()).unwrap();

    // 3. FIXED: Replaced 'public_list' with the actual bounds 
    // (or whatever public data you need to prove against)
    env::commit(&lower_bound);
    env::commit(&upper_bound);
    env::commit(&final_root_digest);
}