use std::fs;
use std::path::Path;
use tonic::{transport::Server, Request, Response, Status};
use sha2::{Sha256, Digest as ShaDigest};

use methods::{MEMBERSHIP_ELF, MEMBERSHIP_ID, RANGE_ELF, RANGE_ID};
use risc0_zkvm::{get_prover_server, ExecutorEnv, ProverOpts, Receipt, sha::Digest};

pub mod prover {
    tonic::include_proto!("prover");
}

use prover::prover_service_server::{ProverService, ProverServiceServer};
use prover::{ProofRequest, ProofResponse, VerifyRequest, VerifyResponse, proof_request::ProofData};

pub struct ProverHost;

impl ProverHost {
    fn save_receipt(receipt: &Receipt) -> String {
        let receipt_bytes = bincode::serialize(receipt).unwrap();
        let mut hasher = Sha256::new();
        hasher.update(&receipt_bytes);
        let id = hex::encode(hasher.finalize());
        let dir = "receipts";
        if !Path::new(dir).exists() { 
            fs::create_dir_all(dir).expect("Failed to create receipts directory"); 
        }
        fs::write(format!("{}/{}.bin", dir, id), &receipt_bytes).expect("Failed to write receipt");
        id
    }
}

#[tonic::async_trait]
impl ProverService for ProverHost {
    async fn generate_proof(&self, request: Request<ProofRequest>) -> Result<Response<ProofResponse>, Status> {
        let req = request.into_inner();
        
        match req.proof_data {
            Some(ProofData::Membership(m)) => {
                println!("🚀 Received Membership Proof Request (Path length: {})", m.all_leaves.len());
                
                // We pass the root as a string because the Guest is reconstructing the 
                // root as a hex string to stay 1:1 compatible with your Go logic.
                let env = ExecutorEnv::builder()
                    .write(&m.actual_value).unwrap()
                    .write(&m.actual_salt).unwrap() 
                    .write(&m.all_leaves).unwrap()  // This is now the Siblings Path
                    .write(&m.public_list).unwrap()
                    .write(&m.public_root).unwrap() // The target root hex string
                    .build()
                    .map_err(|e| Status::internal(format!("Env build failed: {}", e)))?;

                println!("  -> Starting Prover (CPU/GPU Auto-fallback)...");
                let prover = get_prover_server(&ProverOpts::default())
                    .map_err(|e| Status::internal(format!("Prover init failed: {}", e)))?;
                
                let prove_info = prover.prove(env, MEMBERSHIP_ELF)
                    .map_err(|e| Status::internal(format!("Proving failed: {}", e)))?;
                
                let cycles = prove_info.stats.total_cycles;
                println!("  ✅ Membership Proof Generated! Cycles: {}", cycles);
                
                let receipt_bytes = bincode::serialize(&prove_info.receipt)
                    .map_err(|e| Status::internal(format!("Serialization failed: {}", e)))?;
                
                let receipt_id = Self::save_receipt(&prove_info.receipt);
                
                Ok(Response::new(ProofResponse {
                    receipt_id,
                    cycles: cycles as u32,
                    receipt_bytes,
                }))
            },

            Some(ProofData::Range(r)) => {
                println!("🚀 Received Range Proof Request ({} leaves)", r.all_leaves.len());

                let env = ExecutorEnv::builder()
                    .write(&r.actual_value).unwrap()
                    .write(&r.actual_salt).unwrap() 
                    .write(&r.all_leaves).unwrap() 
                    .write(&r.lower_bound).unwrap()
                    .write(&r.upper_bound).unwrap()
                    .write(&r.public_root).unwrap()
                    .build()
                    .map_err(|e| Status::internal(format!("Env build failed: {}", e)))?;

                println!("  -> Starting Prover...");
                let prover = get_prover_server(&ProverOpts::default())
                    .map_err(|e| Status::internal(format!("Prover init failed: {}", e)))?;
                
                let prove_info = prover.prove(env, RANGE_ELF)
                    .map_err(|e| Status::internal(format!("Proving failed: {}", e)))?;
                
                let cycles = prove_info.stats.total_cycles;
                println!("  ✅ Range Proof Generated! Cycles: {}", cycles);
                
                let receipt_bytes = bincode::serialize(&prove_info.receipt)
                    .map_err(|e| Status::internal(format!("Serialization failed: {}", e)))?;
                
                let receipt_id = Self::save_receipt(&prove_info.receipt);
                
                Ok(Response::new(ProofResponse {
                    receipt_id,
                    cycles: cycles as u32,
                    receipt_bytes,
                }))
            },

            _ => Err(Status::unimplemented("Proof type not implemented")),
        }
    }

    async fn verify_proof(&self, request: Request<VerifyRequest>) -> Result<Response<VerifyResponse>, Status> {
        let req = request.into_inner();
        let path = format!("receipts/{}.bin", req.receipt_id);
        
        let bytes = fs::read(path).map_err(|_| Status::not_found("Receipt file not found"))?;
        let receipt: Receipt = bincode::deserialize(&bytes)
            .map_err(|e| Status::internal(format!("Receipt deserialization failed: {}", e)))?;
        
        // Try verifying against both known Image IDs
        let result = receipt.verify(MEMBERSHIP_ID)
            .or_else(|_| receipt.verify(RANGE_ID));

        match result {
            Ok(_) => {
                println!("✅ Verification Successful for receipt {}", req.receipt_id);
                Ok(Response::new(VerifyResponse { valid: true, error_msg: "".into() }))
            },
            Err(e) => {
                println!("❌ Verification Failed: {}", e);
                Ok(Response::new(VerifyResponse { valid: false, error_msg: e.to_string() }))
            },
        }
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = "127.0.0.1:50051".parse()?; 
    println!("✅ ZKP Prover gRPC Server running on {}", addr);

    // 100MB message limit
    let max_msg_size = 100 * 1024 * 1024; 

    Server::builder()
        .initial_connection_window_size(Some(max_msg_size as u32))
        .initial_stream_window_size(Some(max_msg_size as u32))
        .add_service(
            ProverServiceServer::new(ProverHost)
                .max_decoding_message_size(max_msg_size)
                .max_encoding_message_size(max_msg_size)
        )
        .serve(addr)
        .await?;

    Ok(())
}