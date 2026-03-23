use std::fs;
use std::path::Path;
use tonic::{transport::Server, Request, Response, Status};
use sha2::{Sha256, Digest as ShaDigest};

use methods::{MEMBERSHIP_ELF, MEMBERSHIP_ID};
use risc0_zkvm::{get_prover_server, ExecutorEnv, ProverOpts, Receipt, sha::Digest};

pub mod prover {
    tonic::include_proto!("prover");
}

use prover::prover_service_server::{ProverService, ProverServiceServer};
use prover::{ProofRequest, ProofResponse, VerifyRequest, VerifyResponse, proof_request::ProofData};

pub struct ProverHost;

impl ProverHost {
    fn hex_to_digest(hex_str: &str) -> Result<Digest, Status> {
        let bytes = hex::decode(hex_str).map_err(|_| Status::invalid_argument("Invalid hex"))?;
        let array: [u8; 32] = bytes.try_into().map_err(|_| Status::invalid_argument("Not 32 bytes"))?;
        Ok(Digest::from(array))
    }

    fn save_receipt(receipt: &Receipt) -> String {
        let receipt_bytes = bincode::serialize(receipt).unwrap();
        let mut hasher = Sha256::new();
        hasher.update(&receipt_bytes);
        let id = hex::encode(hasher.finalize());
        let dir = "receipts";
        if !Path::new(dir).exists() { fs::create_dir_all(dir).unwrap(); }
        fs::write(format!("{}/{}.bin", dir, id), &receipt_bytes).unwrap();
        id
    }
}

#[tonic::async_trait]
impl ProverService for ProverHost {
    async fn generate_proof(&self, request: Request<ProofRequest>) -> Result<Response<ProofResponse>, Status> {
        let req = request.into_inner();
        
        match req.proof_data {
            Some(ProofData::Membership(m)) => {
                let root = Self::hex_to_digest(&m.public_root)?;

                let env = ExecutorEnv::builder()
                    .write(&m.actual_value).unwrap()
                    .write(&m.actual_salt).unwrap() 
                    .write(&m.all_leaves).unwrap() 
                    .write(&m.public_list).unwrap()
                    .write(&root).unwrap()
                    .build()
                    .map_err(|e| Status::internal(e.to_string()))?;

                let prover = get_prover_server(&ProverOpts::fast()).unwrap();
                let prove_info = prover.prove(env, MEMBERSHIP_ELF)
                    .map_err(|e| Status::internal(format!("Prover failed: {}", e)))?;
                
                let receipt_id = Self::save_receipt(&prove_info.receipt);
                
                Ok(Response::new(ProofResponse {
                    receipt_id,
                    cycles: prove_info.stats.total_cycles as u32,
                    receipt_bytes: bincode::serialize(&prove_info.receipt).unwrap(),
                }))
            },
            _ => Err(Status::unimplemented("Not implemented")),
        }
    }

    async fn verify_proof(&self, request: Request<VerifyRequest>) -> Result<Response<VerifyResponse>, Status> {
        let req = request.into_inner();
        let path = format!("receipts/{}.bin", req.receipt_id);
        let bytes = fs::read(path).map_err(|_| Status::not_found("Receipt not found"))?;
        let receipt: Receipt = bincode::deserialize(&bytes).unwrap();
        
        match receipt.verify(MEMBERSHIP_ID) {
            Ok(_) => Ok(Response::new(VerifyResponse { valid: true, error_msg: "".into() })),
            Err(e) => Ok(Response::new(VerifyResponse { valid: false, error_msg: e.to_string() })),
        }
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = "[::1]:50051".parse()?;
    println!("✅ Server on {}", addr);
    Server::builder().add_service(ProverServiceServer::new(ProverHost)).serve(addr).await?;
    Ok(())
}