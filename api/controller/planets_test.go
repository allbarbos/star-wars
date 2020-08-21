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
			wantStatusCode: 200,
			wantBody:       `[{"id":"5f2c891e9a9e070b1ef2e28c","name":"Alderaan","climate":"temperate","terrain":"grasslands, mountains","totalFilms":2}]`,
		},
		{
			name:           "when get planet with an invalid limit parameter",
			uri:            "http://t.test/?limit=a",
			wantStatusCode: 400,
			wantBody:       `{"error":"limit is invalid"}`,
		},
		{
			name:           "when get planet with an invalid skip parameter",
			uri:            "http://t.test/?skip=a",
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
			wantStatusCode: 200,
			wantBody:       `[{"id":"5f2c891e9a9e070b1ef2e28c","name":"Alderaan","climate":"temperate","terrain":"grasslands, mountains","totalFilms":2}]`,
		},
		{
			name:           "when get non-existent planet",
			uri:            "http://t.test/?search=test",
			findName:       "test",
			errPlanet:      handler.NotFound{Message: "planet not found"},
			wantStatusCode: 400,
			wantBody:       `{"error":"planet not found"}`,
		},
		{
			name:           "when an error happens",
			uri:            "http://t.test/?limit=1&skip=0",
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
	t.Parallel()

	type test struct {
		name           string
		idParam        string
		planet         *entity.Planet
		errPlanet      error
		wantStatusCode int
		wantBody       string
	}

	tests := []test{
		{
			name:    "happy path",
			idParam: "5f29e53f2939a742014a04af",
			planet: &entity.Planet{
				ID:         "5f29e53f2939a742014a04af",
				Name:       "Tatooine",
				Climate:    "arid",
				Terrain:    "desert",
				TotalFilms: 5,
			},
			wantStatusCode: 200,
			wantBody:       `{"id":"5f29e53f2939a742014a04af","name":"Tatooine","climate":"arid","terrain":"desert","totalFilms":5}`,
		},
		{
			name:           "error",
			idParam:        "NotFound",
			errPlanet:      handler.NotFound{Message: "planet not found"},
			wantStatusCode: 404,
			wantBody:       `{"error":"planet not found"}`,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = []gin.Param{{Key: "id", Value: tt.idParam}}
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			srvMock := mock_planet.NewMockService(ctrl)
			srvMock.EXPECT().FindByID(gomock.Any(), tt.idParam).Return(tt.planet, tt.errPlanet)

			Planets{
				Srv: srvMock,
			}.ByID(c)

			assert.Equal(t, tt.wantStatusCode, w.Code)
			assert.Equal(t, tt.wantBody, w.Body.String())
		})
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()

	type test struct {
		name           string
		idParam        string
		planet         *entity.Planet
		err            error
		wantStatusCode int
		wantBody       string
	}

	tests := []test{
		{
			name:           "happy path",
			idParam:        "5f29e53f2939a742014a04af",
			wantStatusCode: 200,
			wantBody:       ``,
		},
		{
			name:           "error",
			idParam:        "5f29e53f2939a742014a04af",
			err:            handler.InternalServer{Message: "error"},
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
			c.Params = []gin.Param{{Key: "id", Value: tt.idParam}}
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			srvMock := mock_planet.NewMockService(ctrl)
			srvMock.EXPECT().Delete(gomock.Any(), tt.idParam).Return(tt.err)

			Planets{
				Srv: srvMock,
			}.Delete(c)

			assert.Equal(t, tt.wantStatusCode, w.Code)
			assert.Equal(t, tt.wantBody, w.Body.String())
		})
	}
}

func TestPost(t *testing.T) {
	t.Parallel()

	type test struct {
		name           string
		body           string
		planet         *entity.Planet
		err            error
		wantStatusCode int
		wantBody       string
	}

	tests := []test{
		{
			name: "happy path",
			body: `{"name":"Kamino","climate":"temperate","terrain":"ocean"}`,
			planet: &entity.Planet{
				Name:    "Kamino",
				Climate: "temperate",
				Terrain: "ocean",
			},
			wantStatusCode: 201,
		},
		{
			name:           "when invalid payload",
			body:           ``,
			wantStatusCode: 400,
			wantBody:       `{"error":"body is invalid"}`,
		},
		{
			name:           "when invalid fields",
			body:           `{"name":"","climate":"temperate","terrain":"ocean"}`,
			wantStatusCode: 400,
			wantBody:       `{"error":"name, climate and terrain is required"}`,
		},
		{
			name: "when internal error",
			body: `{"name":"Kamino","climate":"temperate","terrain":"ocean"}`,
			planet: &entity.Planet{
				Name:    "Kamino",
				Climate: "temperate",
				Terrain: "ocean",
			},
			err:            handler.BadRequest{Message: "planet already registered"},
			wantStatusCode: 400,
			wantBody:       `{"error":"name, climate and terrain is required"}`,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("POST", "/planets", bytes.NewBufferString(tt.body))

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			srvMock := mock_planet.NewMockService(ctrl)

			if tt.planet != nil || tt.err != nil {
				srvMock.EXPECT().Save(gomock.Any(), tt.planet).Return(tt.err)
			}

			Planets{
				Srv: srvMock,
			}.Post(c)

			assert.Equal(t, tt.wantStatusCode, w.Code)
		})
	}
}
