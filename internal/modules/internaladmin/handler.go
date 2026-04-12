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

func (h *Handler) GetPlanConfig(c *gin.Context) {
	res, err := h.svc.GetPlanConfig()
	if err != nil {
		response.Error(c, 500, err.Error())
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
