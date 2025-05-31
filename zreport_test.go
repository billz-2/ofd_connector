package ofdconnector

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/billz-2/ofd_connector/pkg/httpclient"
	mock_httpclient "github.com/billz-2/ofd_connector/pkg/httpclient/mock"
	"github.com/stretchr/testify/assert"
	"github.com/test-go/testify/require"
	"go.uber.org/mock/gomock"
)

func TestZreportOpenSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	httpClient := mock_httpclient.NewMockHTTPClient(ctrl)

	const factoryID = "12342131231223123123"
	createdAtTime := "2025-05-31 12:04:00"
	reqBody, err := json.Marshal(dateTime{DateTime: createdAtTime})
	require.NoError(t, err)
	req, err := httpclient.NewHTTPRequest(
		"localhost:1234/FiscalDrive/ZReport/Open/"+factoryID,
		http.MethodPost,
		ContentTypeUrlEncoded,
		reqBody,
		nil,
	)
	require.NoError(t, err)

	httpClient.EXPECT().Request(gomock.Any(), req).
		Return(&httpclient.HTTPResponse{
			Body:       []byte("OK"),
			StatusCode: 200,
		}).Times(1)

	ofd := &ofdConnector{
		httpClient:     httpClient,
		serviceAddress: "localhost:1234",
		factoryID:      factoryID,
	}

	err = ofd.OpenZreport(ctx, createdAtTime)
	require.NoError(t, err)
}

func TestZreportOpenFailInvalidTime(t *testing.T) {
	ctrl := gomock.NewController(t)
	httpClient := mock_httpclient.NewMockHTTPClient(ctrl)
	const factoryID = "12342131231223123123"
	createdAtTime := "2025-05-31T12:04:00"

	ofd := &ofdConnector{
		httpClient:     httpClient,
		serviceAddress: "localhost:1234",
		factoryID:      factoryID,
	}

	err := ofd.OpenZreport(ctx, createdAtTime)
	require.Error(t, err)
	assert.ErrorContains(t, err, "invalid time format")
}

func TestZreportOpenFailExternal(t *testing.T) {
	ctrl := gomock.NewController(t)
	httpClient := mock_httpclient.NewMockHTTPClient(ctrl)
	const factoryID = "12342131231223123123"
	createdAtTime := "2025-05-31 12:04:00"
	reqBody, err := json.Marshal(dateTime{DateTime: createdAtTime})
	require.NoError(t, err)
	req, err := httpclient.NewHTTPRequest(
		"localhost:1234/FiscalDrive/ZReport/Open/"+factoryID,
		http.MethodPost,
		ContentTypeUrlEncoded,
		reqBody,
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
			StatusCode: 400,
		}).Times(1)
	ofd := &ofdConnector{
		httpClient:     httpClient,
		serviceAddress: "localhost:1234",
		factoryID:      factoryID,
	}

	err = ofd.OpenZreport(ctx, createdAtTime)
	require.Error(t, err)
	assert.ErrorContains(t, err, bodyResponse.Reason)

}

func TestZreportCloseSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	httpClient := mock_httpclient.NewMockHTTPClient(ctrl)

	const factoryID = "12342131231223123123"
	createdAtTime := "2025-05-31 12:04:00"
	reqBody, err := json.Marshal(dateTime{DateTime: createdAtTime})
	require.NoError(t, err)
	req, err := httpclient.NewHTTPRequest(
		"localhost:1234/FiscalDrive/ZReport/Close/"+factoryID,
		http.MethodPost,
		ContentTypeUrlEncoded,
		reqBody,
		nil,
	)
	require.NoError(t, err)

	httpClient.EXPECT().Request(gomock.Any(), req).
		Return(&httpclient.HTTPResponse{
			Body:       []byte("OK"),
			StatusCode: 200,
		}).Times(1)

	ofd := &ofdConnector{
		httpClient:     httpClient,
		serviceAddress: "localhost:1234",
		factoryID:      factoryID,
	}

	err = ofd.CloseZreport(ctx, createdAtTime)
	require.NoError(t, err)
}

