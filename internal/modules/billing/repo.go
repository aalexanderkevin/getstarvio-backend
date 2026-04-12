package billing

import (
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/aalexanderkevin/getstarvio-backend/internal/models"
)

type Repo struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) FindBusinessByUser(userID string) (*models.Business, error) {
	var b models.Business
	err := r.db.Where("user_id = ?", userID).First(&b).Error
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *Repo) FindWalletByBusiness(businessID string) (*models.Wallet, error) {
	var w models.Wallet
	err := r.db.Where("business_id = ?", businessID).First(&w).Error
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *Repo) FindPlanConfig(businessID string) (*models.PlanConfig, error) {
	var p models.PlanConfig
	err := r.db.Where("business_id = ?", businessID).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repo) ListTransactions(businessID string, limit int) ([]models.BillingTransaction, error) {
	if limit <= 0 {
		limit = 200
	}
	var out []models.BillingTransaction
	err := r.db.Where("business_id = ?", businessID).Order("created_at desc").Limit(limit).Find(&out).Error
	return out, err
}

func (r *Repo) UpdateWallet(businessID string, payload map[string]interface{}) error {
	payload["updated_at"] = time.Now().UTC()
	return r.db.Model(&models.Wallet{}).Where("business_id = ?", businessID).Updates(payload).Error
}

func (r *Repo) InsertTransaction(tx models.BillingTransaction) error {
	return r.db.Create(&tx).Error
}

func (r *Repo) CreateTopupOrder(o models.TopupOrder) error {
	return r.db.Create(&o).Error
}

func (r *Repo) FindTopupOrderByExternalID(externalID string) (*models.TopupOrder, error) {
	var o models.TopupOrder
	err := r.db.Where("external_id = ?", externalID).First(&o).Error
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *Repo) FindTopupOrderByInvoiceID(invoiceID string) (*models.TopupOrder, error) {
	var o models.TopupOrder
	err := r.db.Where("invoice_id = ?", invoiceID).First(&o).Error
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *Repo) FindBusinessByID(businessID string) (*models.Business, error) {
	var b models.Business
	err := r.db.Where("id = ?", businessID).First(&b).Error
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *Repo) UpdatePlanConfig(businessID string, p map[string]interface{}) error {
	p["updated_at"] = time.Now().UTC()
	return r.db.Model(&models.PlanConfig{}).Where("business_id = ?", businessID).Updates(p).Error
}

func (r *Repo) ProcessPaidTopup(orderID string, businessID string, credits int, rawPayload string, paidAt time.Time, tx models.BillingTransaction) error {
	return r.db.Transaction(func(db *gorm.DB) error {
		var order models.TopupOrder
		if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", orderID).First(&order).Error; err != nil {
			return err
		}
		if strings.ToLower(order.Status) == "paid" {
			return nil
		}

		if err := db.Model(&models.TopupOrder{}).Where("id = ?", orderID).Updates(map[string]interface{}{
			"status":      "paid",
			"paid_at":     paidAt,
			"raw_payload": rawPayload,
			"updated_at":  time.Now().UTC(),
		}).Error; err != nil {
			return err
		}

		var w models.Wallet
		if err := db.Clauses(clause.Locking{Strength: "UPDATE"}).Where("business_id = ?", businessID).First(&w).Error; err != nil {
			return err
		}
		newTopup := w.TopupCreditsLeft + credits
		if err := db.Model(&models.Wallet{}).Where("business_id = ?", businessID).Updates(map[string]interface{}{
			"topup_credits_left": newTopup,
			"updated_at":         time.Now().UTC(),
		}).Error; err != nil {
			return err
		}

		tx.BalanceAfter = w.WelcomeCreditsLeft + w.SubCreditsLeft + newTopup
		if err := db.Create(&tx).Error; err != nil {
			return err
		}

		return nil
	})
}
