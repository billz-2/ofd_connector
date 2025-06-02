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

func TestGetTXIDSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	httpClient := mock_httpclient.NewMockHTTPClient(ctrl)
	const (
		factoryID = "12342131231223123123"
		txIDRes   = int64(2)
	)
	saleParams := SaleParams{
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
	}
	saleInfoBody, err := json.Marshal(saleParams)
	require.NoError(t, err)
	req, err := httpclient.NewHTTPRequest(
		"localhost:1234/FiscalDrive/Receipt/GetTXID/"+factoryID,
		http.MethodPost,
		constants.ContentTypeJSON,
		saleInfoBody,
		nil,
	)
	require.NoError(t, err)

	body, err := json.Marshal(txIDRes)
	require.NoError(t, err)
	httpClient.EXPECT().Request(gomock.Any(), req).
		Return(&httpclient.HTTPResponse{
			Body:       body,
			StatusCode: http.StatusOK,
		}).Times(1)
	receipt := &receipt{
		httpClient:     httpClient,
		serviceAddress: "localhost:1234",
		factoryID:      factoryID,
	}

	gotTxID, err := receipt.GetTXID(ctx, saleParams)
	require.NoError(t, err)
	require.Equal(t, txIDRes, gotTxID)
}

func TestGetTXIDFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	httpClient := mock_httpclient.NewMockHTTPClient(ctrl)
	const (
		factoryID = "12342131231223123123"
	)
	saleParams := SaleParams{}
	saleInfoBody, err := json.Marshal(saleParams)
	require.NoError(t, err)
	req, err := httpclient.NewHTTPRequest(
		"localhost:1234/FiscalDrive/Receipt/GetTXID/"+factoryID,
		http.MethodPost,
		constants.ContentTypeJSON,
		saleInfoBody,
		nil,
	)
	require.NoError(t, err)

	bodyResponse := errorResponse{
		Reason: "no card found",
		Type:   "errors.errorString",
	}

	body, err := json.Marshal(bodyResponse)
	require.NoError(t, err)
	httpClient.EXPECT().Request(gomock.Any(), req).
		Return(&httpclient.HTTPResponse{
			Body:       body,
			StatusCode: http.StatusNotFound,
		}).Times(1)
	receipt := &receipt{
		httpClient:     httpClient,
		serviceAddress: "localhost:1234",
		factoryID:      factoryID,
	}

	_, err = receipt.GetTXID(ctx, saleParams)
	require.Error(t, err)
	assert.ErrorContains(t, err, bodyResponse.Reason)
}

func TestRegisterTXIDSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	httpClient := mock_httpclient.NewMockHTTPClient(ctrl)
	const (
		factoryID = "12342131231223123123"
		txID      = int64(2)
	)
	txIDRes := ReceiptInfo{
		DateTime:   "2023-01-01 00:00:00",
		TerminalID: "1234567890",
		ReceiptSeq: 12,
		FiscalSign: "00000000000",
		QRCodeURL:  "ADDRESS.com/qr-code.png",
	}
	regTxIDReq := RegisterTXIDReq{TXID: txID}
	regTxIDReqBody, err := json.Marshal(regTxIDReq)
	require.NoError(t, err)
	req, err := httpclient.NewHTTPRequest(
		"localhost:1234/FiscalDrive/Receipt/RegisterTXID/"+factoryID,
		http.MethodPost,
		constants.ContentTypeUrlEncoded,
		regTxIDReqBody,
		nil,
	)
	require.NoError(t, err)
	body, err := json.Marshal(txIDRes)
	require.NoError(t, err)
	httpClient.EXPECT().Request(gomock.Any(), req).
		Return(&httpclient.HTTPResponse{
			Body:       body,
			StatusCode: http.StatusOK,
		}).Times(1)
	receipt := &receipt{
		httpClient:     httpClient,
		serviceAddress: "localhost:1234",
		factoryID:      factoryID,
	}

	receiptInfo, err := receipt.RegisterTXID(ctx, txID)
	require.NoError(t, err)
	assert.Equal(t, txIDRes.DateTime, receiptInfo.DateTime)
	assert.Equal(t, txIDRes.TerminalID, receiptInfo.TerminalID)
	assert.Equal(t, txIDRes.ReceiptSeq, receiptInfo.ReceiptSeq)
	assert.Equal(t, txIDRes.FiscalSign, receiptInfo.FiscalSign)
	assert.Equal(t, txIDRes.QRCodeURL, receiptInfo.QRCodeURL)

}

func TestRegisterTXIDFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	httpClient := mock_httpclient.NewMockHTTPClient(ctrl)
	const (
		factoryID = "12342131231223123123"
		txID      = int64(2)
	)
	regTxIDReq := RegisterTXIDReq{TXID: txID}
	regTxIDReqBody, err := json.Marshal(regTxIDReq)
	require.NoError(t, err)
	req, err := httpclient.NewHTTPRequest(
		"localhost:1234/FiscalDrive/Receipt/RegisterTXID/"+factoryID,
		http.MethodPost,
		constants.ContentTypeUrlEncoded,
		regTxIDReqBody,
		nil,
	)
	require.NoError(t, err)
	bodyResponse := errorResponse{
		Reason: "no card found",
		Type:   "errors.errorString",
	}
	body, err := json.Marshal(bodyResponse)
	require.NoError(t, err)
	httpClient.EXPECT().Request(gomock.Any(), req).
		Return(&httpclient.HTTPResponse{
			Body:       body,
			StatusCode: http.StatusNotFound,
		}).Times(1)

	receipt := &receipt{
		httpClient:     httpClient,
		serviceAddress: "localhost:1234",
		factoryID:      factoryID,
	}
	_, err = receipt.RegisterTXID(ctx, txID)
	require.Error(t, err)
	assert.ErrorContains(t, err, bodyResponse.Reason)
}

func TestGetReceiptInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	httpClient := mock_httpclient.NewMockHTTPClient(ctrl)
	const (
		factoryID = "12342131231223123123"
		index     = uint32(0)
	)
	receiptInfoRes := ReceiptFullInfo{
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
	}

	indexData := indexInfo{Index: index}
	indexBody, err := json.Marshal(indexData)
	require.NoError(t, err)
	req, err := httpclient.NewHTTPRequest(
		"localhost:1234/FiscalDrive/Receipt/Info/"+factoryID,
		http.MethodPost,
		constants.ContentTypeUrlEncoded,
		indexBody,
		nil,
	)
	require.NoError(t, err)

	body, err := json.Marshal(receiptInfoRes)
	require.NoError(t, err)
	httpClient.EXPECT().Request(gomock.Any(), req).
		Return(&httpclient.HTTPResponse{
			Body:       body,
			StatusCode: http.StatusOK,
		}).Times(1)

	receipt := &receipt{
		httpClient:     httpClient,
		serviceAddress: "localhost:1234",
		factoryID:      factoryID,
	}

	receiptInfo, err := receipt.GetReceiptInfo(ctx, index)
	require.NoError(t, err)
	assert.Equal(t, receiptInfoRes.TerminalID, receiptInfo.TerminalID)
	assert.Equal(t, receiptInfoRes.Time, receiptInfo.Time)
	assert.Equal(t, receiptInfoRes.ReceiptSeq, receiptInfo.ReceiptSeq)
	assert.Equal(t, receiptInfoRes.FiscalSign, receiptInfo.FiscalSign)
	assert.Equal(t, receiptInfoRes.ItemsCount, receiptInfo.ItemsCount)
	assert.Equal(t, receiptInfoRes.ItemsHash, receiptInfo.ItemsHash)
	assert.Equal(t, receiptInfoRes.TotalVAT, receiptInfo.TotalVAT)
	assert.Equal(t, receiptInfoRes.ReceivedCash, receiptInfo.ReceivedCash)
	assert.Equal(t, receiptInfoRes.ReceivedCard, receiptInfo.ReceivedCard)
	assert.Equal(t, receiptInfoRes.ReceiptType, receiptInfo.ReceiptType)
	assert.Equal(t, receiptInfoRes.OperationType, receiptInfo.OperationType)
}

func TestGetReceiptInfoFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	httpClient := mock_httpclient.NewMockHTTPClient(ctrl)
	const (
		factoryID = "12342131231223123123"
		index     = uint32(0)
	)
	indexData := indexInfo{Index: index}
	indexBody, err := json.Marshal(indexData)
	require.NoError(t, err)
	req, err := httpclient.NewHTTPRequest(
		"localhost:1234/FiscalDrive/Receipt/Info/"+factoryID,
		http.MethodPost,
		constants.ContentTypeUrlEncoded,
		indexBody,
		nil,
	)
	require.NoError(t, err)
	bodyResponse := errorResponse{
		Reason: "no card found",
		Type:   "errors.errorString",
	}
	body, err := json.Marshal(bodyResponse)
	require.NoError(t, err)
	httpClient.EXPECT().Request(gomock.Any(), req).
		Return(&httpclient.HTTPResponse{
			Body:       body,
			StatusCode: http.StatusNotFound,
		}).Times(1)
	receipt := &receipt{
		httpClient:     httpClient,
		serviceAddress: "localhost:1234",
		factoryID:      factoryID,
	}

	_, err = receipt.GetReceiptInfo(ctx, index)
	require.Error(t, err)
	assert.ErrorContains(t, err, bodyResponse.Reason)

}
