package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/aalexanderkevin/getstarvio-backend/internal/app"
	"github.com/aalexanderkevin/getstarvio-backend/internal/modules/reminder"
)

func registerReminderRoutes(authed *gin.RouterGroup, c *app.Container) {
	repo := reminder.NewRepo(c.DB)
	svc := reminder.NewService(repo, c.Meta, c.Cfg.Meta)
	h := reminder.NewHandler(svc)

	authed.GET("/reminders/log", h.Log)
	authed.POST("/reminders/:id/retry", h.Retry)
	authed.GET("/dashboard/summary", h.DashboardSummary)
}
