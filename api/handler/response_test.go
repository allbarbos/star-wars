package handler

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestResponseSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	ResponseSuccess(200, "ok", c)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "\"ok\"", w.Body.String())
}

func TestResponseError_BadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	ResponseError(BadRequest{Message: "bad request"}, c)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "{\"error\":\"bad request\"}", w.Body.String())
}

func TestResponseError_InternalServer(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	ResponseError(InternalServer{Message: "internal server error"}, c)

	assert.Equal(t, 500, w.Code)
	assert.Equal(t, "{\"error\":\"internal server error\"}", w.Body.String())
}
