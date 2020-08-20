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
	t.Parallel()

	type test struct {
		name           string
		uri            string
		findName       string
		planet         *entity.Planet
		planets        *[]entity.Planet
		errPlanet      error
		errPlanets     error
		wantStatusCode int
		wantBody       string
	}

	tests := []test{
		{
			name: "when get planet with a limit parameter",
			uri:  "http://t.test/?limit=1",
			planets: &[]entity.Planet{
				{
					ID:         "5f2c891e9a9e070b1ef2e28c",
					Name:       "Alderaan",
					Climate:    "temperate",
					Terrain:    "grasslands, mountains",
					TotalFilms: 2,
				},
			},
			planet:         nil,
			errPlanet:      nil,
			errPlanets:     nil,
			wantStatusCode: 200,
			wantBody:       `[{"id":"5f2c891e9a9e070b1ef2e28c","name":"Alderaan","climate":"temperate","terrain":"grasslands, mountains","totalFilms":2}]`,
		},
		{
			name:           "when get planet with an invalid limit parameter",
			uri:            "http://t.test/?limit=a",
			planets:        nil,
			planet:         nil,
			errPlanet:      nil,
			errPlanets:     nil,
			wantStatusCode: 400,
			wantBody:       `{"error":"limit is invalid"}`,
		},
		{
			name:           "when get planet with an invalid skip parameter",
			uri:            "http://t.test/?skip=a",
			planets:        nil,
			planet:         nil,
			errPlanet:      nil,
			errPlanets:     nil,
			wantStatusCode: 400,
			wantBody:       `{"error":"skip is invalid"}`,
		},
		{
			name:     "when get planet with a search parameter",
			uri:      "http://t.test/?search=Alderaan",
			findName: "Alderaan",
			planet: &entity.Planet{
				ID:         "5f2c891e9a9e070b1ef2e28c",
				Name:       "Alderaan",
				Climate:    "temperate",
				Terrain:    "grasslands, mountains",
				TotalFilms: 2,
			},
			errPlanet:      nil,
			planets:        nil,
			errPlanets:     nil,
			wantStatusCode: 200,
			wantBody:       `[{"id":"5f2c891e9a9e070b1ef2e28c","name":"Alderaan","climate":"temperate","terrain":"grasslands, mountains","totalFilms":2}]`,
		},
		{
			name:           "when get non-existent planet",
			uri:            "http://t.test/?search=test",
			findName:       "test",
			planet:         nil,
			errPlanet:      handler.NotFound{Message: "planet not found"},
			planets:        nil,
			errPlanets:     nil,
			wantStatusCode: 400,
			wantBody:       `{"error":"planet not found"}`,
		},
		{
			name:           "when an error happens",
			uri:            "http://t.test/?limit=1&skip=0",
			planets:        nil,
			planet:         nil,
			errPlanet:      nil,
			errPlanets:     handler.InternalServer{Message: "error"},
			wantStatusCode: 500,
			wantBody:       `{"error":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("GET", tt.uri, nil)
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			srvMock := mock_planet.NewMockService(ctrl)

			var limit, skip int64
			limit = 1
			skip = 0

			if tt.planets != nil || tt.errPlanets != nil {

				srvMock.EXPECT().FindAll(gomock.Any(), limit, skip).Return(tt.planets, tt.errPlanets)
			}

			if tt.findName != "" {
				srvMock.EXPECT().FindByName(gomock.Any(), tt.findName).Return(tt.planet, tt.errPlanet)
			}

			Planets{
				Srv: srvMock,
			}.All(c)

			assert.Equal(t, tt.wantStatusCode, w.Code)
			assert.Equal(t, tt.wantBody, w.Body.String())
		})
	}
}

func TestByID(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		pathParam := gin.Param{Key: "id", Value: "5f29e53f2939a742014a04af"}
		c.Params = []gin.Param{pathParam}
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		srvMock := mock_planet.NewMockService(ctrl)
		srvMock.EXPECT().FindByID(gomock.Any(), "5f29e53f2939a742014a04af").Return(
			&entity.Planet{
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
	})

	t.Run("error", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		pathParam := gin.Param{Key: "id", Value: "NotFound"}
		c.Params = []gin.Param{pathParam}
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		srvMock := mock_planet.NewMockService(ctrl)

		srvMock.EXPECT().FindByID(gomock.Any(), "NotFound").Return(nil, handler.NotFound{Message: "planet not found"})

		Planets{
			Srv: srvMock,
		}.ByID(c)

		assert.Equal(t, 404, w.Code)
		assert.Equal(t, `{"error":"planet not found"}`, w.Body.String())
	})
}

func TestDelete(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		pathParam := gin.Param{Key: "id", Value: "5f29e53f2939a742014a04af"}
		c.Params = []gin.Param{pathParam}
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		srvMock := mock_planet.NewMockService(ctrl)
		srvMock.EXPECT().Delete(gomock.Any(), "5f29e53f2939a742014a04af").Return(nil)

		Planets{
			Srv: srvMock,
		}.Delete(c)

		assert.Equal(t, 200, w.Code)
	})

	t.Run("error", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		pathParam := gin.Param{Key: "id", Value: "5f29e53f2939a742014a04af"}
		c.Params = []gin.Param{pathParam}
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		srvMock := mock_planet.NewMockService(ctrl)
		srvMock.EXPECT().Delete(gomock.Any(), "5f29e53f2939a742014a04af").Return(handler.InternalServer{Message: "error"})

		Planets{
			Srv: srvMock,
		}.Delete(c)

		assert.Equal(t, 500, w.Code)
		assert.Equal(t, `{"error":"internal server error"}`, w.Body.String())
	})
}

func TestPost(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body := bytes.NewBufferString(`{"name":"Kamino","climate":"temperate","terrain":"ocean"}`)
		c.Request, _ = http.NewRequest("POST", "/planets", body)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		srvMock := mock_planet.NewMockService(ctrl)

		srvMock.EXPECT().Save(
			gomock.Any(),
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
	})

	t.Run("InvalidPayload", func(t *testing.T) {
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
	})

	t.Run("InvalidFields", func(t *testing.T) {
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
	})

	t.Run("InternalError", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body := bytes.NewBufferString(`{"name":"Kamino","climate":"temperate","terrain":"ocean"}`)
		c.Request, _ = http.NewRequest("POST", "/planets", body)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		srvMock := mock_planet.NewMockService(ctrl)

		srvMock.EXPECT().Save(
			gomock.Any(),
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
	})
}
