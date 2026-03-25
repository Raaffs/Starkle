package zkp

import (
	"github.com/Suy56/ProofChain/internal/models"
)


// FieldLeaf stores the necessary data for a single attribute leaf node

// SaltedCertificate is the data structure that the Issuer sends to the Requestor.
// It contains all field data and salts, allowing the Requestor to reconstruct the tree.
type SaltedCertificate struct {
	SaltedFields map[string]models.LeafFields `json:"salted_fields"`
}


type ZKProof interface {
	New() 
	GenerateRootProof(c models.CertificateBase[string]) (Hash, SaltedCertificate, error)
}

// Proof contains the components needed for a third-party verifier.
// This is returned by the Disclose function.
type ProofVerification struct {
	RootHash    Hash     `json:"root_hash"`    // The committed hash the verifier checks against
	Attribute   string   `json:"attribute"`    // The name of the field being disclosed
	Value       string   `json:"value"`        // The disclosed field value
	Salt        string   `json:"salt"`         // The salt used to generate the leaf hash
	MerkleProof []Hash   `json:"merkle_proof"` // The ordered list of sibling hashes for verification
}
