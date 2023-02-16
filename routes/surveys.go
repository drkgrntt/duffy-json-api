package routes

import (
	"github.com/drkgrntt/duffy-json-api/controllers"
	// "github.com/drkgrntt/duffy-json-api/middleware"
	"github.com/gin-gonic/gin"
)

type SurveyRouteController struct {
	surveyController controllers.SurveyController
}

func NewRouteSurveyController(surveysController controllers.SurveyController) SurveyRouteController {
	return SurveyRouteController{surveysController}
}

func (c *SurveyRouteController) SurveyRoutes(rg *gin.RouterGroup) {
	router := rg.Group("surveys")
	router.GET("/", c.surveyController.GetSurveyResults)
	router.GET("/latest-timestamp", c.surveyController.GetLatestSurveyTimestamp)

	group := router.Group("tallies")
	group.GET("/by-date", c.surveyController.GetSurveysGroupedByDate)
}
