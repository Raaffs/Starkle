package zkp

import (
	"encoding/binary"
	"fmt"
	"slices"
	"sort"

	"github.com/Suy56/ProofChain/internal/models"
)

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


//Merkle Proof is generated on issuer's side and root of merkle tree is anchored on chain
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

	for _, key := range keys {
		leaf := id.FieldLeaves[key]
		id.LeafHashes = append(id.LeafHashes, Hash(leaf.Hash))
	}

	id.RootHash = calculateMerkleRoot(id.LeafHashes)

	//salted certificate is sent to requestor 
	return id.RootHash, saltedCert, nil
}

func GenerateMerklePath(leaves []string, targetHash string) (string, []string) {
    // DO NOT sort(current) here. 
    // The leaves must remain in the same order as they were 
    // when calculateMerkleRoot was called.
    current := make([]string, len(leaves))
    copy(current, leaves)

    var path []string
    
    for len(current) > 1 {
        // Handle odd number of nodes
        if len(current)%2 != 0 {
            current = append(current, current[len(current)-1])
        }

        var nextLevel []string
        for i := 0; i < len(current); i += 2 {
            h1, h2 := current[i], current[i+1]
            
            // Deterministic sorting of PAIRS is fine and matches calculateMerkleRoot
            first, second := h1, h2
            if h1 > h2 {
                first, second = h2, h1
            }

            parent := string(HashData([]byte(first), []byte(second)))
            nextLevel = append(nextLevel, parent)

            // Correct target tracking
            if h1 == targetHash {
                path = append(path, h2)
                targetHash = parent
            } else if h2 == targetHash {
                path = append(path, h1)
                targetHash = parent
            }
        }
        current = nextLevel
    }

    return current[0], path
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
	// We cast to uint32 and use Little Endian to maintain parity with the Rust implementation.
	// The range-prover in Rust serializes integers as 4-byte Little Endian values 
	// before hashing, so the byte representation must be identical here.
		buf := make([]byte, 4)
		binary.LittleEndian.PutUint32(buf, uint32(v))
	default:
		// Unsupported type, return false or handle as needed
		return false
	}
    
	//check if disclosed hash are present in the hash fields user provided
    found := slices.Contains(p.MerkleProof, disclosedLeafHash)
    
    if !found {
        return false 
    }

    calculatedRoot := calculateMerkleRoot(p.MerkleProof)

    return calculatedRoot == expectedRoot
}

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

func calculateLeafHash(value any, salt string) (Hash,error) {
    var leafHash Hash
    switch v := value.(type) {
    case string:
        leafHash = HashData([]byte(v), []byte(salt))
    case int:
        buf := make([]byte, 4)
        binary.LittleEndian.PutUint32(buf, uint32(v))
        leafHash = HashData(buf, []byte(salt))
    case uint32:
        buf := make([]byte, 4)
        binary.LittleEndian.PutUint32(buf, v)
        leafHash = HashData(buf, []byte(salt))
    case nil:
        leafHash = HashData([]byte(""), []byte(salt))
    default:
        // Default to string conversion if type is unknown
        return "",fmt.Errorf("invalid data type used for hashing leaf")
    }
    return leafHash,nil
}