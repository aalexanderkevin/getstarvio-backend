package billing

import (
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"

	"github.com/aalexanderkevin/getstarvio-backend/internal/http/middleware"
	"github.com/aalexanderkevin/getstarvio-backend/internal/http/response"
)

type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Summary(c *gin.Context) {
	res, err := h.svc.Summary(middleware.UserID(c))
	if response.FetchErrorOrEmpty(c, err) {
		return
	}
	response.Success(c, res)
}

func (h *Handler) History(c *gin.Context) {
	res, err := h.svc.History(middleware.UserID(c))
	if response.FetchErrorOrEmpty(c, err) {
		return
	}
	response.Success(c, res)
}

func (h *Handler) ActivateSubscription(c *gin.Context) {
	if err := h.svc.ActivateSubscription(middleware.UserID(c)); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, map[string]bool{"ok": true})
}

func (h *Handler) CancelSubscription(c *gin.Context) {
	if err := h.svc.CancelSubscription(middleware.UserID(c)); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, map[string]bool{"ok": true})
}

func (h *Handler) Checkout(c *gin.Context) {
	var req CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	res, err := h.svc.CreateTopupCheckout(middleware.UserID(c), req)
	if err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, res)
}

func (h *Handler) XenditWebhook(c *gin.Context) {
	if !h.svc.ValidateWebhookToken(c.GetHeader("X-Callback-Token")) {
		response.Error(c, 401, "invalid callback token")
		return
	}

	payloadRaw, _ := io.ReadAll(c.Request.Body)
	var p XenditWebhookPayload
	if err := json.Unmarshal(payloadRaw, &p); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	if err := h.svc.HandleXenditWebhook(p, string(payloadRaw)); err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.Success(c, map[string]bool{"ok": true})
}

func (h *Handler) GetPlanConfig(c *gin.Context) {
	res, err := h.svc.GetPlanConfig(middleware.UserID(c))
	if response.FetchErrorOrEmpty(c, err) {
		return
	}
	response.Success(c, res)
}

func (h *Handler) UpdatePlanConfig(c *gin.Context) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	if err := h.svc.UpdatePlanConfig(middleware.UserID(c), data); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, map[string]bool{"ok": true})
}
