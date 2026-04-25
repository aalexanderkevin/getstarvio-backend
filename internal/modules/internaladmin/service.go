package internaladmin

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/aalexanderkevin/getstarvio-backend/internal/config"
	"github.com/aalexanderkevin/getstarvio-backend/internal/models"
	"github.com/aalexanderkevin/getstarvio-backend/internal/modules/shared"
	"github.com/aalexanderkevin/getstarvio-backend/internal/platform/meta"
)

type Service struct {
	repo *Repo
	cfg  config.Config
	meta *meta.Client
}

func NewService(repo *Repo, cfg config.Config, metaClient *meta.Client) *Service {
	return &Service{repo: repo, cfg: cfg, meta: metaClient}
}

func (s *Service) GetPlanConfig() (PlanConfigResponse, error) {
	p, err := s.repo.GetPrimaryPlanConfig()
	if err != nil {
		return PlanConfigResponse{}, err
	}

	return PlanConfigResponse{
		BusinessID: p.BusinessID,
		FreeBonus:  p.FreeBonus,
		SubCredits: p.SubCredits,
		SubPrice:   p.SubPrice,
		TopupPrice: p.TopupPrice,
		Tier1Price: p.Tier1Price,
		Tier1Creds: p.Tier1Credits,
		Tier2Price: p.Tier2Price,
		Tier2Creds: p.Tier2Credits,
		Tier3Price: p.Tier3Price,
		Tier3Creds: p.Tier3Credits,
	}, nil
}

func (s *Service) UpdatePlanConfig(data map[string]interface{}) error {
	p, err := s.repo.GetPrimaryPlanConfig()
	if err != nil {
		return err
	}

	payload := map[string]interface{}{}
	assignIfInt(payload, "free_bonus", data, "free_bonus", "freeBonus")
	assignIfInt(payload, "sub_credits", data, "sub_credits", "subCredits")
	assignIfInt(payload, "sub_price", data, "sub_price", "subPrice")
	assignIfInt(payload, "topup_price", data, "topup_price", "topupPrice")
	assignIfInt(payload, "tier1_price", data, "tier1_price", "tier1Price")
	assignIfInt(payload, "tier1_credits", data, "tier1_credits", "tier1Credits")
	assignIfInt(payload, "tier2_price", data, "tier2_price", "tier2Price")
	assignIfInt(payload, "tier2_credits", data, "tier2_credits", "tier2Credits")
	assignIfInt(payload, "tier3_price", data, "tier3_price", "tier3Price")
	assignIfInt(payload, "tier3_credits", data, "tier3_credits", "tier3Credits")

	if len(payload) == 0 {
		return fmt.Errorf("empty payload")
	}

	return s.repo.UpdatePlanConfig(p.BusinessID, payload)
}

const (
	internalAccessTokenType  = "internal_admin"
	internalRefreshTokenType = "internal_admin_refresh"
)

var validWATemplateCategories = map[string]struct{}{
	"UTILITY":        {},
	"MARKETING":      {},
	"AUTHENTICATION": {},
}

var validWATemplateLanguages = map[string]struct{}{
	"id":    {},
	"en_US": {},
	"ms_MY": {},
}

var validWATemplateStatuses = map[string]struct{}{
	"DRAFT":    {},
	"PENDING":  {},
	"APPROVED": {},
	"REJECTED": {},
	"PAUSED":   {},
	"FLAGGED":  {},
}

type internalAccessClaims struct {
	InternalAdminID string `json:"internal_admin_id"`
	TokenType       string `json:"token_type"`
	jwt.RegisteredClaims
}

type internalRefreshClaims struct {
	InternalAdminID string `json:"internal_admin_id"`
	TokenType       string `json:"token_type"`
	jwt.RegisteredClaims
}

func (s *Service) Login(req LoginRequest) (*LoginResponse, error) {
	email := strings.TrimSpace(req.Email)
	password := req.Password
	if email == "" || password == "" {
		return nil, fmt.Errorf("email and password are required")
	}

	admin, err := s.repo.FindAdminByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}
	if !admin.IsActive {
		return nil, fmt.Errorf("admin account is disabled")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	accessToken, refreshToken, err := s.issueTokens(admin.ID)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	_ = s.repo.TouchAdminLastLogin(admin.ID, now)

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Admin: InternalAdminBasicDTO{
			ID:    admin.ID,
			Name:  admin.Name,
			Email: admin.Email,
		},
	}, nil
}

