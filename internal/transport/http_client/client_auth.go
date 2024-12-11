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

// logs in on server and retrieves access token
func Login(account *domain.UserAccountData) (*domain.AuthResponse, error) {

	err := app.ValidateAccountData(account)
	if err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(account)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s://%s/login", domain.PROTOCOL_SCHEME, app.Cc.Endpoint)
	resp, err := PostData(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("server response: %s", body)
	}

	ar := new(domain.AuthResponse)
	if err := json.NewDecoder(resp.Body).Decode(ar); err != nil {
		fmt.Println("Error decoding response:", err)
		return nil, err
	}

	// return loginResponse.Token, loginResponse.PublicKey
	return ar, nil
}
