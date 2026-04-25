package reminder

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/aalexanderkevin/getstarvio-backend/internal/config"
	"github.com/aalexanderkevin/getstarvio-backend/internal/models"
	"github.com/aalexanderkevin/getstarvio-backend/internal/modules/shared"
	"github.com/aalexanderkevin/getstarvio-backend/internal/platform/meta"
	"gorm.io/gorm"
)

var ErrMetaWebhookUnauthorized = errors.New("meta webhook unauthorized")

type Service struct {
	repo *Repo
	meta *meta.Client
	cfg  config.MetaConfig
}

func NewService(repo *Repo, metaClient *meta.Client, cfg config.MetaConfig) *Service {
	return &Service{repo: repo, meta: metaClient, cfg: cfg}
}

func (s *Service) Log(userID, status string, limit int) ([]map[string]interface{}, error) {
	biz, err := s.repo.FindBusinessByUser(userID)
	if err != nil {
		return nil, err
	}
	rows, err := s.repo.ListReminderLogs(biz.ID, status, limit)
	if err != nil {
		return nil, err
	}

	out := make([]map[string]interface{}, 0, len(rows))
	for _, r := range rows {
		var sentAt interface{}
		if r.SentAt != nil {
			sentAt = r.SentAt.Format(time.RFC3339)
		}
		out = append(out, map[string]interface{}{
			"id":            r.ID,
			"businessId":    r.BusinessID,
			"customerId":    r.CustomerID,
			"categoryId":    r.CategoryID,
			"cxName":        r.CxName,
			"svcName":       r.SvcName,
			"categoryName":  r.SvcName,
			"scheduledAt":   r.ScheduledAt.Format(time.RFC3339),
			"sentAt":        sentAt,
			"status":        r.Status,
			"kredit":        r.Kredit,
			"errorReason":   r.ErrorReason,
			"retryCount":    r.RetryCount,
			"metaMessageId": r.MetaMessageID,
		})
	}
	return out, nil
}

func (s *Service) Retry(userID, reminderID string) error {
	biz, err := s.repo.FindBusinessByUser(userID)
	if err != nil {
		return err
	}
	return s.repo.RetryReminder(biz.ID, reminderID)
}

func (s *Service) DashboardSummary(userID string) (map[string]interface{}, error) {
	biz, err := s.repo.FindBusinessByUser(userID)
	if err != nil {
		return nil, err
	}

	settings, err := s.repo.FindSettingsByBusiness(biz.ID)
	if err != nil {
		return nil, err
	}
	loc, err := time.LoadLocation(settings.Timezone)
	if err != nil {
		loc = time.FixedZone("Asia/Jakarta", 7*60*60)
	}
	now := time.Now().UTC()
	from := shared.StartOfDay(now, loc).UTC()
	to := shared.EndOfDay(now, loc).UTC()

	totalCustomers, err := s.repo.CountCustomers(biz.ID)
	if err != nil {
		return nil, err
	}
	pending, err := s.repo.CountRemindersByStatus(biz.ID, models.ReminderStatusPending)
	if err != nil {
		return nil, err
	}
	todaySent, err := s.repo.CountSentBetween(biz.ID, from, to)
	if err != nil {
		return nil, err
	}
	todayFailed, err := s.repo.CountFailedBetween(biz.ID, from, to)
	if err != nil {
		return nil, err
	}
	wallet, err := s.repo.FindWalletByBusiness(biz.ID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"date":             now.In(loc).Format("2006-01-02"),
		"timezone":         settings.Timezone,
		"totalCustomers":   totalCustomers,
		"pendingReminders": pending,
		"sentToday":        todaySent,
		"failedToday":      todayFailed,
		"credits": map[string]int{
			"welcome":      wallet.WelcomeCreditsLeft,
			"subscription": wallet.SubCreditsLeft,
			"topup":        wallet.TopupCreditsLeft,
			"total":        wallet.WelcomeCreditsLeft + wallet.SubCreditsLeft + wallet.TopupCreditsLeft,
		},
		"trialEndsAt":        wallet.TrialEndsAt.Format(time.RFC3339),
		"subscriptionStatus": wallet.SubscriptionStatus,
	}, nil
}