func (s *Service) Refresh(req RefreshRequest) (*RefreshResponse, error) {
	if req.RefreshToken == "" {
		return nil, fmt.Errorf("refresh token required")
	}

	claims := &internalRefreshClaims{}
	parsed, err := jwt.ParseWithClaims(req.RefreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWT.Secret), nil
	})
	if err != nil || !parsed.Valid {
		return nil, fmt.Errorf("invalid refresh token")
	}
	if claims.TokenType != internalRefreshTokenType || claims.InternalAdminID == "" {
		return nil, fmt.Errorf("invalid refresh token")
	}

	th := shared.HashToken(req.RefreshToken)
	rt, err := s.repo.FindInternalRefreshToken(th)
	if err != nil {
		return nil, fmt.Errorf("refresh token not found")
	}
	if rt.RevokedAt != nil || rt.ExpiresAt.Before(time.Now().UTC()) {
		return nil, fmt.Errorf("refresh token expired")
	}

	if err := s.repo.RevokeInternalRefreshToken(th, time.Now().UTC()); err != nil {
		return nil, err
	}

	admin, err := s.repo.FindAdminByID(claims.InternalAdminID)
	if err != nil || !admin.IsActive {
		return nil, fmt.Errorf("admin not found or inactive")
	}

	accessToken, refreshToken, err := s.issueTokens(admin.ID)
	if err != nil {
		return nil, err
	}
	return &RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) Logout(req LogoutRequest) error {
	if req.RefreshToken == "" {
		return nil
	}
	return s.repo.RevokeInternalRefreshToken(shared.HashToken(req.RefreshToken), time.Now().UTC())
}

func (s *Service) ListDefaultCategories() ([]DefaultCategoryListItem, error) {
	cats, err := s.repo.ListDefaultCategories()
	if err != nil {
		return nil, err
	}
	out := make([]DefaultCategoryListItem, 0, len(cats))
	for _, c := range cats {
		out = append(out, DefaultCategoryListItem{
			ID:           c.ID,
			Name:         c.Name,
			Category:     c.Category,
			Status:       c.Status,
			Icon:         c.Icon,
			Interval:     c.IntervalDays,
			TemplateID:   c.TemplateID,
			TemplateBody: c.TemplateBody,
			ExampleBody:  c.ExampleBody,
			IsActive:     c.IsActive,
		})
	}
	return out, nil
}

