package ofdconnector

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/billz-2/ofd_connector/internal/constants"
	"github.com/billz-2/ofd_connector/internal/httpclient"
)

const (
	fiscalDriveListEndpoint = "/FiscalDrive/List"
)

type FiscalDriveReaderInfo struct {
	ReaderName    string `json:"ReaderName"`
	ATR           string `json:"ATR"`
	Description   string `json:"Description"`
	FactoryID     string `json:"FactoryID"`
	AppletVersion string `json:"AppletVersion"`
}

type errorResponse struct {
	Reason string `json:"Reason"`
	Type   string `json:"Type"`
}

type FiscalDriveLister interface {
	ListFiscalDrives(context.Context) ([]FiscalDriveReaderInfo, error)
}

// ofdConnector implements the OfdConnector interface
type fiscalDriveLister struct {
	serviceAddress string
	httpClient     httpclient.HTTPClient
}

// NewFiscalDriveLister returns FiscalDriveLister that lists available fiscal drive readers in the system
func NewFiscalDriveLister(config OfdConnectorConfig) (FiscalDriveLister, error) {
	if config.ServiceAddress == "" {
		return nil, fmt.Errorf("invalid url address")
	}

	httpClient := httpclient.NewHTTPClient(
		config.RequestTimeOutSeconds,
	)

	return &fiscalDriveLister{
		serviceAddress: config.ServiceAddress,
		httpClient:     httpClient,
	}, nil
}

// ListFiscalDrives returns list of fiscal drive readers
func (o fiscalDriveLister) ListFiscalDrives(ctx context.Context) ([]FiscalDriveReaderInfo, error) {
	// Implementation for /FiscalDrive/List endpoint
	req, err := httpclient.NewHTTPRequest(
		o.serviceAddress+fiscalDriveListEndpoint,
		http.MethodPost,
		constants.ContentTypeJSON,
		nil,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %s", err.Error()) // TODO: Add error code
	}

	resp := o.httpClient.Request(ctx, req)
	if resp.StatusCode != http.StatusOK {
		errorResp := errorResponse{}
		jsonErr := json.Unmarshal(resp.Body, &errorResp)
		if jsonErr != nil {
			return nil, fmt.Errorf("error unmarshalling body: %s", resp.Body) // TODO: Add error code
		}
		return nil, fmt.Errorf("error: %s", errorResp.Reason)
	}

	var result []FiscalDriveReaderInfo
	err = json.Unmarshal(resp.Body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
