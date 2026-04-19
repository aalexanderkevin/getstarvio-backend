package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/aalexanderkevin/getstarvio-backend/internal/config"
	"github.com/aalexanderkevin/getstarvio-backend/internal/http/response"
)

type InternalClaims struct {
	InternalAdminID string `json:"internal_admin_id"`
	TokenType       string `json:"token_type"`
	jwt.RegisteredClaims
}

func InternalAuth() gin.HandlerFunc {
	cfg := config.Instance()
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			response.Error(c, 401, "missing bearer token")
			c.Abort()
			return
		}

		tok := strings.TrimPrefix(h, "Bearer ")
		claims := &InternalClaims{}
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
		if claims.InternalAdminID == "" || claims.TokenType != "internal_admin" {
			response.Error(c, 401, "invalid internal admin token")
			c.Abort()
			return
		}

		c.Set("internal_admin_id", claims.InternalAdminID)
		c.Next()
	}
}
