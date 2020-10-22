package gotruevault

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSearchTypeMarshalling(t *testing.T) {
	expected := `{
  "test": {
    "type": "eq",
    "value": "0001-01-01T00:00:00Z"
  },
  "test1": {
    "type": "in",
    "value": [
      "test1",
      "test2"
    ]
  },
  "test2": {
    "type": "not_in",
    "value": [
      "test3",
      "test4"
    ]
  },
  "test3": {
    "type": "wildcard",
    "value": "wildcard*"
  },
  "test4": {
    "type": "range",
    "value": {
      "gt": 3,
      "lt": 5
    }
  }
}`

	e := map[string]SearchType{
		"test": Eq{
			Value: Time{time.Time{}},
		},
		"test1": In{
			Value: SearchValues{
				String{"test1"},
				String{"test2"},
			},
		},
		"test2": NotIn{
			Value: SearchValues{
				String{"test3"},
				String{"test4"},
			},
		},
		"test3": Wildcard{
			Value: String{"wildcard*"},
		},
		"test4": Range{
			Value: RangeValue{
				Gt: 3,
				Lt: 5,
			},
		},
	}

	b, err := json.MarshalIndent(e, "", "  ")

	assert.Nil(t, err)
	assert.Equal(t, string(b), expected)
}
