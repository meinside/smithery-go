// smithery.go

// Package smithery provides functions for using smithery.ai APIs.
package smithery

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/fatih/color"
	mcpc "github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/meinside/version-go"
)

const (
	clientName = `meinside/smithery-go`
)

// Client struct
type Client struct {
	apiToken string

	Verbose bool
}

// NewClient returns a new client with given `apiToken`.
func NewClient(apiToken string) *Client {
	return &Client{
		apiToken: apiToken,
	}
}

// ListServers lists servers.
//
// https://smithery.ai/docs/use/registry#list-servers
func (c *Client) ListServers(
	ctx context.Context,
	opts ...RequestOptionListServers,
) (result ResponseServers, err error) {
	// apply options
	params := make(map[string]any)
	for _, opt := range opts {
		params = opt(params)
	}

	// send request and read the response
	var req *http.Request
	if req, err = c.httpGetRequest(
		ctx,
		`https://registry.smithery.ai/servers`,
		params,
	); err == nil {
		client := httpClient()

		var body []byte
		if body, err = c.readHTTPResponse(client, req); err == nil {
			if err = json.Unmarshal(body, &result); err == nil {
				return result, nil
			}
		}
	}

	return ResponseServers{}, err
}

// WithQuery builds a request option for `ListServers` with given `query`.
func WithQuery(query string) RequestOptionListServers {
	return func(params reqParams) reqParams {
		params["q"] = query
		return params
	}
}

// WithPage builds a reuqest option for `ListServers` with given `page`.
func WithPage(page uint) RequestOptionListServers {
	return func(params reqParams) reqParams {
		params["page"] = page
		return params
	}
}

// WithPageSize builds a request option for `ListServers` with given `pageSize`.
func WithPageSize(pageSize uint) RequestOptionListServers {
	return func(params reqParams) reqParams {
		params["pageSize"] = pageSize
		return params
	}
}

// GetServer gets a server.
//
// https://smithery.ai/docs/use/registry#get-server
func (c *Client) GetServer(
	ctx context.Context,
	qualifiedName string,
) (result ResponseServer, err error) {
	// send request and read the response
	var req *http.Request
	if req, err = c.httpGetRequest(
		ctx,
		fmt.Sprintf(`https://registry.smithery.ai/servers/%s`, qualifiedName),
		nil,
	); err == nil {
		client := httpClient()

		var body []byte
		if body, err = c.readHTTPResponse(client, req); err == nil {
			if err = json.Unmarshal(body, &result); err == nil {
				return result, nil
			}
		}
	}

	return ResponseServer{}, err
}

// ConnectWithProfileID connects to server with given `profileID` and `serverName`.
// Returned client(`closer`) should be closed manually after use.
//
// https://smithery.ai/docs/use/connect#using-a-profile-recommended
func (c *Client) ConnectWithProfileID(
	ctx context.Context,
	profileID string,
	serverName string,
) (closer *mcpc.Client, err error) {
	var u *url.URL
	if u, err = url.Parse(fmt.Sprintf(
		`https://server.smithery.ai/%[1]s/mcp`,
		serverName,
	)); err == nil {
		u.RawQuery = getParams(map[string]any{
			"api_key": c.apiToken,
			"profile": profileID,
		}).Encode()

		return c.connect(ctx, u)
	}

	return
}

// ConnectManually connects to server with given `url` and `config`.
// Returned client(`closer`) should be closed manually after use.
//
// https://smithery.ai/docs/use/connect#manual-configuration
func (c *Client) ConnectManually(
	ctx context.Context,
	serverURL string,
	config map[string]any,
) (closer *mcpc.Client, err error) {
	var conf []byte
	if conf, err = json.Marshal(config); err == nil {
		var u *url.URL
		if u, err = url.Parse(serverURL); err == nil {
			u.RawQuery = getParams(map[string]any{
				"api_key": c.apiToken,
				"config":  base64.StdEncoding.EncodeToString(conf),
			}).Encode()

			return c.connect(ctx, u)
		}
	}

	return
}

// connect to server, start, initialize, and return the client
func (c *Client) connect(ctx context.Context, url *url.URL) (closer *mcpc.Client, err error) {
	var streamable *transport.StreamableHTTP
	if streamable, err = transport.NewStreamableHTTP(
		url.String(),
		transport.WithHTTPBasicClient(httpClient()),
	); err == nil {
		closer = mcpc.NewClient(
			streamable,
			mcpc.WithClientCapabilities(mcp.ClientCapabilities{}),
		)

		if err = closer.Start(ctx); err == nil {
			var initialized *mcp.InitializeResult
			if initialized, err = closer.Initialize(ctx, mcp.InitializeRequest{
				Params: mcp.InitializeParams{
					ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
					ClientInfo: mcp.Implementation{
						Name:    clientName,
						Version: version.Build(version.OS | version.Architecture),
					},
				},
			}); err == nil {
				c.verbose(">>> initialized server: %s", prettify(initialized))

				return closer, nil
			}
		}
	}

	return nil, err
}

// print message to stdout if verbose mode is enabled
func (c *Client) verbose(
	format string,
	a ...any,
) {
	if c.Verbose {
		if !strings.HasSuffix(format, "\n") {
			format += "\n"
		}

		printColored(
			color.FgYellow,
			format,
			a...,
		)
	}
}
