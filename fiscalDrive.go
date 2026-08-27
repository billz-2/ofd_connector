package ofdconnector

import (
	"context"
	"encoding/json"
	"errors"
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
	gateway           gateway.Client
	fiscalDriveLister FiscalDriveLister
}

// ofdConnector implements the OfdConnector interface
type fiscalDrive struct {
	gateway           gateway.Client
	fiscalDriveLister FiscalDriveLister
}

func newFiscalDrive(config fiscalDriveConfig) FiscalDriveI {
	return &fiscalDrive{
		gateway:           config.gateway,
		fiscalDriveLister: config.fiscalDriveLister,
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
	fiscalDrivesConnected, err := f.fiscalDriveLister.ListFiscalDrives(ctx)
	if err != nil {
		return FiscalDriveInfo{}, fmt.Errorf("error listing fiscal drives: %w", err)
	}
	// validate if only one smartCard is connected
	if len(fiscalDrivesConnected) == 0 {
		return FiscalDriveInfo{}, errors.New("no smartCard connected")
	}
	// refresh to the factoryID to the first connected fiscal drive,
	// when multiple fiscal drives are connected, the first one will be used
	if len(fiscalDrivesConnected) >= 1 {
		f.gateway.SetFactoryID(fiscalDrivesConnected[0].FactoryID)
	}

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
		if errorResp.Reason == "" {
			errorResp.Reason = "unknown error body:" + string(resp.Body)
		}
		return FiscalDriveInfo{}, fmt.Errorf("failed to get fiscal drive info: %s", errorResp.Reason)
	}

	fiscalDriveInfo := FiscalDriveInfo{}
	if err := json.Unmarshal(resp.Body, &fiscalDriveInfo); err != nil {
		return FiscalDriveInfo{}, fmt.Errorf("error unmarshalling response: %s", err.Error())
	}

	return fiscalDriveInfo, nil
}
