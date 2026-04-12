package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/aalexanderkevin/getstarvio-backend/internal/app"
	"github.com/aalexanderkevin/getstarvio-backend/internal/modules/internaladmin"
)

func registerInternalAdminRoutes(internal *gin.RouterGroup, c *app.Container) {
	repo := internaladmin.NewRepo(c.DB)
	svc := internaladmin.NewService(repo)
	h := internaladmin.NewHandler(svc)

	internal.GET("/plan-config", h.GetPlanConfig)
	internal.PUT("/plan-config", h.UpdatePlanConfig)
}
