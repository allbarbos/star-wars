package controller

import (
	"star-wars/api/handler"
	"star-wars/entity"
	"star-wars/planet"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Planets controller
type Planets struct {
	Srv planet.Service
}

// All get planets
func (p Planets) All(c *gin.Context) {
	limit, err :=  strconv.ParseInt(c.DefaultQuery("limit", "3"), 10, 64)
	if err != nil {
		handler.ResponseError(
			handler.BadRequest{
				Message: "limit is invalid",
			},
			c,
		)
		return
	}

	skip, err := strconv.ParseInt(c.DefaultQuery("skip", "0"), 10, 64)
	if err != nil {
		handler.ResponseError(
				handler.BadRequest{
				Message: "skip is invalid",
			},
			c,
		)
		return
	}

	planets, err := p.Srv.FindAll(limit, skip)

	if err != nil {
		handler.ResponseError(err, c)
		return
	}

	handler.ResponseSuccess(200, planets, c)
}

// ByName get planet
func (p Planets) ByName(c *gin.Context) {
	name := c.Param("name")
	planet, err := p.Srv.FindByName(name)

	if err == nil {
		handler.ResponseSuccess(200, planet, c)
	} else {
		handler.ResponseError(err, c)
	}
}

// ByID get planet
func (p Planets) ByID(c *gin.Context) {
	id := c.Param("id")
	planet, err := p.Srv.FindByID(id)

	if err == nil {
		handler.ResponseSuccess(200, planet, c)
	} else {
		handler.ResponseError(err, c)
	}
}

// Delete planet
func (p Planets) Delete(c *gin.Context) {
	id := c.Param("id")
	err := p.Srv.Delete(id)

	if err == nil {
		handler.ResponseSuccess(200, nil, c)
	} else {
		handler.ResponseError(err, c)
	}
}

// Post save planet
func (p Planets) Post(c *gin.Context) {
	var planet entity.Planet
	err := c.BindJSON(&planet)

	if err != nil {
		handler.ResponseError(
			handler.BadRequest{
				Message: "body is invalid",
			},
			c,
		)
		return
	}

	if planet.IsEmpty([]string{"Name", "Climate", "Terrain"}) {
		handler.ResponseError(
			handler.BadRequest{
				Message: "name, climate and terrain is required",
			},
			c,
		)
		return
	}

	err = p.Srv.Save(&planet)

	if err != nil {
		handler.ResponseError(err, c)
		return
	}

	handler.ResponseSuccess(201, planet, c)
}
