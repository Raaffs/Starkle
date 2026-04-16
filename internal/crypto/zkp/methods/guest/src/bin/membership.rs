#![no_main]
#![no_std]

extern crate alloc;
use alloc::string::String;
use alloc::vec::Vec;
use core::fmt::Write;
use risc0_zkvm::guest::env;
use risc0_zkvm::sha::{Impl, Sha256, Digest};

risc0_zkvm::guest::entry!(main);

fn digest_to_hex(digest: &Digest) -> String {
    let mut s = String::with_capacity(64);
    for &b in digest.as_bytes() {
        write!(s, "{:02x}", b).unwrap();
    }
    s
}

fn main() {
    let value: String = env::read();
    let salt_hex: String = env::read();
    let siblings: Vec<String> = env::read(); // Only ~3-8 strings
    let public_list: Vec<String> = env::read();
    let expected_root_hex: String = env::read();

    // 1. Logic Check
    assert!(public_list.iter().any(|x| x == &value), "Value not in public list");

    // 2. Initial Leaf Hash: sha256(value + salt_hex)
    let mut hasher_input = Vec::new();
    hasher_input.extend_from_slice(value.as_bytes());
    hasher_input.extend_from_slice(salt_hex.as_bytes());
    let leaf_digest = *Impl::hash_bytes(&hasher_input);
    let mut current_hash_hex = digest_to_hex(&leaf_digest);
    env::log(&current_hash_hex);

    // 3. Climb the tree using siblings
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

    // 4. Final Root Validation
    assert_eq!(current_hash_hex, expected_root_hex, "Merkle Root Mismatch!");

    env::commit(&public_list);
    let root_bytes = hex_decode(&current_hash_hex);
    env::commit(&Digest::try_from(root_bytes.as_slice()).unwrap());
}

fn hex_decode(s: &str) -> Vec<u8> {
    let mut res = Vec::with_capacity(s.len() / 2);
    let bytes = s.as_bytes();
    for i in (0..s.len()).step_by(2) {
        let hi = char::from(bytes[i]).to_digit(16).unwrap() as u8;
        let lo = char::from(bytes[i+1]).to_digit(16).unwrap() as u8;
        res.push((hi << 4) | lo);
    }
    res
}