package ofdconnector

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/billz-2/ofd_connector/internal/constants"
	"github.com/billz-2/ofd_connector/internal/httpclient"
	mock_httpclient "github.com/billz-2/ofd_connector/internal/httpclient/mock"
	"github.com/stretchr/testify/assert"
	"github.com/test-go/testify/require"
	"go.uber.org/mock/gomock"
)

func TestListFiscalDrives(t *testing.T) {
	tests := []struct {
		name           string
		responseBody   interface{}
		responseStatus int
		expectedError  string
		expectedDrives []FiscalDriveReaderInfo
	}{
		{
			name: "success",
			responseBody: []FiscalDriveReaderInfo{
				{
					FactoryID:     "12342131231223123123",
					ReaderName:    "reader1",
					ATR:           "1a8f800180318065b08503010101030201040105",
					AppletVersion: "0400",
				},
				{
					FactoryID:     "25123123123123123126",
					ReaderName:    "reader2",
					ATR:           "3b8f800180318065b08503010101030201040105",
					AppletVersion: "0200",
				},
			},
			responseStatus: 200,
			expectedDrives: []FiscalDriveReaderInfo{
				{
					FactoryID:     "12342131231223123123",
					ReaderName:    "reader1",
					ATR:           "1a8f800180318065b08503010101030201040105",
					AppletVersion: "0400",
				},
				{
					FactoryID:     "25123123123123123126",
					ReaderName:    "reader2",
					ATR:           "3b8f800180318065b08503010101030201040105",
					AppletVersion: "0200",
				},
			},
		},
		{
			name: "failure",
			responseBody: errorResponse{
				Reason: "no card connected",
				Type:   "errors.errorString",
			},
			responseStatus: 400,
			expectedError:  "no card connected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			httpClient := mock_httpclient.NewMockHTTPClient(ctrl)

			req, err := httpclient.NewHTTPRequest(
				"localhost:1234/FiscalDrive/List",
				http.MethodPost,
				constants.ContentTypeJSON,
				nil,
				nil,
			)
			require.NoError(t, err)

			body, err := json.Marshal(tt.responseBody)
			require.NoError(t, err)
			httpClient.EXPECT().Request(gomock.Any(), req).
				Return(&httpclient.HTTPResponse{
					Body:       body,
					StatusCode: tt.responseStatus,
				}).Times(1)

			fdLister := &fiscalDriveLister{
				httpClient:     httpClient,
				serviceAddress: "localhost:1234",
			}

			got, err := fdLister.ListFiscalDrives(ctx)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.expectedError)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedDrives, got)
			}
		})
	}
}