func (s *Service) VerifyMetaWebhook(mode, verifyToken, challenge string) (string, error) {
	if mode != "subscribe" {
		return "", ErrMetaWebhookUnauthorized
	}
	if s.cfg.WebhookVerifyToken == "" {
		return "", fmt.Errorf("META_WEBHOOK_VERIFY_TOKEN is not configured")
	}
	if verifyToken != s.cfg.WebhookVerifyToken {
		return "", ErrMetaWebhookUnauthorized
	}
	if challenge == "" {
		return "", fmt.Errorf("hub.challenge is required")
	}
	return challenge, nil
}

func (s *Service) HandleMetaWebhook(raw []byte, signature string) error {
	if len(raw) == 0 {
		return fmt.Errorf("empty payload")
	}
	if err := s.validateMetaWebhookSignature(raw, signature); err != nil {
		return err
	}
	var payload MetaWebhookPayload
	if err := json.Unmarshal(raw, &payload); err != nil {
		return fmt.Errorf("invalid meta payload: %w", err)
	}
	if payload.Object != "whatsapp_business_account" {
		return fmt.Errorf("unsupported meta object: %s", payload.Object)
	}

	for _, entry := range payload.Entry {
		for _, change := range entry.Changes {
			if strings.TrimSpace(change.Field) != "message_template_status_update" {
				continue
			}

			var upd MetaTemplateStatusUpdate
			if err := json.Unmarshal(change.Value, &upd); err != nil {
				return fmt.Errorf("invalid message_template_status_update payload: %w", err)
			}

			metaTemplateID, err := parseMetaTemplateID(upd.MessageTemplateIDRaw)
			if err != nil {
				return fmt.Errorf("invalid message_template_id: %w", err)
			}

			status := strings.ToUpper(strings.TrimSpace(upd.Event))
			if err := s.repo.SetCategoryEnabledByMetaTemplateID(metaTemplateID, status); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Service) validateMetaWebhookSignature(raw []byte, header string) error {
	if strings.TrimSpace(s.cfg.AppSecret) == "" {
		return nil
	}
	h := strings.TrimSpace(header)
	if h == "" || !strings.HasPrefix(h, "sha256=") {
		return ErrMetaWebhookUnauthorized
	}

	mac := hmac.New(sha256.New, []byte(s.cfg.AppSecret))
	_, _ = mac.Write(raw)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(expected), []byte(h)) {
		return ErrMetaWebhookUnauthorized
	}
	return nil
}

func (s *Service) RunWorker(ctx context.Context, pollInterval time.Duration) error {
	if pollInterval <= 0 {
		pollInterval = 30 * time.Second
	}

	s.runCycle(ctx)

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.Canceled) {
				return nil
			}
			return ctx.Err()
		case <-ticker.C:
			s.runCycle(ctx)
		}
	}
}

func (s *Service) runCycle(ctx context.Context) {
	now := time.Now().UTC()
	enqueued, err := s.enqueueDueReminders(now)
	if err != nil {
		fmt.Printf("worker enqueue error: %v\n", err)
	}

	sent, failed, err := s.dispatchDueReminders(ctx, 200)
	if err != nil {
		fmt.Printf("worker dispatch error: %v\n", err)
	}
	if enqueued > 0 || sent > 0 || failed > 0 {
		fmt.Printf("worker cycle enqueued=%d sent=%d failed=%d\n", enqueued, sent, failed)
	}
}

