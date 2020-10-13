package go_truevault

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"time"

	uuid "github.com/satori/go.uuid"
)

type (
	// SearchValue ...
	SearchValue interface {
		searchValue()
	}

	// SearchValues ...
	SearchValues []SearchValue

	// String implements SearchValue that holds a string value
	String struct {
		Value string
	}

	// Float64 implements SearchValue that holds a float64 value
	Float64 struct {
		Value float64
	}

	// Int implements SearchValue that holds an int value
	Int struct {
		Value int
	}

	// Time implements SearchValue that holds an time.Time value
	Time struct {
		Value time.Time
	}

	// RangeValue implements SearchValue that holds range values
	RangeValue struct {
		Gt  float64 `json:"gt,omitempty"`
		Gte float64 `json:"gte,omitempty"`
		Lt  float64 `json:"lt,omitempty"`
		Lte float64 `json:"lte,omitempty"`
	}
)

func (i String) searchValue()     {}
func (i Float64) searchValue()    {}
func (i Int) searchValue()        {}
func (i Time) searchValue()       {}
func (i RangeValue) searchValue() {}

// MarshalJSON ...
func (i String) MarshalJSON() (data []byte, err error) {
	return json.Marshal(i.Value)
}

// MarshalJSON ...
func (i Float64) MarshalJSON() (data []byte, err error) {
	return json.Marshal(i.Value)
}

// MarshalJSON ...
func (i Int) MarshalJSON() (data []byte, err error) {
	return json.Marshal(i.Value)
}

// MarshalJSON ...
func (i Time) MarshalJSON() (data []byte, err error) {
	return json.Marshal(i.Value)
}

type (
	// SearchType ...
	SearchType interface {
		searchType()
	}

	// SearchTypes ...
	SearchTypes []SearchType

	// Eq ...
	Eq struct {
		Value         SearchValue
		CaseSensitive bool
	}

	// In ...
	In struct {
		Value         []SearchValue
		CaseSensitive bool
	}

	// Not ...
	Not struct {
		Value         SearchValue
		CaseSensitive bool
	}

	// NotIn ...
	NotIn struct {
		Value         []SearchValue
		CaseSensitive bool
	}

	// Wildcard ...
	Wildcard struct {
		Value         SearchValue
		CaseSensitive bool
	}

	// Range ...
	Range struct {
		Value RangeValue
	}
)

func (i In) searchType()       {}
func (i Eq) searchType()       {}
func (i Not) searchType()      {}
func (i NotIn) searchType()    {}
func (i Wildcard) searchType() {}
func (i Range) searchType()    {}

// MarshalJSON ...
func (i In) MarshalJSON() (data []byte, err error) {
	return json.Marshal(struct {
		Type          string        `json:"type,omitempty"`
		Value         []SearchValue `json:"value,omitempty"`
		CaseSensitive bool          `json:"case_sensitive,omitempty"`
	}{
		Type:          "in",
		Value:         i.Value,
		CaseSensitive: i.CaseSensitive,
	})
}

// MarshalJSON ...
func (i Eq) MarshalJSON() (data []byte, err error) {
	return json.Marshal(struct {
		Type          string      `json:"type,omitempty"`
		Value         SearchValue `json:"value,omitempty"`
		CaseSensitive bool        `json:"case_sensitive,omitempty"`
	}{
		Type:          "eq",
		Value:         i.Value,
		CaseSensitive: i.CaseSensitive,
	})
}

// MarshalJSON ...
func (i Not) MarshalJSON() (data []byte, err error) {
	return json.Marshal(struct {
		Type          string      `json:"type,omitempty"`
		Value         SearchValue `json:"value,omitempty"`
		CaseSensitive bool        `json:"case_sensitive,omitempty"`
	}{
		Type:          "not",
		Value:         i.Value,
		CaseSensitive: i.CaseSensitive,
	})
}

// MarshalJSON ...
func (i NotIn) MarshalJSON() (data []byte, err error) {
	return json.Marshal(struct {
		Type          string        `json:"type,omitempty"`
		Value         []SearchValue `json:"value,omitempty"`
		CaseSensitive bool          `json:"case_sensitive,omitempty"`
	}{
		Type:          "not_in",
		Value:         i.Value,
		CaseSensitive: i.CaseSensitive,
	})
}

// MarshalJSON ...
func (i Wildcard) MarshalJSON() (data []byte, err error) {
	return json.Marshal(struct {
		Type          string      `json:"type,omitempty"`
		Value         SearchValue `json:"value,omitempty"`
		CaseSensitive bool        `json:"case_sensitive,omitempty"`
	}{
		Type:          "wildcard",
		Value:         i.Value,
		CaseSensitive: i.CaseSensitive,
	})
}

// MarshalJSON ...
func (i Range) MarshalJSON() (data []byte, err error) {
	return json.Marshal(struct {
		Type  string     `json:"type,omitempty"`
		Value RangeValue `json:"value,omitempty"`
	}{
		Type:  "range",
		Value: i.Value,
	})
}

type (
	// FilterType ...
	FilterType string

	// SortOrder ...
	SortOrder string

	// SearchDocument ...
	SearchDocument struct {
		Document   string    `json:"document,omitempty"`
		DocumentID uuid.UUID `json:"document_id,omitempty"`
		OwnerID    uuid.UUID `json:"owner_id,omitempty"`
	}

	// SearchDocuments ...
	SearchDocuments []SearchDocument

	// SearchFilter ...
	SearchFilter struct {
		Filter     map[string]SearchType  `json:"filter,omitempty"`
		FilterType FilterType             `json:"filter_type,omitempty"`
		Page       int                    `json:"page,omitempty"`
		PerPage    int                    `json:"per_page,omitempty"`
		Sort       []map[string]SortOrder `json:"sort,omitempty"`
		SchemaID   uuid.UUID              `json:"schema_id,omitempty"`
	}

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

func (r *SearchDocument) DecodeDocument(v interface{}) error {
	decodeString, err := base64.StdEncoding.DecodeString(r.Document)
	if err != nil {
		return err
	}
	return json.NewDecoder(bytes.NewReader(decodeString)).Decode(v)
}

const (
	And FilterType = "and"
	Or  FilterType = "or"

	AscSort  SortOrder = "asc"
	DescSort SortOrder = "desc"
)
