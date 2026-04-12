package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/aalexanderkevin/getstarvio-backend/internal/config"
	"github.com/aalexanderkevin/getstarvio-backend/internal/http/response"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func Auth() gin.HandlerFunc {
	cfg := config.Instance()
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			response.Error(c, 401, "missing bearer token")
			c.Abort()
			return
		}

		tok := strings.TrimPrefix(h, "Bearer ")
		claims := &Claims{}
		parsed, err := jwt.ParseWithClaims(tok, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWT.Secret), nil
		})
		if err != nil || !parsed.Valid {
			response.Error(c, 401, "invalid token")
			c.Abort()
			return
		}
		if claims.ExpiresAt == nil || claims.ExpiresAt.Time.Before(time.Now()) {
			response.Error(c, 401, "token expired")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Next()
	}
}

func UserID(c *gin.Context) string {
	v, _ := c.Get("user_id")
	s, _ := v.(string)
	return s
}
