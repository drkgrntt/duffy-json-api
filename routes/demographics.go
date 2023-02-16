package routes

import (
	"github.com/drkgrntt/duffy-json-api/controllers"
	// "github.com/drkgrntt/duffy-json-api/middleware"
	"github.com/gin-gonic/gin"
)

type DemographicRouteController struct {
	demographicController controllers.DemographicController
}

func NewRouteDemographicController(demographicController controllers.DemographicController) DemographicRouteController {
	return DemographicRouteController{demographicController}
}

func (c *DemographicRouteController) DemographicRoutes(rg *gin.RouterGroup) {
	router := rg.Group("demographics")
	router.GET("/", c.demographicController.GetDemographics)
}
