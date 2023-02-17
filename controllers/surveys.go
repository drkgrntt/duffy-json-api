package controllers

import (
	"net/http"
	"time"

	"github.com/drkgrntt/duffy-json-api/models"
	"github.com/drkgrntt/duffy-json-api/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SurveyController struct {
	DB *gorm.DB
}

func NewSurveyController(DB *gorm.DB) SurveyController {
	return SurveyController{DB}
}

func (c *SurveyController) GetLatestSurveyTimestamp(ctx *gin.Context) {
	var survey models.Survey

	c.DB.Select("primaryId").Order("primaryId DESC").First(&survey)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": survey.Timestamp})
}

func (c *SurveyController) GetSurveyResults(ctx *gin.Context) {
	days, skip := utils.GetDaysAndSkip(ctx)

	var surveys []models.Survey
	earliest := time.Now().AddDate(0, 0, (-1 * days))
	latest := time.Now().AddDate(0, 0, (-1 * skip))

	c.DB.Where("demo_date > ?", earliest).
		Where("demo_date < ?", latest).
		Order("demo_date DESC").Find(&surveys)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"surveys": surveys}})
}

type GetSurveysGroupedByDateResponse struct {
	FirstBuyer   uint `json:"First Time"`
	ReturnBuyer  uint `json:"Returning"`
	NYC          uint `json:"NYC"`
	NYCSuburbs   uint `json:"NYC Suburbs"`
	OtherUS      uint `json:"Other U.S."`
	OtherCountry uint `json:"Other Country"`
}

func (c *SurveyController) GetSurveysGroupedByDate(ctx *gin.Context) {
	days, skip := utils.GetDaysAndSkip(ctx)

	var surveys []models.Survey
	earliest := time.Now().AddDate(0, 0, (-1 * days))
	latest := time.Now().AddDate(0, 0, (-1 * skip))

	c.DB.Where("demo_date > ?", earliest).
		Where("demo_date < ?", latest).
		Order("demo_date DESC").Find(&surveys)

	response := make(map[string]GetSurveysGroupedByDateResponse)

	for _, survey := range surveys {
		date := utils.FormatDate(survey.Date)
		_, ok := response[date]
		if !ok {
			response[date] = GetSurveysGroupedByDateResponse{}
		}
		val := response[date]

		switch survey.Noob {
		case "Yes":
			val.FirstBuyer++
		case "No":
			val.ReturnBuyer++
		}

		switch survey.Residence {
		case "NYC":
			val.NYC++
		case "NYC Suburbs":
			val.NYCSuburbs++
		case "Other U.S.":
			val.OtherUS++
		case "Other Country":
			val.OtherCountry++
		}

		response[date] = val
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"surveys": response}})
}
