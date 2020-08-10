package controller

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"star-wars/api/handler"
	"star-wars/entity"
	"star-wars/planet/mock_planet"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestByName(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	pathParam := gin.Param{Key: "name", Value: "Tatooine"}
	c.Params = []gin.Param{pathParam}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	srvMock := mock_planet.NewMockService(ctrl)
	srvMock.EXPECT().FindByName("Tatooine").Return(
		entity.Planet{
			ID:         "5f29e53f2939a742014a04af",
			Name:       "Tatooine",
			Climate:    "arid",
			Terrain:    "desert",
			TotalFilms: 5,
		},
		nil,
	)

	Planets{
		Srv: srvMock,
	}.ByName(c)

	assert.Equal(t, 200, w.Code)
	assert.Equal(
		t,
		`{"id":"5f29e53f2939a742014a04af","name":"Tatooine","climate":"arid","terrain":"desert","totalFilms":5}`,
		w.Body.String(),
	)
}

func TestByName_Error(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	pathParam := gin.Param{Key: "name", Value: "NotFound"}
	c.Params = []gin.Param{pathParam}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	srvMock := mock_planet.NewMockService(ctrl)

	srvMock.EXPECT().FindByName("NotFound").Return(entity.Planet{}, handler.NotFound{Message: "planet not found"})

	Planets{
		Srv: srvMock,
	}.ByName(c)

	assert.Equal(t, 404, w.Code)
	assert.Equal(t, `{"error":"planet not found"}`, w.Body.String())
}

func TestByID(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	pathParam := gin.Param{Key: "id", Value: "5f29e53f2939a742014a04af"}
	c.Params = []gin.Param{pathParam}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	srvMock := mock_planet.NewMockService(ctrl)
	srvMock.EXPECT().FindByID("5f29e53f2939a742014a04af").Return(
		entity.Planet{
			ID:         "5f29e53f2939a742014a04af",
			Name:       "Tatooine",
			Climate:    "arid",
			Terrain:    "desert",
			TotalFilms: 5,
		},
		nil,
	)

	Planets{
		Srv: srvMock,
	}.ByID(c)

	assert.Equal(t, 200, w.Code)
	assert.Equal(
		t,
		`{"id":"5f29e53f2939a742014a04af","name":"Tatooine","climate":"arid","terrain":"desert","totalFilms":5}`,
		w.Body.String(),
	)
}

func TestByID_Error(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	pathParam := gin.Param{Key: "id", Value: "NotFound"}
	c.Params = []gin.Param{pathParam}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	srvMock := mock_planet.NewMockService(ctrl)

	srvMock.EXPECT().FindByID("NotFound").Return(entity.Planet{}, handler.NotFound{Message: "planet not found"})

	Planets{
		Srv: srvMock,
	}.ByID(c)

	assert.Equal(t, 404, w.Code)
	assert.Equal(t, `{"error":"planet not found"}`, w.Body.String())
}

func TestDelete(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	pathParam := gin.Param{Key: "id", Value: "5f29e53f2939a742014a04af"}
	c.Params = []gin.Param{pathParam}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	srvMock := mock_planet.NewMockService(ctrl)
	srvMock.EXPECT().Delete("5f29e53f2939a742014a04af").Return(nil)

	Planets{
		Srv: srvMock,
	}.Delete(c)

	assert.Equal(t, 200, w.Code)
}

func TestDelete_Error(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	pathParam := gin.Param{Key: "id", Value: "5f29e53f2939a742014a04af"}
	c.Params = []gin.Param{pathParam}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	srvMock := mock_planet.NewMockService(ctrl)
	srvMock.EXPECT().Delete("5f29e53f2939a742014a04af").Return(handler.InternalServer{Message: "error"})

	Planets{
		Srv: srvMock,
	}.Delete(c)

	assert.Equal(t, 500, w.Code)
	assert.Equal(t, `{"error":"internal server error"}`, w.Body.String())
}

