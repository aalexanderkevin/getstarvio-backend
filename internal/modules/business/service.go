package business

import (
	"fmt"
	"strings"
	"time"

	"github.com/aalexanderkevin/getstarvio-backend/internal/models"
	"github.com/aalexanderkevin/getstarvio-backend/internal/modules/shared"
)

type Service struct{ repo *Repo }

func NewService(repo *Repo) *Service { return &Service{repo: repo} }

func (s *Service) GetBootstrap(userID string) (map[string]interface{}, error) {
	biz, err := s.repo.FindBusinessByUser(userID)
	if err != nil {
		return nil, err
	}
	settings, err := s.repo.FindSettings(biz.ID)
	if err != nil {
		return nil, err
	}
	wallet, err := s.repo.FindWallet(biz.ID)
	if err != nil {
		return nil, err
	}
	planCfg, _ := s.repo.FindPlanConfig(biz.ID)
	cats, err := s.repo.ListCategories(biz.ID)
	if err != nil {
		return nil, err
	}
	customers, err := s.repo.ListCustomers(biz.ID)
	if err != nil {
		return nil, err
	}

	ids := make([]string, 0, len(customers))
	for _, c := range customers {
		ids = append(ids, c.ID)
	}
	services, err := s.repo.ListCustomerServices(ids)
	if err != nil {
		return nil, err
	}
	serviceMap := map[string][]map[string]interface{}{}
	for _, svc := range services {
		serviceMap[svc.CustomerID] = append(serviceMap[svc.CustomerID], map[string]interface{}{
			"name": svc.ServiceName,
			"icon": svc.ServiceIcon,
			"date": svc.LastVisitAt.Format(time.RFC3339),
			"days": svc.IntervalDays,
		})
	}

	rems, err := s.repo.ListReminders(biz.ID, 1000)
	if err != nil {
		return nil, err
	}

	catDTO := make([]map[string]interface{}, 0, len(cats))
	for _, c := range cats {
		catDTO = append(catDTO, map[string]interface{}{
			"id":           c.ID,
			"name":         c.Name,
			"icon":         c.Icon,
			"interval":     c.IntervalDays,
			"templateId":   c.TemplateID,
			"templateBody": c.TemplateBody,
			"metaTemplateId": c.MetaTemplateID,
			"isEnabled":    c.IsEnabled,
		})
	}

	cxDTO := make([]map[string]interface{}, 0, len(customers))
	for _, cx := range customers {
		cxDTO = append(cxDTO, map[string]interface{}{
			"id":       cx.ID,
			"name":     cx.Name,
			"wa":       cx.WA,
			"via":      cx.Via,
			"services": serviceMap[cx.ID],
		})
	}

	remDTO := make([]map[string]interface{}, 0, len(rems))
	for _, r := range rems {
		var sentAt interface{}
		if r.SentAt != nil {
			sentAt = r.SentAt.Format(time.RFC3339)
		} else {
			sentAt = nil
		}
		remDTO = append(remDTO, map[string]interface{}{
			"id":          r.ID,
			"cxId":        r.CustomerID,
			"cxName":      r.CxName,
			"svc":         r.SvcName,
			"scheduledAt": r.ScheduledAt.Format(time.RFC3339),
			"sentAt":      sentAt,
			"status":      r.Status,
			"kredit":      r.Kredit,
		})
	}

	plan := "free"
	if wallet.SubscriptionStatus == models.SubscriptionStatusActive {
		plan = "subscriber"
	}

	var subRenewsAt interface{}
	if wallet.SubscriptionEnds != nil {
		subRenewsAt = wallet.SubscriptionEnds.Format(time.RFC3339)
	} else {
		subRenewsAt = nil
	}

	planConfigDTO := map[string]interface{}{}
	if planCfg != nil {
		planConfigDTO = map[string]interface{}{
			"freeBonus":  planCfg.FreeBonus,
			"subCredits": planCfg.SubCredits,
			"subPrice":   planCfg.SubPrice,
			"topupPrice": planCfg.TopupPrice,
			"tiers": []map[string]int{
				{"price": planCfg.Tier1Price, "credits": planCfg.Tier1Credits},
				{"price": planCfg.Tier2Price, "credits": planCfg.Tier2Credits},
				{"price": planCfg.Tier3Price, "credits": planCfg.Tier3Credits},
			},
		}
	}

	return map[string]interface{}{
		"DATA_VERSION":       4,
		"bizName":            biz.BizName,
		"bizType":            biz.BizType,
		"bizSlug":            biz.BizSlug,
		"adminName":          biz.AdminName,
		"adminEmail":         biz.AdminEmail,
		"ownerWa":            biz.OwnerWA,
		"waNum":              biz.WANum,
		"metaWabaId":         biz.MetaWABAID,
		"metaAccessTokenConfigured": strings.TrimSpace(biz.MetaAccessToken) != "",
		"timezone":           biz.Timezone,
		"country":            biz.Country,
		"plan":               plan,
		"subCreditsLeft":     wallet.SubCreditsLeft,
		"subCreditsMax":      wallet.SubCreditsMax,
		"topupCreditsLeft":   wallet.TopupCreditsLeft,
		"subRenewsAt":        subRenewsAt,
		"remLeft":            wallet.WelcomeCreditsLeft + wallet.SubCreditsLeft + wallet.TopupCreditsLeft,
		"remMax":             wallet.SubCreditsMax,
		"defaultInterval":    settings.DefaultInterval,
		"automationEnabled":  settings.AutomationEnabled,
		"sendTime":           settings.SendTime,
		"autoTopup":          settings.AutoTopupEnabled,
		"autoTopupThreshold": settings.AutoTopupThreshold,
		"autoTopupPackage":   settings.AutoTopupPackageID,
		"billingNotifs": map[string]bool{
			"lowCredit":       settings.BillingNotifLow,
			"criticalCredit":  settings.BillingNotifCritical,
			"subLow":          settings.BillingNotifSubLow,
			"renewalReminder": settings.BillingNotifPreRenew,
		},
		"cats":       catDTO,
		"customers":  cxDTO,
		"reminders":  remDTO,
		"planConfig": planConfigDTO,
	}, nil
}

