package ofdconnector

import (
	"fmt"

	"github.com/billz-2/ofd_connector/internal/gateway"
	"github.com/billz-2/ofd_connector/internal/httpclient"
)

// OfdConnector interface defines the contract for OFD service operations
type OfdConnector interface {
	ZReport() ZReportI
	Receipt() ReceiptI
}

type OfdConnectorConfig struct {
	ServiceAddress        string
	RequestTimeOutSeconds int
	FactoryID             string
}

// ofdConnector implements the OfdConnector interface
type ofdConnector struct {
	zReport ZReportI
	receipt ReceiptI
}

// New creates a new instance of OfdConnector
func New(config OfdConnectorConfig) (OfdConnector, error) {
	if config.ServiceAddress == "" {
		return nil, fmt.Errorf("invalid url address")
	}
	if config.FactoryID == "" {
		return nil, fmt.Errorf("invalid FactoryID")
	}

	httpClient := httpclient.NewHTTPClient(config.RequestTimeOutSeconds)
	gateway := gateway.New(gateway.Config{
		ServiceAddress: config.ServiceAddress,
		FactoryID:      config.FactoryID,
		HttpClient:     httpClient,
	})
	zReport := newZReport(zReportConfig{
		gateway: gateway,
	})
	receipt := newReceipt(receiptConfig{
		gateway: gateway,
	})

	return &ofdConnector{
		zReport: zReport,
		receipt: receipt,
	}, nil
}

func (o *ofdConnector) ZReport() ZReportI {
	return o.ZReport()
}

func (o *ofdConnector) Receipt() ReceiptI {
	return o.Receipt()
}
