package category

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/aalexanderkevin/getstarvio-backend/internal/models"
)

type Service struct{ repo *Repo }

func NewService(repo *Repo) *Service { return &Service{repo: repo} }

func (s *Service) ListDefault() ([]map[string]interface{}, error) {
	cats, err := s.repo.ListDefaults()
	if err != nil {
		return nil, err
	}
	out := make([]map[string]interface{}, 0, len(cats))
	for _, c := range cats {
		out = append(out, map[string]interface{}{
			"id":           c.ID,
			"name":         c.Name,
			"icon":         c.Icon,
			"interval":     c.IntervalDays,
			"templateId":   c.TemplateID,
			"templateBody": c.TemplateBody,
		})
	}
	return out, nil
}

func (s *Service) List(userID string) ([]map[string]interface{}, error) {
	biz, err := s.repo.FindBusinessByUser(userID)
	if err != nil {
		return nil, err
	}
	cats, err := s.repo.List(biz.ID)
	if err != nil {
		return nil, err
	}
	out := make([]map[string]interface{}, 0, len(cats))
	for _, c := range cats {
		out = append(out, map[string]interface{}{
			"id":           c.ID,
			"name":         c.Name,
			"icon":         c.Icon,
			"interval":     c.IntervalDays,
			"templateId":   c.TemplateID,
			"templateBody": c.TemplateBody,
			"isEnabled":    c.IsEnabled,
		})
	}
	return out, nil
}

func (s *Service) Create(userID string, req CreateCategoryRequest) error {
	if req.Name == "" || req.Interval <= 0 {
		return fmt.Errorf("name and interval are required")
	}
	biz, err := s.repo.FindBusinessByUser(userID)
	if err != nil {
		return err
	}
	enabled := true
	if req.IsEnabled != nil {
		enabled = *req.IsEnabled
	}
	cat := models.Category{
		ID: uuid.NewString(), BusinessID: biz.ID, Name: req.Name, Icon: req.Icon,
		IntervalDays: req.Interval, TemplateID: req.TemplateID, TemplateBody: req.TemplateBody, IsEnabled: enabled,
	}
	if cat.Icon == "" {
		cat.Icon = "✨"
	}
	if cat.TemplateID == "" {
		cat.TemplateID = "reminder_return"
	}
	if cat.TemplateBody == "" {
		cat.TemplateBody = "Hai {customer_name}, sudah waktunya untuk {service} di {business_name}."
	}
	return s.repo.Create(cat)
}

func (s *Service) Update(userID, id string, req UpdateCategoryRequest) error {
	biz, err := s.repo.FindBusinessByUser(userID)
	if err != nil {
		return err
	}
	if _, err := s.repo.FindByID(biz.ID, id); err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("category not found")
		}
		return err
	}
	payload := map[string]interface{}{}
	if req.Name != nil {
		payload["name"] = *req.Name
	}
	if req.Icon != nil {
		payload["icon"] = *req.Icon
	}
	if req.Interval != nil {
		payload["interval_days"] = *req.Interval
	}
	if req.TemplateID != nil {
		payload["template_id"] = *req.TemplateID
	}
	if req.TemplateBody != nil {
		payload["template_body"] = *req.TemplateBody
	}
	if req.IsEnabled != nil {
		payload["is_enabled"] = *req.IsEnabled
	}
	if len(payload) == 0 {
		return nil
	}
	return s.repo.Update(biz.ID, id, payload)
}

func (s *Service) Delete(userID, id string) error {
	biz, err := s.repo.FindBusinessByUser(userID)
	if err != nil {
		return err
	}
	return s.repo.Delete(biz.ID, id)
}
