package routes

import (
	"github.com/drkgrntt/duffy-json-api/controllers"
	"github.com/gin-gonic/gin"
)

type ShowRouteController struct {
	showController controllers.ShowController
}

func NewShowRouteController(showController controllers.ShowController) ShowRouteController {
	return ShowRouteController{showController}
}

func (c *ShowRouteController) ShowRoutes(rg *gin.RouterGroup) {
	router := rg.Group("shows")
	router.GET("/", c.showController.GetProductions)
	router.GET("/names", c.showController.GetNames)

	tallies := router.Group("tallies")
	tallies.GET("/price-ranges", c.showController.GetPriceRanges)
	tallies.GET("/average-discounts", c.showController.GetAverageDiscounts)
	tallies.GET("/percent-perfs-at-tkts", c.showController.GetPercentageAtTkts)
}