func TestZreportCloseFailInvalidTime(t *testing.T) {
	ctrl := gomock.NewController(t)
	httpClient := mock_httpclient.NewMockHTTPClient(ctrl)
	const factoryID = "12342131231223123123"
	createdAtTime := "2025-05-31T12:04:00"
	ofd := &ofdConnector{
		httpClient:     httpClient,
		serviceAddress: "localhost:1234",
		factoryID:      factoryID,
	}

	err := ofd.CloseZreport(ctx, createdAtTime)
	require.Error(t, err)
	assert.ErrorContains(t, err, "invalid time format")
}

func TestZReportCloseFailExternal(t *testing.T) {
	ctrl := gomock.NewController(t)
	httpClient := mock_httpclient.NewMockHTTPClient(ctrl)
	const factoryID = "12342131231223123123"
	createdAtTime := "2025-05-31 12:04:00"
	reqBody, err := json.Marshal(dateTime{DateTime: createdAtTime})
	require.NoError(t, err)
	req, err := httpclient.NewHTTPRequest(
		"localhost:1234/FiscalDrive/ZReport/Close/"+factoryID,
		http.MethodPost,
		ContentTypeUrlEncoded,
		reqBody,
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
			StatusCode: 400,
		}).Times(1)
	ofd := &ofdConnector{
		httpClient:     httpClient,
		serviceAddress: "localhost:1234",
		factoryID:      factoryID,
	}

	err = ofd.CloseZreport(ctx, createdAtTime)
	require.Error(t, err)
	assert.ErrorContains(t, err, bodyResponse.Reason)
}

func TestZReportInfoSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	httpClient := mock_httpclient.NewMockHTTPClient(ctrl)

	const factoryID = "12342131231223123123"

	indexData := indexInfo{Index: 0}
	indexBody, err := json.Marshal(indexData)
	require.NoError(t, err)
	req, err := httpclient.NewHTTPRequest(
		"localhost:1234/FiscalDrive/ZReport/Info/"+factoryID,
		http.MethodGet,
		ContentTypeUrlEncoded,
		indexBody,
		nil,
	)
	require.NoError(t, err)

	expectedResponse := ZReportInfo{
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
	}

	responseBody, err := json.Marshal(expectedResponse)
	require.NoError(t, err)

	httpClient.EXPECT().Request(gomock.Any(), req).
		Return(&httpclient.HTTPResponse{
			Body:       responseBody,
			StatusCode: 200,
		}).Times(1)

	ofd := &ofdConnector{
		httpClient:     httpClient,
		serviceAddress: "localhost:1234",
		factoryID:      factoryID,
	}

	info, err := ofd.ZReportInfo(ctx, 0)
	require.NoError(t, err)
	assert.Equal(t, expectedResponse.OpenTime, info.OpenTime)
	assert.Equal(t, expectedResponse.CloseTime, info.CloseTime)
	assert.Equal(t, expectedResponse.TerminalID, info.TerminalID)
	assert.Equal(t, expectedResponse.TotalSaleCount, info.TotalSaleCount)
	assert.Equal(t, expectedResponse.TotalRefundCount, info.TotalRefundCount)
	assert.Equal(t, expectedResponse.TotalCash, info.TotalCash)
	assert.Equal(t, expectedResponse.TotalCard, info.TotalCard)
	assert.Equal(t, expectedResponse.TotalVAT, info.TotalVAT)
	assert.Equal(t, expectedResponse.FirstReceiptSeq, info.FirstReceiptSeq)
	assert.Equal(t, expectedResponse.LastReceiptSeq, info.LastReceiptSeq)

}

func TestZReportFailExternal(t *testing.T) {
	ctrl := gomock.NewController(t)
	httpClient := mock_httpclient.NewMockHTTPClient(ctrl)
	const factoryID = "12342131231223123123"
	indexData := indexInfo{Index: 0}
	indexBody, err := json.Marshal(indexData)
	require.NoError(t, err)
	req, err := httpclient.NewHTTPRequest(
		"localhost:1234/FiscalDrive/ZReport/Info/"+factoryID,
		http.MethodGet,
		ContentTypeUrlEncoded,
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
			StatusCode: 400,
		}).Times(1)
	ofd := &ofdConnector{
		httpClient:     httpClient,
		serviceAddress: "localhost:1234",
		factoryID:      factoryID,
	}

	_, err = ofd.ZReportInfo(ctx, 0)
	require.Error(t, err)
	assert.ErrorContains(t, err, bodyResponse.Reason)
}
