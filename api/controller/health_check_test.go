package controller

import (
	"net/http/httptest"
	"star-wars/planet/mock_planet"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck( t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dbMock := mock_planet.NewMockRepository(ctrl)
	dbMock.EXPECT().Ping().Return("ok")

	HealthCheckController{
		DB: dbMock,
	}.HealthCheck(c)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"status\":\"ok\",\"dependencies\":{\"mongoDb\":\"ok\"}}", w.Body.String())
}

func TestHealthCheck_Error( t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dbMock := mock_planet.NewMockRepository(ctrl)
	dbMock.EXPECT().Ping().Return("error")

	HealthCheckController{
		DB: dbMock,
	}.HealthCheck(c)

	assert.Equal(t, 500, w.Code)
	assert.Equal(t, "{\"status\":\"error\",\"dependencies\":{\"mongoDb\":\"error\"}}", w.Body.String())
}
