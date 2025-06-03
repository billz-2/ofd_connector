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

type OfdConnectorConfigs struct {
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
func New(configs OfdConnectorConfigs) (OfdConnector, error) {
	if configs.ServiceAddress == "" {
		return nil, fmt.Errorf("invalid url address")
	}
	if configs.FactoryID == "" {
		return nil, fmt.Errorf("invalid FactoryID")
	}

	httpClient := httpclient.NewHTTPClient(configs.RequestTimeOutSeconds)
	gateway := gateway.New(gateway.Configs{
		ServiceAddress: configs.ServiceAddress,
		FactoryID:      configs.FactoryID,
		HttpClient:     httpClient,
	})
	zReport := newZReport(zReportConfigs{
		gateway: gateway,
	})
	receipt := newReceipt(receiptConfigs{
		Gateway: gateway,
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
