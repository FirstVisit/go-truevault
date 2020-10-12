package bingboop

import (
	"bytes"
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
	httpClient *http.Client
	urlBuilder URLBuilder
	apiKey     string
}

// NewClient ...
func NewClient(h *http.Client, a string) Client {
	return &trueVaultClient{
		httpClient: h,
		urlBuilder: &DefaultURLBuilder{},
		apiKey:     base64.StdEncoding.EncodeToString([]byte(a + ":")),
	}
}

func (c *trueVaultClient) newRequest(ctx context.Context, method, path, contentType string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Basic "+c.apiKey)
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

func (c *trueVaultClient) SearchDocument(ctx context.Context, vaultID string, filter SearchFilter) (SearchDocumentResult, error) {
	var result SearchDocumentResult
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(filter)

	if err != nil {
		return SearchDocumentResult{}, err
	}

	path := c.urlBuilder.SearchDocumentURL(vaultID)

	req, err := c.newRequest(ctx, http.MethodPost, path, contentTypeApplicationJSON, buf)

	if err != nil {
		return SearchDocumentResult{}, err
	}

	err = c.do(req, &result)

	return result, err
}
