package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/aalexanderkevin/getstarvio-backend/internal/app"
	"github.com/aalexanderkevin/getstarvio-backend/internal/modules/auth"
)

func registerAuthRoutes(api *gin.RouterGroup, c *app.Container) {
	repo := auth.NewRepo(c.DB)
	svc := auth.NewService(repo, c.Cfg)
	h := auth.NewHandler(svc)

	g := api.Group("/auth")
	g.POST("/google/login", h.GoogleLogin)
	g.POST("/refresh", h.Refresh)
	g.POST("/logout", h.Logout)
}
