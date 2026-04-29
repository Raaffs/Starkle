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
    /// Saves the receipt to a specific directory with a descriptive filename.
    /// base_path: Absolute path from the request (e.g., /home/user/...)
    /// filename: Constraint-based name (e.g., "list_item1_item2" or "lb10_ub20")
    fn save_receipt(receipt: &Receipt, base_path: &str, filename: &str) -> Result<String, Status> {
        let receipt_bytes = bincode::serialize(receipt)
            .map_err(|e| Status::internal(format!("Serialization failed: {}", e)))?;
        
        let path = Path::new(base_path);
        if !path.exists() { 
            fs::create_dir_all(path)
                .map_err(|e| Status::internal(format!("Failed to create directory: {}", e)))?; 
        }

        // sanitize filename slightly to ensure no OS conflicts
        let safe_filename = filename.replace(|c: char| !c.is_alphanumeric() && c != '_' && c != '-', "_");
        let full_path = path.join(format!("{}.bin", safe_filename));
        
        fs::write(&full_path, &receipt_bytes)
            .map_err(|e| Status::internal(format!("Failed to write receipt: {}", e)))?;
        
        Ok(full_path.to_string_lossy().into_owned())
    }
}

#[tonic::async_trait]
impl ProverService for ProverHost {
    async fn generate_proof(&self, request: Request<ProofRequest>) -> Result<Response<ProofResponse>, Status> {
        let req = request.into_inner();
        
        match req.proof_data {
            Some(ProofData::Membership(m)) => {
                println!("🚀 Membership Request. Path: {}, List items: {}", m.path, m.public_list.len());
                
                let env = ExecutorEnv::builder()
                    .write(&m.actual_value).unwrap()
                    .write(&m.actual_salt).unwrap() 
                    .write(&m.siblings).unwrap()  
                    .write(&m.public_list).unwrap()
                    .write(&m.public_root).unwrap() 
                    .build()
                    .map_err(|e| Status::internal(format!("Env build failed: {}", e)))?;

                let prover = get_prover_server(&ProverOpts::default())
                    .map_err(|e| Status::internal(format!("Prover init failed: {}", e)))?;
                
                let prove_info = prover.prove(env, MEMBERSHIP_ELF)
                    .map_err(|e| Status::internal(format!("Proving failed: {}", e)))?;
                
                // Filename based on public list (joined items)
                let filename = if m.public_list.is_empty() {
                    "empty_list".to_string()
                } else {
                    m.public_list.join("_")
                };

                let receipt_path = Self::save_receipt(&prove_info.receipt, &m.path, &filename)?;
                
                Ok(Response::new(ProofResponse {
                    receipt_id: receipt_path,
                    cycles: prove_info.stats.total_cycles as u32,
                    receipt_bytes: bincode::serialize(&prove_info.receipt).unwrap(),
                }))
            },

            Some(ProofData::Range(r)) => {
                println!("🚀 Range Request. Path: {}, Range: {} - {}", r.path, r.lower_bound, r.upper_bound);

                let env = ExecutorEnv::builder()
                    .write(&r.actual_value).unwrap()
                    .write(&r.actual_salt).unwrap() 
                    .write(&r.siblings).unwrap() 
                    .write(&r.lower_bound).unwrap()
                    .write(&r.upper_bound).unwrap()
                    .write(&r.public_root).unwrap()
                    .build()
                    .map_err(|e| Status::internal(format!("Env build failed: {}", e)))?;

                let prover = get_prover_server(&ProverOpts::default())
                    .map_err(|e| Status::internal(format!("Prover init failed: {}", e)))?;
                
                let prove_info = prover.prove(env, RANGE_ELF)
                    .map_err(|e| Status::internal(format!("Proving failed: {}", e)))?;
                
                // Filename based on lower/upper bounds
                let filename = format!("lb{}_ub{}", r.lower_bound, r.upper_bound);
                let receipt_path = Self::save_receipt(&prove_info.receipt, &r.path, &filename)?;
                
                Ok(Response::new(ProofResponse {
                    receipt_id: receipt_path,
                    cycles: prove_info.stats.total_cycles as u32,
                    receipt_bytes: bincode::serialize(&prove_info.receipt).unwrap(),
                }))
            },

            _ => Err(Status::unimplemented("Proof type not implemented")),
        }
    }

    async fn verify_proof(&self, request: Request<VerifyRequest>) -> Result<Response<VerifyResponse>, Status> {
        let req = request.into_inner();
        
        // receipt_id is now the full path returned by generate_proof
        let bytes = fs::read(&req.receipt_id)
            .map_err(|_| Status::not_found(format!("Receipt file not found at {}", req.receipt_id)))?;
            
        let receipt: Receipt = bincode::deserialize(&bytes)
            .map_err(|e| Status::internal(format!("Receipt deserialization failed: {}", e)))?;
        
        let result = receipt.verify(MEMBERSHIP_ID)
            .or_else(|_| receipt.verify(RANGE_ID));

        match result {
            Ok(_) => {
                println!("✅ Verified: {}", req.receipt_id);
                Ok(Response::new(VerifyResponse { valid: true, error_msg: "".into() }))
            },
            Err(e) => {
                println!("❌ Failed: {}", e);
                Ok(Response::new(VerifyResponse { valid: false, error_msg: e.to_string() }))
            },
        }
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = "127.0.0.1:50051".parse()?; 
    println!("✅ ZKP Prover gRPC Server running on {}", addr);

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