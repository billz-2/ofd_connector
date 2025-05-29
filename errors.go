package ofdconnector

import "fmt"

var (
	ErrorFactoryIDNotSet   = fmt.Errorf("factoryID is not set")
	ErrorInvalidUrlAddress = fmt.Errorf("invalid url address")
)
