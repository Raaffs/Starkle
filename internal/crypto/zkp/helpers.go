package zkp

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/Suy56/ProofChain/internal/models"
	"github.com/Suy56/ProofChain/internal/utils"
)

type Hash string

// SaltCertificate iterates over the models.CertificateData, salts each field, and returns the map of salted leaves.
func SaltCertificate(c models.CertificateBase[string]) (SaltedCertificate, error) {
	fieldMap := make(map[string]string)
	var keys []string

	for k, v := range utils.Walk(c) {
		strVal := fmt.Sprint(v)
		fieldMap[k] = strVal
		keys = append(keys, k)
	}

	// Sort keys alphabetically to ensure deterministic processing order
	sort.Strings(keys)

	saltedFields := make(map[string]models.LeafFields)

	for _, key := range keys {
		value := fieldMap[key]

		salt, err := utils.GenerateSalt()
		if err != nil {
			return SaltedCertificate{}, err
		}

		// Leaf Hash Logic: Hash(Value + Salt)
		leafHash := HashData([]byte(value), []byte(salt))

		leaf := models.LeafFields{
			Key:   key,
			Value: value,
			Salt:  salt,
			Hash:  models.Hash(leafHash),
		}

		saltedFields[key] = leaf
	}

	return SaltedCertificate{SaltedFields: saltedFields}, nil
}
// HashData performs a SHA256 hash on the provided inputs.
func HashData(data ...[]byte) Hash {
	h := sha256.New()
	for _, d := range data {
		h.Write(d)
	}
	return Hash(hex.EncodeToString(h.Sum(nil)))
}
