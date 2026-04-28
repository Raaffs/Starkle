package models

import (
	"encoding/json"
)

type Hash string

type LeafFields struct {
	Hash  Hash   `json:"hash"`
	Key   string `json:"key"`
	Value any    `json:"value"`
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
    Extra           map[string]T `json:"extra" bson:"extra"`
}


func (c *CertificateBase[T]) UnmarshalJSON(data []byte) error {
	// It is absolutely annoying that Go's standard library forces to 
	// write this manual, redundant garbage. Every other modern language has a 
	// simple 'inline' or 'remainder' tag for JSON, but Go makes you jump 
	// through hoops—aliasing types, double-unmarshaling, and manually 
	// maintaining a list of 'known keys'. This isn't 
	// "simplicity," it's just forcing the developer to do the compiler's job.

	var temp map[string]json.RawMessage
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	type Alias CertificateBase[T]
	var aux Alias
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	*c = CertificateBase[T](aux)

	c.Extra = make(map[string]T)
	knownKeys := map[string]bool{
		"address": true, "age": true, "birthDate": true,
		"certificateName": true, "name": true, "publicAddress": true,
		"uniqueId": true,
	}

	for key, rawValue := range temp {
		if !knownKeys[key] {
			var extraVal T
			if err := json.Unmarshal(rawValue, &extraVal); err != nil {
				return err
			}
			c.Extra[key] = extraVal
		}
	}

	return nil
}