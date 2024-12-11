// universal Message-Response structures to support multiple transports
package domain

import (
	"bytes"
	"io"
	"net/http"
)

// universal Message data structure for client-server communication
type Message struct {
	Address     string
	ContentType string
	Body        *bytes.Buffer
	Metadata    map[string]string
}

// NewMessage constructor for Message
func NewMessage() *Message {
	return &Message{
		Body:     new(bytes.Buffer),       // init Body as a new bytes.Buffer
		Metadata: make(map[string]string), // init the map
	}
}

// universal Response data structure for client-server communication
type Response struct {
	StatusCode    int
	Status        string
	ContentLength int64
	Body          *bytes.Buffer
	Metadata      map[string][]string
}

// NewResponse constructor for Response
func NewResponse() *Response {
	return &Response{
		Body:     new(bytes.Buffer),         // init Body as a new bytes.Buffer
		Metadata: make(map[string][]string), // init the map
	}
}

// converts the http.Response object to Response
func NewResponseFromHTTP(r *http.Response) (*Response, error) {
	res := &Response{
		StatusCode:    r.StatusCode,
		Status:        r.Status,
		ContentLength: r.ContentLength,
		Body:          new(bytes.Buffer),         // init Body as a new bytes.Buffer
		Metadata:      make(map[string][]string), // init the map
	}

	//read out body
	if _, err := io.Copy(res.Body, r.Body); err != nil {
		return nil, err
	}
	defer r.Body.Close()

	//read out Header copying all the slices
	for key, values := range r.Header {
		res.Metadata[key] = append([]string(nil), values...)
	}

	return res, nil
}
