package controller

import (
	"star-wars/api/handler"
	"star-wars/planet"

	"github.com/gin-gonic/gin"
)

type PlanetsController struct {
	Srv planet.Service
}

// GetByName find
func (p PlanetsController) GetByName(c *gin.Context) {
	name := c.Param("name")
	planet, err := p.Srv.FindByName(name)

	if err == nil {
		handler.ResponseSuccess(200, planet, c)
	} else {
		handler.ResponseError(err, c)
	}
}
