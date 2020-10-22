package gotruevault

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	//ContentTypeApplicationJSON application/json mime type
	ContentTypeApplicationJSON = "application/json"
)

var (
	// ErrUnauthorized ...
	ErrUnauthorized = errors.New("error: authorization failed")

	// ErrServerError ...
	ErrServerError = errors.New("error: server error")

	// ErrBadRequest ...
	ErrBadRequest = errors.New("error: bad request")
)

//URLBuilder is the interface for building URLs
//go:generate mockery --name URLBuilder
type URLBuilder interface {
	SearchDocumentURL(vaultID string) string
}

//DefaultURLBuilder implements URLBuilder interface
type DefaultURLBuilder struct{}

// SearchDocumentURL ...
func (t *DefaultURLBuilder) SearchDocumentURL(vaultID string) string {
	return fmt.Sprintf("https://api.truevault.com/v1/vaults/%s/search", vaultID)
}

//Client contains the base http requirements to make requests to TrueVault
type Client struct {
	URLBuilder    URLBuilder
	httpClient    *http.Client
	authorization string
}

// New creates a new Client instance
func New(h *http.Client, ub URLBuilder, accessTokenOrKey string) Client {
	return Client{
		httpClient:    h,
		URLBuilder:    ub,
		authorization: buildAuthorizationValue(accessTokenOrKey),
	}
}

// NewDefaultClient creates a Client with the default URLBuilder
func NewDefaultClient(h *http.Client, accessTokenOrKey string) Client {
	return New(h, &DefaultURLBuilder{}, accessTokenOrKey)
}

// WithNewAccessTokenOrKey creates a  new Client instance with new Access Token or API key
func (c *Client) WithNewAccessTokenOrKey(accessTokenOrKey string) Client {
	return New(c.httpClient, c.URLBuilder, accessTokenOrKey)
}

func buildAuthorizationValue(key string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(key+":"))
}

// NewRequest builds an http.Request that contains the Authorization and Content-Type header
func (c *Client) NewRequest(ctx context.Context, method, path, contentType string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", c.authorization)
	req.Header.Set("Content-Type", contentType)

	return req, nil
}

// Do sends an HTTP request and returns an HTTP response, following policy (e.g. redirects, cookies, auth) as configured on the client
func (c *Client) Do(req *http.Request, v interface{}) error {
	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	switch res.StatusCode {
	case http.StatusUnauthorized:
		return ErrUnauthorized
	case http.StatusInternalServerError:
		return ErrServerError
	case http.StatusBadRequest:
		return ErrBadRequest
	}

	return json.NewDecoder(res.Body).Decode(v)
}
