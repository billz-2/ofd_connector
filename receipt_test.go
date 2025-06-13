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

func TestGetTXID(t *testing.T) {
	const (
		factoryID      = "12342131231223123123"
		serviceAddress = "localhost:1234"
	)

	tests := []struct {
		name           string
		saleParams     SaleParams
		responseBody   any
		responseStatus int
		expectError    bool
		errorContains  string
		expectedTxID   int64
	}{
		{
			name: "success",
			saleParams: SaleParams{
				ExtraInfo: ExtraInfo{
					CarNumber:         "АА000АА",
					CardType:          1,
					CashedOutFromCard: 10542,
					PINFL:             "123456789012",
					PPTID:             "83",
					PhoneNumber:       "9998905741148",
					QRPaymentID:       "12345678901234567890",
					QRPaymentProvider: 30,
					TIN:               "1234567890",
				},
				Items: []Item{
					{
						Name:        "Product 1",
						Barcode:     "123456789",
						Labels:      []string{"label1", "label2"},
						SPIC:        "12345",
						Units:       1,
						PackageCode: "PKG001",
						OwnerType:   1,
						Price:       1000,
						VATPercent:  12,
						VAT:         120,
						Amount:      1000,
						Discount:    0,
						Other:       0,
						CommissionInfo: &CommissionInfo{
							TIN:   "1234567890",
							PINFL: "123456789012",
						},
					},
				},
				Location: Location{
					Latitude:  12.345,
					Longitude: 67.890,
				},
				Operation:    0,
				ReceivedCard: 10000,
				ReceivedCash: 0,
				RefundInfo: RefundInfo{
					TerminalID: "1234567890",
					ReceiptSeq: 12,
					DateTime:   "2023-01-01 00:00:00",
					FiscalSign: "00000000000",
				},
				Time: "2023-01-01 00:00:00",
				Type: 0,
			},
			responseBody:   int64(2),
			responseStatus: http.StatusOK,
			expectError:    false,
			expectedTxID:   2,
		},
		{
			name:       "failure",
			saleParams: SaleParams{},
			responseBody: errorResponse{
				Reason: "no card found",
				Type:   "errors.errorString",
			},
			responseStatus: http.StatusNotFound,
			expectError:    true,
			errorContains:  "no card found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			httpClient := mock_httpclient.NewMockHTTPClient(ctrl)

			saleInfoBody, err := json.Marshal(tt.saleParams)
			require.NoError(t, err)

			req, err := httpclient.NewHTTPRequest(
				serviceAddress+"/FiscalDrive/Receipt/GetTXID/"+factoryID,
				http.MethodPost,
				constants.ContentTypeJSON,
				saleInfoBody,
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
			receipt := &receipt{
				gateway: gateway,
			}

			gotTxID, err := receipt.GetTXID(ctx, tt.saleParams)
			if tt.expectError {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.errorContains)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expectedTxID, gotTxID)
		})
	}
}
func TestRegisterTXID(t *testing.T) {
	const (
		factoryID      = "12342131231223123123"
		serviceAddress = "localhost:1234"
		txID           = int64(2)
	)

	tests := []struct {
		name           string
		responseBody   any
		responseStatus int
		expectError    bool
		errorContains  string
		expectedInfo   ReceiptInfo
	}{
		{
			name: "success",
			responseBody: ReceiptInfo{
				DateTime:   "2023-01-01 00:00:00",
				TerminalID: "1234567890",
				ReceiptSeq: 12,
				FiscalSign: "00000000000",
				QRCodeURL:  "ADDRESS.com/qr-code.png",
			},
			responseStatus: http.StatusOK,
			expectError:    false,
			expectedInfo: ReceiptInfo{
				DateTime:   "2023-01-01 00:00:00",
				TerminalID: "1234567890",
				ReceiptSeq: 12,
				FiscalSign: "00000000000",
				QRCodeURL:  "ADDRESS.com/qr-code.png",
			},
		},
		{
			name: "failure",
			responseBody: errorResponse{
				Reason: "no card found",
				Type:   "errors.errorString",
			},
			responseStatus: http.StatusNotFound,
			expectError:    true,
			errorContains:  "no card found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			httpClient := mock_httpclient.NewMockHTTPClient(ctrl)

			regTxIDReq := txIDReq{TXID: txID}
			regTxIDReqBody, err := json.Marshal(regTxIDReq)
			require.NoError(t, err)

			req, err := httpclient.NewHTTPRequest(
				serviceAddress+"/FiscalDrive/Receipt/RegisterTXID/"+factoryID,
				http.MethodPost,
				constants.ContentTypeUrlEncoded,
				regTxIDReqBody,
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
			receipt := &receipt{
				gateway: gateway,
			}

			receiptInfo, err := receipt.RegisterTXID(ctx, txID)
			if tt.expectError {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.errorContains)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedInfo.DateTime, receiptInfo.DateTime)
			assert.Equal(t, tt.expectedInfo.TerminalID, receiptInfo.TerminalID)
			assert.Equal(t, tt.expectedInfo.ReceiptSeq, receiptInfo.ReceiptSeq)
			assert.Equal(t, tt.expectedInfo.FiscalSign, receiptInfo.FiscalSign)
			assert.Equal(t, tt.expectedInfo.QRCodeURL, receiptInfo.QRCodeURL)
		})
	}
}
func TestGetReceiptInfo(t *testing.T) {
	const (
		factoryID      = "12342131231223123123"
		serviceAddress = "localhost:1234"
		index          = uint32(0)
	)

	tests := []struct {
		name           string
		responseBody   any
		responseStatus int
		expectError    bool
		errorContains  string
		expectedInfo   ReceiptFullInfo
	}{
		{
			name: "success",
			responseBody: ReceiptFullInfo{
				Time:          "2023-01-01 00:00:00",
				TerminalID:    "1234567890",
				ReceiptSeq:    12,
				FiscalSign:    "00000000000",
				ItemsCount:    12,
				ItemsHash:     "213123r324123213",
				TotalVAT:      1221,
				ReceivedCash:  12,
				ReceivedCard:  34123,
				ReceiptType:   "Purchase",
				OperationType: "Cash",
			},
			responseStatus: http.StatusOK,
			expectError:    false,
			expectedInfo: ReceiptFullInfo{
				Time:          "2023-01-01 00:00:00",
				TerminalID:    "1234567890",
				ReceiptSeq:    12,
				FiscalSign:    "00000000000",
				ItemsCount:    12,
				ItemsHash:     "213123r324123213",
				TotalVAT:      1221,
				ReceivedCash:  12,
				ReceivedCard:  34123,
				ReceiptType:   "Purchase",
				OperationType: "Cash",
			},
		},
		{
			name: "failure",
			responseBody: errorResponse{
				Reason: "no card found",
				Type:   "errors.errorString",
			},
			responseStatus: http.StatusNotFound,
			expectError:    true,
			errorContains:  "no card found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			httpClient := mock_httpclient.NewMockHTTPClient(ctrl)

			indexData := indexInfo{Index: index}
			indexBody, err := json.Marshal(indexData)
			require.NoError(t, err)

			req, err := httpclient.NewHTTPRequest(
				serviceAddress+"/FiscalDrive/Receipt/Info/"+factoryID,
				http.MethodPost,
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
				HttpClient:     httpClient,
				ServiceAddress: serviceAddress,
				FactoryID:      factoryID,
			})
			receipt := &receipt{
				gateway: gateway,
			}

			receiptInfo, err := receipt.GetReceiptInfo(ctx, index)
			if tt.expectError {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.errorContains)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedInfo.TerminalID, receiptInfo.TerminalID)
			assert.Equal(t, tt.expectedInfo.Time, receiptInfo.Time)
			assert.Equal(t, tt.expectedInfo.ReceiptSeq, receiptInfo.ReceiptSeq)
			assert.Equal(t, tt.expectedInfo.FiscalSign, receiptInfo.FiscalSign)
			assert.Equal(t, tt.expectedInfo.ItemsCount, receiptInfo.ItemsCount)
			assert.Equal(t, tt.expectedInfo.ItemsHash, receiptInfo.ItemsHash)
			assert.Equal(t, tt.expectedInfo.TotalVAT, receiptInfo.TotalVAT)
			assert.Equal(t, tt.expectedInfo.ReceivedCash, receiptInfo.ReceivedCash)
			assert.Equal(t, tt.expectedInfo.ReceivedCard, receiptInfo.ReceivedCard)
			assert.Equal(t, tt.expectedInfo.ReceiptType, receiptInfo.ReceiptType)
			assert.Equal(t, tt.expectedInfo.OperationType, receiptInfo.OperationType)
		})
	}
}

