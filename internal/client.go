package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client wraps the Fan Token Intel REST API.
type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

// NewClient creates a Client. apiKey may be empty for public endpoints.
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL: baseURL,
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) do(req *http.Request) ([]byte, error) {
	if c.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.APIKey)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "fti-cli/1.0")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode >= 400 {
		var apiErr struct {
			Detail string `json:"detail"`
		}
		if e := json.Unmarshal(body, &apiErr); e == nil && apiErr.Detail != "" {
			return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, apiErr.Detail)
		}
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	return body, nil
}

// Get performs a GET request. If out is non-nil the body is JSON-decoded into it.
// The raw body bytes are always returned.
func (c *Client) Get(path string, params url.Values, out interface{}) ([]byte, error) {
	endpoint := c.BaseURL + path
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	body, err := c.do(req)
	if err != nil {
		return nil, err
	}

	if out != nil {
		if err := json.Unmarshal(body, out); err != nil {
			return nil, fmt.Errorf("parsing response: %w", err)
		}
	}
	return body, nil
}

// Post performs a POST request with a JSON payload.
func (c *Client) Post(path string, payload interface{}, out interface{}) ([]byte, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.BaseURL+path, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	body, err := c.do(req)
	if err != nil {
		return nil, err
	}

	if out != nil {
		if err := json.Unmarshal(body, out); err != nil {
			return nil, fmt.Errorf("parsing response: %w", err)
		}
	}
	return body, nil
}
