package models

type Document struct {
	Shahash           string `bson:"shahash" json:"shahash"`
	EncryptedDocument []byte `bson:"encryptedDocument" json:"encryptedDocument"`
	PublicAddress     string `bson:"publicAddress" json:"publicAddress"`
}

type CertificateBase[T any] struct {
    Address         T            `json:"address" bson:"address"`
    Age             T            `json:"age" bson:"age"`
    BirthDate       T            `json:"birthDate" bson:"birthDate"`
    CertificateName T            `json:"certificateName" bson:"certificateName"`
    Name            T            `json:"name" bson:"name"`
    PublicAddress   T            `json:"publicAddress" bson:"publicAddress"`
    UniqueID        T            `json:"uniqueId" bson:"uniqueId"`
    Extra           map[string]T `json:"extra" bson:",inline"`
}

type CertificateData = CertificateBase[string]

