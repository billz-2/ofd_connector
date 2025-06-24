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

func TestZreportOpen(t *testing.T) {
	tests := []struct {
		name          string
		createdAtTime string
		responseBody  interface{}
		responseCode  int
		expectedError string
	}{
		{
			name:          "success",
			createdAtTime: "2025-05-31 12:04:00",
			responseBody:  "OK",
			responseCode:  200,
			expectedError: "",
		},
		{
			name:          "invalid time format",
			createdAtTime: "2025-05-31T12:04:00",
			responseBody:  nil,
			responseCode:  0,
			expectedError: "invalid time format",
		},
		{
			name:          "external error",
			createdAtTime: "2025-05-31 12:04:00",
			responseBody: errorResponse{
				Reason: "no card found",
				Type:   "errors.errorString",
			},
			responseCode:  400,
			expectedError: "no card found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			httpClient := mock_httpclient.NewMockHTTPClient(ctrl)
			const factoryID = "12342131231223123123"

			gateway := gateway.New(gateway.Config{
				ServiceAddress: "localhost:1234",
				FactoryID:      factoryID,
				HttpClient:     httpClient,
			})
			zReport := &zReport{
				gateway: gateway,
			}

			if tt.responseBody != nil {
				reqBody, err := json.Marshal(dateTime{DateTime: tt.createdAtTime})
				require.NoError(t, err)
				req, err := httpclient.NewHTTPRequest(
					"localhost:1234/FiscalDrive/ZReport/Open/"+factoryID,
					http.MethodPost,
					constants.ContentTypeUrlEncoded,
					reqBody,
					nil,
				)
				require.NoError(t, err)

				var responseBody []byte
				responseBody, err = json.Marshal(tt.responseBody)
				require.NoError(t, err)

				httpClient.EXPECT().Request(gomock.Any(), req).
					Return(&httpclient.HTTPResponse{
						Body:       responseBody,
						StatusCode: tt.responseCode,
					}).Times(1)
			}

			err := zReport.OpenZreport(ctx, tt.createdAtTime)
			if tt.expectedError != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestZReportClose(t *testing.T) {
	tests := []struct {
		name          string
		createdAtTime string
		responseBody  interface{}
		responseCode  int
		expectedError string
	}{
		{
			name:          "success",
			createdAtTime: "2025-05-31 12:04:00",
			responseBody:  "OK",
			responseCode:  200,
			expectedError: "",
		},
		{
			name:          "invalid time format",
			createdAtTime: "2025-05-31T12:04:00",
			responseBody:  nil,
			responseCode:  0,
			expectedError: "invalid time format",
		},
		{
			name:          "external error",
			createdAtTime: "2025-05-31 12:04:00",
			responseBody: errorResponse{
				Reason: "no card found",
				Type:   "errors.errorString",
			},
			responseCode:  400,
			expectedError: "no card found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			httpClient := mock_httpclient.NewMockHTTPClient(ctrl)
			const factoryID = "12342131231223123123"

			gateway := gateway.New(gateway.Config{
				ServiceAddress: "localhost:1234",
				FactoryID:      factoryID,
				HttpClient:     httpClient,
			})
			zReport := &zReport{
				gateway: gateway,
			}

			if tt.responseBody != nil {
				reqBody, err := json.Marshal(dateTime{DateTime: tt.createdAtTime})
				require.NoError(t, err)
				req, err := httpclient.NewHTTPRequest(
					"localhost:1234/FiscalDrive/ZReport/Close/"+factoryID,
					http.MethodPost,
					constants.ContentTypeUrlEncoded,
					reqBody,
					nil,
				)
				require.NoError(t, err)

				var responseBody []byte
				responseBody, err = json.Marshal(tt.responseBody)
				require.NoError(t, err)

				httpClient.EXPECT().Request(gomock.Any(), req).
					Return(&httpclient.HTTPResponse{
						Body:       responseBody,
						StatusCode: tt.responseCode,
					}).Times(1)
			}

			err := zReport.CloseZreport(ctx, tt.createdAtTime)
			if tt.expectedError != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
func TestZReportInfo(t *testing.T) {
	tests := []struct {
		name             string
		responseBody     interface{}
		responseStatus   int
		expectedError    string
		expectedResponse *ZReportInfo
	}{
		{
			name: "success",
			responseBody: ZReportInfo{
				OpenTime:         "2023-05-31 12:04:00",
				CloseTime:        "2023-05-31 13:04:00",
				TerminalID:       "TERM123",
				TotalSaleCount:   10,
				TotalRefundCount: 2,
				TotalCash: TotalAmount{
					Sale:   1000,
					Refund: 12,
				},
				TotalCard: TotalAmount{
					Sale:   2000,
					Refund: 11,
				},
				TotalVAT: TotalAmount{
					Sale:   100,
					Refund: 12,
				},
				FirstReceiptSeq: 1001,
				LastReceiptSeq:  1012,
			},
			responseStatus: 200,
			expectedError:  "",
			expectedResponse: &ZReportInfo{
				OpenTime:         "2023-05-31 12:04:00",
				CloseTime:        "2023-05-31 13:04:00",
				TerminalID:       "TERM123",
				TotalSaleCount:   10,
				TotalRefundCount: 2,
				TotalCash: TotalAmount{
					Sale:   1000,
					Refund: 12,
				},
				TotalCard: TotalAmount{
					Sale:   2000,
					Refund: 11,
				},
				TotalVAT: TotalAmount{
					Sale:   100,
					Refund: 12,
				},
				FirstReceiptSeq: 1001,
				LastReceiptSeq:  1012,
			},
		},
		{
			name: "external error",
			responseBody: errorResponse{
				Reason: "no card found",
				Type:   "errors.errorString",
			},
			responseStatus:   400,
			expectedError:    "no card found",
			expectedResponse: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			httpClient := mock_httpclient.NewMockHTTPClient(ctrl)

			const factoryID = "12342131231223123123"
			indexData := indexInfo{Index: 0}
			indexBody, err := json.Marshal(indexData)
			require.NoError(t, err)

			req, err := httpclient.NewHTTPRequest(
				"localhost:1234/FiscalDrive/ZReport/Info/"+factoryID,
				http.MethodGet,
				constants.ContentTypeUrlEncoded,
				indexBody,
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
				ServiceAddress: "localhost:1234",
				FactoryID:      factoryID,
				HttpClient:     httpClient,
			})
			zReport := &zReport{
				gateway: gateway,
			}

			info, err := zReport.GetZReportInfo(ctx, 0)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedResponse.OpenTime, info.OpenTime)
				assert.Equal(t, tt.expectedResponse.CloseTime, info.CloseTime)
				assert.Equal(t, tt.expectedResponse.TerminalID, info.TerminalID)
				assert.Equal(t, tt.expectedResponse.TotalSaleCount, info.TotalSaleCount)
				assert.Equal(t, tt.expectedResponse.TotalRefundCount, info.TotalRefundCount)
				assert.Equal(t, tt.expectedResponse.TotalCash, info.TotalCash)
				assert.Equal(t, tt.expectedResponse.TotalCard, info.TotalCard)
				assert.Equal(t, tt.expectedResponse.TotalVAT, info.TotalVAT)
				assert.Equal(t, tt.expectedResponse.FirstReceiptSeq, info.FirstReceiptSeq)
				assert.Equal(t, tt.expectedResponse.LastReceiptSeq, info.LastReceiptSeq)
			}
		})
	}
}

func TestSyncZReport(t *testing.T) {
	const (
		factoryID      = "12342131231223123123"
		serviceAddress = "localhost:1234"
	)
	tests := []struct {
		name           string
		itemsCount     uint16
		responseBody   any
		responseStatus int
		expectError    bool
		errorContains  string
		mockFunc       func(httpClient *mock_httpclient.MockHTTPClient, req *httpclient.HTTPRequest, response any, responseStatus int)
	}{
		{
			name:           "success",
			itemsCount:     uint16(2),
			responseBody:   nil,
			responseStatus: http.StatusOK,
			expectError:    false,
			mockFunc: func(httpClient *mock_httpclient.MockHTTPClient, req *httpclient.HTTPRequest, response any, responseStatus int) {
				responseBody, err := json.Marshal(response)
				require.NoError(t, err)
				httpClient.EXPECT().Request(gomock.Any(), req).
					Return(&httpclient.HTTPResponse{
						Body:       responseBody,
						StatusCode: responseStatus,
					}).Times(1)
			},
		},
		{
			name:          "fail itemsCount validation: high value",
			itemsCount:    uint16(33),
			expectError:   true,
			errorContains: "itemsCount must be in range [1, 32]",
		},
		{
			name:          "fail itemsCount validation: zero value",
			itemsCount:    uint16(0),
			expectError:   true,
			errorContains: "itemsCount must be in range [1, 32]",
		},
		{
			name:       "failure",
			itemsCount: uint16(2),
			responseBody: errorResponse{
				Reason: "internal error",
				Type:   "errors.errorString",
			},
			responseStatus: http.StatusInternalServerError,
			expectError:    true,
			errorContains:  "internal error",
			mockFunc: func(httpClient *mock_httpclient.MockHTTPClient, req *httpclient.HTTPRequest, response any, responseStatus int) {
				responseBody, err := json.Marshal(response)
				require.NoError(t, err)
				httpClient.EXPECT().Request(gomock.Any(), req).
					Return(&httpclient.HTTPResponse{
						Body:       responseBody,
						StatusCode: responseStatus,
					}).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			httpClient := mock_httpclient.NewMockHTTPClient(ctrl)
			reqBody, err := json.Marshal(itemsCountReq{ItemsCount: tt.itemsCount})
			require.NoError(t, err)
			req, err := httpclient.NewHTTPRequest(
				serviceAddress+zReportSyncEndpoint+factoryID,
				http.MethodPost,
				constants.ContentTypeUrlEncoded,
				reqBody,
				nil,
			)
			require.NoError(t, err)

			if tt.mockFunc != nil {
				tt.mockFunc(httpClient, req, tt.responseBody, tt.responseStatus)
			}

			gateway := gateway.New(gateway.Config{
				HttpClient:     httpClient,
				ServiceAddress: serviceAddress,
				FactoryID:      factoryID,
			})
			zReport := &zReport{
				gateway: gateway,
			}
			err = zReport.SyncZReports(ctx, tt.itemsCount)
			if tt.expectError {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.errorContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
