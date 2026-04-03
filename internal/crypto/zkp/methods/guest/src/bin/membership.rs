#![no_main]
#![no_std]

extern crate alloc;
use alloc::string::String;
use alloc::vec::Vec;
use core::fmt::Write;
use risc0_zkvm::guest::env;
use risc0_zkvm::sha::{Impl, Sha256, Digest};

risc0_zkvm::guest::entry!(main);

/// Replicates Go's hex.EncodeToString
fn to_hex_string(bytes: &[u8]) -> String {
    let mut s = String::with_capacity(bytes.len() * 2);
    for &b in bytes {
        write!(s, "{:02x}", b).unwrap();
    }
    s
}

/// Decodes hex string to bytes safely in no_std
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

/// Replicates Go's HashData([]byte(h1), []byte(h2)) 
/// where h1 and h2 are hex strings.
fn hash_data_go_style(h1: &str, h2: &str) -> Digest {
    let mut combined = String::with_capacity(h1.len() + h2.len());
    combined.push_str(h1);
    combined.push_str(h2);
    *Impl::hash_bytes(combined.as_bytes())
}

fn main() {
    let value: String = env::read();
    let salt_hex: String = env::read();
    let all_leaves_hex: Vec<String> = env::read();
    let public_list: Vec<String> = env::read();
    // let public_root: Digest = env::read();

    // Public List Check
    assert!(public_list.iter().any(|x| x == &value), "Value not in public list");

    // Leaf Calculation: matches Go's HashData([]byte(value), []byte(saltHexString))
    // Go's salt is already a hex string — hash the hex string bytes directly, do NOT decode
    let mut hasher_input = Vec::new();
    hasher_input.extend_from_slice(value.as_bytes());
    hasher_input.extend_from_slice(salt_hex.as_bytes()); // ← KEY FIX: was decode_hex(&salt_hex)

    let raw_leaf_hash = *Impl::hash_bytes(&hasher_input);
    let calculated_leaf_hex = to_hex_string(raw_leaf_hash.as_bytes());

    assert!(
        all_leaves_hex.iter().any(|l| l == &calculated_leaf_hex),
        "Leaf mismatch!"
    );

    // Merkle Root Reconstruction (unchanged)
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

    env::commit(&public_list);
    env::commit(&final_root_digest);
}