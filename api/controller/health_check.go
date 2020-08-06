package controller

import (
	"star-wars/entity"
	"star-wars/planet"

	"github.com/gin-gonic/gin"
)

// HealthCheckController health of application
type HealthCheckController struct {
	DB planet.Repository
}

// HealthCheck returns application health
func (h HealthCheckController) HealthCheck(c *gin.Context) {
	hc := entity.HealthCheck{
		Status: "ok",
		Dependencies: entity.Dependencies{
			MongoDB: h.DB.Ping(),
		},
	}

	status := hc.CheckDependencies()
	c.JSON(status, hc)
}
