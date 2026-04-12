package xendit

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/segmentio/ksuid"

	"github.com/aalexanderkevin/getstarvio-backend/internal/config"
)

type Client struct {
	cfg        config.XenditConfig
	httpClient *http.Client
}

type CreateInvoiceInput struct {
	ExternalID         string
	Amount             int
	PayerEmail         string
	Description        string
	SuccessRedirectURL string
	FailureRedirectURL string
}

type CreateInvoiceOutput struct {
	InvoiceID   string
	ExternalID  string
	InvoiceURL  string
	Status      string
	RawResponse string
}

func NewClient(cfg config.XenditConfig) *Client {
	return &Client{cfg: cfg, httpClient: &http.Client{Timeout: 10 * time.Second}}
}

func (c *Client) CreateInvoice(ctx context.Context, in CreateInvoiceInput) (CreateInvoiceOutput, error) {
	if c.cfg.APIKey == "" {
		return CreateInvoiceOutput{
			InvoiceID:   "mock-" + ksuid.New().String(),
			ExternalID:  in.ExternalID,
			InvoiceURL:  "https://checkout.mock.xendit.local/" + in.ExternalID,
			Status:      "PENDING",
			RawResponse: "{}",
		}, nil
	}

	payload := map[string]any{
		"external_id": in.ExternalID,
		"amount":      in.Amount,
		"description": in.Description,
	}
	if in.PayerEmail != "" {
		payload["payer_email"] = in.PayerEmail
	}
	successURL := in.SuccessRedirectURL
	if successURL == "" {
		successURL = c.cfg.SuccessRedirect
	}
	if successURL != "" {
		payload["success_redirect_url"] = successURL
	}
	failureURL := in.FailureRedirectURL
	if failureURL == "" {
		failureURL = c.cfg.FailureRedirect
	}
	if failureURL != "" {
		payload["failure_redirect_url"] = failureURL
	}

	b, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.xendit.co/v2/invoices", bytes.NewBuffer(b))
	if err != nil {
		return CreateInvoiceOutput{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	cred := base64.StdEncoding.EncodeToString([]byte(c.cfg.APIKey + ":"))
	req.Header.Set("Authorization", "Basic "+cred)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return CreateInvoiceOutput{}, err
	}
	defer resp.Body.Close()

	rb, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return CreateInvoiceOutput{}, fmt.Errorf("xendit create invoice failed (%d): %s", resp.StatusCode, string(rb))
	}

	var out struct {
		ID         string `json:"id"`
		ExternalID string `json:"external_id"`
		InvoiceURL string `json:"invoice_url"`
		Status     string `json:"status"`
	}
	if err := json.Unmarshal(rb, &out); err != nil {
		return CreateInvoiceOutput{}, err
	}

	return CreateInvoiceOutput{
		InvoiceID:   out.ID,
		ExternalID:  out.ExternalID,
		InvoiceURL:  out.InvoiceURL,
		Status:      out.Status,
		RawResponse: string(rb),
	}, nil
}

func (c *Client) ValidateCallbackToken(token string) bool {
	if c.cfg.CallbackToken == "" {
		return true
	}
	return token == c.cfg.CallbackToken
}
