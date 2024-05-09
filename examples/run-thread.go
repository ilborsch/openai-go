package main

import (
	"fmt"
	"github.com/ilborsch/openai-go/openai"
)

func main() {
	// Let's assume you have already created an assistant and successfully attached a few files to it.
	// In this example you will learn how to:
	// 1. Create a thread
	// 2. Add a message to the thread
	// 3. Create a run for the thread.

	APIKey := "your api key"
	assistantID := "your assistant id"

	// initialize client
	client := openai.New(APIKey)

	threadID, err := client.Threads.Create()
	if err != nil {
		// handle error
	}

	userPrompt := "Hi! What can you do?"
	err = client.Messages.AddMessageToThread(threadID, userPrompt)
	if err != nil {
		// handle error
	}

	runID, err := client.Runs.Create(threadID, assistantID)
	if err != nil {
		// handle error :)
	}

	fmt.Println(runID)

}
