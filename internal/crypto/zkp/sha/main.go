package main

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

func main() {
	// 1. DATA TO HASH
	// The Private Birth Year (uint32)
	var birthYear uint32 = 1998
	
	// A real 32-byte salt (hex) - commonly used for identity masking
	saltHex := "a1b2c3d4e5f60718293a4b5c6d7e8f90a1b2c3d4e5f60718293a4b5c6d7e8f90"
	salt, _ := hex.DecodeString(saltHex)

	// 2. CONVERSION TO BYTES
	// Rust's u32::to_le_bytes() produces 4 bytes in Little-Endian order.
	yearBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(yearBytes, birthYear)

	// 3. GENERATE THE COMMITMENT
	// We hash (Year + Salt) in that specific order.
	hasher := sha256.New()
	hasher.Write(yearBytes)
	hasher.Write(salt)
	commitment := hasher.Sum(nil)

	// 4. PRINT FORMATTED RESULTS FOR RUST
	fmt.Println("--- COPY THESE INTO YOUR RUST HOST ---")
	fmt.Printf("let birth_year: u32 = %d;\n", birthYear)
	
	fmt.Print("let salt: [u8; 32] = [")
	for i, b := range salt {
		fmt.Printf("0x%02x", b)
		if i < 31 { fmt.Print(", ") }
	}
	fmt.Println("];")

	fmt.Printf("\n// This is your 'expected_commitment' for the Host/Guest\n")
	fmt.Print("let commitment: [u8; 32] = [")
	for i, b := range commitment {
		fmt.Printf("0x%02x", b)
		if i < 31 { fmt.Print(", ") }
	}
	fmt.Println("];")
	
	fmt.Printf("\nHex Commitment (for your DB): %x\n", commitment)
}