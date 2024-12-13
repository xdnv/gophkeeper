// provides message posting functions
package http_client

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"internal/app"
	"internal/domain"
)

// Sends command to a server
func ExecServerCommand(ctx context.Context, body *bytes.Buffer) (*domain.Response, error) {
	m := domain.NewMessage()
	m.ContentType = "application/json"
	m.Body = body

	m.Address = fmt.Sprintf("%s://%s/command", domain.PROTOCOL_SCHEME, app.Cc.Endpoint)

	// set metadata for extended posting
	m.Metadata["Content-Type"] = m.ContentType
	m.Metadata["Accept-Encoding"] = "gzip"

	// // set real client IP
	// if agentIP != "" {
	// 	m.Metadata["X-Real-IP"] = agentIP
	// }

	// //optionally encrypt message
	// _, err := encryptMessage(m)
	// if err != nil {
	// 	return nil, err
	// }

	//compress message
	if app.Cc.UseCompression {
		err := compressMessage(m)
		if err != nil {
			return nil, err
		}
	}

	return PostMessageExtended(ctx, m)
}

// func encryptMessage(m *domain.Message) (bool, error) {
// 	if !cryptor.CanEncrypt() {
// 		return false, nil
// 	}

// 	msg, err := cryptor.Encrypt(m.Body.Bytes())
// 	if err != nil {
// 		return false, err
// 	}

// 	m.Body.Reset()             // Clear the buffer
// 	_, err = m.Body.Write(msg) // Write encrypted data back to buffer
// 	if err != nil {
// 		return false, err
// 	}

// 	m.Metadata["X-Encrypted"] = "true"

// 	return true, nil
// }

func compressMessage(m *domain.Message) error {
	var buf bytes.Buffer

	g := gzip.NewWriter(&buf)
	if _, err := g.Write(m.Body.Bytes()); err != nil {
		return err
	}
	if err := g.Close(); err != nil {
		return err
	}

	m.Body.Reset()                      // Clear the buffer
	_, err := m.Body.Write(buf.Bytes()) // Write compressed data back to buffer
	if err != nil {
		return err
	}

	m.Metadata["Content-Encoding"] = "gzip"

	return nil
}