func (s *Service) enqueueDueReminders(now time.Time) (int, error) {
	businesses, err := s.repo.ListBusinesses()
	if err != nil {
		return 0, err
	}

	created := 0
	for _, biz := range businesses {
		settings, err := s.repo.FindSettingsByBusiness(biz.ID)
		if err != nil {
			continue
		}
		loc, err := time.LoadLocation(settings.Timezone)
		if err != nil {
			loc = time.FixedZone("Asia/Jakarta", 7*60*60)
		}
		hour, minute := parseSendTime(settings.SendTime)

		services, err := s.repo.ListSchedulableServices(biz.ID)
		if err != nil {
			continue
		}

		for _, svc := range services {
			if svc.IntervalDays <= 0 {
				continue
			}
			if !svc.CategoryEnabled {
				continue
			}

			scheduledAt := dueAtWithSendTime(svc.LastVisitAt, svc.IntervalDays, hour, minute, loc)
			if scheduledAt.After(now) {
				continue
			}

			exists, err := s.repo.ReminderExists(biz.ID, svc.CustomerID, svc.ServiceName, scheduledAt)
			if err != nil || exists {
				continue
			}

			if err := s.repo.CreateReminder(models.Reminder{
				ID:          uuid.NewString(),
				BusinessID:  biz.ID,
				CustomerID:  svc.CustomerID,
				CategoryID:  svc.CategoryID,
				CxName:      svc.CustomerName,
				SvcName:     svc.ServiceName,
				ScheduledAt: scheduledAt,
				Status:      models.ReminderStatusPending,
				Kredit:      1,
			}); err == nil {
				created++
			}
		}
	}

	return created, nil
}

func (s *Service) dispatchDueReminders(ctx context.Context, limit int) (int, int, error) {
	pending, err := s.repo.ListDuePending(limit)
	if err != nil {
		return 0, 0, err
	}

	sent := 0
	failed := 0
	for _, rem := range pending {
		dctx, err := s.repo.GetDispatchContext(rem.ID)
		if err != nil {
			_ = s.repo.MarkReminderFailed(rem.ID, "dispatch context not found")
			failed++
			continue
		}

		wallet, err := s.repo.FindWalletByBusiness(rem.BusinessID)
		if err != nil {
			_ = s.repo.MarkReminderFailed(rem.ID, "wallet not found")
			failed++
			continue
		}

		if can, reason := canSend(*wallet, time.Now().UTC()); !can {
			_ = s.repo.MarkReminderFailed(rem.ID, reason)
			failed++
			continue
		}

		templateName := "reminder_return"
		if dctx.Category != nil && dctx.Category.TemplateID != "" {
			templateName = dctx.Category.TemplateID
		}
		templateVariables := []string{"customer_name", "days_since_last_visit", "service_name", "business_name"}
		if waTpl, err := s.repo.FindWATemplateByMetaTemplateName(templateName); err == nil {
			var fromDB []string
			if err := json.Unmarshal([]byte(waTpl.BodyExample), &fromDB); err == nil {
				if len(fromDB) > 0 {
					templateVariables = fromDB
				}
			}
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			_ = s.repo.MarkReminderFailed(rem.ID, fmt.Sprintf("fetch wa template failed: %v", err))
			failed++
			continue
		}

		parameters, err := resolveTemplateVariables(templateVariables, dctx, time.Now().UTC())
		if err != nil {
			_ = s.repo.MarkReminderFailed(rem.ID, err.Error())
			failed++
			continue
		}

		to := shared.NormalizePhone(dctx.Customer.PhoneNumber, "62")
		if to == "" {
			_ = s.repo.MarkReminderFailed(rem.ID, "customer wa is empty")
			failed++
			continue
		}

		sendCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		metaID, err := s.meta.SendTemplate(sendCtx, meta.SendTemplateInput{
			To:           to,
			TemplateName: templateName,
			LanguageCode: "id",
			Parameters:   parameters,
			AccessToken:  dctx.Business.MetaAccessToken,
			RefID:        rem.ID,
		})
		cancel()
		if err != nil {
			_ = s.repo.MarkReminderFailed(rem.ID, err.Error())
			failed++
			continue
		}

		if err := s.repo.MarkReminderSentAndDeduct(rem.ID, metaID); err != nil {
			failed++
			continue
		}
		sent++
	}

	return sent, failed, nil
}

func parseSendTime(v string) (int, int) {
	parts := strings.Split(v, ":")
	if len(parts) != 2 {
		return 9, 0
	}
	h, errH := strconv.Atoi(parts[0])
	m, errM := strconv.Atoi(parts[1])
	if errH != nil || errM != nil {
		return 9, 0
	}
	if h < 0 || h > 23 {
		h = 9
	}
	if m < 0 || m > 59 {
		m = 0
	}
	return h, m
}

