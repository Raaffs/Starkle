package zkp

import (
	"encoding/binary"
	"slices"
	"sort"

	"github.com/Suy56/ProofChain/internal/models"
)

// MerkleProof implements the ZKProof interface and holds the committed data state.
type MerkleProof struct {
	RootHash   Hash
	FieldLeaves map[string]models.LeafFields // Map for O(1) lookup during Disclosure
	LeafHashes []Hash                // Ordered list for Merkle Tree construction
}

func NewMerkleProof() *MerkleProof {
	mp:=&MerkleProof{}
	mp.New()
	return mp
}

func (id *MerkleProof) New()  {
	id.RootHash = ""
	id.FieldLeaves = make(map[string]models.LeafFields)
	id.LeafHashes = make([]Hash, 0)
}

// GenerateRootProof (Issuer side)
func (id *MerkleProof) GenerateRootProof(c models.CertificateBase[any]) (Hash, SaltedCertificate, error) {
	saltedCert, err := SaltCertificate(c)
	if err != nil {
		return "", SaltedCertificate{}, err
	}

	var keys []string
	for key := range saltedCert.SaltedFields {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	id.FieldLeaves = saltedCert.SaltedFields
	id.LeafHashes = make([]Hash, 0, len(keys))

	//  Build the ordered LeafHashes list
	for _, key := range keys {
		leaf := id.FieldLeaves[key]
		id.LeafHashes = append(id.LeafHashes, Hash(leaf.Hash))
	}

	id.RootHash = calculateMerkleRoot(id.LeafHashes)

	// The Issuer sends the Root Hash (to Blockchain) and the SaltedCertificate (to Requestor)
	return id.RootHash, saltedCert, nil
}



// VerifyProof checks if a disclosed proof matches the expected root hash.
// This runs on the client/verifier side.
func VerifyProof(p ProofVerification, expectedRoot Hash) bool {
    //Re-calculate the leaf hash for the field we are checking
	var disclosedLeafHash Hash
	switch v:= p.Value.(type) {
	case string:
		// No conversion needed for string values
		disclosedLeafHash = HashData([]byte(v), []byte(p.Salt))
	case int:
		// Convert int to byte slice for hashing
		buf := make([]byte, 4)
		binary.LittleEndian.PutUint32(buf, uint32(v))
	default:
		// Unsupported type, return false or handle as needed
		return false
	}
    
	//check if disclosed hash are present in the hash fields user provided
    found := slices.Contains(p.MerkleProof, disclosedLeafHash)
    
    if !found {
        // If the hash isn't in the list, the user is passing invalid proof 
        return false 
    }

    calculatedRoot := calculateMerkleRoot(p.MerkleProof)

    return calculatedRoot == expectedRoot && calculatedRoot == p.RootHash
}

// calculateMerkleRoot calculates the root hash from an ordered list of leaf hashes.
func calculateMerkleRoot(leaves []Hash) Hash {
	if len(leaves) == 0 {
		return ""
	}
	if len(leaves) == 1 {
		return leaves[0]
	}

	currentLeaves := make([]Hash, len(leaves))
	copy(currentLeaves, leaves)
	if len(currentLeaves)%2 != 0 {
		currentLeaves = append(currentLeaves, currentLeaves[len(currentLeaves)-1])
	}

	var nextLevel []Hash
	for i := 0; i < len(currentLeaves); i += 2 {
		h1 := currentLeaves[i]
		h2 := currentLeaves[i+1]
		// Sort hashes before concatenating to ensure canonical parent hash
		if h1 < h2 {
			nextLevel = append(nextLevel, HashData([]byte(h1), []byte(h2)))
		} else {
			nextLevel = append(nextLevel, HashData([]byte(h2), []byte(h1)))
		}
	}

	return calculateMerkleRoot(nextLevel)
}

