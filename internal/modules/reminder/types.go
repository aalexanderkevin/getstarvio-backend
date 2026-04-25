package reminder

import "encoding/json"

type MetaWebhookPayload struct {
	Object string             `json:"object"`
	Entry  []MetaWebhookEntry `json:"entry"`
}

type MetaWebhookEntry struct {
	ID      string              `json:"id"`
	Changes []MetaWebhookChange `json:"changes"`
}

type MetaWebhookChange struct {
	Field string          `json:"field"`
	Value json.RawMessage `json:"value"`
}

type MetaWebhookValue struct {
	MessagingProduct string               `json:"messaging_product"`
	Metadata         MetaWebhookMetadata  `json:"metadata"`
	Messages         []MetaWebhookMessage `json:"messages"`
	Statuses         []MetaWebhookStatus  `json:"statuses"`
}

type MetaTemplateStatusUpdate struct {
	Event                string          `json:"event"`
	MessageTemplateIDRaw json.RawMessage `json:"message_template_id"`
	MessageTemplateName  string          `json:"message_template_name"`
	MessageTemplateLang  string          `json:"message_template_language"`
	MessageTemplateCat   string          `json:"message_template_category"`
	Reason               string          `json:"reason"`
}

type MetaTemplateCategoryUpdate struct {
	MessageTemplateIDRaw    json.RawMessage `json:"message_template_id"`
	MessageTemplateName     string          `json:"message_template_name"`
	MessageTemplateLang     string          `json:"message_template_language"`
	MessageNewCategory      string          `json:"new_category"`
	MessagePreviousCategory string          `json:"previous_category"`
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
