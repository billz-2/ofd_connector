package ofdconnector

import (
	"fmt"

	"github.com/billz-2/ofd_connector/internal/httpclient"
)

// OfdConnector interface defines the contract for OFD service operations
type OfdConnector interface {
	ZReport() ZReportI
}

type OfdConnectorConfigs struct {
	ServiceAddress        string
	RequestTimeOutSeconds int
	FactoryID             string
}

// ofdConnector implements the OfdConnector interface
type ofdConnector struct {
	serviceAddress string
	httpClient     httpclient.HTTPClient
	factoryID      string
	zReport        ZReportI
}

// New creates a new instance of OfdConnector
func New(configs OfdConnectorConfigs) (OfdConnector, error) {
	if configs.ServiceAddress == "" {
		return nil, fmt.Errorf("invalid url address")
	}
	if configs.FactoryID == "" {
		return nil, fmt.Errorf("invalid FactoryID")
	}

	httpClient := httpclient.NewHTTPClient(configs.RequestTimeOutSeconds)
	zReport := newZReport(zReportConfigs{
		ServiceAddress: configs.ServiceAddress,
		FactoryID:      configs.FactoryID,
		HttpClient:     httpClient,
	})

	return &ofdConnector{
		serviceAddress: configs.ServiceAddress,
		httpClient:     httpClient,
		zReport:        zReport,
		factoryID:      configs.FactoryID,
	}, nil
}

func (o *ofdConnector) ZReport() ZReportI {
	return o.ZReport()
}
