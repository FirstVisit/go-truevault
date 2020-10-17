package client

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

type Client struct {
	httpClient    *http.Client
	UrlBuilder    URLBuilder
	authorization string
}

// NewClient creates a TrueVault client
func NewClient(h *http.Client, ub URLBuilder, accessTokenOrKey string) Client {
	return Client{
		httpClient:    h,
		UrlBuilder:    ub,
		authorization: buildAuthorizationValue(accessTokenOrKey),
	}
}

// NewDefaultClient creates a TrueVault client with the default URLBuilder
func NewDefaultClient(h *http.Client, accessTokenOrKey string) Client {
	return NewClient(h, &defaultURLBuilder{}, accessTokenOrKey)
}

// WithNewAccessTokenOrKey creates a  new Client instance with new Access Token or API key
func (c *Client) WithNewAccessTokenOrKey(accessTokenOrKey string) Client {
	return Client{
		httpClient:    c.httpClient,
		UrlBuilder:    c.UrlBuilder,
		authorization: buildAuthorizationValue(accessTokenOrKey),
	}
}

func buildAuthorizationValue(key string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(key+":"))
}

func (c *Client) NewRequest(ctx context.Context, method, path, contentType string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", c.authorization)
	req.Header.Set("Content-Type", contentType)

	return req, nil
}

func (c *Client) Do(req *http.Request, v interface{}) error {
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
