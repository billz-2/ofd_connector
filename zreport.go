package ofdconnector

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/billz-2/ofd_connector/pkg/httpclient"
)

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
func (o *ofdConnector) OpenZreport(ctx context.Context, createdTime string) error {
	if o.factoryID == "" {
		return fmt.Errorf("factoryID cannot be empty")
	}

	if _, err := time.Parse(TimeFormat, createdTime); err != nil && createdTime != "now" {
		return fmt.Errorf("invalid time format, can be 'now' or in format %s", TimeFormat)
	}

	bodyBytes, err := json.Marshal(dateTime{DateTime: createdTime})
	if err != nil {
		return fmt.Errorf("error marshalling body: %s", err.Error())
	}

	endpoint := fmt.Sprintf("%s/FiscalDrive/ZReport/Open/%s", o.serviceAddress, o.factoryID)
	req, err := httpclient.NewHTTPRequest(
		endpoint,
		http.MethodPost,
		ContentTypeUrlEncoded,
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
func (o *ofdConnector) CloseZreport(ctx context.Context, closedTime string) error {
	if o.factoryID == "" {
		return fmt.Errorf("factoryID cannot be empty")
	}

	if _, err := time.Parse(TimeFormat, closedTime); err != nil && closedTime != "now" {
		return fmt.Errorf("invalid time format, can be 'now' or in format %s", TimeFormat)
	}

	bodyBytes, err := json.Marshal(dateTime{DateTime: closedTime})
	if err != nil {
		return fmt.Errorf("error marshalling body: %s", err.Error())
	}

	endpoint := fmt.Sprintf("%s/FiscalDrive/ZReport/Close/%s", o.serviceAddress, o.factoryID)
	req, err := httpclient.NewHTTPRequest(
		endpoint,
		http.MethodPost,
		ContentTypeUrlEncoded,
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
