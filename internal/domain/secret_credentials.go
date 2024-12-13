package domain

import (
	"fmt"
)

// Secret structure for datatype "credentials"
type SecretCredentials struct {
	Address  string `db:"address" json:"address"`   // URL to apply credentials
	Login    string `db:"login" json:"login"`       // login credential
	Password string `db:"password" json:"password"` // password credential
}

// String representation for "credentials" datatype
func (k SecretCredentials) ToString() string {
	return fmt.Sprintf("Address: %s\nLogin: %s\nPassword: %s", k.Address, k.Login, k.Password)
}
