package ofdconnector

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/billz-2/ofd_connector/internal/constants"
	"github.com/billz-2/ofd_connector/internal/httpclient"
)

type ReceiptI interface {
	GetTXID(ctx context.Context, req SaleParams) (int64, error)
	RegisterTXID(ctx context.Context, txID int64) (ReceiptInfo, error)
}

type receiptConfigs struct {
	ServiceAddress string
	FactoryID      string
	HttpClient     httpclient.HTTPClient
}

// ofdConnector implements the OfdConnector interface
type receipt struct {
	serviceAddress string
	httpClient     httpclient.HTTPClient
	factoryID      string
}

func newReceipt(configs receiptConfigs) ReceiptI {
	return &receipt{
		serviceAddress: configs.ServiceAddress,
		httpClient:     configs.HttpClient,
		factoryID:      configs.FactoryID,
	}
}

// SaleParams represents the parameters for a sale operation
type SaleParams struct {
	ReceivedCash int64      `json:"ReceivedCash"`
	ReceivedCard int64      `json:"ReceivedCard"`
	Time         string     `json:"Time"`
	Type         int        `json:"Type"`
	Operation    int        `json:"Operation"`
	Location     Location   `json:"Location"`
	Items        []Item     `json:"Items"`
	ExtraInfo    ExtraInfo  `json:"ExtraInfo"`
	RefundInfo   RefundInfo `json:"RefundInfo"`
}

type RefundInfo struct {
	TerminalID string `json:"TerminalID"` // Fiscal module serial number where the refunded receipt was registered
	ReceiptSeq uint64 `json:"ReceiptSeq"` // Number of the receipt being refunded
	DateTime   string `json:"DateTime"`   // Date and time of the refunded receipt (format YYYYMMDDHHMMSS)
	FiscalSign string `json:"FiscalSign"` // Fiscal sign of the receipt
}

type Location struct {
	Latitude  float64 `json:"Latitude"`
	Longitude float64 `json:"Longitude"`
}

type CommissionInfo struct {
	TIN   string `json:"TIN"`
	PINFL string `json:"PINFL"`
}

type Item struct {
	Name           string          `json:"Name"`
	Barcode        string          `json:"Barcode"`
	Labels         []string        `json:"Labels"`
	SPIC           string          `json:"SPIC"`
	Units          int64           `json:"Units"`
	PackageCode    string          `json:"PackageCode"`
	OwnerType      int             `json:"OwnerType"`
	Price          int64           `json:"Price"`
	VATPercent     int             `json:"VATPercent"`
	VAT            int64           `json:"VAT"`
	Amount         int64           `json:"Amount"`
	Discount       int64           `json:"Discount"`
	Other          int64           `json:"Other"`
	CommissionInfo *CommissionInfo `json:"CommissionInfo,omitempty"`
}

type ExtraInfo struct {
	CarNumber         string `json:"CarNumber"`
	CardType          int    `json:"CardType"`
	CashedOutFromCard int64  `json:"CashedOutFromCard"`
	PhoneNumber       string `json:"PhoneNumber"`
	QRPaymentID       string `json:"QRPaymentID"`
	QRPaymentProvider int    `json:"QRPaymentProvider"`
	PPTID             string `json:"PPTID"`
	PINFL             string `json:"PINFL"`
	TIN               string `json:"TIN"`
}

// ReceiptInfo returned from ofd, response of RegisterTXID success
type ReceiptInfo struct {
	TerminalID string `json:"TerminalID"`
	ReceiptSeq uint64 `json:"ReceiptSeq"`
	DateTime   string `json:"DateTime"`
	FiscalSign string `json:"FiscalSign"`
	QRCodeURL  string `json:"QRCodeURL"`
}

type ReceiptFullInfo struct {
	TerminalID    string `json:"TerminalID"`
	ReceiptSeq    int    `json:"ReceiptSeq"`
	Time          string `json:"Time"`
	FiscalSign    string `json:"FiscalSign"`
	ReceiptType   string `json:"ReceiptType"`
	OperationType string `json:"OperationType"`
	ReceivedCash  int64  `json:"ReceivedCash"`
	ReceivedCard  int64  `json:"ReceivedCard"`
	TotalVAT      int64  `json:"TotalVAT"`
	ItemsCount    int    `json:"ItemsCount"`
}

type TotalAmount struct {
	Sale   int64 `json:"Sale"`
	Refund int64 `json:"Refund"`
}

type RegisterTXIDReq struct {
	TXID int64 `json:"TXID"`
}

// GetTXID returns the txID for a sale
func (r *receipt) GetTXID(ctx context.Context, params SaleParams) (int64, error) {
	// Prepare the request body
	bodyBytes, err := json.Marshal(params)
	if err != nil {
		return 0, fmt.Errorf("error marshalling body: %s", err.Error())
	}

	endpoint := fmt.Sprintf("%s/FiscalDrive/Receipt/GetTXID/%s", r.serviceAddress, r.factoryID)
	req, err := httpclient.NewHTTPRequest(
		endpoint,
		http.MethodPost,
		constants.ContentTypeJSON,
		bodyBytes,
		nil,
	)
	if err != nil {
		return 0, fmt.Errorf("error creating request: %s", err.Error())
	}

	resp := r.httpClient.Request(ctx, req)
	if resp.StatusCode != http.StatusOK {
		errorResp := errorResponse{}
		if err = json.Unmarshal(resp.Body, &errorResp); err != nil {
			return 0, fmt.Errorf("error unmarshalling error response: %s responseBody: %s",
				err.Error(),
				string(resp.Body),
			)
		}
		return 0, fmt.Errorf("failed to get txID: %s", errorResp.Reason)
	}

	// Parse the response
	var txID int64
	err = json.Unmarshal(resp.Body, &txID)
	if err != nil {
		return 0, fmt.Errorf("error unmarshalling response: %s", err.Error())
	}

	return txID, nil
}

func (r *receipt) RegisterTXID(ctx context.Context, txID int64) (ReceiptInfo, error) {

	txIDInfo := RegisterTXIDReq{TXID: txID}
	reqBody, err := json.Marshal(txIDInfo)
	if err != nil {
		return ReceiptInfo{}, fmt.Errorf("error marshalling request body: %s", err.Error())
	}
	endpoint := fmt.Sprintf(
		"%s/FiscalDrive/Receipt/RegisterTXID/%s", r.serviceAddress, r.factoryID)
	req, err := httpclient.NewHTTPRequest(
		endpoint,
		http.MethodPost,
		constants.ContentTypeUrlEncoded,
		reqBody,
		nil,
	)
	if err != nil {
		return ReceiptInfo{}, fmt.Errorf("error creating request: %s", err.Error())
	}

	resp := r.httpClient.Request(ctx, req)
	if resp.StatusCode != http.StatusOK {
		errorResp := errorResponse{}
		if err = json.Unmarshal(resp.Body, &errorResp); err != nil {
			return ReceiptInfo{}, fmt.Errorf("error unmarshalling error response: %s responseBody: %s",
				err.Error(),
				string(resp.Body),
			)
		}
		return ReceiptInfo{}, fmt.Errorf("failed to register txID: %s", errorResp.Reason)
	}

	receiptInfo := ReceiptInfo{}
	err = json.Unmarshal(resp.Body, &receiptInfo)
	if err != nil {
		return ReceiptInfo{}, fmt.Errorf("error unmarshalling response: %s", err.Error())
	}

	return receiptInfo, nil
}
