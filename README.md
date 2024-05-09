# Go OpenAI API

## Introduction

`openai-go` is an API for Go language that provides intuitive and easy-to-understand interface to communicate with OpenAI API.

The library abstracts complexity of OpenAI REST API and provides user-friendly API to interact with instead.

## Install

Run this command to install latest version of library:

```bash
  go install github.com/ilborsch/openai-go-old@latest
```

## Examples

See [examples/](https://github.com/ilborsch/openai-go/tree/main/examples) for a variety of examples.

**As easy as:**

```go

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
	if err != nil {}

	// create an assistant
	assistantID, err := client.Assistants.Create(
		"My New Assistant!",       // name
		"Be polite and friendly!", // instructions
		vsID,                      // vector store ID
		nil,                       // tools to use (file-search by default)
	)
	if err != nil {}
	
	fmt.Println(assistantID) 
	// Output: asst_jLkZWgH4uio1zM0HoW2MXbL
}

```

## Credits

* Sasha Draganov for inspiration and help

We'll be more than happy to see [your contributions](./CONTRIBUTING.md)!

## Authors

- [@ilborsch](https://www.github.com/ilborsch)


## License

[MIT](https://choosealicense.com/licenses/mit/)

