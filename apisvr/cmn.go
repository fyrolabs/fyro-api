package apisvr

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Respond success with JSON data
func Respond(c *gin.Context, data any) {
	c.JSON(http.StatusOK, data)
}

// RespondOK success with no data
func RespondOK(c *gin.Context) {
	c.Status(http.StatusOK)
}
