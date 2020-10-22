package document

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"

	gotruevault "github.com/FirstVisit/go-truevault"

	uuid "github.com/satori/go.uuid"
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
		Info          SearchDocumentResultInfo `json:"info,omitempty"`
		Documents     SearchDocuments          `json:"documents,omitempty"`
		Result        string                   `json:"result,omitempty"`
		TransactionID uuid.UUID                `json:"transaction_id,omitempty"`
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

// TrueVaultDocument implements the Document interface
type TrueVaultDocument struct {
	*gotruevault.Client
}

// New creates a new document service
func New(client gotruevault.Client) Document {
	return &TrueVaultDocument{&client}
}

// SearchDocument https://docs.truevault.com/documentsearch#search-documents
func (r *TrueVaultDocument) SearchDocument(ctx context.Context, vaultID string, filter gotruevault.SearchOption) (SearchDocumentResult, error) {
	var result SearchDocumentResult
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(filter); err != nil {
		return SearchDocumentResult{}, err
	}

	path := r.URLBuilder.SearchDocumentURL(vaultID)

	req, err := r.NewRequest(ctx, http.MethodPost, path, gotruevault.ContentTypeApplicationJSON, buf)

	if err != nil {
		return SearchDocumentResult{}, err
	}

	err = r.Do(req, &result)

	return result, err
}
