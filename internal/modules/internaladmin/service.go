package internaladmin

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/aalexanderkevin/getstarvio-backend/internal/config"
	"github.com/aalexanderkevin/getstarvio-backend/internal/models"
	"github.com/aalexanderkevin/getstarvio-backend/internal/modules/shared"
)

type Service struct {
	repo *Repo
	cfg  config.Config
}

func NewService(repo *Repo, cfg config.Config) *Service {
	return &Service{repo: repo, cfg: cfg}
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
