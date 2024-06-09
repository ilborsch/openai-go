package main

import (
	"fmt"
	"github.com/ilborsch/openai-go/openai"
	"io"
	"os"
)

func main() {
	// In this example you will learn how to:
	// 1. Create a vector store
	// 2. Attach files to the vector store
	// 3. Create an assistant with vector store

	APIKey := "your api key"

	// init client
	client := openai.New(APIKey)

	// create a vector store
	vectorStoreID, err := client.CreateVectorStore("Crypto scam guides storage")
	// or alternatively (thanks to Go syntax sugar)
	// vectorStoreID, err := client.VectorStores.Create(...)
	if err != nil {
		// handle error
	}

	fileName := "some-data.txt"
	file, _ := os.Open(fileName)
	fileData, _ := io.ReadAll(file)

	// upload the file to OpenAI portal, you can reference it later within OpenAI by fileID
	fileID, err := client.UploadFile(fileName, fileData)
	if err != nil {
		// handle error
	}

	// attach file to the vector store you created before
	err = client.AddVectorStoreFile(vectorStoreID, fileID)
	if err != nil {
		// handle error
	}

	// create an assistant
	assistantID, err := client.CreateAssistant(
		"AI Tutor", // name
		"You are an AI tutor. Be polite and friendly!", // instructions
		vectorStoreID, // vector store ID
		nil,           // tools

		//[]assistants.Tool{
		//	{
		//		Type: assistants.ToolFileSearch,
		//	},
		//},
		// alternatively you could specify tools like so... But... Really???
	)
	// file search tool is specified by default
	// therefore, for the sake of simplicity we can leave this argument as nil

	fmt.Println(assistantID)
	// Output: asst_jLkZWgH4uio1zM0HoW2MXbL
}
