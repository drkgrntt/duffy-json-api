package main

import (
	"log"

	"github.com/drkgrntt/duffy-json-api/controllers"
	"github.com/drkgrntt/duffy-json-api/database"
	"github.com/drkgrntt/duffy-json-api/routes"
	"github.com/drkgrntt/duffy-json-api/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	server                     *gin.Engine
	ShowController             controllers.ShowController
	ShowRouteController        routes.ShowRouteController
	SurveyController           controllers.SurveyController
	SurveyRouteController      routes.SurveyRouteController
	HotelPriceController       controllers.HotelPriceController
	HotelPriceRouteController  routes.HotelPriceRouteController
	DemographicController      controllers.DemographicController
	DemographicRouteController routes.DemographicRouteController
)

func init() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	database.ConnectDB(&config)
	database.ConnectSurveyDB(&config)
	database.ConnectHpfDB(&config)

	ShowController = controllers.NewShowController(database.GetDatabase())
	ShowRouteController = routes.NewRouteShowController(ShowController)

	SurveyController = controllers.NewSurveyController(database.GetSurveyDatabase())
	SurveyRouteController = routes.NewRouteSurveyController(SurveyController)

	HotelPriceController = controllers.NewHotelPriceController(database.GetHpfDatabase())
	HotelPriceRouteController = routes.NewRouteHotelPriceController(HotelPriceController)

	DemographicController = controllers.NewDemographicController(database.GetDatabase())
	DemographicRouteController = routes.NewRouteDemographicController(DemographicController)

	log.Println("Server is running in", config.Environment, "mode")
	if config.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	server = gin.Default()
}

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	// Allow all origins
	corsConfig.AllowAllOrigins = true

	server.Use(cors.New(corsConfig))

	router := server.Group("/")

	ShowRouteController.ShowRoutes(router)
	SurveyRouteController.SurveyRoutes(router)
	HotelPriceRouteController.HotelPriceRoutes(router)
	DemographicRouteController.DemographicRoutes(router)

	log.Fatal(server.Run(":" + config.ServerPort))
}
