package billing

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/aalexanderkevin/getstarvio-backend/internal/models"
	"github.com/aalexanderkevin/getstarvio-backend/internal/platform/xendit"
)

type Service struct {
	repo   *Repo
	xendit *xendit.Client
}

func NewService(repo *Repo, x *xendit.Client) *Service {
	return &Service{repo: repo, xendit: x}
}

func (s *Service) ValidateWebhookToken(token string) bool {
	return s.xendit.ValidateCallbackToken(token)
}

func (s *Service) Summary(userID string) (map[string]interface{}, error) {
	biz, err := s.repo.FindBusinessByUser(userID)
	if err != nil {
		return nil, err
	}
	wallet, err := s.repo.FindWalletByBusiness(biz.ID)
	if err != nil {
		return nil, err
	}
	planCfg, err := s.repo.FindPlanConfig(biz.ID)
	if err != nil {
		return nil, err
	}

	plan := "free"
	if wallet.SubscriptionStatus == models.SubscriptionStatusActive {
		plan = "subscriber"
	}

	return map[string]interface{}{
		"plan":               plan,
		"subscriptionStatus": wallet.SubscriptionStatus,
		"trialEndsAt":        wallet.TrialEndsAt.Format(time.RFC3339),
		"subscriptionEndsAt": nullableTime(wallet.SubscriptionEnds),
		"welcomeCreditsLeft": wallet.WelcomeCreditsLeft,
		"subCreditsLeft":     wallet.SubCreditsLeft,
		"topupCreditsLeft":   wallet.TopupCreditsLeft,
		"subCreditsMax":      wallet.SubCreditsMax,
		"remLeft":            wallet.WelcomeCreditsLeft + wallet.SubCreditsLeft + wallet.TopupCreditsLeft,
		"planConfig": map[string]interface{}{
			"freeBonus":  planCfg.FreeBonus,
			"subCredits": planCfg.SubCredits,
			"subPrice":   planCfg.SubPrice,
			"topupPrice": planCfg.TopupPrice,
			"tiers": []map[string]int{
				{"price": planCfg.Tier1Price, "credits": planCfg.Tier1Credits},
				{"price": planCfg.Tier2Price, "credits": planCfg.Tier2Credits},
				{"price": planCfg.Tier3Price, "credits": planCfg.Tier3Credits},
			},
		},
	}, nil
}

func (s *Service) History(userID string) ([]map[string]interface{}, error) {
	biz, err := s.repo.FindBusinessByUser(userID)
	if err != nil {
		return nil, err
	}
	txs, err := s.repo.ListTransactions(biz.ID, 300)
	if err != nil {
		return nil, err
	}
	out := make([]map[string]interface{}, 0, len(txs))
	for _, tx := range txs {
		out = append(out, map[string]interface{}{
			"id":           tx.ID,
			"type":         tx.Type,
			"label":        tx.Label,
			"delta":        tx.Delta,
			"balanceAfter": tx.BalanceAfter,
			"note":         tx.Note,
			"createdAt":    tx.CreatedAt.Format(time.RFC3339),
		})
	}
	return out, nil
}

func (s *Service) ActivateSubscription(userID string) error {
	biz, err := s.repo.FindBusinessByUser(userID)
	if err != nil {
		return err
	}
	wallet, err := s.repo.FindWalletByBusiness(biz.ID)
	if err != nil {
		return err
	}
	planCfg, err := s.repo.FindPlanConfig(biz.ID)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	ends := now.AddDate(0, 0, 30)
	newSub := wallet.SubCreditsLeft + planCfg.SubCredits
	newBal := wallet.WelcomeCreditsLeft + newSub + wallet.TopupCreditsLeft

	if err := s.repo.UpdateWallet(biz.ID, map[string]interface{}{
		"subscription_status":     models.SubscriptionStatusActive,
		"subscription_started_at": now,
		"subscription_ends_at":    ends,
		"sub_credits_left":        newSub,
		"sub_credits_max":         planCfg.SubCredits,
	}); err != nil {
		return err
	}

	return s.repo.InsertTransaction(models.BillingTransaction{
		ID: uuid.NewString(), BusinessID: biz.ID, Type: "subscription", Label: "Subscription",
		Delta: planCfg.SubCredits, BalanceAfter: newBal,
		Note: fmt.Sprintf("Aktivasi subscriber — Rp %d/bulan (%d kredit)", planCfg.SubPrice, planCfg.SubCredits),
	})
}

