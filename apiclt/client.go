package apiclt

import (
	"bytes"
	"encoding/json"
	"io"
	"maps"
	"net/http"
	"net/url"
)

type Client struct {
	baseURL string
	headers http.Header
}

type APIResponse struct {
	StatusCode int
	Data       []byte
	APIError   *APIError
}

func (r *APIResponse) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

type APIError struct {
	Name    string `json:"name"`
	Code    string `json:"code"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func (e *APIError) Error() string {
	return e.Message
}

type ClientOpts struct {
	BaseURL     string
	Port        uint
	Headers     http.Header
	BearerToken string
}

func NewClient(opts ClientOpts) Client {
	headers := http.Header{}
	headers.Add("Content-Type", "application/json")

	// Apply custom headers
	maps.Copy(headers, opts.Headers)

	if opts.BearerToken != "" {
		headers.Add("Authorization", "Bearer "+opts.BearerToken)
	}

	c := Client{
		baseURL: opts.BaseURL,
		headers: headers,
	}

	return c
}

type FetchOpts struct {
	Path    string
	Query   map[string]string
	Method  string
	Body    any
	Headers http.Header
}

func (c *Client) Fetch(opts FetchOpts) (*APIResponse, error) {
	reqURL, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, err
	}

	reqURL.Path, err = url.JoinPath(reqURL.Path, opts.Path)
	if err != nil {
		return nil, err
	}

	query := reqURL.Query()
	for k, v := range opts.Query {
		query.Add(k, v)
	}
	reqURL.RawQuery = query.Encode()

	method := "GET"
	if opts.Method != "" {
		method = opts.Method
	}

	var body []byte
	switch v := opts.Body.(type) {
	case nil:
		body = nil
	case json.RawMessage:
		body = v // If already raw JSON, assign directly
	default:
		// All other types need to be serialized
		serialized, err := json.Marshal(opts.Body)
		if err != nil {
			return nil, err
		}
		body = serialized
	}
	bodyRdr := bytes.NewReader(body)

	req, err := http.NewRequest(method, reqURL.String(), bodyRdr)
	if err != nil {
		return nil, err
	}

	// Apply headers
	headers := c.headers.Clone()
	maps.Copy(headers, opts.Headers)
	req.Header = headers

	httpRes, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpRes.Body.Close()

	data, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return nil, err
	}

	apiRes := APIResponse{
		StatusCode: httpRes.StatusCode,
		Data:       data,
	}

	if !apiRes.IsSuccess() {
		// Try to unmarshal the error data
		var apiErr APIError
		if err := json.Unmarshal(data, &apiErr); err != nil {
			return nil, err
		}

		apiRes.APIError = &apiErr
	}
	return &apiRes, nil
}

// Perform is a convenience function that fetches JSON data and unmarshals it into a struct.
func Perform[T any](client *Client, opts FetchOpts) (*T, error) {
	res, err := client.Fetch(opts)
	if err != nil {
		return nil, err
	}

	var resource T
	if err := json.Unmarshal(res.Data, &resource); err != nil {
		return nil, err
	}

	return &resource, nil
}
