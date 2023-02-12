package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/drkgrntt/duffy-json-api/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SurveysController struct {
	DB *gorm.DB
}

func NewSurveysController(DB *gorm.DB) SurveysController {
	return SurveysController{DB}
}

func (c *SurveysController) GetLatestSurveyTimestamp(ctx *gin.Context) {
	var survey models.Survey

	c.DB.Select("primaryId").Order("primaryId DESC").First(&survey)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": survey.Timestamp})
}

func (c *SurveysController) GetSurveyResults(ctx *gin.Context) {
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
			skip = 7
			err = nil
		}
	}

	var surveys []models.Survey
	earliest := time.Now().AddDate(0, 0, (-1 * days))
	latest := time.Now().AddDate(0, 0, (-1 * skip))

	c.DB.Where("demo_date > ?", earliest).
		Where("demo_date < ?", latest).
		Order("demo_date DESC").Find(&surveys)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"surveys": surveys}})
}
