package routes

import (
	"net/http"
	"os"
	"star-wars/api/controller"
	"star-wars/health"
	"star-wars/planet"

	"github.com/gin-gonic/gin"
)

func Config() *gin.Engine {
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.Use(configCors)

	router.GET("/health-check", healthCheckBuilder())
	// router.GET("/planets", indexBuilder())
	// router.POST("/planets", shortenerBuilder())
	// router.GET("/planets/:id", indexBuilder())
	router.GET("/planets/:name", planetsByNameBuilder())
	// router.DELETE("/planets/:id", indexBuilder())

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

func healthCheckBuilder() gin.HandlerFunc {
	return controller.HealthCheckController{
		DB: health.NewRepository(),
	}.HealthCheck
}

func planetsByNameBuilder() gin.HandlerFunc {
	r := planet.NewRepository()
	s := planet.NewService(r)
	return controller.PlanetsController{
		Srv: s,
	}.GetByName
}
