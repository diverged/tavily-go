package utils

import (
	"encoding/json"
	"strings"

	"github.com/pkoukk/tiktoken-go"
)

// GetMaxItemsFromList returns a JSON string of items from a list, limited by max tokens
const DefaultEncoding = "cl100k_base" // This is the encoding used by GPT-3.5 and GPT-4 models

func GetTotalTokensFromString(s string, encodingName string) (int, error) {
	tke, err := tiktoken.GetEncoding(encodingName)
	if err != nil {
		return 0, err
	}
	tokens := tke.Encode(s, nil, nil)
	return len(tokens), nil
}

func GetMaxTokensFromString(s string, maxTokens int, encodingName string) (string, error) {
	tke, err := tiktoken.GetEncoding(encodingName)
	if err != nil {
		return "", err
	}
	tokens := tke.Encode(s, nil, nil)
	if len(tokens) <= maxTokens {
		return s, nil
	}
	decodedTokens := tke.Decode(tokens[:maxTokens])
	return strings.TrimSpace(string(decodedTokens)), nil
}

func GetMaxItemsFromList(items []interface{}, maxTokens int) (string, error) {
	tke, err := tiktoken.GetEncoding(DefaultEncoding)
	if err != nil {
		return "", err
	}

	result := make([]interface{}, 0)
	currentTokens := 2 // Start with 2 tokens for the opening and closing brackets of the JSON array

	for _, item := range items {
		itemJSON, err := json.Marshal(item)
		if err != nil {
			return "", err
		}

		itemTokens := tke.Encode(string(itemJSON), nil, nil)
		if len(result) > 0 {
			currentTokens++ // Add 1 token for the comma separator
		}

		if currentTokens+len(itemTokens) > maxTokens {
			break
		}

		result = append(result, item)
		currentTokens += len(itemTokens)
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(resultJSON), nil
}
