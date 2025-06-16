package ofdconnector

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/billz-2/ofd_connector/internal/constants"
	"github.com/billz-2/ofd_connector/internal/gateway"
)

const (
	fiscalDriveInfoEndpoint = "/FiscalDrive/Info/"
)

type FiscalDriveI interface {
	FiscalDriveInfo(ctx context.Context) (FiscalDriveInfo, error)
}

type fiscalDriveConfig struct {
	gateway gateway.Client
}

// ofdConnector implements the OfdConnector interface
type fiscalDrive struct {
	gateway gateway.Client
}

func newFiscalDrive(config fiscalDriveConfig) FiscalDriveI {
	return &fiscalDrive{
		gateway: config.gateway,
	}
}

type MemoryInfo struct {
	AvailablePersistentMemory int64 `json:"AvailablePersistentMemory"`
	AvailableResetMemory      int64 `json:"AvailableResetMemory"`
	AvailableDeselectMemory   int64 `json:"AvailableDeselectMemory"`
}

type FiscalDriveInfo struct {
	AppletVersion string     `json:"AppletVersion"`
	TerminalID    string     `json:"TerminalID"`
	SyncChallenge string     `json:"SyncChallenge"`
	Locked        bool       `json:"Locked"`
	JCREVersion   string     `json:"JCREVersion"`
	POSLocked     bool       `json:"POSLocked"`
	POSAuth       bool       `json:"POSAuth"`
	MemoryInfo    MemoryInfo `json:"MemoryInfo"`
}

func (f *fiscalDrive) FiscalDriveInfo(ctx context.Context) (FiscalDriveInfo, error) {
	endpoint := f.gateway.FactoryEndpoint(fiscalDriveInfoEndpoint)
	resp, err := f.gateway.HTTPRequest(
		ctx,
		endpoint,
		http.MethodPost,
		constants.ContentTypeJSON,
		nil,
		nil,
	)
	if err != nil {
		return FiscalDriveInfo{}, fmt.Errorf("error creating request: %s", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		errorResp := errorResponse{}
		if err := json.Unmarshal(resp.Body, &errorResp); err != nil {
			return FiscalDriveInfo{}, fmt.Errorf("error unmarshalling error response: %s responseBody: %s",
				err.Error(),
				string(resp.Body),
			)
		}

		return FiscalDriveInfo{}, fmt.Errorf("failed to get fiscal drive info: %s", errorResp.Reason)
	}

	fiscalDriveInfo := FiscalDriveInfo{}
	if err := json.Unmarshal(resp.Body, &fiscalDriveInfo); err != nil {
		return FiscalDriveInfo{}, fmt.Errorf("error unmarshalling response: %s", err.Error())
	}

	return fiscalDriveInfo, nil
}
