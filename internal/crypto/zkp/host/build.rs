fn main() -> Result<(), Box<dyn std::error::Error>> {
    // Tells cargo to re-run this if the proto changes
    println!("cargo:rerun-if-changed=../proto/prover.proto");
    
    tonic_build::configure()
        .build_server(true)
        .compile(
            &["../proto/prover.proto"], 
            &["../proto"]
        )?;
    Ok(())
}