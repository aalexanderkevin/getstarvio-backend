package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/aalexanderkevin/getstarvio-backend/internal/app"
	"github.com/aalexanderkevin/getstarvio-backend/internal/modules/billing"
)

func registerBillingRoutes(authed *gin.RouterGroup, c *app.Container) {
	repo := billing.NewRepo(c.DB)
	svc := billing.NewService(repo, c.Xendit)
	h := billing.NewHandler(svc)

	authed.GET("/billing/summary", h.Summary)
	authed.GET("/billing/history", h.History)
	authed.POST("/billing/subscription/activate", h.ActivateSubscription)
	authed.POST("/billing/subscription/cancel", h.CancelSubscription)
	authed.POST("/billing/topup/checkout", h.Checkout)
}
