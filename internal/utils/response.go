package utils

import "github.com/gin-gonic/gin"

type Response struct {
	Success bool        `json:"success"`
	Error *Error     `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func ResponseSuccess(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Data: data,
	})
}

func ResponseError(c *gin.Context, statusCode int, errCode, errMsg string) {
	c.JSON(statusCode, Response{
		Success: false,
		Error: &Error{
			Code: errCode,
			Message: errMsg,
		},
	})
}

