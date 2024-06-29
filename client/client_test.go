package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/diverged/tavily-go/models"
	"github.com/stretchr/testify/assert"
)

func TestNewTavilyClient(t *testing.T) {
	apiKey := "test-api-key"
	client := NewTavilyClient(apiKey)

	assert.NotNil(t, client)
	assert.Equal(t, "https://api.tavily.com/search", client.BaseURL)
	assert.Equal(t, apiKey, client.APIKey)
	assert.NotNil(t, client.HTTPClient)
}

func TestSearch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var req models.SearchRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)

		assert.Equal(t, "test query", req.Query)
		assert.Equal(t, "test-api-key", req.APIKey)

		resp := models.SearchResponse{
			Results: []models.SearchResult{
				{Title: "Test Result", URL: "https://example.com", Content: "Test content", Score: 0.9},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := &TavilyClient{
		BaseURL:    server.URL,
		APIKey:     "test-api-key",
		HTTPClient: server.Client(),
	}

	req := models.SearchRequest{Query: "test query"}
	resp, err := client.Search(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Results, 1)
	assert.Equal(t, "Test Result", resp.Results[0].Title)
}

func TestGetSearchContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := models.SearchResponse{
			Results: []models.SearchResult{
				{URL: "https://example.com", Content: "Test content"},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := &TavilyClient{
		BaseURL:    server.URL,
		APIKey:     "test-api-key",
		HTTPClient: server.Client(),
	}

	context, err := client.GetSearchContext("test query", "auto", 100)

	assert.NoError(t, err)
	assert.Contains(t, context, "https://example.com")
	assert.Contains(t, context, "Test content")
}

func TestQnASearch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := models.SearchResponse{
			Answer: "Test answer",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := &TavilyClient{
		BaseURL:    server.URL,
		APIKey:     "test-api-key",
		HTTPClient: server.Client(),
	}

	answer, err := client.QnASearch("test question", "auto")

	assert.NoError(t, err)
	assert.Equal(t, "Test answer", answer)
}

func TestGetCompanyInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := models.SearchResponse{
			Results: []models.SearchResult{
				{Title: "Company News", URL: "https://example.com/news", Content: "News content", Score: 0.9},
				{Title: "Company Finance", URL: "https://example.com/finance", Content: "Finance content", Score: 0.8},
				{Title: "Company News 2", URL: "https://example.com/news2", Content: "More news content", Score: 0.7},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := &TavilyClient{
		BaseURL:    server.URL,
		APIKey:     "test-api-key",
		HTTPClient: server.Client(),
	}

	results, err := client.GetCompanyInfo("test company", "auto", 2)

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "Company News", results[0].Title)
	assert.Equal(t, "Company Finance", results[1].Title)
}
