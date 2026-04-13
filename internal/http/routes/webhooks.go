package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/aalexanderkevin/getstarvio-backend/internal/app"
	"github.com/aalexanderkevin/getstarvio-backend/internal/modules/billing"
	"github.com/aalexanderkevin/getstarvio-backend/internal/modules/reminder"
)

func registerWebhookRoutes(api *gin.RouterGroup, c *app.Container) {
	billingRepo := billing.NewRepo(c.DB)
	billingSvc := billing.NewService(billingRepo, c.Xendit)
	billingHandler := billing.NewHandler(billingSvc)

	reminderRepo := reminder.NewRepo(c.DB)
	reminderSvc := reminder.NewService(reminderRepo, c.Meta, c.Cfg.Meta)
	reminderHandler := reminder.NewHandler(reminderSvc)

	webhooks := api.Group("/webhooks")
	webhooks.POST("/xendit", billingHandler.XenditWebhook)
	webhooks.GET("/meta", reminderHandler.MetaWebhook)
	webhooks.POST("/meta", reminderHandler.MetaWebhook)
}
