package document

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"

	gotruevault "github.com/FirstVisit/go-truevault"
	"github.com/FirstVisit/go-truevault/client"
	"github.com/google/uuid"
)

type (
	// SearchDocument ...
	SearchDocument struct {
		Document   string    `json:"document,omitempty"`
		DocumentID uuid.UUID `json:"document_id,omitempty"`
		OwnerID    uuid.UUID `json:"owner_id,omitempty"`
	}

	// SearchDocuments ...
	SearchDocuments []SearchDocument

	// SearchDocumentResultInfo ...
	SearchDocumentResultInfo struct {
		PerPage          int `json:"per_page,omitempty"`
		CurrentPage      int `json:"current_page,omitempty"`
		NumPage          int `json:"num_page,omitempty"`
		TotalResultCount int `json:"total_result_count,omitempty"`
	}

	// SearchDocumentResult ...
	SearchDocumentResult struct {
		Info          SearchDocumentResultInfo
		Documents     SearchDocuments
		Result        string
		TransactionID uuid.UUID
	}
)

// DecodeDocument ...
func (r *SearchDocument) DecodeDocument(v interface{}) error {
	decodeString, err := base64.StdEncoding.DecodeString(r.Document)
	if err != nil {
		return err
	}
	return json.NewDecoder(bytes.NewReader(decodeString)).Decode(v)
}

// Document ...
//go:generate mockery --name Document
type Document interface {
	SearchDocument(ctx context.Context, vaultID string, filter gotruevault.SearchOption) (SearchDocumentResult, error)
}

type defaultDocumentService struct {
	*client.Client
}

// New creates a new document service
func New(client client.Client) Document {
	return &defaultDocumentService{
		Client: &client,
	}
}

// SearchDocument https://docs.truevault.com/documentsearch#search-documents
func (r *defaultDocumentService) SearchDocument(ctx context.Context, vaultID string, filter gotruevault.SearchOption) (SearchDocumentResult, error) {
	var result SearchDocumentResult
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(filter); err != nil {
		return SearchDocumentResult{}, err
	}

	path := r.UrlBuilder.SearchDocumentURL(vaultID)

	req, err := r.NewRequest(ctx, http.MethodPost, path, client.ContentTypeApplicationJSON, buf)

	if err != nil {
		return SearchDocumentResult{}, err
	}

	err = r.Do(req, &result)

	return result, err
}