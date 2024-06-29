package utils

import (
	"encoding/json"
	"strings"

	"github.com/diverged/tavily-go/config"
)

// GetTotalTokensFromString calculates the total number of tokens in a string
// Note: This is a simplified version and may not be as accurate as the Python tiktoken library
func GetTotalTokensFromString(s string, encodingName string) int {
	// Simple approximation: count words and punctuation
	return len(strings.Fields(s)) + strings.Count(s, ".") + strings.Count(s, ",") + strings.Count(s, "!") + strings.Count(s, "?")
}

// GetMaxTokensFromString extracts a substring with a maximum number of tokens
func GetMaxTokensFromString(s string, maxTokens int, encodingName string) string {
	words := strings.Fields(s)
	result := []string{}
	count := 0

	for _, word := range words {
		wordTokens := GetTotalTokensFromString(word, encodingName)
		if count+wordTokens > maxTokens {
			break
		}
		result = append(result, word)
		count += wordTokens
	}

	return strings.Join(result, " ")
}

// GetMaxItemsFromList returns a JSON string of items from a list, limited by max tokens
func GetMaxItemsFromList(data []interface{}, maxTokens int) (string, error) {
	result := []string{}
	currentTokens := 0

	for _, item := range data {
		itemJSON, err := json.Marshal(item)
		if err != nil {
			return "", err
		}

		itemStr := string(itemJSON)
		newTotalTokens := currentTokens + GetTotalTokensFromString(itemStr, config.DefaultModelEncoding)

		if newTotalTokens > maxTokens {
			break
		}

		result = append(result, itemStr)
		currentTokens = newTotalTokens
	}

	return "[" + strings.Join(result, ",") + "]", nil
}
