package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/aalexanderkevin/getstarvio-backend/internal/app"
	"github.com/aalexanderkevin/getstarvio-backend/internal/modules/business"
)

func registerBusinessRoutes(authed *gin.RouterGroup, c *app.Container) {
	repo := business.NewRepo(c.DB)
	svc := business.NewService(repo)
	h := business.NewHandler(svc)

	authed.GET("/me/bootstrap", h.Bootstrap)
	authed.PUT("/business/profile", h.UpdateProfile)
	authed.PUT("/business/whatsapp", h.UpdateWhatsApp)
	authed.PUT("/business/settings", h.UpdateSettings)
}
