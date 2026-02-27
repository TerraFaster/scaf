package templates

import (
	"net/http"
	"time"
)

func newHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
		// TLS verification is enforced by default in Go's http.Client
	}
}
