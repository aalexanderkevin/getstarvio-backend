package internaladmin

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string                `json:"accessToken"`
	RefreshToken string                `json:"refreshToken"`
	Admin        InternalAdminBasicDTO `json:"admin"`
}

type InternalAdminBasicDTO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type RefreshResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type DefaultCategoryListItem struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Category     string  `json:"category"`
	Status       string  `json:"status"`
	Icon         *string `json:"icon"`
	Interval     *int    `json:"interval"`
	TemplateID   string  `json:"templateId"`
	TemplateBody string  `json:"templateBody"`
	ExampleBody  string  `json:"exampleBody"`
	IsActive     bool    `json:"isActive"`
}

type CreateDefaultCategoryRequest struct {
	Name         string  `json:"name"`
	Category     string  `json:"category"`
	Status       *string `json:"status"`
	Icon         *string `json:"icon"`
	Interval     *int    `json:"interval"`
	TemplateID   string  `json:"templateId"`
	TemplateBody string  `json:"templateBody"`
	ExampleBody  string  `json:"exampleBody"`
	IsActive     *bool   `json:"isActive"`
}

type PlanConfigResponse struct {
	BusinessID string `json:"businessId"`
	FreeBonus  int    `json:"freeBonus"`
	SubCredits int    `json:"subCredits"`
	SubPrice   int    `json:"subPrice"`
	TopupPrice int    `json:"topupPrice"`
	Tier1Price int    `json:"tier1Price"`
	Tier1Creds int    `json:"tier1Credits"`
	Tier2Price int    `json:"tier2Price"`
	Tier2Creds int    `json:"tier2Credits"`
	Tier3Price int    `json:"tier3Price"`
	Tier3Creds int    `json:"tier3Credits"`
}

type WATemplateItem struct {
	ID                 string                             `json:"id"`
	MetaTemplateName   string                             `json:"metaTemplateName"`
	TemplateAlias      string                             `json:"templateAlias"`
	Category           string                             `json:"category"`
	Language           string                             `json:"language"`
	Status             string                             `json:"status"`
	Body               string                             `json:"body"`
	BodyExample        []string                           `json:"bodyExample"`
	BodyExamplePreview []WATemplateBodyExamplePreviewItem `json:"bodyExamplePreview"`
	MetaTemplateID     string                             `json:"metaTemplateId"`
	CreatedAt          string                             `json:"createdAt"`
	UpdatedAt          string                             `json:"updatedAt"`
}

type WATemplateBodyExamplePreviewItem struct {
	Key    string `json:"key"`
	Sample string `json:"sample"`
}

type CreateWATemplateRequest struct {
	MetaTemplateName string   `json:"metaTemplateName"`
	TemplateAlias    string   `json:"templateAlias"`
	Category         string   `json:"category"`
	Language         string   `json:"language"`
	Status           string   `json:"status"`
	Body             string   `json:"body"`
	BodyExample      []string `json:"bodyExample"`
}

type UpdateWATemplateRequest struct {
	MetaTemplateName *string   `json:"metaTemplateName"`
	TemplateAlias    *string   `json:"templateAlias"`
	Category         *string   `json:"category"`
	Language         *string   `json:"language"`
	Status           *string   `json:"status"`
	Body             *string   `json:"body"`
	BodyExample      *[]string `json:"bodyExample"`
	MetaTemplateID   *string   `json:"metaTemplateId"`
}

type WATemplateVariableOption struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Description string `json:"description"`
	Sample      string `json:"sample"`
}
