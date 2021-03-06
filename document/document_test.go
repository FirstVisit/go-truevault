package document

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	gotruevault "github.com/FirstVisit/go-truevault"
	_clientMock "github.com/FirstVisit/go-truevault/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSearchDocument_DecodeDocument(t *testing.T) {
	type document struct {
		X int
		Y string
	}

	tests := []struct {
		name     string
		doc      string
		expected document
	}{
		{name: "int fields", doc: "eyJYIjoxMjM0fQ==", expected: document{X: 1234}},
		{name: "string fields", doc: "eyJZIjoiVEVTVElORyJ9", expected: document{Y: "TESTING"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result document
			r := &SearchDocument{
				Document: tt.doc,
			}
			assert.Nil(t, r.DecodeDocument(&result))
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_trueVaultClient_SearchDocument_ReturnsCorrectStatusCode(t *testing.T) {
	// Unauthorized
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))

	urlBuilder := new(_clientMock.URLBuilder)
	urlBuilder.On("SearchDocumentURL", mock.Anything).Once().Return(ts.URL)
	service := New(gotruevault.New(http.DefaultClient, urlBuilder, ""))
	_, err := service.SearchDocument(context.TODO(), "vaultID", gotruevault.SearchOption{})
	assert.Equal(t, err, gotruevault.ErrUnauthorized)

	// Internal Server Error
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	urlBuilder = new(_clientMock.URLBuilder)
	urlBuilder.On("SearchDocumentURL", mock.Anything).Once().Return(ts.URL)
	service = New(gotruevault.New(http.DefaultClient, urlBuilder, ""))
	_, err = service.SearchDocument(context.TODO(), "vaultID", gotruevault.SearchOption{})
	assert.Equal(t, err, gotruevault.ErrServerError)
}

func Test_trueVaultClient_SearchDocument_ReturnsSearchResult(t *testing.T) {
	expectedResult := SearchDocumentResult{Result: "Success"}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(expectedResult)
	}))

	urlBuilder := new(_clientMock.URLBuilder)
	urlBuilder.On("SearchDocumentURL", mock.Anything).Once().Return(ts.URL)
	service := New(gotruevault.New(http.DefaultClient, urlBuilder, ""))
	result, err := service.SearchDocument(context.TODO(), "testing", gotruevault.SearchOption{})
	assert.Nil(t, err)
	assert.Equal(t, result, expectedResult)
}