func (s *Service) CreateDefaultCategory(req CreateDefaultCategoryRequest) (map[string]interface{}, error) {
	name := strings.TrimSpace(req.Name)
	category := strings.TrimSpace(req.Category)
	status := "PENDING"
	if req.Status != nil {
		trimmedStatus := strings.ToUpper(strings.TrimSpace(*req.Status))
		if trimmedStatus != "" {
			status = trimmedStatus
		}
	}
	templateID := strings.TrimSpace(req.TemplateID)
	templateBody := strings.TrimSpace(req.TemplateBody)
	exampleBody := strings.TrimSpace(req.ExampleBody)
	if name == "" || templateID == "" || templateBody == "" || exampleBody == "" {
		return nil, fmt.Errorf("name, templateId, templateBody, and exampleBody are required")
	}
	if category == "" {
		category = "UTILITY"
	}
	if req.Interval != nil && *req.Interval <= 0 {
		return nil, fmt.Errorf("interval must be greater than 0")
	}

	var vars []string
	if err := json.Unmarshal([]byte(exampleBody), &vars); err != nil || len(vars) == 0 {
		return nil, fmt.Errorf("exampleBody must be valid JSON string array")
	}

	var icon *string
	if req.Icon != nil {
		trimmed := strings.TrimSpace(*req.Icon)
		if trimmed != "" {
			icon = &trimmed
		}
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	m := models.DefaultCategory{
		ID:           uuid.NewString(),
		Name:         name,
		Category:     category,
		Status:       status,
		Icon:         icon,
		IntervalDays: req.Interval,
		TemplateID:   templateID,
		TemplateBody: templateBody,
		ExampleBody:  exampleBody,
		IsActive:     isActive,
	}
	if err := s.repo.CreateDefaultCategory(m); err != nil {
		return nil, err
	}
	return map[string]interface{}{"ok": true, "id": m.ID}, nil
}

func (s *Service) ListWATemplates(category, status, metaTemplateName string) ([]WATemplateItem, error) {
	category = strings.ToUpper(strings.TrimSpace(category))
	status = strings.ToUpper(strings.TrimSpace(status))
	metaTemplateName = strings.TrimSpace(metaTemplateName)

	rows, err := s.repo.ListWATemplates(category, status, metaTemplateName)
	if err != nil {
		return nil, err
	}
	out := make([]WATemplateItem, 0, len(rows))
	for _, row := range rows {
		item, err := toWATemplateItem(row)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, nil
}

func (s *Service) GetWATemplate(id string) (*WATemplateItem, error) {
	row, err := s.repo.FindWATemplateByID(strings.TrimSpace(id))
	if err != nil {
		return nil, err
	}
	item, err := toWATemplateItem(*row)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *Service) CreateWATemplate(req CreateWATemplateRequest) (map[string]interface{}, error) {
	metaTemplateName := strings.TrimSpace(req.MetaTemplateName)
	templateAlias := strings.TrimSpace(req.TemplateAlias)
	category, err := normalizeWATemplateCategory(req.Category)
	if err != nil {
		return nil, err
	}
	language, err := normalizeWATemplateLanguage(req.Language)
	if err != nil {
		return nil, err
	}
	status, err := normalizeWATemplateStatus(req.Status)
	if err != nil {
		return nil, err
	}
	if status != "DRAFT" && status != "PENDING" {
		return nil, fmt.Errorf("status on create must be DRAFT or PENDING")
	}
	body := strings.TrimSpace(req.Body)
	if metaTemplateName == "" || templateAlias == "" || body == "" {
		return nil, fmt.Errorf("metaTemplateName, templateAlias, and body are required")
	}
	bodyExampleJSON, bodyExampleKeys, err := marshalBodyExample(req.BodyExample)
	if err != nil {
		return nil, err
	}
	if err := validateTemplatePlaceholders(body, len(bodyExampleKeys)); err != nil {
		return nil, err
	}

	row := models.WATemplate{
		ID:               uuid.NewString(),
		MetaTemplateName: metaTemplateName,
		TemplateAlias:    templateAlias,
		Category:         category,
		Language:         language,
		Status:           status,
		Body:             body,
		BodyExample:      bodyExampleJSON,
	}

	if status == "PENDING" {
		if s.meta == nil {
			return nil, fmt.Errorf("meta client not configured")
		}
		wabaID := strings.TrimSpace(s.cfg.Meta.WABAID)
		accessToken := strings.TrimSpace(s.cfg.Meta.AccessToken)
		if wabaID == "" || accessToken == "" {
			return nil, fmt.Errorf("meta credentials are not configured in env")
		}

		metaRes, err := s.meta.CreateTemplate(context.Background(), meta.CreateTemplateInput{
			Name:                metaTemplateName,
			WABAID:              wabaID,
			AccessToken:         accessToken,
			Category:            category,
			Language:            language,
			BodyText:            body,
			ExampleBodyTextVars: templateVariableSamples(bodyExampleKeys),
			RefID:               row.ID,
		})
		if err != nil {
			return nil, err
		}
		row.MetaTemplateID = strings.TrimSpace(metaRes.ID)
		if row.MetaTemplateID == "" {
			return nil, fmt.Errorf("meta create template response missing id")
		}
		metaStatus := strings.ToUpper(strings.TrimSpace(metaRes.Status))
		if _, ok := validWATemplateStatuses[metaStatus]; ok {
			row.Status = metaStatus
		}
	}

	if err := s.repo.CreateWATemplate(row); err != nil {
		return nil, err
	}
	return map[string]interface{}{"ok": true, "id": row.ID}, nil
}

func (s *Service) UpdateWATemplate(id string, req UpdateWATemplateRequest) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("id is required")
	}
	row, err := s.repo.FindWATemplateByID(id)
	if err != nil {
		return err
	}

	payload := map[string]interface{}{}
	var existingKeys []string
	if err := json.Unmarshal([]byte(row.BodyExample), &existingKeys); err != nil {
		return fmt.Errorf("invalid bodyExample in db: %w", err)
	}

	nextStatus := strings.ToUpper(strings.TrimSpace(row.Status))
	effectiveMetaTemplateName := strings.TrimSpace(row.MetaTemplateName)
	effectiveCategory := strings.TrimSpace(row.Category)
	effectiveLanguage := strings.TrimSpace(row.Language)
	effectiveBody := strings.TrimSpace(row.Body)
	effectiveBodyExampleKeys := append([]string(nil), existingKeys...)

	if req.MetaTemplateName != nil {
		v := strings.TrimSpace(*req.MetaTemplateName)
		if v == "" {
			return fmt.Errorf("metaTemplateName cannot be empty")
		}
		payload["meta_template_name"] = v
		effectiveMetaTemplateName = v
	}
	if req.TemplateAlias != nil {
		v := strings.TrimSpace(*req.TemplateAlias)
		if v == "" {
			return fmt.Errorf("templateAlias cannot be empty")
		}
		payload["template_alias"] = v
	}
	if req.Category != nil {
		v, err := normalizeWATemplateCategory(*req.Category)
		if err != nil {
			return err
		}
		payload["category"] = v
		effectiveCategory = v
	}
	if req.Language != nil {
		v, err := normalizeWATemplateLanguage(*req.Language)
		if err != nil {
			return err
		}
		payload["language"] = v
		effectiveLanguage = v
	}
	if req.Status != nil {
		v, err := normalizeWATemplateStatus(*req.Status)
		if err != nil {
			return err
		}
		payload["status"] = v
		nextStatus = v
	}
	if req.Body != nil {
		v := strings.TrimSpace(*req.Body)
		if v == "" {
			return fmt.Errorf("body cannot be empty")
		}
		payload["body"] = v
		effectiveBody = v
	}
	if req.BodyExample != nil {
		v, keys, err := marshalBodyExample(*req.BodyExample)
		if err != nil {
			return err
		}
		payload["body_example"] = v
		effectiveBodyExampleKeys = keys
	}
	if req.MetaTemplateID != nil {
		payload["meta_template_id"] = strings.TrimSpace(*req.MetaTemplateID)
	}

	if err := validateTemplatePlaceholders(effectiveBody, len(effectiveBodyExampleKeys)); err != nil {
		return err
	}

	if strings.EqualFold(strings.TrimSpace(row.Status), "DRAFT") && nextStatus == "PENDING" {
		if s.meta == nil {
			return fmt.Errorf("meta client not configured")
		}
		wabaID := strings.TrimSpace(s.cfg.Meta.WABAID)
		accessToken := strings.TrimSpace(s.cfg.Meta.AccessToken)
		if wabaID == "" || accessToken == "" {
			return fmt.Errorf("meta credentials are not configured in env")
		}

		metaRes, err := s.meta.CreateTemplate(context.Background(), meta.CreateTemplateInput{
			Name:                effectiveMetaTemplateName,
			WABAID:              wabaID,
			AccessToken:         accessToken,
			Category:            effectiveCategory,
			Language:            effectiveLanguage,
			BodyText:            effectiveBody,
			ExampleBodyTextVars: templateVariableSamples(effectiveBodyExampleKeys),
			RefID:               row.ID,
		})
		if err != nil {
			return err
		}
		metaTemplateID := strings.TrimSpace(metaRes.ID)
		if metaTemplateID == "" {
			return fmt.Errorf("meta create template response missing id")
		}
		payload["meta_template_id"] = metaTemplateID
		metaStatus := strings.ToUpper(strings.TrimSpace(metaRes.Status))
		if _, ok := validWATemplateStatuses[metaStatus]; ok {
			payload["status"] = metaStatus
		}
	}

	if len(payload) == 0 {
		return nil
	}
	return s.repo.UpdateWATemplate(id, payload)
}

