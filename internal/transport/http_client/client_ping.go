// provides ping functions
package http_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"internal/app"
	"internal/domain"
	"net/http"
)

// Pings server
func Ping() (*domain.Response, error) {

	const message = "ping"
	jsonData, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	m := domain.NewMessage()
	m.Address = fmt.Sprintf("%s://%s/ping", domain.PROTOCOL_SCHEME, app.Cc.Endpoint)
	m.ContentType = "application/json"
	m.Body = bytes.NewBuffer(jsonData)

	// set metadata for extended posting
	m.Metadata["Content-Type"] = m.ContentType
	m.Metadata["Authorization"] = "Bearer " + app.Cc.AuthToken

	r, err := PostMessageExtended(context.Background(), m)
	if err != nil {
		return nil, err
	}

	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server response: %s", r.Status)
	}

	return r, nil
}
