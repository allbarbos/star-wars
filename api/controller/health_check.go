package controller

import (
	"context"
	"star-wars/entity"
	"star-wars/planet"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthCheck controller
type HealthCheck struct {
	DB planet.Repository
}

// HealthCheck returns application health
func (h HealthCheck) HealthCheck(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	hc := entity.HealthCheck{
		Status: "ok",
		Dependencies: entity.Dependencies{
			MongoDB: h.DB.Ping(ctx),
		},
	}

	status := hc.CheckDependencies()
	c.JSON(status, hc)
}
