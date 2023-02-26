package utils

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/now"
)

func GetDaysAndSkip(ctx *gin.Context) (int, int) {
	days := 7
	skip := 0
	var err error

	daysQuery := ctx.Query("days")
	skipQuery := ctx.Query("skip")

	if daysQuery != "" {
		days, err = strconv.Atoi(daysQuery)
		if err != nil {
			days = 7
			err = nil
		}
	}
	if skipQuery != "" {
		skip, err = strconv.Atoi(skipQuery)
		if err != nil {
			skip = 0
			err = nil
		}
	}

	return days, skip
}

func GetEarliestAndLatest(ctx *gin.Context) (earliest time.Time, latest time.Time) {
	days, skip := GetDaysAndSkip(ctx)

	location, _ := time.LoadLocation("America/New_York")

	earliest = now.BeginningOfDay().In(location)
	latest = now.EndOfDay().In(location)

	earliest = earliest.AddDate(0, 0, (-1 * (days - 1)))
	latest = latest.AddDate(0, 0, (-1 * skip))

	return
}
