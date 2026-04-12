package billing

type CheckoutRequest struct {
	PackageID string `json:"packageId"`
}

type XenditWebhookPayload struct {
	ID         string `json:"id"`
	ExternalID string `json:"external_id"`
	Status     string `json:"status"`
	PaidAt     string `json:"paid_at"`
}
