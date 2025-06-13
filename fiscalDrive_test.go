package ofdconnector

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/billz-2/ofd_connector/internal/constants"
	"github.com/billz-2/ofd_connector/internal/gateway"
	"github.com/billz-2/ofd_connector/internal/httpclient"
	mock_httpclient "github.com/billz-2/ofd_connector/internal/httpclient/mock"
	"github.com/stretchr/testify/assert"
	"github.com/test-go/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetFiscalDrive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	const (
		factoryID      = "12342131231223123123"
		serviceAddress = "localhost:1234"
	)

	tests := []struct {
		name           string
		responseBody   any
		responseStatus int
		expectError    bool
		errorContains  string
		expectedInfo   FiscalDriveInfo
	}{
		{
			name: "success",
			responseBody: FiscalDriveInfo{
				AppletVersion: "2.0.0",
				TerminalID:    "TERM123",
				SyncChallenge: "challenge123",
				Locked:        false,
				JCREVersion:   "3.0.4",
				POSLocked:     false,
				POSAuth:       true,
				MemoryInfo: MemoryInfo{
					AvailablePersistentMemory: 200,
					AvailableResetMemory:      100,
					AvailableDeselectMemory:   23,
				},
			},
			responseStatus: http.StatusOK,
			expectError:    false,
			expectedInfo: FiscalDriveInfo{
				AppletVersion: "2.0.0",
				TerminalID:    "TERM123",
				SyncChallenge: "challenge123",
				Locked:        false,
				JCREVersion:   "3.0.4",
				POSLocked:     false,
				POSAuth:       true,
				MemoryInfo: MemoryInfo{
					AvailablePersistentMemory: 200,
					AvailableResetMemory:      100,
					AvailableDeselectMemory:   23,
				},
			},
		},
		{
			name: "failure",
			responseBody: errorResponse{
				Reason: "fiscal drive not found",
				Type:   "errors.errorString",
			},
			responseStatus: http.StatusNotFound,
			expectError:    true,
			errorContains:  "fiscal drive not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			httpClient := mock_httpclient.NewMockHTTPClient(ctrl)

			req, err := httpclient.NewHTTPRequest(
				serviceAddress+"/FiscalDrive/Info/"+factoryID,
				http.MethodGet,
				constants.ContentTypeJSON,
				nil,
				nil,
			)
			require.NoError(t, err)

			responseBody, err := json.Marshal(tt.responseBody)
			require.NoError(t, err)

			httpClient.EXPECT().Request(gomock.Any(), req).
				Return(&httpclient.HTTPResponse{
					Body:       responseBody,
					StatusCode: tt.responseStatus,
				}).Times(1)

			gateway := gateway.New(gateway.Config{
				HttpClient:     httpClient,
				ServiceAddress: serviceAddress,
				FactoryID:      factoryID,
			})
			fiscalDrive := &fiscalDrive{
				gateway: gateway,
			}

			info, err := fiscalDrive.FiscalDriveInfo(ctx)
			if tt.expectError {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.errorContains)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedInfo.AppletVersion, info.AppletVersion)
			assert.Equal(t, tt.expectedInfo.TerminalID, info.TerminalID)
			assert.Equal(t, tt.expectedInfo.SyncChallenge, info.SyncChallenge)
			assert.Equal(t, tt.expectedInfo.Locked, info.Locked)
			assert.Equal(t, tt.expectedInfo.JCREVersion, info.JCREVersion)
			assert.Equal(t, tt.expectedInfo.POSLocked, info.POSLocked)
			assert.Equal(t, tt.expectedInfo.POSAuth, info.POSAuth)
			assert.Equal(t, tt.expectedInfo.MemoryInfo, info.MemoryInfo)
		})
	}
}
