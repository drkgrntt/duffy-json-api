package routes

import (
	"github.com/drkgrntt/duffy-json-api/controllers"
	"github.com/gin-gonic/gin"
)

type SalesRouteController struct {
	salesController controllers.SalesController
}

func NewSalesRouteController(salesController controllers.SalesController) SalesRouteController {
	return SalesRouteController{salesController}
}

func (c *SalesRouteController) SalesRoutes(rg *gin.RouterGroup) {
	router := rg.Group("sales")

	tallies := router.Group("tallies")
	tallies.GET("/tickets", c.salesController.GetTicketSales)
}
