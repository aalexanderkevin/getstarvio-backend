package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/aalexanderkevin/getstarvio-backend/internal/config"
	"github.com/aalexanderkevin/getstarvio-backend/internal/http/response"
)

func InternalToken() gin.HandlerFunc {
	cfg := config.Instance()
	return func(c *gin.Context) {
		tok := c.GetHeader("X-Internal-Token")
		if tok == "" || tok != cfg.Internal.Token {
			response.Error(c, 401, "invalid internal token")
			c.Abort()
			return
		}
		c.Next()
	}
}
