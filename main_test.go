package ofdconnector

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/test-go/testify/require"
)

type HTTPResponse struct {
	Body       []byte
	StatusCode int
	Error      error
}

var (
	ctx context.Context
)

func TestMain(m *testing.M) {
	ctx = context.Background()
	m.Run()
}

func TestNew(t *testing.T) {
	ofd, err := New(OfdConnectorConfigs{
		ServiceAddress:        "localhost:1232",
		RequestTimeOutSeconds: 10,
	})
	require.NoError(t, err)
	assert.NotNil(t, ofd)
}

func TestNewInvalidAddress(t *testing.T) {
	ofd, err := New(OfdConnectorConfigs{
		ServiceAddress:        "",
		RequestTimeOutSeconds: 10,
	})
	require.Error(t, err)
	require.Nil(t, ofd)

	assert.True(t, strings.Contains(err.Error(), "invalid url address"))
}

func TestSetFactoryID(t *testing.T) {
	ofd := ofdConnector{"localhost:1232", nil, ""}
	expectedID := "test123"

	ofd.SetFactoryID(expectedID)
	assert.Equal(t, expectedID, ofd.factoryID)
}
