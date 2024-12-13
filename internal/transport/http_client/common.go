// provides user security functions
package http_client

import (
	"bytes"
	"context"
	"internal/domain"
	"net/http"
)

// simple HTTP post function
func PostData(address string, contentType string, body *bytes.Buffer) (*http.Response, error) {
	resp, err := http.Post(address, contentType, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// simple HTTP post function based on Message input format
func PostMessage(m *domain.Message) (*domain.Response, error) {
	resp, err := http.Post(m.Address, m.ContentType, m.Body)
	if err != nil {
		return nil, err
	}

	res, err := domain.NewResponseFromHTTP(resp)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// extended HTTP post function based on Message input format, supports headers & more
func PostMessageExtended(ctx context.Context, m *domain.Message) (*domain.Response, error) {

	r, err := http.NewRequestWithContext(ctx, "POST", m.Address, m.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	r.Close = true //whether to close the connection after replying to this request (for servers) or after sending the request (for clients).

	// set HTTP headers
	for k, v := range m.Metadata {
		r.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}

	res, err := domain.NewResponseFromHTTP(resp)
	if err != nil {
		return nil, err
	}

	return res, nil
}
