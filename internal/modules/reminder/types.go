package reminder

type MetaWebhookPayload struct {
	Object string             `json:"object"`
	Entry  []MetaWebhookEntry `json:"entry"`
}

type MetaWebhookEntry struct {
	ID      string              `json:"id"`
	Changes []MetaWebhookChange `json:"changes"`
}

type MetaWebhookChange struct {
	Field string           `json:"field"`
	Value MetaWebhookValue `json:"value"`
}

type MetaWebhookValue struct {
	MessagingProduct string               `json:"messaging_product"`
	Metadata         MetaWebhookMetadata  `json:"metadata"`
	Messages         []MetaWebhookMessage `json:"messages"`
	Statuses         []MetaWebhookStatus  `json:"statuses"`
}

type MetaWebhookMetadata struct {
	DisplayPhoneNumber string `json:"display_phone_number"`
	PhoneNumberID      string `json:"phone_number_id"`
}

type MetaWebhookMessage struct {
	ID   string `json:"id"`
	From string `json:"from"`
	Type string `json:"type"`
}

type MetaWebhookStatus struct {
	ID          string `json:"id"`
	RecipientID string `json:"recipient_id"`
	Status      string `json:"status"`
	Timestamp   string `json:"timestamp"`
}
