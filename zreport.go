package ofdconnector

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/billz-2/ofd_connector/internal/constants"
	"github.com/billz-2/ofd_connector/internal/gateway"
	"github.com/billz-2/ofd_connector/internal/validators"
	"github.com/go-playground/validator/v10"
)

const (
	zReportCloseEndpoint     = "/FiscalDrive/ZReport/Close/"
	zReportInfoEndpoint      = "/FiscalDrive/ZReport/Info/"
	zReportOpenEndpoint      = "/FiscalDrive/ZReport/Open/"
	zReportSyncEndpoint      = "/DataBase/Files/Sync/ZReports/"
	fiscalMemoryInfoEndpoint = "/FiscalDrive/FiscalMemory/Info/"
)

type ZReportI interface {
	OpenZreport(ctx context.Context, createdTime string) error
	CloseZreport(ctx context.Context, closedTime string) error
	GetCurrentZReportInfo(ctx context.Context) (ZReportInfo, error)
	SyncZReports(ctx context.Context, itemsCount uint16) error
}

type zReportConfig struct {
	gateway gateway.Client
}

// ofdConnector implements the OfdConnector interface
type zReport struct {
	gateway gateway.Client
}

func newZReport(config zReportConfig) ZReportI {
	return &zReport{
		gateway: config.gateway,
	}
}

type indexInfo struct {
	Index uint32 `json:"Index"`
}

type fiscalMemoryInfoResp struct {
	ZReportsCount uint32 `json:"ZReportsCount"`
}

type ZReportInfo struct {
	ZReportIndex     uint32      `json:"ZReportIndex"`
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
			return fmt.Errorf(
				"invalid time format, can be 'now' or in format %s, err: %w",
				constants.TimeFormat,
				validateErr,
			)
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
			return fmt.Errorf(
				"invalid time format, can be 'now' or in format %s, err: %w",
				constants.TimeFormat,
				validateErr,
			)
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

func (o zReport) getFiscalMemoryInfo(ctx context.Context) (fiscalMemoryInfoResp, error) {
	endpoint := o.gateway.FactoryEndpoint(fiscalMemoryInfoEndpoint)
	resp, err := o.gateway.HTTPRequest(
		ctx,
		endpoint,
		http.MethodGet,
		constants.ContentTypeJSON,
		nil,
		nil,
	)
	if err != nil {
		return fiscalMemoryInfoResp{}, fmt.Errorf("error fetching fiscal memory info: %s", err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		errorResp := errorResponse{}
		if err := json.Unmarshal(resp.Body, &errorResp); err != nil {
			return fiscalMemoryInfoResp{}, fmt.Errorf("error unmarshalling error response: %s responseBody: %s",
				err.Error(),
				string(resp.Body),
			)
		}
		return fiscalMemoryInfoResp{}, fmt.Errorf("failed to get fiscal memory info: %s", errorResp.Reason)
	}

	info := fiscalMemoryInfoResp{}
	if err := json.Unmarshal(resp.Body, &info); err != nil {
		return fiscalMemoryInfoResp{}, fmt.Errorf("error unmarshalling fiscal memory info: %s", err.Error())
	}
	return info, nil
}

// GetCurrentZReportInfo returns the last closed ZReport info for the fiscal drive.
// The index param is kept for backward compatibility but is no longer used;
// the index is derived from FiscalMemory/Info as ZReportsCount-1.
func (o zReport) GetCurrentZReportInfo(ctx context.Context) (ZReportInfo, error) {
	memInfo, err := o.getFiscalMemoryInfo(ctx)
	if err != nil {
		return ZReportInfo{}, err
	}

	zReportIndex := memInfo.ZReportsCount - 1

	bodyBytes, err := json.Marshal(indexInfo{Index: zReportIndex})
	if err != nil {
		return ZReportInfo{}, fmt.Errorf("error marshalling body: %s", err.Error())
	}

	endpoint := o.gateway.FactoryEndpoint(zReportInfoEndpoint)
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
	zReportInfo.ZReportIndex = zReportIndex

	return zReportInfo, nil
}

func (o zReport) SyncZReports(ctx context.Context, itemsCount uint16) error {
	req := itemsCountReq{
		ItemsCount: itemsCount,
	}
	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		return errors.New("itemsCount must be in range [1, 32]")
	}

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("error marshalling body: %s", err.Error())
	}

	endpoint := o.gateway.FactoryEndpoint(zReportSyncEndpoint)
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
		return fmt.Errorf("failed to sync Z reports: %s", errorResp.Reason)
	}

	return nil
}
