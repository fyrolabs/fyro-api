package fapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Data  any   `json:"data,omitempty"`
	Error error `json:"error,omitempty"`
}

func Respond(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Data: data,
	})
}

func RespondOK(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Data: gin.H{
			"result": "ok",
		},
	})
}
