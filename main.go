package main

import (
	"log"
	"net/http"

	"github.com/drkgrntt/duffy-json-api/controllers"
	"github.com/drkgrntt/duffy-json-api/database"
	"github.com/drkgrntt/duffy-json-api/routes"
	"github.com/drkgrntt/duffy-json-api/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	server               *gin.Engine
	ShowsController      controllers.ShowsController
	ShowsRouteController routes.ShowsRouteController
)

func init() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	database.ConnectDB(&config)

	ShowsController = controllers.NewShowsController(database.DB)
	ShowsRouteController = routes.NewRouteShowsController(ShowsController)

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

	router := server.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		message := "Yo, pickem is up and running!"
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})

	ShowsRouteController.ShowsRoute(router)

	log.Fatal(server.Run(":" + config.ServerPort))
}
