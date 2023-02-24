package converter_util

import (
	"time"
)

func ToDateTime(datetime string) (time.Time, error) {
	date, err := time.Parse(time.RFC3339, datetime)
	return date, err
}
