package utils

import (
	"time"
)

func FormatDate(date time.Time) string {
	return date.Format("01-02-2006")
}
