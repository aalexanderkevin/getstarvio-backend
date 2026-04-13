package reminder

import (
	"errors"
	"io"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/aalexanderkevin/getstarvio-backend/internal/http/middleware"
	"github.com/aalexanderkevin/getstarvio-backend/internal/http/response"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Log(c *gin.Context) {
	limit := 200
	if q := c.Query("limit"); q != "" {
		if v, err := strconv.Atoi(q); err == nil {
			limit = v
		}
	}

	res, err := h.svc.Log(middleware.UserID(c), c.Query("status"), limit)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.Success(c, res)
}

func (h *Handler) Retry(c *gin.Context) {
	if err := h.svc.Retry(middleware.UserID(c), c.Param("id")); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.Success(c, map[string]bool{"ok": true})
}

func (h *Handler) DashboardSummary(c *gin.Context) {
	res, err := h.svc.DashboardSummary(middleware.UserID(c))
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.Success(c, res)
}

func (h *Handler) MetaWebhook(c *gin.Context) {
	if c.Request.Method == "GET" {
		challenge, err := h.svc.VerifyMetaWebhook(
			c.Query("hub.mode"),
			c.Query("hub.verify_token"),
			c.Query("hub.challenge"),
		)
		if err != nil {
			if errors.Is(err, ErrMetaWebhookUnauthorized) {
				c.String(401, "unauthorized")
				return
			}
			c.String(400, err.Error())
			return
		}
		c.String(200, challenge)
		return
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(400, "invalid payload")
		return
	}
	if err := h.svc.HandleMetaWebhook(body, c.GetHeader("X-Hub-Signature-256")); err != nil {
		if errors.Is(err, ErrMetaWebhookUnauthorized) {
			c.String(401, "unauthorized")
			return
		}
		c.String(400, err.Error())
		return
	}
	c.String(200, "EVENT_RECEIVED")
}
