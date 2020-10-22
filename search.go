package gotruevault

import (
	"encoding/json"
	"time"

	uuid "github.com/satori/go.uuid"
)

type (
	// SearchValue an interface that represents an search value
	SearchValue interface {
		searchValue()
	}

	// SearchValues alias for a list of SearchValue
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

// MarshalJSON returns JSON encoding of String
func (i String) MarshalJSON() (data []byte, err error) {
	return json.Marshal(i.Value)
}

// MarshalJSON returns JSON encoding of Float64
func (i Float64) MarshalJSON() (data []byte, err error) {
	return json.Marshal(i.Value)
}

// MarshalJSON returns JSON encoding of Int
func (i Int) MarshalJSON() (data []byte, err error) {
	return json.Marshal(i.Value)
}

// MarshalJSON returns JSON encoding of Time
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

	singleSearchValueJSON struct {
		Type          string      `json:"type,omitempty"`
		Value         SearchValue `json:"value,omitempty"`
		CaseSensitive bool        `json:"case_sensitive,omitempty"`
	}

	multiSearchValueJSON struct {
		Type          string        `json:"type,omitempty"`
		Value         []SearchValue `json:"value,omitempty"`
		CaseSensitive bool          `json:"case_sensitive,omitempty"`
	}

	rangeSearchValueJSON struct {
		Type  string     `json:"type,omitempty"`
		Value RangeValue `json:"value,omitempty"`
	}
)

func (i In) searchType()       {}
func (i Eq) searchType()       {}
func (i Not) searchType()      {}
func (i NotIn) searchType()    {}
func (i Wildcard) searchType() {}
func (i Range) searchType()    {}

// MarshalJSON returns JSON encoding of In
func (i In) MarshalJSON() (data []byte, err error) {
	return json.Marshal(multiSearchValueJSON{
		Type:          "in",
		Value:         i.Value,
		CaseSensitive: i.CaseSensitive,
	})
}

// MarshalJSON returns JSON encoding of Eq
func (i Eq) MarshalJSON() (data []byte, err error) {
	return json.Marshal(singleSearchValueJSON{
		Type:          "eq",
		Value:         i.Value,
		CaseSensitive: i.CaseSensitive,
	})
}

// MarshalJSON returns JSON encoding of Not
func (i Not) MarshalJSON() (data []byte, err error) {
	return json.Marshal(singleSearchValueJSON{
		Type:          "not",
		Value:         i.Value,
		CaseSensitive: i.CaseSensitive,
	})
}

// MarshalJSON returns JSON encoding of NotIn
func (i NotIn) MarshalJSON() (data []byte, err error) {
	return json.Marshal(multiSearchValueJSON{
		Type:          "not_in",
		Value:         i.Value,
		CaseSensitive: i.CaseSensitive,
	})
}

// MarshalJSON returns JSON encoding of Wildcard
func (i Wildcard) MarshalJSON() (data []byte, err error) {
	return json.Marshal(singleSearchValueJSON{
		Type:          "wildcard",
		Value:         i.Value,
		CaseSensitive: i.CaseSensitive,
	})
}

// MarshalJSON returns JSON encoding of Range
func (i Range) MarshalJSON() (data []byte, err error) {
	return json.Marshal(rangeSearchValueJSON{
		Type:  "range",
		Value: i.Value,
	})
}

type (
	// FilterType ...
	FilterType string

	// SortOrder ...
	SortOrder string

	// SearchOption ...
	SearchOption struct {
		Filter     map[string]SearchType  `json:"filter,omitempty"`
		FilterType FilterType             `json:"filter_type,omitempty"`
		Page       int                    `json:"page,omitempty"`
		PerPage    int                    `json:"per_page,omitempty"`
		Sort       []map[string]SortOrder `json:"sort,omitempty"`
		SchemaID   uuid.UUID              `json:"schema_id,omitempty"`
	}
)

const (
	// And ...
	And FilterType = "and"
	// Or ...
	Or FilterType = "or"

	// Asc ...
	Asc SortOrder = "asc"
	// Desc ...
	Desc SortOrder = "desc"
)
