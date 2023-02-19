package routes

import (
	"github.com/drkgrntt/duffy-json-api/controllers"
	"github.com/gin-gonic/gin"
)

type DemographicRouteController struct {
	demographicController controllers.DemographicController
}

func NewDemographicRouteController(demographicController controllers.DemographicController) DemographicRouteController {
	return DemographicRouteController{demographicController}
}

func (c *DemographicRouteController) DemographicRoutes(rg *gin.RouterGroup) {
	router := rg.Group("demographics")

	tallies := router.Group("tallies")
	tallies.GET("/", c.demographicController.GetTallies)
	tallies.GET("/domestic", c.demographicController.GetDomesticTallies)
	tallies.GET("/international", c.demographicController.GetInternationalTallies)
}
