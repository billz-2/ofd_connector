package ofdconnector

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.udevs.io/billz/ofd_connector/pkg/httpclient"
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

// FiscalDriveList returns list of fiscal drive readers
func (o *ofdConnector) FiscalDriveList(ctx context.Context) ([]FiscalDriveReaderInfo, error) {
	// Implementation for /FiscalDrive/List endpoint
	req, err := httpclient.NewHTTPRequest(
		o.serviceAddress+"/FiscalDrive/List",
		http.MethodPost,
		ContentTypeJSON,
		nil,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %s", err.Error()) // TODO: Add error code
	}

	resp := o.httpClient.Request(ctx, req)
	if resp.StatusCode != 200 {
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
