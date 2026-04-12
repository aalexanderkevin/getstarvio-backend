package category

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/aalexanderkevin/getstarvio-backend/internal/models"
	"github.com/aalexanderkevin/getstarvio-backend/internal/platform/meta"
)

type Service struct {
	repo *Repo
	meta *meta.Client
}

func NewService(repo *Repo, metaClient *meta.Client) *Service {
	return &Service{repo: repo, meta: metaClient}
}

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
			"id":             c.ID,
			"name":           c.Name,
			"icon":           c.Icon,
			"interval":       c.IntervalDays,
			"templateId":     c.TemplateID,
			"templateBody":   c.TemplateBody,
			"metaTemplateId": c.MetaTemplateID,
			"isEnabled":      c.IsEnabled,
		})
	}
	return out, nil
}

func (s *Service) Create(userID string, req CreateCategoryRequest) (map[string]interface{}, error) {
	if req.Name == "" || req.Interval <= 0 {
		return nil, fmt.Errorf("name and interval are required")
	}
	biz, err := s.repo.FindBusinessByUser(userID)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(biz.MetaWABAID) == "" {
		return nil, fmt.Errorf("business metaWabaId is required")
	}
	if strings.TrimSpace(biz.MetaAccessToken) == "" {
		return nil, fmt.Errorf("business metaAccessToken is required")
	}
	enabled := true
	if req.IsEnabled != nil {
		enabled = *req.IsEnabled
	}

	templateName := req.TemplateID
	if templateName == "" {
		templateName = defaultTemplateName(req.Name)
	}
	templateBody := req.TemplateBody
	if templateBody == "" {
		templateBody = "Halo {{1}}! Sudah {{2}} hari sejak {{3}} terakhir kamu di {{4}}. Yuk balik lagi — kami tunggu! 😊"
	}

	metaTemplate, err := s.meta.CreateTemplate(context.Background(), meta.CreateTemplateInput{
		Name:                templateName,
		WABAID:              biz.MetaWABAID,
		AccessToken:         biz.MetaAccessToken,
		Category:            "UTILITY",
		Language:            "id",
		BodyText:            templateBody,
		ExampleBodyTextVars: []string{"Pelanggan", strconv.Itoa(req.Interval), req.Name, biz.BizName},
	})
	if err != nil {
		return nil, err
	}
	if strings.EqualFold(metaTemplate.Status, "REJECTED") {
		return nil, fmt.Errorf("meta template rejected")
	}

	cat := models.Category{
		ID:             uuid.NewString(),
		BusinessID:     biz.ID,
		Name:           req.Name,
		Icon:           req.Icon,
		IntervalDays:   req.Interval,
		TemplateID:     templateName,
		TemplateBody:   templateBody,
		IsEnabled:      enabled,
		MetaTemplateID: metaTemplate.ID,
	}
	if cat.Icon == "" {
		cat.Icon = "✨"
	}
	if err := s.repo.Create(cat); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ok":             true,
		"categoryId":     cat.ID,
		"metaTemplateId": metaTemplate.ID,
		"metaStatus":     metaTemplate.Status,
	}, nil
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

func defaultTemplateName(categoryName string) string {
	s := strings.ToLower(strings.TrimSpace(categoryName))
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			return r
		}
		return -1
	}, s)
	if s == "" {
		s = "category_reminder"
	}
	return fmt.Sprintf("%s_%d", s, time.Now().Unix())
}
