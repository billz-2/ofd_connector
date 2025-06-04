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

func TestListFiscalDrivesSuccess(t *testing.T) {
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

	body, err := json.Marshal([]FiscalDriveReaderInfo{
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
	})
	require.NoError(t, err)
	httpClient.EXPECT().Request(gomock.Any(), req).
		Return(&httpclient.HTTPResponse{
			Body:       body,
			StatusCode: 200,
		}).Times(1)

	fdLister := &fiscalDriveLister{
		httpClient:     httpClient,
		serviceAddress: "localhost:1234",
	}

	got, err := fdLister.ListFiscalDrives(ctx)
	require.NoError(t, err)
	require.Len(t, got, 2)
	assert.Equal(t, "12342131231223123123", got[0].FactoryID)
	assert.Equal(t, "reader1", got[0].ReaderName)
	assert.Equal(t, "1a8f800180318065b08503010101030201040105", got[0].ATR)
	assert.Equal(t, "0400", got[0].AppletVersion)
	assert.Equal(t, "25123123123123123126", got[1].FactoryID)
	assert.Equal(t, "reader2", got[1].ReaderName)
	assert.Equal(t, "3b8f800180318065b08503010101030201040105", got[1].ATR)
	assert.Equal(t, "0200", got[1].AppletVersion)
}

func TestListFiscalDrives_Fail(t *testing.T) {
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

	body, err := json.Marshal(errorResponse{
		Reason: "no card connected",
		Type:   "errors.errorString",
	})
	require.NoError(t, err)
	httpClient.EXPECT().Request(gomock.Any(), req).
		Return(&httpclient.HTTPResponse{
			Body:       body,
			StatusCode: 400,
		}).Times(1)

	fdLister := &fiscalDriveLister{
		httpClient:     httpClient,
		serviceAddress: "localhost:1234",
	}

	got, err := fdLister.ListFiscalDrives(ctx)
	require.Error(t, err)
	require.Nil(t, got)
	assert.ErrorContains(t, err, "no card connected")
}
