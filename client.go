package gotruevault

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const (
	contentTypeApplicationJSON = "application/json"
)

var (
	// ErrUnauthorized ...
	ErrUnauthorized = errors.New("error: authorization failed")

	// ErrServerError ...
	ErrServerError = errors.New("error: server error")

	// ErrBadRequest ...
	ErrBadRequest = errors.New("error: bad request")
)

// Client is the interface that wraps http calls to TrueVault
//go:generate mockery --name Client
type Client interface {
	SearchDocument(ctx context.Context, vaultID string, filter SearchFilter) (SearchDocumentResult, error)
}

type trueVaultClient struct {
	httpClient    *http.Client
	urlBuilder    URLBuilder
	authorization string
}

// NewClient creates a TrueVault client
func NewClient(h *http.Client, accessTokenOrKey string) Client {
	return &trueVaultClient{
		httpClient:    h,
		urlBuilder:    &DefaultURLBuilder{},
		authorization: buildAuthorizationValue(accessTokenOrKey),
	}
}

// WithNewAccessTokenOrKey creates a  new Cient instance with new Access Token or API key
func (c *trueVaultClient) WithNewAccessTokenOrKey(accessTokenOrKey string) Client {
	return &trueVaultClient{
		httpClient:    c.httpClient,
		urlBuilder:    c.urlBuilder,
		authorization: buildAuthorizationValue(accessTokenOrKey),
	}
}

func buildAuthorizationValue(key string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(key+":"))
}

func (c *trueVaultClient) newRequest(ctx context.Context, method, path, contentType string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", c.authorization)
	req.Header.Set("Content-Type", contentType)

	return req, nil
}

func (c *trueVaultClient) do(req *http.Request, v interface{}) error {
	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		_ = res.Body.Close()
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
