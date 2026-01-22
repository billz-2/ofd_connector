gen-mock: 
	mockgen -destination=/Users/chris/go/github.com/billz/ofd_connector/mock/receipt.go -package=mock_ofdconnector github.com/billz-2/ofd_connector ReceiptI 
	mockgen -destination=/Users/chris/go/github.com/billz/ofd_connector/mock/zreport.go -package=mock_ofdconnector github.com/billz-2/ofd_connector ZReportI
	mockgen -destination=/Users/chris/go/github.com/billz/ofd_connector/mock/fiscalDriveList.go -package=mock_ofdconnector github.com/billz-2/ofd_connector FiscalDriveLister