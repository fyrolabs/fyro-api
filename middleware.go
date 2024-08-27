package fapi

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// No errors
		errCount := len(c.Errors)
		if errCount == 0 {
			return
		}

		var lastErr *ResponseError
		for i, ginErr := range c.Errors {
			var respError *ResponseError
			switch {
			case !errors.As(ginErr.Err, &respError):
				respError = ServerError(ginErr.Err)
			}

			if respError.Name == "ServerError" {
				LogError(respError)
			}

			if i == errCount-1 {
				lastErr = respError
			}
		}

		c.JSON(lastErr.Status, Response{
			Error: lastErr,
		})
	}
}

func LogError(err *ResponseError) {
	if err.Name != "ServerError" {
		return // Don't log non-server errors
	}

	// .Message is not shown to public, log internally
	log.Printf("[SERVER_ERROR]: %s\n", err.PrivateMessage)
}
