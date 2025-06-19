# smithery-go

A go library for using [smithery](https://smithery.ai/) APIs, built with [mark3labs/mcp-go](https://github.com/mark3labs/mcp-go).

## Usage

```go
package main

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/meinside/smithery-go"
)

const (
	apiToken  = `your-smithery-api-token`
	profileID = `your-smithery-profile-id`

	qualifiedName = `exa` // https://smithery.ai/server/exa
)

func main() {
	client := smithery.NewClient(apiToken)

	// list servers,
	if servers, err := client.ListServers(context.TODO()); err == nil {
		fmt.Printf("> servers = %+v\n", servers)
	} else {
		fmt.Printf("* failed to list servers: %s\n", err)
	}

	// get server,
	if server, err := client.GetServer(context.TODO(), qualifiedName); err == nil {
		fmt.Printf("> server for qualified name '%s': %+v\n", qualifiedName, server)
	} else {
		fmt.Printf("* failed to get server: %s\n", err)
	}

	// connect,
	if conn, err := client.ConnectWithProfileID(
		context.TODO(),
		profileID,
		qualifiedName,
	); err == nil {
		defer conn.Close()	// NOTE: do not forget to close it after use

		// do various things with the connection,
		// eg. get tools from the connection, generate with your LLM, and call the corresponding tools

		// list tools,
		if tools, err := conn.ListTools(
			context.TODO(),
			mcp.ListToolsRequest{},
		); err == nil {
			fmt.Printf("> tools = %+v\n", tools)

			// do something with your LLM,
			fnNameFromYourLLM := "web_search_exa"
			fnArgsFromYourLLM := map[string]any{
				"query": "shoebill",
			}

			// and call the corresponding tool with your function arguments
			if result, err := conn.CallTool(
				context.TODO(),
				mcp.CallToolRequest{
					Request: mcp.Request{
						Method: "tools/call",
					},
					Params: mcp.CallToolParams{
						Name:      fnNameFromYourLLM,
						Arguments: fnArgsFromYourLLM,
					},
				},
			); err == nil {
				fmt.Printf("> call tool result: %+v\n", result)
			} else {
				fmt.Printf("* failed to call tool: %s", err)
			}
		} else {
			fmt.Printf("* failed to list tools: %s\n", err)
		}
	} else {
		fmt.Printf("* failed to connect to server: %s\n", err)
	}
}
```

## License

MIT

