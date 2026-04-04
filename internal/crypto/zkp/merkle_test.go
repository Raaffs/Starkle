package zkp

import (
	"log"
	"sort"
	"testing"

	"github.com/Suy56/ProofChain/internal/models"
)

func TestVerifyProof_Scenarios(t *testing.T) {
	// 1. Setup Initial Data
	input := models.CertificateBase[any]{
		Name:            "alice",
		CertificateName: "Master of Electrical and ElectronicsEngineering",
		BirthDate:       "1990-01-01",
		Address:         1,
		Age: 		   	 23,
		PublicAddress: 	"fefefefe",
		Extra:           map[string]any{"MembershipID": "1225789"},
	}

	merkle := NewMerkleProof()
	root, saltedCert, err := merkle.GenerateRootProof(input)
	if err != nil {
		t.Fatalf("Failed to generate root: %v", err)
	}
	log.Println("salted cert: ",saltedCert)
	extractHashes := func(s SaltedCertificate) []Hash {
		var h []Hash
		var keys []string
		for field := range s.SaltedFields {
			keys = append(keys, field)
		}
		sort.Strings(keys)
		for _, key := range keys {
			h = append(h, Hash(s.SaltedFields[key].Hash))
		}
		return h
	}

	allLeaves := extractHashes(saltedCert)

	t.Run("Valid Disclosure", func(t *testing.T) {
		p := ProofVerification{
			RootHash:    root,
			Attribute:   "Name",
			Value:       saltedCert.SaltedFields["Name"].Value,
			Salt:        saltedCert.SaltedFields["Name"].Salt,
			MerkleProof: allLeaves,
		}
		if ok := VerifyProof(p, root); !ok {
			t.Errorf("Expected valid proof to pass")
		}
	})

	t.Run("Tampered Value", func(t *testing.T) {
		p := ProofVerification{
			RootHash:    root,
			Attribute:   "Name",
			Value:       "bob", // Tampered
			Salt:        saltedCert.SaltedFields["Name"].Salt,
			MerkleProof: allLeaves,
		}
		if ok := VerifyProof(p, root); ok {
			t.Errorf("Expected failure: Value does not match any leaf hash")
		}
	})

	t.Run("Tampered Salt", func(t *testing.T) {
		p := ProofVerification{
			RootHash:    root,
			Attribute:   "Name",
			Value:       saltedCert.SaltedFields["Name"].Value,
			Salt:        "malicious_salt", // Tampered
			MerkleProof: allLeaves,
		}
		if ok := VerifyProof(p, root); ok {
			t.Errorf("Expected failure: Hash(Value+FakeSalt) not in leaf set")
		}
	})

	t.Run("Tampered Leaf Set", func(t *testing.T) {
		// Create a copy of leaves and swap one out
		maliciousLeaves := make([]Hash, len(allLeaves))
		copy(maliciousLeaves, allLeaves)
		maliciousLeaves[0] = Hash("00000000000000000000000000000000")

		p := ProofVerification{
			RootHash:    root,
			Attribute:   "Name",
			Value:       saltedCert.SaltedFields["Name"].Value,
			Salt:        saltedCert.SaltedFields["Name"].Salt,
			MerkleProof: maliciousLeaves,
		}
		if ok := VerifyProof(p, root); ok {
			t.Errorf("Expected failure: Tampered leaf set should produce wrong root")
		}
	})

	t.Run("Malicious institute and requestor roots", func(t *testing.T) {
		p := ProofVerification{
			//fake root
			RootHash:    Hash("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"),
			Attribute:   "Name",
			Value:       saltedCert.SaltedFields["Name"].Value,
			Salt:        saltedCert.SaltedFields["Name"].Salt,
			MerkleProof: allLeaves,
		}
		if ok := VerifyProof(p, Hash("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")); ok {
			t.Errorf("Expected failure: Proof is valid but root hash doesn't match expectedRoot")
		}
	})
}