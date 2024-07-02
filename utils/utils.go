package utils

import (
	"encoding/json"
	"strings"

	"github.com/pkoukk/tiktoken-go"

	"github.com/diverged/tavily-go/config"
)

// GetTotalTokensFromString calculates the total number of tokens in a given string
// using the specified encoding model.
func GetTotalTokensFromString(s string, encodingModel string) (int, error) {
	tkm, err := tiktoken.EncodingForModel(encodingModel)
	if err != nil {
		return 0, err
	}
	tokens := tkm.Encode(s, nil, nil)
	return len(tokens), nil
}

// GetMaxTokensFromString truncates a given string to a maximum number of tokens
// using the specified encoding model. It returns the truncated string and an error
// if the encoding model is not found or if tokenization fails.
func GetMaxTokensFromString(s string, maxTokens int, encodingModel string) (string, error) {
	tkm, err := tiktoken.EncodingForModel(encodingModel)
	if err != nil {
		return "", err
	}
	tokens := tkm.Encode(s, nil, nil)
	if len(tokens) <= maxTokens {
		return s, nil
	}
	decodedTokens := tkm.Decode(tokens[:maxTokens])
	return strings.TrimSpace(string(decodedTokens)), nil
}

// GetMaxItemsFromList takes a list of items and returns a JSON string containing
// as many items as possible without exceeding the specified maximum number of tokens.
// It returns an error if JSON marshaling fails or if the encoding model is not found.
func GetMaxItemsFromList(items []interface{}, maxTokens int) (string, error) {
	tkm, err := tiktoken.EncodingForModel(config.DefaultModelEncoding)
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

		itemTokens := tkm.Encode(string(itemJSON), nil, nil)
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
