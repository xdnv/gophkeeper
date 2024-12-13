package domain

import (
	"fmt"
)

// Secret structure for datatype "creditcard"
type SecretCreditcard struct {
	CardNumber     string `db:"card_number" json:"card_number"`         // Card number in 4x4 space-delimited blocks of digits
	ExpirationDate string `db:"expiration_date" json:"expiration_date"` // Expiration date in "MM/YY" format
	SecurityCode   string `db:"security_code" json:"security_code"`     // Security (CVC/CVV/CVC2/etc.) code, 3 digits
}

// String representation for "creditcard" datatype
func (k SecretCreditcard) ToString() string {
	return fmt.Sprintf("Card Number: %s\nExpiration Date: %s\nCVV: %s", k.CardNumber, k.ExpirationDate, k.SecurityCode)
}
