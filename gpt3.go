package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	gogpt "github.com/sashabaranov/go-gpt3"
)

func gpt3() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	c := gogpt.NewClient(os.Getenv("KEY"))
	ctx := context.Background()

	req := gogpt.CompletionRequest{
		MaxTokens: 5,
		Prompt:    "Lorem ipsum",
	}
	resp, err := c.CreateCompletion(ctx, "ada", req)
	if err != nil {
		return
	}
	fmt.Println(resp.Choices[0].Text)

	searchReq := gogpt.SearchRequest{
		Documents: []string{"White House", "hospital", "school"},
		Query:     "the president",
	}
	searchResp, err := c.Search(ctx, "ada", searchReq)
	if err != nil {
		return
	}
	fmt.Println(searchResp.SearchResults)
}
