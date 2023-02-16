package controllers

import (
	"net/http"
	"time"

	"github.com/drkgrntt/duffy-json-api/models"
	"github.com/drkgrntt/duffy-json-api/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DemographicController struct {
	DB *gorm.DB
}

func NewDemographicController(DB *gorm.DB) DemographicController {
	return DemographicController{DB}
}

type GetDemographicsResponse struct {
	Domestic      uint `json:"domestic"`
	International uint `json:"international"`
}

func (c *DemographicController) GetDemographics(ctx *gin.Context) {
	days, skip := utils.GetDaysAndSkip(ctx)

	var analytics []models.Analytic
	earliest := time.Now().AddDate(0, 0, (-1 * days))
	latest := time.Now().AddDate(0, 0, (-1 * skip))

	c.DB.Select("country, created_at").Where("created_at > ?", earliest).
		Where("created_at < ?", latest).
		Order("created_at DESC").Find(&analytics)

	response := make(map[string]GetDemographicsResponse)

	for _, analytic := range analytics {
		_, ok := response[utils.FormatDate(analytic.CreatedAt)]
		if !ok {
			response[utils.FormatDate(analytic.CreatedAt)] = GetDemographicsResponse{}
		}
		val := response[utils.FormatDate(analytic.CreatedAt)]

		switch analytic.Country {
		case "US":
			val.Domestic++
		default:
			val.International++
		}

		response[utils.FormatDate(analytic.CreatedAt)] = val
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"demogrpahics": response}})
}
