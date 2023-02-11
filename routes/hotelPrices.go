package routes

import (
	"github.com/drkgrntt/duffy-json-api/controllers"
	// "github.com/drkgrntt/duffy-json-api/middleware"
	"github.com/gin-gonic/gin"
)

type HotelPricesRouteController struct {
	surveysController controllers.HotelPricesController
}

func NewRouteHotelPriceController(surveysController controllers.HotelPricesController) HotelPricesRouteController {
	return HotelPricesRouteController{surveysController}
}

func (c *HotelPricesRouteController) HotelPriceRoute(rg *gin.RouterGroup) {
	router := rg.Group("hotel-prices")
	router.GET("/", c.surveysController.GetHotelPrices)
}
