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
	server                 *gin.Engine
	ShowsController        controllers.ShowsController
	ShowsRouteController   routes.ShowsRouteController
	SurveysController      controllers.SurveysController
	SurveysRouteController routes.SurveyRouteController
)

func init() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	database.ConnectDB(&config)
	database.ConnectSurveyDB(&config)

	ShowsController = controllers.NewShowsController(database.GetDatabase())
	ShowsRouteController = routes.NewRouteShowsController(ShowsController)

	SurveysController = controllers.NewSurveysController(database.GetSurveyDatabase())
	SurveysRouteController = routes.NewRouteSurveyController(SurveysController)

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

	ShowsRouteController.ShowsRoute(router)
	SurveysRouteController.SurveyRoute(router)

	log.Fatal(server.Run(":" + config.ServerPort))
}
