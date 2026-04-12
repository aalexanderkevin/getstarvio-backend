package meta

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/segmentio/ksuid"

	"github.com/aalexanderkevin/getstarvio-backend/internal/config"
)

type Client struct {
	cfg        config.MetaConfig
	httpClient *http.Client
}

type SendTemplateInput struct {
	To           string
	TemplateName string
	LanguageCode string
	Parameters   []string
}

func NewClient(cfg config.MetaConfig) *Client {
	return &Client{
		cfg:        cfg,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) SendTemplate(ctx context.Context, in SendTemplateInput) (string, error) {
	if c.cfg.AccessToken == "" || c.cfg.PhoneNumberID == "" {
		return "mock-" + ksuid.New().String(), nil
	}

	if in.LanguageCode == "" {
		in.LanguageCode = "id"
	}

	components := make([]map[string]any, 0)
	if len(in.Parameters) > 0 {
		params := make([]map[string]string, 0, len(in.Parameters))
		for _, p := range in.Parameters {
			params = append(params, map[string]string{"type": "text", "text": p})
		}
		components = append(components, map[string]any{
			"type":       "body",
			"parameters": params,
		})
	}

	payload := map[string]any{
		"messaging_product": "whatsapp",
		"to":                in.To,
		"type":              "template",
		"template": map[string]any{
			"name":       in.TemplateName,
			"language":   map[string]string{"code": in.LanguageCode},
			"components": components,
		},
	}

	b, _ := json.Marshal(payload)
	url := fmt.Sprintf("https://graph.facebook.com/%s/%s/messages", c.cfg.APIVersion, c.cfg.PhoneNumberID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.cfg.AccessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("meta send failed (%d): %s", resp.StatusCode, string(body))
	}

	var out struct {
		Messages []struct {
			ID string `json:"id"`
		} `json:"messages"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return "", err
	}
	if len(out.Messages) == 0 || out.Messages[0].ID == "" {
		return "", fmt.Errorf("meta response missing message id")
	}

	return out.Messages[0].ID, nil
}
