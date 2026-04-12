package business

import (
	"github.com/gin-gonic/gin"

	"github.com/aalexanderkevin/getstarvio-backend/internal/http/middleware"
	"github.com/aalexanderkevin/getstarvio-backend/internal/http/response"
)

type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Bootstrap(c *gin.Context) {
	res, err := h.svc.GetBootstrap(middleware.UserID(c))
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.Success(c, res)
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	if err := h.svc.UpdateProfile(middleware.UserID(c), req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, map[string]bool{"ok": true})
}

func (h *Handler) UpdateWhatsApp(c *gin.Context) {
	var req UpdateWhatsAppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	if err := h.svc.UpdateWhatsApp(middleware.UserID(c), req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, map[string]bool{"ok": true})
}

func (h *Handler) UpdateSettings(c *gin.Context) {
	var req UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	if err := h.svc.UpdateSettings(middleware.UserID(c), req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, map[string]bool{"ok": true})
}
