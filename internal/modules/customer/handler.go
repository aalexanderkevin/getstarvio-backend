package customer

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/aalexanderkevin/getstarvio-backend/internal/http/middleware"
	"github.com/aalexanderkevin/getstarvio-backend/internal/http/response"
)

type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) List(c *gin.Context) {
	page, err := parsePositiveIntQuery(c, "page", 1)
	if err != nil {
		response.Error(c, 400, "invalid page parameter")
		return
	}
	limit, err := parsePositiveIntQuery(c, "limit", 20)
	if err != nil {
		response.Error(c, 400, "invalid limit parameter")
		return
	}
	if limit > 100 {
		limit = 100
	}

	res, err := h.svc.List(
		middleware.UserID(c),
		c.Query("q"),
		c.Query("status"),
		c.DefaultQuery("sort", "urgent"),
		c.Query("date"),
		page,
		limit,
	)
	if response.FetchErrorOrEmpty(c, err) {
		return
	}
	response.SuccessWithPaginationAndStatusCount(c, res.Data, res.Pagination, res.StatusCount)
}

func parsePositiveIntQuery(c *gin.Context, key string, fallback int) (int, error) {
	raw := c.Query(key)
	if raw == "" {
		return fallback, nil
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return 0, err
	}
	if v < 1 {
		return 0, fmt.Errorf("%s must be >= 1", key)
	}
	return v, nil
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
	res, err := h.svc.CheckinLookup(middleware.UserID(c), req.PhoneNumber)
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
