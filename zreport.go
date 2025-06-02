package ofdconnector

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/billz-2/ofd_connector/internal/constants"
	"github.com/billz-2/ofd_connector/internal/httpclient"
	"github.com/billz-2/ofd_connector/internal/validators"
)

type ZReportI interface {
	OpenZreport(ctx context.Context, createdTime string) error
	CloseZreport(ctx context.Context, closedTime string) error
	GetZReportInfo(ctx context.Context, index uint32) (ZReportInfo, error)
}

type zReportConfigs struct {
	ServiceAddress string
	FactoryID      string
	HttpClient     httpclient.HTTPClient
}

// ofdConnector implements the OfdConnector interface
type zReport struct {
	serviceAddress string
	httpClient     httpclient.HTTPClient
	factoryID      string
}

func newZReport(configs zReportConfigs) ZReportI {
	return &zReport{
		serviceAddress: configs.ServiceAddress,
		httpClient:     configs.HttpClient,
		factoryID:      configs.FactoryID,
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
	if createdTime != "now" && !validators.IsValidateTimeFormat(createdTime) {
		return fmt.Errorf("invalid time format, can be 'now' or in format %s", constants.TimeFormat)
	}

	bodyBytes, err := json.Marshal(dateTime{DateTime: createdTime})
	if err != nil {
		return fmt.Errorf("error marshalling body: %s", err.Error())
	}

	endpoint := fmt.Sprintf("%s/FiscalDrive/ZReport/Open/%s", o.serviceAddress, o.factoryID)
	req, err := httpclient.NewHTTPRequest(
		endpoint,
		http.MethodPost,
		constants.ContentTypeUrlEncoded,
		bodyBytes,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error creating request: %s", err.Error())
	}

	resp := o.httpClient.Request(ctx, req)
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
	if closedTime != "now" && !validators.IsValidateTimeFormat(closedTime) {
		return fmt.Errorf("invalid time format, can be 'now' or in format %s", constants.TimeFormat)
	}

	bodyBytes, err := json.Marshal(dateTime{DateTime: closedTime})
	if err != nil {
		return fmt.Errorf("error marshalling body: %s", err.Error())
	}

	endpoint := fmt.Sprintf("%s/FiscalDrive/ZReport/Close/%s", o.serviceAddress, o.factoryID)
	req, err := httpclient.NewHTTPRequest(
		endpoint,
		http.MethodPost,
		constants.ContentTypeUrlEncoded,
		bodyBytes,
		nil,
	)
	if err != nil {
		return fmt.Errorf("error creating request: %s", err.Error())
	}

	resp := o.httpClient.Request(ctx, req)
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
	if index < 0 {
		return ZReportInfo{}, fmt.Errorf("index cannot be negative")
	}

	bodyBytes, err := json.Marshal(indexInfo{Index: index})
	if err != nil {
		return ZReportInfo{}, fmt.Errorf("error marshalling body: %s", err.Error())
	}

	endpoint := fmt.Sprintf("%s/FiscalDrive/ZReport/Info/%s", o.serviceAddress, o.factoryID)
	req, err := httpclient.NewHTTPRequest(
		endpoint,
		http.MethodGet,
		constants.ContentTypeUrlEncoded,
		bodyBytes,
		nil,
	)
	if err != nil {
		return ZReportInfo{}, fmt.Errorf("error creating request: %s", err.Error())
	}

	resp := o.httpClient.Request(ctx, req)
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
