package utils

import (
	"time"
)

func FormatDate(date time.Time) string {
	location, _ := time.LoadLocation("America/New_York")
	return date.In(location).Format("01-02-2006")
}
