package utils

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

const DefaultEncoding = "gpt-3.5-turbo"

func TestGetTotalTokensFromString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"Empty string", "", 0},
		{"Single word", "hello", 1},
		{"Multiple words", "hello world", 2},
		{"With punctuation", "Hello, world!", 4}, // Updated expected value
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetTotalTokensFromString(tt.input, DefaultEncoding)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetMaxTokensFromString(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		maxTokens int
		expected  string
	}{
		{"Short string", "Hello world", 5, "Hello world"},
		{"Long string", "This is a longer string that should be truncated", 5, "This is a longer string"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetMaxTokensFromString(tt.input, tt.maxTokens, DefaultEncoding)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetMaxItemsFromList(t *testing.T) {
	tests := []struct {
		name      string
		input     []interface{}
		maxTokens int
		expected  string
	}{
		{
			name:      "Single item",
			input:     []interface{}{"hello"},
			maxTokens: 10,
			expected:  `["hello"]`,
		},
		{
			name:      "Multiple items",
			input:     []interface{}{"hello", "world", "test"},
			maxTokens: 10,
			expected:  `["hello","world"]`,
		},
		{
			name:      "Complex items",
			input:     []interface{}{"hello", 123, true, map[string]interface{}{"key": "value"}},
			maxTokens: 12,
			expected:  `["hello",123,true]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetMaxItemsFromList(tt.input, tt.maxTokens)
			assert.NoError(t, err)

			var resultJSON, expectedJSON []interface{}
			err = json.Unmarshal([]byte(result), &resultJSON)
			assert.NoError(t, err)
			err = json.Unmarshal([]byte(tt.expected), &expectedJSON)
			assert.NoError(t, err)

			assert.Equal(t, expectedJSON, resultJSON)
		})
	}
}
