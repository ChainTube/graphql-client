package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
)

// PostWithSchema is similar to Post, except it skips decoding the raw json response
// unpacked onto Response. This is used to test extension keys which are not
// available when using Post.
func (p *Client) PostWithSchema(query string, response interface{}, options ...Option) error {
	r, err := p.newRequest(query, options...)
	if err != nil {
		return fmt.Errorf("build: %w", err)
	}

	w := httptest.NewRecorder()
	p.h.ServeHTTP(w, r)

	if w.Code >= http.StatusBadRequest {
		return fmt.Errorf("http %d: %s", w.Code, w.Body.String())
	}

	// use the response object from params because it contains type information
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		return fmt.Errorf("decode: %w", err)
	}

	return nil
}

// MustPostWithSchema is a convenience wrapper around PostWithSchema that automatically panics on error
func (p *Client) MustPostWithSchema(query string, response interface{}, options ...Option) {
	if err := p.PostWithSchema(query, response, options...); err != nil {
		panic(err)
	}
}
