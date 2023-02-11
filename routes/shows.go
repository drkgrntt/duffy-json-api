package routes

import (
	"github.com/drkgrntt/duffy-json-api/controllers"
	// "github.com/drkgrntt/duffy-json-api/middleware"
	"github.com/gin-gonic/gin"
)

type ShowsRouteController struct {
	showsController controllers.ShowsController
}

func NewRouteShowsController(showsController controllers.ShowsController) ShowsRouteController {
	return ShowsRouteController{showsController}
}

func (c *ShowsRouteController) ShowsRoute(rg *gin.RouterGroup) {
	router := rg.Group("shows")
	router.GET("/", c.showsController.GetProductions)
	// router.GET("/:showsId",  c.showsController.GetShows)
}
