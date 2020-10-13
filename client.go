package go_truevault

import (
	"encoding/base64"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// Error API Response. Contains the error message as well as the type of error
type Error struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

// Client makes all the API calls to TrueVault
type Client struct {
	apiKey  string
	accessToken string
	baseURL string
	http    *http.Client
}

// NewClient creates a TrueVault client
func NewClient(apiKey string, accessToken string) Client {
	return Client{
		apiKey: apiKey,
		accessToken: accessToken,
		baseURL: "https://api.truevault.com/",
		http:    &http.Client{},
	}
}

func (c Client) buildRequest(httpMethod string, slug string, data url.Values) (*http.Request, error) {
	var body *strings.Reader
	if data != nil {
		body = strings.NewReader(data.Encode())
	}

	req, err := http.NewRequest(httpMethod, c.baseURL + slug, body)
	if err != nil {
		return nil, err
	}

	key := c.apiKey
	if c.accessToken != "" {
		key = c.accessToken
	}

	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(key + ":")))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

func (c Client) do(req *http.Request) ([]byte, error) {
	resp, err := c.http.Do(req)
	if err != nil {
		return []byte{}, errors.New("failed to call true vault")
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Fatalf("truevault#post failed to close respones\n%+v", err)
		}
	}()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return []byte{}, err
	}

	return body, nil
}

func (c Client) post(slug string, data url.Values) ([]byte, error) {
	req, err := c.buildRequest(http.MethodPost, slug, data)
	if err != nil {
		return []byte{}, err
	}

	return c.do(req)
}

func (c *Client) delete(slug string, data url.Values) ([]byte, error) {
	req, err := c.buildRequest(http.MethodDelete, slug, data)
	if err != nil {
		return []byte{}, err
	}

	return c.do(req)
}

func (c *Client) get(slug string, queryParams url.Values) ([]byte, error) {
	params := "?"
	if queryParams != nil {
		params += queryParams.Encode()
	}

	req, err := http.NewRequest(http.MethodGet, c.baseURL + slug + params, nil)
	if err != nil {
		return []byte{}, err
	}

	return c.do(req)
}