func TestGetDatabaseFilesCount(t *testing.T) {
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
		expectedCount  map[string]int64
	}{
		{
			name:           "success",
			responseBody:   map[string]int64{factoryID: 12},
			responseStatus: http.StatusOK,
			expectError:    false,
			expectedCount:  map[string]int64{factoryID: 12},
		},
		{
			name: "failure",
			responseBody: errorResponse{
				Reason: "no card found",
				Type:   "errors.errorString",
			},
			responseStatus: http.StatusNotFound,
			expectError:    true,
			errorContains:  "no card found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			httpClient := mock_httpclient.NewMockHTTPClient(ctrl)

			statusFilter := statusData{Status: 0}
			statusBody, err := json.Marshal(statusFilter)
			require.NoError(t, err)

			req, err := httpclient.NewHTTPRequest(
				serviceAddress+"/Database/Files/Count",
				http.MethodPost,
				constants.ContentTypeUrlEncoded,
				statusBody,
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
			receipt := &receipt{
				gateway: gateway,
			}

			countRes, err := receipt.GetDatabaseFilesCount(ctx, statusFilter.Status)
			if tt.expectError {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.errorContains)
			} else {
				require.NoError(t, err)
				require.Len(t, countRes, len(tt.expectedCount))
				assert.Equal(t, tt.expectedCount[factoryID], countRes[factoryID])
			}
		})
	}
}
func TestResetDatabaseFilesStatus(t *testing.T) {
	const (
		factoryID      = "12342131231223123123"
		txID           = int64(2)
		serviceAddress = "localhost:1234"
	)

	tests := []struct {
		name           string
		responseBody   any
		responseStatus int
		expectError    bool
		errorContains  string
	}{
		{
			name:           "success",
			responseBody:   nil,
			responseStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name: "failure",
			responseBody: errorResponse{
				Reason: "internal error",
				Type:   "errors.errorString",
			},
			responseStatus: http.StatusInternalServerError,
			expectError:    true,
			errorContains:  "internal error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			httpClient := mock_httpclient.NewMockHTTPClient(ctrl)

			regTxIDReq := txIDReq{TXID: txID}
			regTxIDReqBody, err := json.Marshal(regTxIDReq)
			require.NoError(t, err)

			req, err := httpclient.NewHTTPRequest(
				serviceAddress+databaseFilesStatusReset,
				http.MethodPost,
				constants.ContentTypeUrlEncoded,
				regTxIDReqBody,
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
			receipt := &receipt{
				gateway: gateway,
			}

			err = receipt.ResetDatabaseFilesStatus(ctx, txID)
			if tt.expectError {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.errorContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
