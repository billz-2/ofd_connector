package gateway

import (
	"context"

	"github.com/billz-2/ofd_connector/internal/httpclient"
)

type Client interface {
	HTTPRequest(
		ctx context.Context,
		uri string,
		method string,
		contentType string,
		body []byte,
		headers map[string]string,
	) (*httpclient.HTTPResponse, error)
	// FactoryEndpoint returns the endpoint for the given route with factoryID appended to the end
	FactoryEndpoint(route string) string
}

type gateway struct {
	serviceAddress string
	factoryID      string
	httpClient     httpclient.HTTPClient
}

type Config struct {
	ServiceAddress string
	FactoryID      string
	HttpClient     httpclient.HTTPClient
}

func New(config Config) Client {
	return gateway{
		serviceAddress: config.ServiceAddress,
		httpClient:     config.HttpClient,
		factoryID:      config.FactoryID,
	}
}

func (g gateway) HTTPRequest(
	ctx context.Context,
	uri string,
	method string,
	contentType string,
	body []byte,
	headers map[string]string,
) (*httpclient.HTTPResponse, error) {
	endpoint := g.serviceAddress + uri
	request, err := httpclient.NewHTTPRequest(
		endpoint,
		method,
		contentType,
		body,
		headers,
	)
	if err != nil {
		return nil, err
	}

	return g.httpClient.Request(ctx, request), nil
}

func (g gateway) FactoryEndpoint(route string) string {
	return route + g.factoryID
}
