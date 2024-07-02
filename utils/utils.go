package utils

import (
	"encoding/json"
	"strings"

	"github.com/pkoukk/tiktoken-go"

	"github.com/diverged/tavily-go/config"
)

func GetTotalTokensFromString(s string, encodingModel string) (int, error) {
	tkm, err := tiktoken.EncodingForModel(encodingModel)
	if err != nil {
		return 0, err
	}
	tokens := tkm.Encode(s, nil, nil)
	return len(tokens), nil
}

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
