package api

import (
	"net/http"
	"os"
	"star-wars/api/controller"
	"star-wars/planet"

	"github.com/gin-gonic/gin"
)

func Config() *gin.Engine {
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.Use(configCors)

	router.GET("/health-check", healthCtrl().HealthCheck)
	router.GET("/planets/id/:id", planetsCtrl().GetByID)
	router.GET("/planets/name/:name", planetsCtrl().GetByName)
	router.DELETE("/planets/:id", planetsCtrl().Delete)
	// router.GET("/planets", indexBuilder())
	// router.POST("/planets", shortenerBuilder())

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

func healthCtrl() controller.HealthCheckController {
	return controller.HealthCheckController{
		DB: planet.NewRepository(),
	}
}

func planetsCtrl() controller.PlanetsController {
	r := planet.NewRepository()
	s := planet.NewService(r)
	return controller.PlanetsController{
		Srv: s,
	}
}
