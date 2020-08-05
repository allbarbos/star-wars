package handler

import (
	"net/http"
	"reflect"
	"strings"

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

	if strings.Contains(typeError, "BadRequest") {
		status = http.StatusBadRequest
	}

	c.JSON(status, gin.H{"error": message})
}
