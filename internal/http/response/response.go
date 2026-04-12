package response

import "github.com/gin-gonic/gin"

type Envelope struct {
	Error   bool        `json:"error"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Envelope{Error: false, Data: data})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(201, Envelope{Error: false, Data: data})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(code, Envelope{Error: true, Message: message})
}
