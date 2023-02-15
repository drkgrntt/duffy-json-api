package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
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
