
#![no_main]
#![no_std]

extern  crate alloc;
use alloc::vec::Vec;
use alloc::string::String;
use risc0_zkvm::guest::env;
use risc0_zkvm::sha::{Impl, Sha256, Digest};

risc0_zkvm::guest::entry!(main);

fn main(){
    let value: String = env::read();
    let salt: [u32;8]=env::read();
    let mut all_leaves: Vec<Digest>=env::read();

    let lower_bound: u32=env::read();
    let upper_bound: u32=env::read();

    let public_root: [u32; 8]=env::read();
    let value_num: u32 = value.parse().unwrap();
    assert!(value_num>=lower_bound && value_num<=upper_bound, "Value not in range");


    let disclosed_leaf_hash = *Impl::hash_pair(
        &Digest::from(*Impl::hash_bytes(value.as_bytes())), 
        &Digest::from(salt)
    );
    assert!(all_leaves.contains(&disclosed_leaf_hash), "Hash of value and salt not in list of all leaves");

    let mut current_level= all_leaves;

    while current_level.len()>1{
        if current_level.len()%2==1{
            let last=*current_level.last().unwrap();
            current_level.push(last);
        }
        let mut next_level: Vec<Digest>=Vec::with_capacity(current_level.len()/2);
        for i in (0..current_level.len()).step_by(2){
            let h1=current_level[i];
            let h2=current_level[i+1];

            if h1<h2{
                next_level.push(*Impl::hash_pair(&h1, &h2));
            }else{
                next_level.push(*Impl::hash_pair(&h2, &h1));
            }
        }
        current_level=next_level
    }
    assert_eq!(current_level[0], Digest::from(public_root), "Root mismatch");
    env::commit(&lower_bound);
    env::commit(&upper_bound);
    env::commit(&public_root);
}