package routes

import (
	"github.com/drkgrntt/duffy-json-api/controllers"
	// "github.com/drkgrntt/duffy-json-api/middleware"
	"github.com/gin-gonic/gin"
)

type SurveyRouteController struct {
	surveysController controllers.SurveysController
}

func NewRouteSurveyController(surveysController controllers.SurveysController) SurveyRouteController {
	return SurveyRouteController{surveysController}
}

func (c *SurveyRouteController) SurveyRoute(rg *gin.RouterGroup) {

	router := rg.Group("surveys")
	router.GET("/", c.surveysController.GetSurveyResults)
}
