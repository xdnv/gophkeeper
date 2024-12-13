package domain

import (
	"fmt"
)

// Secret structure for datatype "text"
type SecretText struct {
	Text string `db:"text" json:"text"` // Secret text
}

// String representation for "text" datatype
func (k SecretText) ToString() string {
	return fmt.Sprintf("Text:\n%s", k.Text)
}
