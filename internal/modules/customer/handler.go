package customer

import (
	"github.com/gin-gonic/gin"

	"github.com/aalexanderkevin/getstarvio-backend/internal/http/middleware"
	"github.com/aalexanderkevin/getstarvio-backend/internal/http/response"
)

type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) List(c *gin.Context) {
	res, err := h.svc.List(
		middleware.UserID(c),
		c.Query("q"),
		c.Query("status"),
		c.DefaultQuery("sort", "urgent"),
	)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.Success(c, res)
}

func (h *Handler) Create(c *gin.Context) {
	var req CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	if err := h.svc.Create(middleware.UserID(c), req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Created(c, map[string]bool{"ok": true})
}

func (h *Handler) Update(c *gin.Context) {
	var req UpdateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	if err := h.svc.Update(middleware.UserID(c), c.Param("id"), req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, map[string]bool{"ok": true})
}

func (h *Handler) Delete(c *gin.Context) {
	if err := h.svc.Delete(middleware.UserID(c), c.Param("id")); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, map[string]bool{"ok": true})
}

func (h *Handler) RecordVisit(c *gin.Context) {
	var req VisitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	if err := h.svc.RecordVisit(middleware.UserID(c), req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, map[string]bool{"ok": true})
}

func (h *Handler) CheckinLookup(c *gin.Context) {
	var req CheckinLookupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	res, err := h.svc.CheckinLookup(middleware.UserID(c), req.WA)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, res)
}

func (h *Handler) CheckinSubmit(c *gin.Context) {
	var req CheckinSubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	if err := h.svc.CheckinSubmit(middleware.UserID(c), req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, map[string]bool{"ok": true})
}
