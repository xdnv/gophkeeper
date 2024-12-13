// provides user security functions
package http_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"internal/app"
	"internal/domain"
	"io"
	"net/http"
)

// registers new user if its login is available
func Register(account *domain.UserAccountData) error {

	err := app.ValidateAccountData(account)
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(account)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s://%s/register", domain.PROTOCOL_SCHEME, app.Cc.Endpoint)
	resp, err := PostData(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server response: %s", body)
	}

	// var loginResponse LoginResponse
	// if err := json.NewDecoder(resp.Body).Decode(&loginResponse); err != nil {
	// 	fmt.Println("Error decoding response:", err)
	// 	return "", ""
	// }

	// return loginResponse.Token, loginResponse.PublicKey
	return nil
}
