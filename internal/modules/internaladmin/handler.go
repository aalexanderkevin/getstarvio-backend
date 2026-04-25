package internaladmin

import (
	"github.com/gin-gonic/gin"

	"github.com/aalexanderkevin/getstarvio-backend/internal/http/response"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	res, err := h.svc.Login(req)
	if err != nil {
		response.Error(c, 401, err.Error())
		return
	}
	response.Success(c, res)
}

func (h *Handler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	res, err := h.svc.Refresh(req)
	if err != nil {
		response.Error(c, 401, err.Error())
		return
	}
	response.Success(c, res)
}

func (h *Handler) Logout(c *gin.Context) {
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	if err := h.svc.Logout(req); err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.Success(c, map[string]bool{"ok": true})
}

func (h *Handler) ListDefaultCategories(c *gin.Context) {
	res, err := h.svc.ListDefaultCategories()
	if response.FetchErrorOrEmpty(c, err) {
		return
	}
	response.Success(c, res)
}

func (h *Handler) CreateDefaultCategory(c *gin.Context) {
	var req CreateDefaultCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	res, err := h.svc.CreateDefaultCategory(req)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Created(c, res)
}

func (h *Handler) GetPlanConfig(c *gin.Context) {
	res, err := h.svc.GetPlanConfig()
	if response.FetchErrorOrEmpty(c, err) {
		return
	}
	response.Success(c, res)
}

func (h *Handler) UpdatePlanConfig(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Error(c, 400, err.Error())
		return
	}

	if err := h.svc.UpdatePlanConfig(body); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, map[string]bool{"ok": true})
}

func (h *Handler) ListWATemplates(c *gin.Context) {
	res, err := h.svc.ListWATemplates(
		c.Query("category"),
		c.Query("status"),
		c.Query("metaTemplateName"),
	)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.Success(c, res)
}

func (h *Handler) GetWATemplate(c *gin.Context) {
	res, err := h.svc.GetWATemplate(c.Param("id"))
	if response.FetchErrorOrEmpty(c, err) {
		return
	}
	response.Success(c, res)
}

func (h *Handler) CreateWATemplate(c *gin.Context) {
	var req CreateWATemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	res, err := h.svc.CreateWATemplate(req)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Created(c, res)
}

func (h *Handler) UpdateWATemplate(c *gin.Context) {
	var req UpdateWATemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	if err := h.svc.UpdateWATemplate(c.Param("id"), req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, map[string]bool{"ok": true})
}

func (h *Handler) DeleteWATemplate(c *gin.Context) {
	if err := h.svc.DeleteWATemplate(c.Param("id")); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, map[string]bool{"ok": true})
}

func (h *Handler) ListWATemplateVariables(c *gin.Context) {
	response.Success(c, h.svc.ListWATemplateVariables())
}
