package zkp

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/Suy56/ProofChain/internal/models"
	"github.com/Suy56/ProofChain/internal/utils"
)

type Hash string


// SaltCertificate iterates over the models.CertificateData, salts each field, and returns the map of salted leaves.
func SaltCertificate(c models.CertificateBase[any]) (SaltedCertificate, error) {
	fieldMap := make(map[string]any)
	var keys []string

	for k, v := range utils.Walk(c) {
		fieldMap[k] = v
		keys = append(keys, k)
	}

	// Sort keys alphabetically to ensure deterministic processing order
	sort.Strings(keys)

	saltedFields := make(map[string]models.LeafFields)

	for _, key := range keys {
		value := fieldMap[key]
		var leafHash Hash
		salt, err := utils.GenerateSalt()
		if err != nil {
			return SaltedCertificate{}, err
		}

		switch v := value.(type) {
		case string:
			leafHash = HashData([]byte(v), []byte(salt))
		case int:
	        buf := make([]byte, 4)
        	binary.LittleEndian.PutUint32(buf, uint32(v))
			leafHash = HashData(buf, []byte(salt))
		default:
			return SaltedCertificate{}, fmt.Errorf("unsupported field type for key %s", key)
		}

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

func HashData(data ...[]byte) Hash {
	h := sha256.New()
	for _, d := range data {
		h.Write(d)
	}
	return Hash(hex.EncodeToString(h.Sum(nil)))
}
