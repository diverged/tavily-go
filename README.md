# tavily-go

An unofficial Go client for the Tavily search API, modeled closely on the official [Tavily Python Wrapper](https://github.com/tavily-ai/tavily-python).

## Installation

To install the tavily-go package, use the following command:

```shell
go get github.com/diverged/tavily-go
```

## Usage

Below is a simple example of how to use the tavily-go package:

```go
package main

import (
    "fmt"
    "log"

    tavilygo "github.com/diverged/tavily-go"
    "github.com/diverged/tavily-go/models"
)

func main() {
    // Create a new Tavily client
    client := tavilygo.NewClient("your-api-key-here")

    // Perform a search
    searchReq := models.SearchRequest{
        Query:       "What is the capital of France?",
        SearchDepth: "basic",
    }
    
    response, err := tavilygo.Search(client, searchReq)
    if err != nil {
        log.Fatalf("Error performing search: %v", err)
    }

    // Print the search results
    for _, result := range response.Results {
        fmt.Printf("Title: %s\nURL: %s\n\n", result.Title, result.URL)
    }

    // Perform a Q&A search
    answer, err := tavilygo.QnASearch(client, "What is the capital of France?", "basic")
    if err != nil {
        log.Fatalf("Error performing Q&A search: %v", err)
    }

    fmt.Printf("Answer: %s\n", answer)
}
```

## License

[MIT](https://github.com/diverged/tavily-go/blob/main/LICENSE)
