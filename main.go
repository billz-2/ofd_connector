package ofdconnector

import (
	"context"
	"fmt"

	"github.com/billz-2/ofd_connector/internal/httpclient"
)

// OfdConnector interface defines the contract for OFD service operations
type OfdConnector interface {
	FiscalDriveList(context.Context) ([]FiscalDriveReaderInfo, error)
	SetFactoryID(factoryID string)

	ZReportI
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
		return nil, fmt.Errorf("invalid url address")
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
