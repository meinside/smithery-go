// helpers.go

package smithery

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

const (
	defaultTimeoutSeconds               = 60
	defaultIdleTimeoutSeconds           = 90
	defaultTLSHandshakeTimeoutSeconds   = 5
	defaultResponseHeaderTimeoutSeconds = 30
	defaultExpectContinueTimeoutSeconds = 1
)

// for reusing http client
var _httpClient *http.Client

// helper function for generating a http client
func httpClient() *http.Client {
	if _httpClient == nil {
		_httpClient = &http.Client{
			Timeout: defaultTimeoutSeconds * time.Second,
			Transport: &http.Transport{
				IdleConnTimeout:       defaultIdleTimeoutSeconds * time.Second,
				TLSHandshakeTimeout:   defaultTLSHandshakeTimeoutSeconds * time.Second,
				ResponseHeaderTimeout: defaultResponseHeaderTimeoutSeconds * time.Second,
				ExpectContinueTimeout: defaultExpectContinueTimeoutSeconds * time.Second,
			},
		}
	}
	return _httpClient
}

// helper function for converting get params
func getParams(params map[string]any) url.Values {
	converted := url.Values{}
	for k, v := range params {
		converted.Add(k, fmt.Sprintf("%v", v))
	}
	return converted
}

// helper function for generating a http request
func (c *Client) httpRequest(ctx context.Context, method, url string) (req *http.Request, err error) {
	if req, err = http.NewRequestWithContext(ctx, method, url, nil); err == nil {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.apiToken))
		req.Header.Add("Accept", "application/json")
		return req, nil
	}
	return
}

// helper function for generating a http get request
func (c *Client) httpGetRequest(ctx context.Context, getURL string, params map[string]any) (req *http.Request, err error) {
	getParams := getParams(params)

	var u *url.URL
	if u, err = url.Parse(getURL); err == nil {
		u.RawQuery = getParams.Encode()

		var req *http.Request
		if req, err = c.httpRequest(ctx, "GET", u.String()); err == nil {
			return req, nil
		}
	}
	return
}

// helper function for requesting a http response and reading its body
func (c *Client) readHTTPResponse(
	client *http.Client,
	req *http.Request,
) (body []byte, err error) {
	if c.Verbose {
		if dumped, err := httputil.DumpRequest(req, true); err == nil {
			log.Printf(">>> dump of request:\n\n%s", string(dumped))
		}
	}

	var resp *http.Response
	if resp, err = client.Do(req); err == nil {
		defer resp.Body.Close()

		if c.Verbose {
			if dumped, err := httputil.DumpResponse(resp, true); err == nil {
				log.Printf(">>> dump of response:\n\n%s", string(dumped))
			}
		}

		if resp.StatusCode == 200 {
			if body, err = io.ReadAll(resp.Body); err == nil {
				return body, nil
			}
		} else {
			body, _ = io.ReadAll(resp.Body)
			err = fmt.Errorf("http error %d (%s)", resp.StatusCode, string(body))
		}
	}
	return
}

// prettify given thing in JSON format
func prettify(v any, flatten ...bool) string {
	if len(flatten) > 0 && flatten[0] {
		if bytes, err := json.Marshal(v); err == nil {
			return string(bytes)
		}
	} else {
		if bytes, err := json.MarshalIndent(v, "", "  "); err == nil {
			return string(bytes)
		}
	}
	return fmt.Sprintf("%+v", v)
}
