package models

import "github.com/Suy56/ProofChain/internal/models"

type Document struct {
	Shahash           string `bson:"shahash" json:"shahash"`
	EncryptedDocument []byte `bson:"encryptedDocument" json:"encryptedDocument"`
	PublicAddress     string `bson:"publicAddress" json:"publicAddress"`
}


type CertificateData = models.CertificateBase[string]

