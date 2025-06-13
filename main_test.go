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
	tests := []struct {
		name        string
		config      OfdConnectorConfig
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid config",
			config: OfdConnectorConfig{
				ServiceAddress:        "localhost:1232",
				RequestTimeOutSeconds: 10,
				FactoryID:             "12342131231223123123",
			},
			expectError: false,
		},
		{
			name: "invalid address",
			config: OfdConnectorConfig{
				ServiceAddress:        "",
				RequestTimeOutSeconds: 10,
				FactoryID:             "12342131231223123123",
			},
			expectError: true,
			errorMsg:    "invalid url address",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ofd, err := New(tt.config)

			if tt.expectError {
				require.Error(t, err)
				require.Nil(t, ofd)
				assert.True(t, strings.Contains(err.Error(), tt.errorMsg))
			} else {
				require.NoError(t, err)
				assert.NotNil(t, ofd)
			}
		})
	}
}
