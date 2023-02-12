package routes

import (
	"github.com/drkgrntt/duffy-json-api/controllers"
	// "github.com/drkgrntt/duffy-json-api/middleware"
	"github.com/gin-gonic/gin"
)

type ShowRouteController struct {
	showController controllers.ShowController
}

func NewRouteShowController(showController controllers.ShowController) ShowRouteController {
	return ShowRouteController{showController}
}

func (c *ShowRouteController) ShowRoutes(rg *gin.RouterGroup) {
	router := rg.Group("show")
	router.GET("/", c.showController.GetProductions)
	// router.GET("/:showId",  c.showController.GetShow)
}
