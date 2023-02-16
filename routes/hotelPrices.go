package routes

import (
	"github.com/drkgrntt/duffy-json-api/controllers"
	// "github.com/drkgrntt/duffy-json-api/middleware"
	"github.com/gin-gonic/gin"
)

type HotelPriceRouteController struct {
	hotelPriceController controllers.HotelPriceController
}

func NewRouteHotelPriceController(hotelPriceController controllers.HotelPriceController) HotelPriceRouteController {
	return HotelPriceRouteController{hotelPriceController}
}

func (c *HotelPriceRouteController) HotelPriceRoutes(rg *gin.RouterGroup) {
	router := rg.Group("hotel-prices")
	router.GET("/", c.hotelPriceController.GetHotelPrices)
	router.GET("/today", c.hotelPriceController.GetTodaysAverage)
	router.GET("/week", c.hotelPriceController.GetThisWeeksAverage)
}
