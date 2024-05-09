package main

import (
	"fmt"
	"github.com/ilborsch/openai-go/openai"
)

func main() {
	APIKey := "your-api-key"

	// create a client
	client := openai.New(APIKey)

	// create a store for assistant files
	vsID, err := client.Assistants.VectorStores.Create("Crypto scam data storage")
	if err != nil {
	}

	// create an assistant
	assistantID, err := client.Assistants.Create(
		"My New Assistant!",       // name
		"Be polite and friendly!", // instructions
		vsID,                      // vector store ID
		nil,                       // tools to use (file search is default)
	)
	if err != nil {
	}
	fmt.Println("My assistant ID: " + assistantID)
	// Output: "My assistant ID: asst_jLkZWgH4uio1zM0HoW2MXbL"
}
