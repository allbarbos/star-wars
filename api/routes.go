package api

import (
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

	router.GET("/health-check", healthCtrl().HealthCheck)
	router.GET("/planets", planetsCtrl().All)
	router.GET("/planets/id/:id", planetsCtrl().ByID)
	router.GET("/planets/name/:name", planetsCtrl().ByName)
	router.DELETE("/planets/:id", planetsCtrl().Delete)
	router.POST("/planets", planetsCtrl().Post)

	return router
}

func healthCtrl() controller.HealthCheck {
	return controller.HealthCheck{
		DB: planet.NewRepository(),
	}
}

func planetsCtrl() controller.Planets {
	sw := swapi.New()
	r := planet.NewRepository()
	s := planet.NewService(r, sw)
	return controller.Planets{
		Srv: s,
	}
}