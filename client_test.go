package gotruevault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockURLBuilder struct {
	documentURL string
}

func (r *mockURLBuilder) SearchDocumentURL(string) string {
	return r.documentURL
}

func Test_trueVaultClient_SearchDocument_ReturnsCorrectStatusCode(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))

	client := trueVaultClient{
		httpClient: http.DefaultClient,
		urlBuilder: &mockURLBuilder{documentURL: ts.URL},
	}

	_, err := client.SearchDocument(context.TODO(), "vaultID", SearchFilter{})
	assert.Equal(t, err, ErrUnauthorized)

	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	client = trueVaultClient{
		httpClient: http.DefaultClient,
		urlBuilder: &mockURLBuilder{documentURL: ts.URL},
	}

	_, err = client.SearchDocument(context.TODO(), "vaultID", SearchFilter{})
	assert.Equal(t, err, ErrServerError)
}

func Test_trueVaultClient_SearchDocument_ReturnsSearchResult(t *testing.T) {
	expectedResult := SearchDocumentResult{Result: "Success"}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(expectedResult)
	}))

	client := trueVaultClient{
		httpClient: http.DefaultClient,
		urlBuilder: &mockURLBuilder{documentURL: ts.URL},
	}

	result, err := client.SearchDocument(context.TODO(), "vaultID", SearchFilter{})
	assert.Nil(t, err)
	assert.Equal(t, result, expectedResult)
}
