package handler

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

// ResponseSuccess creates payload
func ResponseSuccess(status int, body interface{}, c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(status, body)
}

// ResponseError creates payload
func ResponseError(err error, c *gin.Context) {
	typeError := reflect.TypeOf(err).String()
	status := http.StatusInternalServerError
	message := err.Error()

	switch typeError {
	case "handler.BadRequest":
		status = http.StatusBadRequest
	case "handler.NotFound":
		status = http.StatusNotFound
	default:
		status = http.StatusInternalServerError
	}

	c.JSON(status, gin.H{"error": message})
}
