package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var (

	// List of characters replaced by javascript code to make queries url-safe.
	base64subs = map[string]string{
		"-": "+",
		"_": "/",
		".": "=",
	}
)

type Error struct {
	Message string `json:"error"`
}

func NewError(err error) Error {
	return Error{err.Error()}
}

func desanitize64(query string) string {
	// Before feeding the string into decoded, we must "reconstruct" the base64 data.
	// Javascript replaces a few characters to be url-safe.
	for olds, news := range base64subs {
		query = strings.Replace(query, olds, news, -1)
	}

	return query
}

// Send a query result to client
func serveResult(c *gin.Context, result interface{}, err interface{}) {
	if err == nil {
		successResponse(c, result)
	} else {
		badRequest(c, err)
	}
}

// Send successful response back to client
func successResponse(c *gin.Context, data interface{}) {
	c.JSON(200, data)
}

// Send an error response back to client
func errorResponse(c *gin.Context, status int, err interface{}) {
	var message interface{}

	switch v := err.(type) {
	case error:
		message = v.Error()
	case string:
		message = v
	default:
		message = v
	}

	c.AbortWithStatusJSON(status, gin.H{"status": status, "error": message})
}

// Send a bad request (http 400) back to client
func badRequest(c *gin.Context, err interface{}) {
	errorResponse(c, http.StatusBadRequest, err)
}

// Send a not found (http 404) back to client
func notFoundf(c *gin.Context, format string, a ...any) {
	errorResponse(c, http.StatusNotFound, fmt.Errorf(format, a...))
}