func (s *Service) DeleteWATemplate(id string) error {
	return s.repo.DeleteWATemplate(strings.TrimSpace(id))
}

func (s *Service) ListWATemplateVariables() []WATemplateVariableOption {
	opts := shared.ListTemplateVariableOptions()
	out := make([]WATemplateVariableOption, 0, len(opts))
	for _, opt := range opts {
		out = append(out, WATemplateVariableOption{
			Key:         opt.Key,
			Label:       opt.Label,
			Description: opt.Description,
			Sample:      opt.Sample,
		})
	}
	return out
}

func toWATemplateItem(row models.WATemplate) (WATemplateItem, error) {
	var bodyExample []string
	if err := json.Unmarshal([]byte(row.BodyExample), &bodyExample); err != nil {
		return WATemplateItem{}, fmt.Errorf("invalid bodyExample in db: %w", err)
	}
	preview := make([]WATemplateBodyExamplePreviewItem, 0, len(bodyExample))
	for _, key := range bodyExample {
		sample, ok := shared.TemplateVariableSampleForKey(key)
		if !ok {
			sample = ""
		}
		preview = append(preview, WATemplateBodyExamplePreviewItem{
			Key:    key,
			Sample: sample,
		})
	}
	return WATemplateItem{
		ID:                 row.ID,
		MetaTemplateName:   row.MetaTemplateName,
		TemplateAlias:      row.TemplateAlias,
		Category:           row.Category,
		Language:           row.Language,
		Status:             row.Status,
		Body:               row.Body,
		BodyExample:        bodyExample,
		BodyExamplePreview: preview,
		MetaTemplateID:     row.MetaTemplateID,
		CreatedAt:          row.CreatedAt.Format(time.RFC3339),
		UpdatedAt:          row.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func marshalBodyExample(input []string) (string, []string, error) {
	keys, err := shared.NormalizeTemplateVariableKeys(input)
	if err != nil {
		return "", nil, err
	}
	b, err := json.Marshal(keys)
	if err != nil {
		return "", nil, err
	}
	return string(b), keys, nil
}

var templatePlaceholderPattern = regexp.MustCompile(`\{\{(\d+)\}\}`)

func validateTemplatePlaceholders(body string, mappingLen int) error {
	matches := templatePlaceholderPattern.FindAllStringSubmatch(body, -1)
	max := 0
	for _, m := range matches {
		if len(m) < 2 {
			continue
		}
		n, err := strconv.Atoi(m[1])
		if err != nil {
			continue
		}
		if n > max {
			max = n
		}
	}
	if max == 0 {
		return nil
	}
	if max != mappingLen {
		return fmt.Errorf("body placeholders count (%d) must match bodyExample keys count (%d)", max, mappingLen)
	}
	return nil
}

func templateVariableSamples(keys []string) []string {
	out := make([]string, 0, len(keys))
	for _, key := range keys {
		sample, ok := shared.TemplateVariableSampleForKey(key)
		if !ok {
			out = append(out, "")
			continue
		}
		out = append(out, sample)
	}
	return out
}

func normalizeWATemplateCategory(v string) (string, error) {
	c := strings.ToUpper(strings.TrimSpace(v))
	if _, ok := validWATemplateCategories[c]; !ok {
		return "", fmt.Errorf("invalid category")
	}
	return c, nil
}

func normalizeWATemplateLanguage(v string) (string, error) {
	lang := strings.TrimSpace(v)
	switch strings.ToLower(strings.ReplaceAll(lang, "-", "_")) {
	case "id":
		lang = "id"
	case "en_us":
		lang = "en_US"
	case "ms_my":
		lang = "ms_MY"
	default:
		return "", fmt.Errorf("invalid language")
	}
	if _, ok := validWATemplateLanguages[lang]; !ok {
		return "", fmt.Errorf("invalid language")
	}
	return lang, nil
}

func normalizeWATemplateStatus(v string) (string, error) {
	status := strings.ToUpper(strings.TrimSpace(v))
	if _, ok := validWATemplateStatuses[status]; !ok {
		return "", fmt.Errorf("invalid status")
	}
	return status, nil
}

func (s *Service) issueTokens(adminID string) (string, string, error) {
	now := time.Now().UTC()
	accExp := now.Add(time.Duration(s.cfg.JWT.AccessTTLMinutes) * time.Minute)
	refExp := now.Add(time.Duration(s.cfg.JWT.RefreshTTLHours) * time.Hour)

	accClaims := internalAccessClaims{
		InternalAdminID: adminID,
		TokenType:       internalAccessTokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   adminID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(accExp),
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accClaims).SignedString([]byte(s.cfg.JWT.Secret))
	if err != nil {
		return "", "", err
	}

	refClaims := internalRefreshClaims{
		InternalAdminID: adminID,
		TokenType:       internalRefreshTokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   adminID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(refExp),
		},
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refClaims).SignedString([]byte(s.cfg.JWT.Secret))
	if err != nil {
		return "", "", err
	}

	if err := s.repo.SaveInternalRefreshToken(models.InternalRefreshToken{
		ID:        uuid.NewString(),
		AdminID:   adminID,
		TokenHash: shared.HashToken(refreshToken),
		ExpiresAt: refExp,
	}); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func assignIfInt(dst map[string]interface{}, dstKey string, src map[string]interface{}, keys ...string) {
	for _, key := range keys {
		v, ok := src[key]
		if !ok {
			continue
		}
		n, err := toInt(v)
		if err != nil {
			continue
		}
		dst[dstKey] = n
		return
	}
}

func toInt(v interface{}) (int, error) {
	switch t := v.(type) {
	case int:
		return t, nil
	case int32:
		return int(t), nil
	case int64:
		return int(t), nil
	case float32:
		return int(t), nil
	case float64:
		return int(t), nil
	case string:
		n, err := strconv.Atoi(t)
		if err != nil {
			return 0, err
		}
		return n, nil
	default:
		return 0, fmt.Errorf("not numeric")
	}
}