func (s *Service) CancelSubscription(userID string) error {
	biz, err := s.repo.FindBusinessByUser(userID)
	if err != nil {
		return err
	}
	wallet, err := s.repo.FindWalletByBusiness(biz.ID)
	if err != nil {
		return err
	}

	newBal := wallet.WelcomeCreditsLeft + wallet.TopupCreditsLeft
	if err := s.repo.UpdateWallet(biz.ID, map[string]interface{}{
		"subscription_status":  models.SubscriptionStatusCancelled,
		"subscription_ends_at": time.Now().UTC(),
		"sub_credits_left":     0,
	}); err != nil {
		return err
	}

	return s.repo.InsertTransaction(models.BillingTransaction{
		ID: uuid.NewString(), BusinessID: biz.ID, Type: "subscription", Label: "Cancel Subscription", Delta: -wallet.SubCreditsLeft,
		BalanceAfter: newBal, Note: "Langganan dibatalkan",
	})
}

func (s *Service) CreateTopupCheckout(userID string, req CheckoutRequest) (map[string]interface{}, error) {
	biz, err := s.repo.FindBusinessByUser(userID)
	if err != nil {
		return nil, err
	}
	wallet, err := s.repo.FindWalletByBusiness(biz.ID)
	if err != nil {
		return nil, err
	}
	if wallet.SubscriptionStatus != models.SubscriptionStatusActive {
		return nil, fmt.Errorf("top-up hanya untuk subscriber aktif")
	}
	plan, err := s.repo.FindPlanConfig(biz.ID)
	if err != nil {
		return nil, err
	}

	pkgID := req.PackageID
	if pkgID == "" {
		pkgID = "p1"
	}
	amount, credits, err := pickPackage(plan, pkgID)
	if err != nil {
		return nil, err
	}

	extID := "topup-" + uuid.NewString()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	inv, err := s.xendit.CreateInvoice(
		ctx,
		xendit.CreateInvoiceInput{
			ExternalID:  extID,
			Amount:      amount,
			PayerEmail:  biz.AdminEmail,
			Description: fmt.Sprintf("Top up %d kredit", credits),
		},
	)
	if err != nil {
		return nil, err
	}

	order := models.TopupOrder{
		ID: uuid.NewString(), BusinessID: biz.ID, ExternalID: extID, InvoiceID: inv.InvoiceID,
		PackageID: pkgID, AmountIDR: amount, Credits: credits, Status: strings.ToLower(inv.Status), CheckoutURL: inv.InvoiceURL, RawPayload: inv.RawResponse,
	}
	if err := s.repo.CreateTopupOrder(order); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"orderId":     order.ID,
		"externalId":  extID,
		"invoiceId":   inv.InvoiceID,
		"checkoutUrl": inv.InvoiceURL,
		"status":      inv.Status,
	}, nil
}

func (s *Service) HandleXenditWebhook(payload XenditWebhookPayload, raw string) error {
	var order *models.TopupOrder
	var err error
	if payload.ExternalID != "" {
		order, err = s.repo.FindTopupOrderByExternalID(payload.ExternalID)
	} else {
		order, err = s.repo.FindTopupOrderByInvoiceID(payload.ID)
	}
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}

	st := strings.ToLower(payload.Status)
	if st != "paid" && st != "settled" {
		return nil
	}
	if strings.ToLower(order.Status) == "paid" {
		return nil
	}

	paidAt := time.Now().UTC()
	if payload.PaidAt != "" {
		if t, err := time.Parse(time.RFC3339, payload.PaidAt); err == nil {
			paidAt = t.UTC()
		}
	}

	tx := models.BillingTransaction{
		ID: uuid.NewString(), BusinessID: order.BusinessID, Type: "topup", Label: "Top Up", Delta: order.Credits,
		Note: fmt.Sprintf("Top Up %d kredit", order.Credits), MetaJSON: raw,
	}
	return s.repo.ProcessPaidTopup(order.ID, order.BusinessID, order.Credits, raw, paidAt, tx)
}

