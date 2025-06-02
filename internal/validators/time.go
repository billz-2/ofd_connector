package validators

import (
	"time"

	"github.com/billz-2/ofd_connector/internal/constants"
)

func IsValidateTimeFormat(dateTime string) bool {
	_, err := time.Parse(constants.TimeFormat, dateTime)
	return err == nil
}
