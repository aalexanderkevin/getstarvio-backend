package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/aalexanderkevin/getstarvio-backend/internal/app"
	"github.com/aalexanderkevin/getstarvio-backend/internal/modules/category"
)

func registerCategoryRoutes(authed *gin.RouterGroup, c *app.Container) {
	repo := category.NewRepo(c.DB)
	svc := category.NewService(repo)
	h := category.NewHandler(svc)

	authed.GET("/categories", h.List)
	authed.POST("/categories", h.Create)
	authed.PATCH("/categories/:id", h.Update)
	authed.DELETE("/categories/:id", h.Delete)
}
