package main

import (
	"bufio"
	"fmt"
	"github.com/ilborsch/openai-go/openai"
	"github.com/ilborsch/openai-go/openai/assistants/runs"
	"os"
	"strings"
	"time"
)

func main() {
	// Let's assume you have already created an assistant and successfully attached a few files to it.
	// This is an example of how you may utilize library's very simple API in your own project
	// In this example you will see how to:
	// 1. Create a thread
	// 2. Add messages to the thread
	// 3. Create runs for the messages
	// 4. Pool for the assistant response

	// In addition, strongly recommend to get familiar with this example before:
	// https://github.com/ilborsch/openai-go/blob/main/examples/run-thread.go

	APIKey := "your api key"
	assistantID := "your assistant id"

	// initialize client
	client := openai.New(APIKey)

	// create a thread
	threadID, err := client.CreateThread()
	if err != nil {
		// handler error
	}

	// Now we are going to enter endless loop to ask for user inputs (until user enters "q")
	// In this loop we:
	// 1. Add the message to the thread
	// 2. Create a run for this message.
	// 3. Pool for assistant response every X seconds (0.2 sec in our example)

	const poolInterval = 200 * time.Millisecond

	for {
		userPrompt := readLine()
		if strings.ToLower(userPrompt) == "q" {
			return
		}

		// add the message to the thread
		err = client.AddMessageToThread(threadID, userPrompt)
		if err != nil {
			// handle error (f.e. wrong threadID passed)
		}

		// create a run for the message
		runID, err := client.CreateRun(threadID, assistantID)
		if err != nil {
			// handle error (f.e. wrong assistantID, threadID or connection error)
		}

		// In a real-world application you don't want to block main application and wait until assistant responds
		// You would typically run response pooling loop in a different goroutine
		// This is a very simple example of how you can do it

		response := make(chan string)
		go getAssistantResponse(response, client, runID, threadID, poolInterval)
		fmt.Printf("Assistant response: %s \n", <-response)
	}
}

func readLine() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your message (enter \"q\" to exit): ")
	input, _ := reader.ReadString('\n')
	return input
}

func getAssistantResponse(
	responseChan chan string,
	client *openai.OpenAI,
	runID string,
	threadID string,
	duration time.Duration,
) {
	runStatus := ""
	// start pooling
	for {
		// get current run status
		run, err := client.GetRun(threadID, runID)
		if err != nil {
			// handle error, for example:
			responseChan <- err.Error()
			close(responseChan)
			return
		}

		runStatus = run.Status
		if runStatus == runs.StatusCompleted {
			// Here we are going to fetch a list of latest thread messages
			// Retrieve latest assistant response from it and quit this goroutine

			// We developed handy function to extract latest assistant response from a list of thread messages
			// in one API call
			// so you don't need to go into details of OpenAI REST API data structure...
			assistantResponse, err := client.LatestAssistantResponse(threadID)
			if err != nil {
				// handle error
			}
			responseChan <- assistantResponse
			close(responseChan)
			return
		}
		// sleep status is still incomplete
		time.Sleep(duration)
	}

	// You may also implement pooling
}
