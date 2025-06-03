package ofdconnector

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/billz-2/ofd_connector/internal/constants"
	"github.com/billz-2/ofd_connector/internal/gateway"
	"github.com/billz-2/ofd_connector/internal/validators"
)

const (
	zReportCloseEndpoint = "/FiscalDrive/ZReport/Close/"
	zReportInfoEndpoint  = "/FiscalDrive/ZReport/Info/"
	zReportOpenEndpoint  = "/FiscalDrive/ZReport/Open/"
)

type ZReportI interface {
	OpenZreport(ctx context.Context, createdTime string) error
	CloseZreport(ctx context.Context, closedTime string) error
	GetZReportInfo(ctx context.Context, index uint32) (ZReportInfo, error)
}

type zReportConfigs struct {
	gateway gateway.GatewayI
}

// ofdConnector implements the OfdConnector interface
type zReport struct {
	gateway gateway.GatewayI
}

func newZReport(configs zReportConfigs) ZReportI {
	return &zReport{
		gateway: configs.gateway,
	}
}

type indexInfo struct {
	Index uint32 `json:"Index"`
}

type ZReportInfo struct {
	TerminalID       string      `json:"TerminalID"`
	OpenTime         string      `json:"OpenTime"`
	CloseTime        string      `json:"CloseTime"`
	TotalSaleCount   int         `json:"TotalSaleCount"`
	TotalRefundCount int         `json:"TotalRefundCount"`
	TotalCash        TotalAmount `json:"TotalCash"`
	TotalCard        TotalAmount `json:"TotalCard"`
	TotalVAT         TotalAmount `json:"TotalVAT"`
	FirstReceiptSeq  int         `json:"FirstReceiptSeq"`
	LastReceiptSeq   int         `json:"LastReceiptSeq"`
}

type dateTime struct {
	DateTime string `json:"DateTime"`
}

// OpenZreport opens a new Zreport for the fiscal drive
// createdTime is in format "2006-01-02 15:04:05" or "now"
func (o zReport) OpenZreport(ctx context.Context, createdTime string) error {
	if createdTime != "now" {
		validateErr := validators.ValidateTimeFormat(createdTime)
		if validateErr != nil {
			return fmt.Errorf("invalid time format, can be 'now' or in format %s, err: %w", constants.TimeFormat, validateErr)
		}
	}

	bodyBytes, err := json.Marshal(dateTime{DateTime: createdTime})
	if err != nil {
		return fmt.Errorf("error marshalling body: %s", err.Error())
	}

	endpoint := o.gateway.FactoryEndpoint(zReportOpenEndpoint)
	resp, err := o.gateway.HTTPRequest(
		ctx,
		endpoint,
		http.MethodPost,
		constants.ContentTypeUrlEncoded,
		bodyBytes,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error creating request: %s", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		errorResp := errorResponse{}
		if err := json.Unmarshal(resp.Body, &errorResp); err != nil {
			return fmt.Errorf("error unmarshalling error response: %s responseBody: %s",
				err.Error(),
				string(resp.Body),
			)
		}
		return fmt.Errorf("failed to open Z report: %s", errorResp.Reason)
	}

	return nil
}

// CloseZreport closes the Zreport for the fiscal drive
// closedTime is in format "2006-01-02 15:04:05" or "now"
func (o zReport) CloseZreport(ctx context.Context, closedTime string) error {
	if closedTime != "now" {
		validateErr := validators.ValidateTimeFormat(closedTime)
		if validateErr != nil {
			return fmt.Errorf("invalid time format, can be 'now' or in format %s, err: %w", constants.TimeFormat, validateErr)
		}
	}

	bodyBytes, err := json.Marshal(dateTime{DateTime: closedTime})
	if err != nil {
		return fmt.Errorf("error marshalling body: %s", err.Error())
	}

	endpoint := o.gateway.FactoryEndpoint(zReportCloseEndpoint)
	resp, err := o.gateway.HTTPRequest(
		ctx,
		endpoint,
		http.MethodPost,
		constants.ContentTypeUrlEncoded,
		bodyBytes,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error creating request: %s", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		errorResp := errorResponse{}
		if err := json.Unmarshal(resp.Body, &errorResp); err != nil {
			return fmt.Errorf("error unmarshalling error response: %s responseBody: %s",
				err.Error(),
				string(resp.Body),
			)
		}
		return fmt.Errorf("failed to close Z report: %s", errorResp.Reason)
	}

	return nil
}

// GetZReportInfo returns the Zreport info for the fiscal drive
// index 0-current zReport, 1-previous zReport, 2-before previous zReport, etc.
func (o zReport) GetZReportInfo(ctx context.Context, index uint32) (ZReportInfo, error) {
	bodyBytes, err := json.Marshal(indexInfo{Index: index})
	if err != nil {
		return ZReportInfo{}, fmt.Errorf("error marshalling body: %s", err.Error())
	}

	endpoint := o.gateway.FactoryEndpoint(receiptInfoEndpoint)
	resp, err := o.gateway.HTTPRequest(
		ctx,
		endpoint,
		http.MethodGet,
		constants.ContentTypeUrlEncoded,
		bodyBytes,
		nil,
	)
	if err != nil {
		return ZReportInfo{}, fmt.Errorf("error creating request: %s", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		errorResp := errorResponse{}
		if err := json.Unmarshal(resp.Body, &errorResp); err != nil {
			return ZReportInfo{}, fmt.Errorf("error unmarshalling error response: %s responseBody: %s",
				err.Error(),
				string(resp.Body),
			)
		}

		return ZReportInfo{}, fmt.Errorf("failed to get Z report info: %s", errorResp.Reason)
	}

	zReportInfo := ZReportInfo{}
	if err := json.Unmarshal(resp.Body, &zReportInfo); err != nil {
		return ZReportInfo{}, fmt.Errorf("error unmarshalling response: %s", err.Error())
	}
	return zReportInfo, nil
}
