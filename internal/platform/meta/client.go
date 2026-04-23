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
	logStore   FacebookLogStore
}

type SendTemplateInput struct {
	To           string
	TemplateName string
	LanguageCode string
	Parameters   []string
	AccessToken  string
	RefID        string
}

type CreateTemplateInput struct {
	Name                string
	WABAID              string
	Category            string
	Language            string
	BodyText            string
	ExampleBodyTextVars []string
	AccessToken         string
	RefID               string
}

type CreateTemplateResult struct {
	ID       string
	Status   string
	Category string
}

func NewClient(cfg config.MetaConfig, stores ...FacebookLogStore) *Client {
	timeoutSeconds := cfg.HTTPTimeoutSeconds
	if timeoutSeconds <= 0 {
		timeoutSeconds = 30
	}
	var store FacebookLogStore
	if len(stores) > 0 {
		store = stores[0]
	}

	return &Client{
		cfg:        cfg,
		httpClient: &http.Client{Timeout: time.Duration(timeoutSeconds) * time.Second},
		logStore:   store,
	}
}

func (c *Client) SendTemplate(ctx context.Context, in SendTemplateInput) (string, error) {
	if c.cfg.PhoneNumberID == "" {
		return "mock-" + ksuid.New().String(), nil
	}
	if in.AccessToken == "" {
		return "", fmt.Errorf("meta access token is required")
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
	req.Header.Set("Authorization", "Bearer "+in.AccessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logAPICall("send_message", url, string(b), err.Error(), 0, in.RefID)
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	c.logAPICall("send_message", url, string(b), string(body), resp.StatusCode, in.RefID)
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

func (c *Client) CreateTemplate(ctx context.Context, in CreateTemplateInput) (*CreateTemplateResult, error) {
	if in.WABAID == "" {
		return nil, fmt.Errorf("meta waba id is required")
	}
	if in.AccessToken == "" {
		return nil, fmt.Errorf("meta access token is required")
	}

	if in.Name == "" {
		return nil, fmt.Errorf("template name is required")
	}
	if in.Category == "" {
		in.Category = "UTILITY"
	}
	if in.Language == "" {
		in.Language = "id"
	}
	if in.BodyText == "" {
		return nil, fmt.Errorf("template body is required")
	}
	if len(in.ExampleBodyTextVars) == 0 {
		in.ExampleBodyTextVars = []string{"Pelanggan", "30", "Facial Treatment", "Celestial Spa & Wellness"}
	}

	payload := map[string]any{
		"name":             in.Name,
		"category":         in.Category,
		"language":         in.Language,
		"parameter_format": "POSITIONAL",
		"components": []map[string]any{
			{
				"type": "BODY",
				"text": in.BodyText,
				"example": map[string]any{
					"body_text": [][]string{in.ExampleBodyTextVars},
				},
			},
		},
	}

	b, _ := json.Marshal(payload)
	url := fmt.Sprintf("https://graph.facebook.com/%s/%s/message_templates", c.cfg.APIVersion, in.WABAID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+in.AccessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logAPICall("create_template", url, string(b), err.Error(), 0, in.RefID)
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	c.logAPICall("create_template", url, string(b), string(body), resp.StatusCode, in.RefID)
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("meta create template failed (%d): %s", resp.StatusCode, string(body))
	}

	var out struct {
		ID       string `json:"id"`
		Status   string `json:"status"`
		Category string `json:"category"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, err
	}
	if out.ID == "" {
		return nil, fmt.Errorf("meta create template response missing id")
	}

	return &CreateTemplateResult{
		ID:       out.ID,
		Status:   out.Status,
		Category: out.Category,
	}, nil
}

func (c *Client) logAPICall(operation, url, requestBody, responseBody string, responseCode int, refID string) {
	if c.logStore == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = c.logStore.CreateFacebookLog(ctx, FacebookLogEntry{
		Operation:    operation,
		URL:          url,
		RequestBody:  requestBody,
		ResponseBody: responseBody,
		ResponseCode: responseCode,
		RefID:        refID,
	})
}
