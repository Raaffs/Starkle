package models

type Hash string


type LeafFields struct {
	Hash  Hash   `json:"hash"` // The salted hash of the value
	Key   string `json:"key"`
	Value string `json:"value"`
	Salt  string `json:"salt"`
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