func TestAll(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "http://t.test/?limit=3&skip=0", nil)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	srvMock := mock_planet.NewMockService(ctrl)

	var limit, skip int64
	limit = 3
	skip = 0

	srvMock.EXPECT().FindAll(limit, skip).Return(
		[]entity.Planet{
			{
				ID:         "5f2c891e9a9e070b1ef2e28c",
				Name:       "Alderaan",
				Climate:    "temperate",
				Terrain:    "grasslands, mountains",
				TotalFilms: 2,
			},
			{
				ID:         "5f2c891e9a9e070b1ef2e28d",
				Name:       "Tatooine",
				Climate:    "arid",
				Terrain:    "desert",
				TotalFilms: 5,
			},
			{
				ID:         "5f2c891e9a9e070b1ef2e28e",
				Name:       "Yavin IV",
				Climate:    "temperate, tropical",
				Terrain:    "jungle, rainforests",
				TotalFilms: 1,
			},
		},
		nil,
	)

	Planets{
		Srv: srvMock,
	}.All(c)

	assert.Equal(t, 200, w.Code)
	assert.Equal(
		t,
		`[{"id":"5f2c891e9a9e070b1ef2e28c","name":"Alderaan","climate":"temperate","terrain":"grasslands, mountains","totalFilms":2},{"id":"5f2c891e9a9e070b1ef2e28d","name":"Tatooine","climate":"arid","terrain":"desert","totalFilms":5},{"id":"5f2c891e9a9e070b1ef2e28e","name":"Yavin IV","climate":"temperate, tropical","terrain":"jungle, rainforests","totalFilms":1}]`,
		w.Body.String(),
	)
}

func TestAll_LimitInvalid(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "http://t.test/?limit=a&skip=0", nil)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	srvMock := mock_planet.NewMockService(ctrl)

	Planets{
		Srv: srvMock,
	}.All(c)

	assert.Equal(t, 400, w.Code)
	assert.Equal(
		t,
		`{"error":"limit is invalid"}`,
		w.Body.String(),
	)
}

func TestAll_SkipInvalid(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "http://t.test/?limit=3&skip=a", nil)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	srvMock := mock_planet.NewMockService(ctrl)

	Planets{
		Srv: srvMock,
	}.All(c)

	assert.Equal(t, 400, w.Code)
	assert.Equal(
		t,
		`{"error":"skip is invalid"}`,
		w.Body.String(),
	)
}

func TestAll_Error(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "http://t.test/?limit=3&skip=0", nil)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	srvMock := mock_planet.NewMockService(ctrl)

	var limit, skip int64
	limit = 3
	skip = 0

	srvMock.EXPECT().FindAll(limit, skip).Return(
		[]entity.Planet{},
		handler.InternalServer{Message: "error"},
	)

	Planets{
		Srv: srvMock,
	}.All(c)

	assert.Equal(t, 500, w.Code)
	assert.Equal(
		t,
		`{"error":"internal server error"}`,
		w.Body.String(),
	)
}

func TestPost(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := bytes.NewBufferString(`{"name":"Kamino","climate":"temperate","terrain":"ocean"}`)
	c.Request, _ = http.NewRequest("POST", "/planets", body)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	srvMock := mock_planet.NewMockService(ctrl)

	srvMock.EXPECT().Save(
		&entity.Planet{
			Name:    "Kamino",
			Climate: "temperate",
			Terrain: "ocean",
		},
	).Return(nil)

	Planets{
		Srv: srvMock,
	}.Post(c)

	assert.Equal(t, 201, w.Code)
}

func TestPost_InvalidPayload(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := bytes.NewBufferString(``)
	c.Request, _ = http.NewRequest("POST", "/planets", body)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	srvMock := mock_planet.NewMockService(ctrl)

	Planets{
		Srv: srvMock,
	}.Post(c)

	assert.Equal(t, 400, w.Code)
	assert.Equal(
		t,
		`{"error":"body is invalid"}`,
		w.Body.String(),
	)
}

func TestPost_InvalidFields(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := bytes.NewBufferString(`{"name":"","climate":"temperate","terrain":"ocean"}`)
	c.Request, _ = http.NewRequest("POST", "/planets", body)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	srvMock := mock_planet.NewMockService(ctrl)

	Planets{
		Srv: srvMock,
	}.Post(c)

	assert.Equal(t, 400, w.Code)
	assert.Equal(
		t,
		`{"error":"name, climate and terrain is required"}`,
		w.Body.String(),
	)
}

func TestPost_InternalError(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := bytes.NewBufferString(`{"name":"Kamino","climate":"temperate","terrain":"ocean"}`)
	c.Request, _ = http.NewRequest("POST", "/planets", body)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	srvMock := mock_planet.NewMockService(ctrl)

	srvMock.EXPECT().Save(
		&entity.Planet{
			Name:    "Kamino",
			Climate: "temperate",
			Terrain: "ocean",
		},
	).Return(handler.BadRequest{Message: "planet already registered"})

	Planets{
		Srv: srvMock,
	}.Post(c)

	assert.Equal(t, 400, w.Code)
	assert.Equal(
		t,
		`{"error":"planet already registered"}`,
		w.Body.String(),
	)
}
