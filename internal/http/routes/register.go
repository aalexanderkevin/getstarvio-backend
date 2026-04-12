package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/aalexanderkevin/getstarvio-backend/internal/app"
	"github.com/aalexanderkevin/getstarvio-backend/internal/http/middleware"
)

func Register(r *gin.Engine, c *app.Container) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/healthz", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"ok":      true,
			"service": "getstarvio-backend",
			"time":    time.Now().UTC().Format(time.RFC3339),
		})
	})

	api := r.Group(c.Cfg.Service.APIPrefix)
	registerAuthRoutes(api, c)
	registerWebhookRoutes(api, c)

	internalGroup := api.Group("/internal")
	internalGroup.Use(middleware.InternalToken())
	registerInternalAdminRoutes(internalGroup, c)

	authed := api.Group("")
	authed.Use(middleware.Auth())
	registerBusinessRoutes(authed, c)
	registerCategoryRoutes(authed, c)
	registerCustomerRoutes(authed, c)
	registerReminderRoutes(authed, c)
	registerBillingRoutes(authed, c)
}