func dueAtWithSendTime(lastVisit time.Time, intervalDays, sendHour, sendMinute int, loc *time.Location) time.Time {
	lv := lastVisit.In(loc)
	dueDate := lv.AddDate(0, 0, intervalDays)
	scheduledLocal := time.Date(
		dueDate.Year(),
		dueDate.Month(),
		dueDate.Day(),
		sendHour,
		sendMinute,
		0,
		0,
		loc,
	)
	return scheduledLocal.UTC()
}

func canSend(w models.Wallet, now time.Time) (bool, string) {
	if now.After(w.TrialEndsAt) && w.SubscriptionStatus != models.SubscriptionStatusActive {
		return false, "trial ended, subscription required"
	}
	if w.WelcomeCreditsLeft+w.SubCreditsLeft+w.TopupCreditsLeft <= 0 {
		return false, "credits empty"
	}
	return true, ""
}

func parseMetaTemplateID(raw json.RawMessage) (string, error) {
	if len(raw) == 0 {
		return "", fmt.Errorf("empty value")
	}

	var asString string
	if err := json.Unmarshal(raw, &asString); err == nil {
		asString = strings.TrimSpace(asString)
		if asString == "" {
			return "", fmt.Errorf("empty string")
		}
		return asString, nil
	}

	var asNumber json.Number
	if err := json.Unmarshal(raw, &asNumber); err == nil {
		s := strings.TrimSpace(asNumber.String())
		if s == "" {
			return "", fmt.Errorf("empty number")
		}
		return s, nil
	}

	return "", fmt.Errorf("unsupported type")
}

func resolveTemplateVariables(keys []string, dctx *DispatchContext, now time.Time) ([]string, error) {
	if len(keys) == 0 {
		return []string{}, nil
	}
	loc, err := time.LoadLocation(strings.TrimSpace(dctx.Business.Timezone))
	if err != nil || loc == nil {
		loc = time.FixedZone("Asia/Jakarta", 7*60*60)
	}
	var lastVisit *time.Time
	if dctx.Category != nil && dctx.Category.IntervalDays > 0 {
		t := dctx.Reminder.ScheduledAt.AddDate(0, 0, -dctx.Category.IntervalDays)
		lastVisit = &t
	}

	out := make([]string, 0, len(keys))
	for _, key := range keys {
		k := strings.TrimSpace(key)
		switch k {
		case "customer_name":
			v := strings.TrimSpace(dctx.Customer.Name)
			if v == "" {
				v = strings.TrimSpace(dctx.Reminder.CxName)
			}
			out = append(out, v)
		case "service_name":
			v := strings.TrimSpace(dctx.Reminder.SvcName)
			if v == "" && dctx.Category != nil {
				v = strings.TrimSpace(dctx.Category.Name)
			}
			out = append(out, v)
		case "business_name":
			out = append(out, strings.TrimSpace(dctx.Business.BizName))
		case "days_since_last_visit":
			if lastVisit == nil {
				return nil, fmt.Errorf("cannot resolve days_since_last_visit: last visit date unavailable")
			}
			days := int(now.In(loc).Sub(lastVisit.In(loc)).Hours() / 24)
			if days < 0 {
				days = 0
			}
			out = append(out, strconv.Itoa(days))
		case "last_visit_date":
			if lastVisit == nil {
				return nil, fmt.Errorf("cannot resolve last_visit_date: last visit date unavailable")
			}
			out = append(out, formatDateID(lastVisit.In(loc)))
		default:
			return nil, fmt.Errorf("unsupported template variable key: %s", k)
		}
	}
	return out, nil
}

func formatDateID(t time.Time) string {
	months := map[time.Month]string{
		time.January:   "Januari",
		time.February:  "Februari",
		time.March:     "Maret",
		time.April:     "April",
		time.May:       "Mei",
		time.June:      "Juni",
		time.July:      "Juli",
		time.August:    "Agustus",
		time.September: "September",
		time.October:   "Oktober",
		time.November:  "November",
		time.December:  "Desember",
	}
	return fmt.Sprintf("%02d %s %d", t.Day(), months[t.Month()], t.Year())
}
