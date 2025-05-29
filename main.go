package ofdconnector

import (
	"context"
	"fmt"

	"gitlab.udevs.io/billz/ofd_connector/pkg/httpclient"
)

// OfdConnector interface defines the contract for OFD service operations
type OfdConnector interface {
	FiscalDriveList(context.Context) ([]FiscalDriveReaderInfo, error)
	SetFactoryID(factoryID string)
	// Add methods as needed
}

type OfdConnectorConfigs struct {
	ServiceAddress        string
	RequestTimeOutSeconds int
}

// ofdConnector implements the OfdConnector interface
type ofdConnector struct {
	serviceAddress string
	httpClient     httpclient.HTTPClient
	factoryID      string
}

// New creates a new instance of OfdConnector
func New(configs OfdConnectorConfigs) (OfdConnector, error) {
	if configs.ServiceAddress == "" {
		return nil, fmt.Errorf("ofdServiceAddress cannot be empty %w", ErrorInvalidUrlAddress)
	}

	httpClient := httpclient.NewHTTPClient(configs.RequestTimeOutSeconds)

	return &ofdConnector{
		serviceAddress: configs.ServiceAddress,
		httpClient:     httpClient,
	}, nil
}

func (o *ofdConnector) SetFactoryID(factoryID string) {
	o.factoryID = factoryID
}
