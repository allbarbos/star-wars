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

// GetByID find
func (p PlanetsController) GetByID(c *gin.Context) {
	id := c.Param("id")
	planet, err := p.Srv.FindByID(id)

	if err == nil {
		handler.ResponseSuccess(200, planet, c)
	} else {
		handler.ResponseError(err, c)
	}
}

// Delete planet
func (p PlanetsController) Delete(c *gin.Context) {
	id := c.Param("id")
	err := p.Srv.Delete(id)

	if err == nil {
		handler.ResponseSuccess(200, nil, c)
	} else {
		handler.ResponseError(err, c)
	}
}
