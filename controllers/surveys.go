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

func (c *SurveyController) GetSurveysGroupedByDate(ctx *gin.Context) {
	days, skip := utils.GetDaysAndSkip(ctx)

	var surveys []models.Survey
	earliest := time.Now().AddDate(0, 0, (-1 * days))
	latest := time.Now().AddDate(0, 0, (-1 * skip))

	c.DB.Where("demo_date > ?", earliest).
		Where("demo_date < ?", latest).
		Order("demo_date DESC").Find(&surveys)

	response := map[string][]models.Survey{}

	for _, survey := range surveys {
		val, ok := response[formatSurveyDate(survey.Date)]
		if !ok {
			response[formatSurveyDate(survey.Date)] = make([]models.Survey, 0)
		}
		response[formatSurveyDate(survey.Date)] = append(val, survey)
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"surveys": response}})
}

func formatSurveyDate(date time.Time) string {
	return date.Format("01-02-2006")
}
