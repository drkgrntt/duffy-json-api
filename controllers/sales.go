package controllers

import (
	"net/http"
	"strings"

	"github.com/drkgrntt/duffy-json-api/models"
	"github.com/drkgrntt/duffy-json-api/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SalesController struct {
	DB *gorm.DB
}

func NewSalesController(DB *gorm.DB) SalesController {
	return SalesController{DB}
}

func (c *SalesController) GetTicketSales(ctx *gin.Context) {
	earliest, latest := utils.GetEarliestAndLatest(ctx)
	locations := ctx.QueryArray("location")

	var salesDays []models.TktsSalesDay

	c.DB.Select("date", "tickets_sold", "location").
		Where("date >= ?", earliest).
		Where("date < ?", latest).
		Find(&salesDays)

	response := make(map[string]map[string]int)

	for _, salesDay := range salesDays {
		date := utils.FormatDate(salesDay.Date)

		_, ok := response[date]
		if !ok {
			response[date] = make(map[string]int)
			response[date]["all"] = 0

			for _, location := range locations {
				response[date][strings.ToLower(location)] = 0
			}
		}

		response[date]["all"] += salesDay.TicketsSold

		for _, location := range locations {
			lowerLocation := strings.ToLower(location)
			if lowerLocation == strings.ToLower(salesDay.Location) {
				response[date][lowerLocation] += salesDay.TicketsSold
			}
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"sales": response}})
}
