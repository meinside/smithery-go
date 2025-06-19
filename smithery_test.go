// test_smithery.go

package smithery

import (
	"context"
	"os"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

// check and return environment variable for given key
func mustHaveEnvVar(t *testing.T, key string) string {
	if value, exists := os.LookupEnv(key); !exists {
		t.Fatalf("no environment variable: %s", key)
	} else {
		return value
	}
	return ""
}

// TestRegistryFunctions tests registry functions.
func TestRegistryFunctions(t *testing.T) {
	client := NewClient(mustHaveEnvVar(t, "API_TOKEN"))
	client.Verbose = os.Getenv("VERBOSE") == "true"

	// test `ListServers`
	if res, err := client.ListServers(
		context.TODO(),
		WithPage(1),
		WithPageSize(10),
		WithQuery("is:verified"),
	); err != nil {
		t.Errorf("failed to list servers: %s", err)
	} else {
		client.verbose("listed servers: %s", prettify(res))
	}

	// test `GetServer`
	if res, err := client.GetServer(
		context.TODO(),
		`@smithery/toolbox`,
	); err != nil {
		t.Errorf("failed to get server: %s", err)
	} else {
		client.verbose("server: %s", prettify(res))
	}
}

// TestConnections tests connections to server.
func TestConnections(t *testing.T) {
	client := NewClient(mustHaveEnvVar(t, "API_TOKEN"))
	client.Verbose = os.Getenv("VERBOSE") == "true"

	// test `ConnectWithProfileID`
	if mcpc, err := client.ConnectWithProfileID(
		context.TODO(),
		mustHaveEnvVar(t, "PROFILE_ID"),
		`exa`,
	); err != nil {
		t.Errorf("failed to connect to server with profile id: %s", err)
	} else {
		defer mcpc.Close() // NOTE: should be closed manually after use

		// list tools,
		if tools, err := mcpc.ListTools(
			context.TODO(),
			mcp.ListToolsRequest{},
		); err != nil {
			t.Errorf("failed to list tools: %s", err)
		} else {
			client.verbose("listed tools: %s", prettify(tools))

			// call tool,
			if result, err := mcpc.CallTool(
				context.TODO(),
				mcp.CallToolRequest{
					Request: mcp.Request{
						Method: `tools/call`,
					},
					Params: mcp.CallToolParams{
						Name: "web_search_exa",
						Arguments: map[string]any{
							"query": "mcp and smithery",
						},
					},
				},
			); err != nil {
				t.Errorf("failed to call tool: %s", err)
			} else {
				client.verbose("call tool result: %s", prettify(result))
			}
		}
	}

	// test `ConnectManually`
	if mcpc, err := client.ConnectManually(
		context.TODO(),
		mustHaveEnvVar(t, "NAVER_SERVER_URL"),
		map[string]any{
			"NAVER_CLIENT_ID":     mustHaveEnvVar(t, "NAVER_CLIENT_ID"),
			"NAVER_CLIENT_SECRET": mustHaveEnvVar(t, "NAVER_CLIENT_SECRET"),
		},
	); err != nil {
		t.Errorf("failed to connect to server with config: %s", err)
	} else {
		defer mcpc.Close() // NOTE: should be closed manually after use

		// list tools,
		if tools, err := mcpc.ListTools(
			context.TODO(),
			mcp.ListToolsRequest{},
		); err != nil {
			t.Errorf("failed to list tools: %s", err)
		} else {
			client.verbose("listed tools: %s", prettify(tools))

			// call tool,
			if result, err := mcpc.CallTool(
				context.TODO(),
				mcp.CallToolRequest{
					Request: mcp.Request{
						Method: `tools/call`,
					},
					Params: mcp.CallToolParams{
						Name: "search_image",
						Arguments: map[string]any{
							"query": "shoebill",
						},
					},
				},
			); err != nil {
				t.Errorf("failed to call tool: %s", err)
			} else {
				client.verbose("call tool result: %s", prettify(result))
			}
		}
	}
}
