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
	"net/url"
	"strings"
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

// Error API Response. Contains the error message as well as the type of error
type Error struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

//URLBuilder is the interface for building URLs
//go:generate mockery --name URLBuilder
type URLBuilder interface {
	SearchDocumentURL(vaultID string) string
	GetUserURL(userId []string) string
	CreateUserURL() string
	ListUserURL(queryParams url.Values) string
	UpdateUserURL(userId string) string
	UpdateUserPasswordURL(userId string) string
	DeleteUserURL(userId string) string
	CreateAccessTokenURL(userId string) string
	CreateApiKeyURL(userId string) string
}

//DefaultURLBuilder implements URLBuilder interface
type DefaultURLBuilder struct{}

// SearchDocumentURL ...
func (t *DefaultURLBuilder) SearchDocumentURL(vaultID string) string {
	return fmt.Sprintf("https://api.truevault.com/v1/vaults/%s/search", vaultID)
}

// GetUserURL returns the TrueVault `Get User` route for the specified user id(s)
func (t *DefaultURLBuilder) GetUserURL(userId []string) string {
	return fmt.Sprintf("https://api.truevault.com/v2/users/"+strings.Join(userId, ","))
}

// CreateUserURL returns the TrueVault `Create User` route
func (t *DefaultURLBuilder) CreateUserURL() string {
	return "https://api.truevault.com/v1/users"
}

// UpdateUserURL returns the TrueVault `Update User` route
func (t *DefaultURLBuilder) UpdateUserURL(userId string) string {
	return "https://api.truevault.com/v1/users/" + userId
}

// UpdateUserPasswordURL returns the TrueVault `Update User Password` route
func (t *DefaultURLBuilder) UpdateUserPasswordURL(userId string) string {
	return "https://api.truevault.com/v1/users/" + userId
}

// DeleteUserURL returns the TrueVault `Delete User` route
func (t *DefaultURLBuilder) DeleteUserURL(userId string) string {
	return "https://api.truevault.com/v1/users/" + userId
}

// CreateAccessTokenURL returns the TrueVault `Create Access Token` route
func (t *DefaultURLBuilder) CreateAccessTokenURL(userId string) string {
	return "https://api.truevault.com/v1/users/" + userId
}

// CreateApiKeyURL returns the TrueVault `Create API Key` route
func (t *DefaultURLBuilder) CreateApiKeyURL(userId string) string {
	return "https://api.truevault.com/v1/users/" + userId + "/api_key"
}

// ListUserURL returns the TrueVault `List User` route
func (t *DefaultURLBuilder) ListUserURL(queryParams url.Values) string {
	params := "?"
	if queryParams != nil {
		params += queryParams.Encode()
	}
	return fmt.Sprintf("https://api.truevault.com/v2/users/?%s", params)
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
		authorization: "Basic " + base64.StdEncoding.EncodeToString([]byte(accessTokenOrKey+":")),
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
