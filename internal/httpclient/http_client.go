package httpclient

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	ErrBadResponseStatusCode = errors.New("bad response status code")
)

type HTTPClient interface {
	Request(ctx context.Context, req *HTTPRequest) *HTTPResponse
}

type httpClient struct {
	client *http.Client
}

func NewHTTPClient(timeoutInSeconds int) HTTPClient {
	return &httpClient{
		client: &http.Client{
			Timeout: time.Duration(timeoutInSeconds) * time.Second,
		},
	}
}

func (s *httpClient) Request(ctx context.Context, req *HTTPRequest) *HTTPResponse {
	var requestBody io.Reader
	if req.Body != nil {
		requestBody = strings.NewReader(string(req.Body))
	}
	httpRequest, err := http.NewRequestWithContext(
		ctx,
		req.Method,
		req.URI,
		requestBody,
	)
	if err != nil {
		return &HTTPResponse{Error: err}
	}
	if req.ContentType != "" {
		httpRequest.Header.Add("Content-Type", req.ContentType)
	}
	if req.Headers != nil {
		for k, v := range req.Headers {
			httpRequest.Header.Add(k, v)
		}
	}

	resp, err := s.client.Do(httpRequest)
	if err != nil {
		return &HTTPResponse{Error: err}
	}
	defer func() {
		er := resp.Body.Close()
		if er != nil {
			fmt.Printf("error: %#v\r\n", er)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &HTTPResponse{
			StatusCode: resp.StatusCode,
			Error:      err,
		}
	}

	if resp.StatusCode != http.StatusOK {
		return &HTTPResponse{
			StatusCode: resp.StatusCode,
			Error:      ErrBadResponseStatusCode,
			Body:       body,
		}
	}

	return &HTTPResponse{
		StatusCode: resp.StatusCode,
		Body:       body,
	}
}

type HTTPRequest struct {
	URI         string
	Method      string
	ContentType string
	Body        []byte
	Headers     map[string]string
}

func NewHTTPRequest(uri, method, contentType string, body []byte, headers map[string]string) (*HTTPRequest, error) {
	if uri == "" {
		return nil, fmt.Errorf("uri parametr cannot be empty")
	}
	if method == "" {
		method = http.MethodGet
	}
	if contentType == "" {
		contentType = "application/json"
	}

	if contentType == "application/x-www-form-urlencoded" && len(body) > 0 {
		// For URL encoded content type, ensure body is properly formatted
		bodyValues := map[string]any{}
		err := json.Unmarshal(body, &bodyValues)
		if err != nil {
			return nil, err
		}

		// Convert body to url encoded format
		formUrlValues := url.Values{}
		for k, v := range bodyValues {
			formUrlValues.Add(k, fmt.Sprintf("%v", v))
		}
		body = []byte(formUrlValues.Encode())
	}

	return &HTTPRequest{
		URI:         uri,
		Method:      method,
		ContentType: contentType,
		Body:        body,
		Headers:     headers,
	}, nil
}

type HTTPResponse struct {
	Body       []byte
	StatusCode int
	Error      error
}

func GetBasicAuthHeader(username, password string) map[string]string {
	// Create basic auth header
	auth := username + ":" + password
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))

	return map[string]string{
		"Authorization": "Basic " + encodedAuth,
	}
}
