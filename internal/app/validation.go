// common APP functions
package app

import (
	"fmt"
	"internal/domain"
	"strings"
)

// Account data quality gate
func ValidateAccountData(account *domain.UserAccountData) error {

	account.Login = strings.TrimSpace(account.Login)
	account.Email = strings.TrimSpace(account.Email)

	if account.Login == "" {
		return fmt.Errorf("login cannot be empty")
	}
	if account.Password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	if account.Password != strings.TrimSpace(account.Password) {
		return fmt.Errorf("please remove leading/trailing spaces from password")
	}
	if account.Email != "" && !domain.IsValidEmail(account.Email) {
		return fmt.Errorf("please use valid email address")
	}

	return nil
}
