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
	FiscalDrive() FiscalDriveI
}

type OfdConnectorConfig struct {
	ServiceAddress        string
	RequestTimeOutSeconds int
	FactoryID             string
}

// ofdConnector implements the OfdConnector interface
type ofdConnector struct {
	zReport     ZReportI
	receipt     ReceiptI
	fiscalDrive FiscalDriveI
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
	fiscalDrive := newFiscalDrive(fiscalDriveConfig{
		gateway: gateway,
	})

	return &ofdConnector{
		zReport:     zReport,
		receipt:     receipt,
		fiscalDrive: fiscalDrive,
	}, nil
}

func (o *ofdConnector) ZReport() ZReportI {
	return o.zReport
}

func (o *ofdConnector) Receipt() ReceiptI {
	return o.receipt
}

func (o *ofdConnector) FiscalDrive() FiscalDriveI {
	return o.fiscalDrive
}
