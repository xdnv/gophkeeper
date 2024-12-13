package domain

import (
	"encoding/json"
	"fmt"
)

// Keeper secret types
const SECRET_CREDENTIALS = "credentials"
const SECRET_CREDITCARD = "creditcard"
const SECRET_TEXT = "text"
const SECRET_BINARY = "binary"

// Secret structure template
type SecretData interface {
	ToString() string
}

// universal human-readble output for Keeper* classes
func KeepReadable(s SecretData) string {
	return s.ToString()
}

// universal secret deserialzer
func KeepDeserialized(class string, data []byte) (*SecretData, error) {
	var object SecretData

	switch class {
	case SECRET_CREDENTIALS:
		object = new(SecretCredentials)
	case SECRET_CREDITCARD:
		object = new(SecretCreditcard)
	case SECRET_TEXT:
		object = new(SecretText)
	case SECRET_BINARY:
		object = new(KeeperBinary)
	default:
		return nil, fmt.Errorf("unknown secret type: %s", class)
	}

	err := json.Unmarshal(data, &object)
	if err != nil {
		return nil, err
	}
	return &object, nil
}
