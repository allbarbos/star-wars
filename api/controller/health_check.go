package controller

import (
	"star-wars/entity"
	"star-wars/planet"

	"github.com/gin-gonic/gin"
)

// HealthCheck controller
type HealthCheck struct {
	DB planet.Repository
}

// HealthCheck returns application health
func (h HealthCheck) HealthCheck(c *gin.Context) {
	hc := entity.HealthCheck{
		Status: "ok",
		Dependencies: entity.Dependencies{
			MongoDB: h.DB.Ping(),
		},
	}

	status := hc.CheckDependencies()
	c.JSON(status, hc)
}