func (s *Service) GetPlanConfig(userID string) (map[string]interface{}, error) {
	biz, err := s.repo.FindBusinessByUser(userID)
	if err != nil {
		return nil, err
	}
	p, err := s.repo.FindPlanConfig(biz.ID)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"freeBonus":  p.FreeBonus,
		"subCredits": p.SubCredits,
		"subPrice":   p.SubPrice,
		"topupPrice": p.TopupPrice,
		"tiers": []map[string]int{
			{"price": p.Tier1Price, "credits": p.Tier1Credits},
			{"price": p.Tier2Price, "credits": p.Tier2Credits},
			{"price": p.Tier3Price, "credits": p.Tier3Credits},
		},
	}, nil
}

func (s *Service) UpdatePlanConfig(userID string, data map[string]interface{}) error {
	biz, err := s.repo.FindBusinessByUser(userID)
	if err != nil {
		return err
	}
	payload := map[string]interface{}{}
	for _, k := range []string{"free_bonus", "sub_credits", "sub_price", "topup_price", "tier1_price", "tier1_credits", "tier2_price", "tier2_credits", "tier3_price", "tier3_credits"} {
		if v, ok := data[k]; ok {
			payload[k] = v
		}
	}
	if len(payload) == 0 {
		return fmt.Errorf("empty payload")
	}
	return s.repo.UpdatePlanConfig(biz.ID, payload)
}

func (s *Service) CanSendReminder(businessID string) (bool, string, error) {
	w, err := s.repo.FindWalletByBusiness(businessID)
	if err != nil {
		return false, "", err
	}
	now := time.Now().UTC()
	if now.After(w.TrialEndsAt) && w.SubscriptionStatus != models.SubscriptionStatusActive {
		return false, "trial ended, subscription required", nil
	}
	if w.WelcomeCreditsLeft+w.SubCreditsLeft+w.TopupCreditsLeft <= 0 {
		return false, "credits empty", nil
	}
	return true, "", nil
}

func (s *Service) DeductReminderCredit(businessID, note string) error {
	w, err := s.repo.FindWalletByBusiness(businessID)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	if now.After(w.TrialEndsAt) && w.SubscriptionStatus != models.SubscriptionStatusActive {
		return fmt.Errorf("cannot deduct credit: trial ended and no active subscription")
	}
	if w.WelcomeCreditsLeft+w.SubCreditsLeft+w.TopupCreditsLeft <= 0 {
		return fmt.Errorf("cannot deduct credit: credits empty")
	}

	payload := map[string]interface{}{}
	if w.WelcomeCreditsLeft > 0 {
		payload["welcome_credits_left"] = w.WelcomeCreditsLeft - 1
	} else if w.SubCreditsLeft > 0 {
		payload["sub_credits_left"] = w.SubCreditsLeft - 1
	} else {
		payload["topup_credits_left"] = w.TopupCreditsLeft - 1
	}

	if err := s.repo.UpdateWallet(businessID, payload); err != nil {
		return err
	}

	nw := w.WelcomeCreditsLeft
	ns := w.SubCreditsLeft
	nt := w.TopupCreditsLeft
	if v, ok := payload["welcome_credits_left"]; ok {
		nw = v.(int)
	}
	if v, ok := payload["sub_credits_left"]; ok {
		ns = v.(int)
	}
	if v, ok := payload["topup_credits_left"]; ok {
		nt = v.(int)
	}

	return s.repo.InsertTransaction(models.BillingTransaction{
		ID: uuid.NewString(), BusinessID: businessID, Type: "usage", Label: "Penggunaan", Delta: -1,
		BalanceAfter: nw + ns + nt,
		Note:         note,
	})
}

func pickPackage(p *models.PlanConfig, id string) (int, int, error) {
	switch id {
	case "p1":
		return p.Tier1Price, p.Tier1Credits, nil
	case "p2":
		return p.Tier2Price, p.Tier2Credits, nil
	case "p3":
		return p.Tier3Price, p.Tier3Credits, nil
	default:
		return 0, 0, fmt.Errorf("invalid packageId")
	}
}

func nullableTime(t *time.Time) interface{} {
	if t == nil {
		return nil
	}
	return t.Format(time.RFC3339)
}

func toJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
