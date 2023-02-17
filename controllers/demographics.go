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

type GetTalliesResponse struct {
	Domestic      uint `json:"domestic"`
	International uint `json:"international"`
}

func (c *DemographicController) GetTallies(ctx *gin.Context) {
	days, skip := utils.GetDaysAndSkip(ctx)

	var analytics []models.Analytic
	earliest := time.Now().AddDate(0, 0, (-1 * days))
	latest := time.Now().AddDate(0, 0, (-1 * skip))

	c.DB.Select("country, created_at").
		Where("created_at > ?", earliest).
		Where("created_at < ?", latest).
		Find(&analytics)

	response := make(map[string]GetTalliesResponse)

	for _, analytic := range analytics {
		date := utils.FormatDate(analytic.CreatedAt)
		_, ok := response[date]
		if !ok {
			response[date] = GetTalliesResponse{}
		}
		val := response[date]

		switch analytic.Country {
		case "US":
			val.Domestic++
		default:
			val.International++
		}

		response[date] = val
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"demogrpahics": response}})
}

func (c *DemographicController) GetDomesticTallies(ctx *gin.Context) {
	days, skip := utils.GetDaysAndSkip(ctx)

	var analytics []models.Analytic
	earliest := time.Now().AddDate(0, 0, (-1 * days))
	latest := time.Now().AddDate(0, 0, (-1 * skip))

	c.DB.Select("state, created_at").
		Where("country = ?", "US").
		Where("created_at > ?", earliest).
		Where("created_at < ?", latest).
		Find(&analytics)

	response := make(map[string]map[string]uint)

	for _, analytic := range analytics {
		date := utils.FormatDate(analytic.CreatedAt)
		_, ok := response[date]
		if !ok {
			response[date] = make(map[string]uint)
		}
		val := response[date]

		val[analytic.State]++

		response[date] = val
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"demogrpahics": response}})
}

func (c *DemographicController) GetInternationalTallies(ctx *gin.Context) {
	days, skip := utils.GetDaysAndSkip(ctx)

	var analytics []models.Analytic
	earliest := time.Now().AddDate(0, 0, (-1 * days))
	latest := time.Now().AddDate(0, 0, (-1 * skip))

	c.DB.Select("country, created_at").
		// Where("country != ?", "US").
		Where("created_at > ?", earliest).
		Where("created_at < ?", latest).
		Find(&analytics)

	response := make(map[string]map[string]uint)

	for _, analytic := range analytics {
		date := utils.FormatDate(analytic.CreatedAt)
		_, ok := response[date]
		if !ok {
			response[date] = make(map[string]uint)
		}
		val := response[date]

		val[utils.GetCountryFromCode(analytic.Country)]++

		response[date] = val
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"demogrpahics": response}})
}
