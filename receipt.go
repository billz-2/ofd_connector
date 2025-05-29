package ofdconnector

// SaleParams represents the parameters for a sale operation
type SaleParams struct {
	ReceivedCash int64     `json:"ReceivedCash"`
	ReceivedCard int64     `json:"ReceivedCard"`
	Time         string    `json:"Time"`
	Type         int       `json:"Type"`
	Operation    int       `json:"Operation"`
	Location     Location  `json:"Location"`
	Items        []Item    `json:"Items"`
	ExtraInfo    ExtraInfo `json:"ExtraInfo"`
}

type Location struct {
	Latitude  float64 `json:"Latitude"`
	Longitude float64 `json:"Longitude"`
}

type CommissionInfo struct {
	TIN string `json:"TIN"`
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
	PhoneNumber       string `json:"PhoneNumber"`
	QRPaymentID       string `json:"QRPaymentID"`
	QRPaymentProvider int    `json:"QRPaymentProvider"`
	CashedOutFromCard int64  `json:"CashedOutFromCard"`
	PPTID             string `json:"PPTID"`
	CardType          int    `json:"CardType"`
}

// ReceiptInfo returned from ofd
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

///here methods to implement receipt endpoints
