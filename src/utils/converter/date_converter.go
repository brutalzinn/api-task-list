package converter_util

import (
	"time"
)

func ToDateTime(datetime string) (time.Time, error) {
	date, err := time.Parse(time.RFC3339, datetime)
	return date, err
}
func ToDateTimeString(datetime time.Time) string {
	return datetime.Format(time.RFC3339)
}
func ToDate(datetime string) (time.Time, error) {
	date, err := time.Parse(time.DateOnly, datetime)
	return date, err
}
func ToDateString(datetime time.Time) string {
	return datetime.Format(time.DateOnly)
}
