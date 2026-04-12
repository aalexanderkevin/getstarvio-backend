package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/aalexanderkevin/getstarvio-backend/internal/app"
	"github.com/aalexanderkevin/getstarvio-backend/internal/modules/customer"
)

func registerCustomerRoutes(authed *gin.RouterGroup, c *app.Container) {
	repo := customer.NewRepo(c.DB)
	svc := customer.NewService(repo)
	h := customer.NewHandler(svc)

	authed.GET("/customers", h.List)
	authed.POST("/customers", h.Create)
	authed.PATCH("/customers/:id", h.Update)
	authed.DELETE("/customers/:id", h.Delete)

	authed.POST("/visits", h.RecordVisit)
	authed.POST("/checkin/lookup", h.CheckinLookup)
	authed.POST("/checkin/submit", h.CheckinSubmit)
}
