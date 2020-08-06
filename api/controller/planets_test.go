package controller

import (
	"net/http/httptest"
	"star-wars/api/handler"
	"star-wars/entity"
	"star-wars/planet/mock_planet"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetByName( t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	pathParam := gin.Param{Key: "name", Value: "Tatooine"}
	c.Params = []gin.Param{pathParam}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	srvMock := mock_planet.NewMockService(ctrl)
	srvMock.EXPECT().FindByName("Tatooine").Return(
		entity.Planet{
			ID: "5f29e53f2939a742014a04af",
			Name: "Tatooine",
			Climate: "arid",
			Terrain: "desert",
			TotalFilms: 5,
		},
		nil,
	)

	PlanetsController{
		Srv: srvMock,
	}.GetByName(c)

	assert.Equal(t, 200, w.Code)
	assert.Equal(
		t,
		`{"id":"5f29e53f2939a742014a04af","name":"Tatooine","climate":"arid","terrain":"desert","totalFilms":5}`,
		w.Body.String(),
	)
}

func TestGetByName_Error( t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	pathParam := gin.Param{Key: "name", Value: "NotFound"}
	c.Params = []gin.Param{pathParam}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	srvMock := mock_planet.NewMockService(ctrl)

	srvMock.EXPECT().FindByName("NotFound").Return(entity.Planet{}, handler.NotFound{ Message: "planet not found" })

	PlanetsController{
		Srv: srvMock,
	}.GetByName(c)

	assert.Equal(t, 404, w.Code)
	assert.Equal(t, `{"error":"planet not found"}`, w.Body.String())
}
