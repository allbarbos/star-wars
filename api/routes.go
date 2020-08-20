package api

import (
	"net/http"
	"os"
	"star-wars/api/controller"
	"star-wars/planet"
	"star-wars/swapi"

	"github.com/gin-gonic/gin"
)

func Config() *gin.Engine {
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.Use(configCors)

	router.GET("/health-check", healthCtrl().HealthCheck)
	router.GET("/planets", planetsCtrl().All)
	router.GET("/planets/:id", planetsCtrl().ByID)
	// router.GET("/planets/name/:name", planetsCtrl().ByName)
	router.POST("/planets", planetsCtrl().Post)
	router.DELETE("/planets/:id", planetsCtrl().Delete)

	return router
}

func configCors(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Content-Type", "application/json")
	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusOK)
	}
}

func healthCtrl() controller.HealthCheck {
	return controller.HealthCheck{
		DB: planet.NewRepository(),
	}
}

func planetsCtrl() controller.Planets {
	s := planet.NewService(planet.NewRepository(), swapi.New())
	return controller.Planets{
		Srv: s,
	}
}
