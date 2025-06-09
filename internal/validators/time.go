package validators

import (
	"time"

	"github.com/billz-2/ofd_connector/internal/constants"
)

func ValidateTimeFormat(dateTime string) error {
	_, err := time.Parse(constants.TimeFormat, dateTime)
	return err
}
