package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/aalexanderkevin/getstarvio-backend/internal/app"
	"github.com/aalexanderkevin/getstarvio-backend/internal/http/middleware"
	"github.com/aalexanderkevin/getstarvio-backend/internal/modules/internaladmin"
)

func registerInternalAdminRoutes(internal *gin.RouterGroup, c *app.Container) {
	repo := internaladmin.NewRepo(c.DB)
	svc := internaladmin.NewService(repo, c.Cfg, c.Meta)
	h := internaladmin.NewHandler(svc)

	internal.POST("/auth/login", h.Login)
	internal.POST("/auth/refresh", h.Refresh)
	internal.POST("/auth/logout", h.Logout)

	authed := internal.Group("")
	authed.Use(middleware.InternalAuth())
	authed.GET("/categories", h.ListDefaultCategories)
	authed.POST("/categories", h.CreateDefaultCategory)
	authed.GET("/wa-templates", h.ListWATemplates)
	authed.GET("/wa-templates/variables", h.ListWATemplateVariables)
	authed.GET("/wa-templates/:id", h.GetWATemplate)
	authed.POST("/wa-templates", h.CreateWATemplate)
	authed.PATCH("/wa-templates/:id", h.UpdateWATemplate)
	authed.DELETE("/wa-templates/:id", h.DeleteWATemplate)
	authed.GET("/plan-config", h.GetPlanConfig)
	authed.PUT("/plan-config", h.UpdatePlanConfig)
}
