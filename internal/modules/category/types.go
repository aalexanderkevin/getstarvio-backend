package category

type CreateCategoryRequest struct {
	Name         string `json:"name"`
	Icon         string `json:"icon"`
	Interval     int    `json:"interval"`
	TemplateID   string `json:"templateId"`
	TemplateBody string `json:"templateBody"`
	IsEnabled    *bool  `json:"isEnabled"`
}

type UpdateCategoryRequest struct {
	Name         *string `json:"name"`
	Icon         *string `json:"icon"`
	Interval     *int    `json:"interval"`
	TemplateID   *string `json:"templateId"`
	TemplateBody *string `json:"templateBody"`
	IsEnabled    *bool   `json:"isEnabled"`
}
