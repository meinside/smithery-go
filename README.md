# smithery-go

A go library for using [smithery](https://smithery.ai/) APIs, built with [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk).

## Usage

```go
package main

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

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
	if session, err := client.ConnectWithProfileID(
		context.TODO(),
		profileID,
		qualifiedName,
	); err == nil {
		// NOTE: do not forget to close it after use
		defer session.Close()

		// do various things with the connection,
		// eg. get tools from the connection, generate with your LLM, and call the corresponding tools

		// list tools,
		if tools, err := session.ListTools(
			context.TODO(),
			&mcp.ListToolsParams{},
		); err == nil {
			fmt.Printf("> tools = %+v\n", tools)

			// do something with your LLM,
			fnNameFromYourLLM := "web_search_exa"
			fnArgsFromYourLLM := map[string]any{
				"query": "shoebill",
			}

			// and call the corresponding tool with your function arguments
			if result, err := session.CallTool(
				context.TODO(),
				&mcp.CallToolParams{
					Name:      fnNameFromYourLLM,
					Arguments: fnArgsFromYourLLM,
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

