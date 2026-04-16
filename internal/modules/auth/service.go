package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/aalexanderkevin/getstarvio-backend/internal/config"
	"github.com/aalexanderkevin/getstarvio-backend/internal/models"
	"github.com/aalexanderkevin/getstarvio-backend/internal/modules/shared"
)

type Service struct {
	repo *Repo
	cfg  config.Config
	hc   *http.Client
}

func NewService(repo *Repo, cfg config.Config) *Service {
	return &Service{repo: repo, cfg: cfg, hc: &http.Client{Timeout: 10 * time.Second}}
}

func (s *Service) LoginWithGoogle(ctx context.Context, req GoogleLoginRequest) (*TokenResponse, error) {
	id, err := s.verifyGoogleIdentity(ctx, req)
	if err != nil {
		return nil, err
	}

	u, err := s.repo.FindUserByGoogleSub(id.Sub)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			now := time.Now().UTC()
			userID := uuid.NewString()
			bizID := uuid.NewString()
			slug := slugify(id.Name)
			if slug == "" {
				slug = "biz-" + strings.ToLower(strings.ReplaceAll(userID[:8], "-", ""))
			}

			user := models.User{ID: userID, GoogleSub: id.Sub, Email: id.Email, Name: id.Name}
			biz := models.Business{
				ID: bizID, UserID: userID, BizName: "Bisnis Saya", BizType: "lainnya", BizSlug: slug,
				AdminName: id.Name, AdminEmail: id.Email, Timezone: "Asia/Jakarta", Country: "ID",
			}
			settings := models.BusinessSettings{
				ID: uuid.NewString(), BusinessID: bizID, Timezone: "Asia/Jakarta", AutomationEnabled: true,
				DefaultInterval: 30, SendTime: "09:00",
				BillingNotifLow: true, BillingNotifCritical: true, BillingNotifSubLow: true, BillingNotifPreRenew: true,
				AutoTopupEnabled: false, AutoTopupThreshold: 10, AutoTopupPackageID: "p1",
			}
			wallet := models.Wallet{
				ID: uuid.NewString(), BusinessID: bizID, TrialStartedAt: now, TrialEndsAt: now.AddDate(0, 0, 20),
				SubscriptionStatus: models.SubscriptionStatusNone, WelcomeCreditsLeft: 100, SubCreditsLeft: 0, TopupCreditsLeft: 0, SubCreditsMax: 250,
			}
			plan := models.PlanConfig{
				ID: uuid.NewString(), BusinessID: bizID, FreeBonus: 100, SubCredits: 250, SubPrice: 250000,
				TopupPrice: 1000, Tier1Price: 250000, Tier1Credits: 300, Tier2Price: 500000, Tier2Credits: 625, Tier3Price: 1000000, Tier3Credits: 1500,
			}
			tx := models.BillingTransaction{
				ID: uuid.NewString(), BusinessID: bizID, Type: "welcome", Label: "Welcome Bonus", Delta: 100,
				BalanceAfter: 100, Note: "Bonus kredit pendaftaran",
			}
			if err := s.repo.BootstrapNewUser(user, biz, settings, wallet, plan, tx); err != nil {
				return nil, err
			}
			u = &user
		} else {
			return nil, err
		}
	}

	return s.issueTokens(u.ID)
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	_ = ctx
	if refreshToken == "" {
		return nil, fmt.Errorf("refresh token required")
	}

	claims := &refreshClaims{}
	parsed, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWT.Secret), nil
	})
	if err != nil || !parsed.Valid {
		return nil, fmt.Errorf("invalid refresh token")
	}

	th := shared.HashToken(refreshToken)
	rt, err := s.repo.FindRefreshToken(th)
	if err != nil {
		return nil, fmt.Errorf("refresh token not found")
	}
	if rt.RevokedAt != nil || rt.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("refresh token expired")
	}

	if err := s.repo.RevokeRefreshToken(th, time.Now().UTC()); err != nil {
		return nil, err
	}

	_, err = s.repo.FindUserByID(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return s.issueTokens(claims.UserID)
}

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	_ = ctx
	if refreshToken == "" {
		return nil
	}
	th := shared.HashToken(refreshToken)
	return s.repo.RevokeRefreshToken(th, time.Now().UTC())
}

type accessClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type refreshClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func (s *Service) issueTokens(userID string) (*TokenResponse, error) {
	now := time.Now().UTC()
	accExp := now.Add(time.Duration(s.cfg.JWT.AccessTTLMinutes) * time.Minute)
	refExp := now.Add(time.Duration(s.cfg.JWT.RefreshTTLHours) * time.Hour)

	accClaims := accessClaims{UserID: userID, RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(accExp), IssuedAt: jwt.NewNumericDate(now), Subject: userID}}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accClaims).SignedString([]byte(s.cfg.JWT.Secret))
	if err != nil {
		return nil, err
	}

	refClaims := refreshClaims{UserID: userID, RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(refExp), IssuedAt: jwt.NewNumericDate(now), Subject: userID}}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refClaims).SignedString([]byte(s.cfg.JWT.Secret))
	if err != nil {
		return nil, err
	}

	if err := s.repo.SaveRefreshToken(models.RefreshToken{ID: uuid.NewString(), UserID: userID, TokenHash: shared.HashToken(refreshToken), ExpiresAt: refExp}); err != nil {
		return nil, err
	}

	return &TokenResponse{AccessToken: accessToken, RefreshToken: refreshToken, UserID: userID}, nil
}

func (s *Service) verifyGoogleIdentity(ctx context.Context, req GoogleLoginRequest) (*googleIdentity, error) {
	if req.IDToken == "" {
		return nil, fmt.Errorf("idToken is required")
	}

	u := "https://oauth2.googleapis.com/tokeninfo?id_token=" + url.QueryEscape(req.IDToken)
	hreq, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.hc.Do(hreq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("google token verification failed: %s", string(b))
	}

	var out struct {
		Sub           string `json:"sub"`
		Email         string `json:"email"`
		Name          string `json:"name"`
		EmailVerified string `json:"email_verified"`
		Aud           string `json:"aud"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	if out.Sub == "" || out.Email == "" {
		return nil, fmt.Errorf("invalid google identity payload")
	}
	if s.cfg.Google.ClientID != "" && out.Aud != "" && out.Aud != s.cfg.Google.ClientID {
		return nil, fmt.Errorf("google audience mismatch")
	}
	if out.Name == "" {
		out.Name = strings.Split(out.Email, "@")[0]
	}

	return &googleIdentity{Sub: out.Sub, Email: out.Email, Name: out.Name}, nil
}

func slugify(v string) string {
	v = strings.ToLower(strings.TrimSpace(v))
	v = strings.ReplaceAll(v, "_", "-")
	v = strings.ReplaceAll(v, " ", "-")
	v = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, v)
	for strings.Contains(v, "--") {
		v = strings.ReplaceAll(v, "--", "-")
	}
	return strings.Trim(v, "-")
}
