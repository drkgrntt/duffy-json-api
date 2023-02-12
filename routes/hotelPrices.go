package routes

import (
	"github.com/drkgrntt/duffy-json-api/controllers"
	// "github.com/drkgrntt/duffy-json-api/middleware"
	"github.com/gin-gonic/gin"
)

type HotelPriceRouteController struct {
	surveysController controllers.HotelPriceController
}

func NewRouteHotelPriceController(surveysController controllers.HotelPriceController) HotelPriceRouteController {
	return HotelPriceRouteController{surveysController}
}

func (c *HotelPriceRouteController) HotelPriceRoutes(rg *gin.RouterGroup) {
	router := rg.Group("hotel-prices")
	router.GET("/", c.surveysController.GetHotelPrices)
	router.GET("/today", c.surveysController.GetTodaysAverage)
	router.GET("/week", c.surveysController.GetThisWeeksAverage)
}
