// types.go

package smithery

type reqParams map[string]any

// RequestOptionListServers for request options of `ListServers`.
type RequestOptionListServers func(reqParams) reqParams

// ResponseServers struct for the result of `ListServers`.
//
// https://smithery.ai/docs/use/registry#response-type
type ResponseServers struct {
	Servers []struct {
		QualifiedName string `json:"qualifiedName"`
		DisplayName   string `json:"displayName"`
		Description   string `json:"description"`
		Homepage      string `json:"homepage"`
		IconURL       string `json:"iconUrl"`
		UseCount      int    `json:"useCount"`
		IsDeployed    bool   `json:"isDeployed"`
		Remote        bool   `json:"remote"`
		CreatedAt     string `json:"createdAt"`
	} `json:"servers"`
	Pagination struct {
		CurrentPage int `json:"currentPage"`
		PageSize    int `json:"pageSize"`
		TotalPages  int `json:"totalPages"`
		TotalCount  int `json:"totalCount"`
	} `json:"pagination"`
}

// ResponseServer struct for the result of `GetServer`.
//
// https://smithery.ai/docs/use/registry#response-type-2
type ResponseServer struct {
	QualifiedName string       `json:"qualifiedName"`
	DisplayName   string       `json:"displayName"`
	Description   string       `json:"description"`
	IconURL       string       `json:"iconUrl,omitempty"`
	Remote        bool         `json:"remote"`
	DeploymentURL string       `json:"deploymentUrl,omitempty"`
	Connections   []Connection `json:"connections"`
	Security      Security     `json:"security"`
	Tools         []Tool       `json:"tools,omitempty"`
}

// ConnectionType enum
type ConnectionType string

// ConnectionType constants
const (
	ConnectionTypeHTTP  ConnectionType = "http"
	ConnectionTypeStdio ConnectionType = "stdio"
)

// Connection struct
type Connection struct {
	Type         ConnectionType `json:"type"`
	URL          string         `json:"url,omitempty"` // for Type == ConnectionTypeHTTP
	ConfigSchema map[string]any `json:"configSchema"`
}

// Security struct
type Security struct {
	ScanPassed bool `json:"scanPassed,omitempty"`
}

// Tool struct
type Tool struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	InputSchema map[string]any `json:"inputSchema"`
}