func (s *Service) UpdateProfile(userID string, req UpdateProfileRequest) error {
	biz, err := s.repo.FindBusinessByUser(userID)
	if err != nil {
		return err
	}

	payload := map[string]interface{}{}
	if req.BizName != "" {
		payload["biz_name"] = req.BizName
	}
	if req.BizType != "" {
		payload["biz_type"] = req.BizType
	}
	if req.BizSlug != "" {
		payload["biz_slug"] = slugify(req.BizSlug)
	}
	if req.Timezone != "" {
		payload["timezone"] = req.Timezone
	}
	if req.Country != "" {
		payload["country"] = strings.ToUpper(req.Country)
	}
	if len(payload) == 0 {
		return nil
	}
	return s.repo.UpdateProfile(biz.ID, payload)
}

func (s *Service) UpdateWhatsApp(userID string, req UpdateWhatsAppRequest) error {
	biz, err := s.repo.FindBusinessByUser(userID)
	if err != nil {
		return err
	}
	owner := biz.OwnerWA
	wanum := biz.WANum
	metaWabaID := biz.MetaWABAID
	metaAccessToken := biz.MetaAccessToken
	if req.OwnerWA != "" {
		owner = shared.NormalizePhone(req.OwnerWA, "62")
	}
	if req.WANum != "" {
		wanum = shared.NormalizePhone(req.WANum, "62")
	}
	if req.MetaWABAID != "" {
		metaWabaID = strings.TrimSpace(req.MetaWABAID)
	}
	if req.MetaAccessToken != "" {
		metaAccessToken = strings.TrimSpace(req.MetaAccessToken)
	}
	return s.repo.UpdateWhatsApp(biz.ID, owner, wanum, metaWabaID, metaAccessToken)
}

func (s *Service) UpdateSettings(userID string, req UpdateSettingsRequest) error {
	biz, err := s.repo.FindBusinessByUser(userID)
	if err != nil {
		return err
	}
	payload := map[string]interface{}{}
	if req.AutomationEnabled != nil {
		payload["automation_enabled"] = *req.AutomationEnabled
	}
	if req.DefaultInterval != nil && *req.DefaultInterval > 0 {
		payload["default_interval"] = *req.DefaultInterval
	}
	if req.SendTime != "" {
		payload["send_time"] = req.SendTime
	}
	if req.Timezone != "" {
		payload["timezone"] = req.Timezone
	}
	if req.BillingNotifLow != nil {
		payload["billing_notif_low"] = *req.BillingNotifLow
	}
	if req.BillingNotifCritical != nil {
		payload["billing_notif_critical"] = *req.BillingNotifCritical
	}
	if req.BillingNotifSubLow != nil {
		payload["billing_notif_sub_low"] = *req.BillingNotifSubLow
	}
	if req.BillingNotifPreRenewal != nil {
		payload["billing_notif_pre_renewal"] = *req.BillingNotifPreRenewal
	}
	if req.AutoTopupEnabled != nil {
		payload["auto_topup_enabled"] = *req.AutoTopupEnabled
	}
	if req.AutoTopupThreshold != nil && *req.AutoTopupThreshold > 0 {
		payload["auto_topup_threshold"] = *req.AutoTopupThreshold
	}
	if req.AutoTopupPackageID != "" {
		payload["auto_topup_package_id"] = req.AutoTopupPackageID
	}
	if len(payload) == 0 {
		return fmt.Errorf("no settings payload")
	}
	return s.repo.UpdateSettings(biz.ID, payload)
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
