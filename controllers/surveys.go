package controllers

import (
	"net/http"
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

func (c *SurveysController) GetSurveyResults(ctx *gin.Context) {
	var surveys []models.Survey
	pastWeek := time.Now().AddDate(0, 0, -7)

	c.DB.Where("demo_date > ?", pastWeek).Order("demo_date DESC").Find(&surveys)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"surveys": surveys}})
}
