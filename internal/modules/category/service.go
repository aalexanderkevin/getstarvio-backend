package category

import (
	"context"
	"encoding/json"
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
			"exampleBody":  c.ExampleBody,
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
	if req.Name == "" || req.Interval <= 0 || strings.TrimSpace(req.DefaultCategoryID) == "" {
		return nil, fmt.Errorf("name, interval, and defaultCategoryId are required")
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

	defCat, err := s.repo.FindDefaultByID(strings.TrimSpace(req.DefaultCategoryID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("default category not found")
		}
		return nil, err
	}

	templateName := defCat.TemplateID
	if templateName == "" {
		templateName = defaultTemplateName(req.Name)
	}
	templateName = concatTemplateNameWithBusiness(templateName, biz.BizName)

	templateBody := strings.TrimSpace(defCat.TemplateBody)
	if templateBody == "" {
		return nil, fmt.Errorf("default category templateBody is empty")
	}

	exampleVars, err := parseExampleBodyVars(defCat.ExampleBody, req.Name, req.Interval, biz.BizName)
	if err != nil {
		return nil, err
	}

	metaTemplate, err := s.meta.CreateTemplate(context.Background(), meta.CreateTemplateInput{
		Name:                templateName,
		WABAID:              biz.MetaWABAID,
		AccessToken:         biz.MetaAccessToken,
		Category:            "UTILITY",
		Language:            "id",
		BodyText:            templateBody,
		ExampleBodyTextVars: exampleVars,
	})
	if err != nil {
		return nil, err
	}
	if strings.EqualFold(metaTemplate.Status, "REJECTED") {
		return nil, fmt.Errorf("meta template rejected")
	}

	icon := strings.TrimSpace(req.Icon)
	if icon == "" {
		icon = strings.TrimSpace(defCat.Icon)
	}
	if icon == "" {
		icon = "✨"
	}

	cat := models.Category{
		ID:             uuid.NewString(),
		BusinessID:     biz.ID,
		Name:           req.Name,
		Icon:           icon,
		IntervalDays:   req.Interval,
		TemplateID:     templateName,
		TemplateBody:   templateBody,
		IsEnabled:      true,
		MetaTemplateID: metaTemplate.ID,
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
	if req.Icon != nil {
		payload["icon"] = *req.Icon
	}
	if req.Interval != nil {
		payload["interval_days"] = *req.Interval
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
	s := normalizeTemplateNamePart(categoryName)
	if s == "" {
		s = "category_reminder"
	}
	return fmt.Sprintf("%s_%d", s, time.Now().Unix())
}

func concatTemplateNameWithBusiness(templateName, businessName string) string {
	base := normalizeTemplateNamePart(templateName)
	if base == "" {
		base = "category_reminder"
	}
	biz := normalizeTemplateNamePart(businessName)
	if biz == "" {
		return base
	}
	return fmt.Sprintf("%s_%s", base, biz)
}

func normalizeTemplateNamePart(v string) string {
	v = strings.ToLower(strings.TrimSpace(v))
	v = strings.ReplaceAll(v, " ", "_")
	v = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			return r
		}
		return -1
	}, v)
	v = strings.Trim(v, "_")
	for strings.Contains(v, "__") {
		v = strings.ReplaceAll(v, "__", "_")
	}
	return v
}

func parseExampleBodyVars(raw, serviceName string, interval int, businessName string) ([]string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil, fmt.Errorf("default category exampleBody is empty")
	}

	var vars []string
	if err := json.Unmarshal([]byte(trimmed), &vars); err != nil {
		return nil, fmt.Errorf("invalid default category exampleBody: %w", err)
	}
	if len(vars) == 0 {
		return nil, fmt.Errorf("default category exampleBody has no values")
	}

	out := make([]string, 0, len(vars))
	for _, v := range vars {
		vv := strings.TrimSpace(v)
		switch strings.ToLower(vv) {
		case "{{interval}}":
			vv = strconv.Itoa(interval)
		case "{{service}}":
			vv = serviceName
		case "{{business}}":
			vv = businessName
		}
		out = append(out, vv)
	}
	return out, nil
}
