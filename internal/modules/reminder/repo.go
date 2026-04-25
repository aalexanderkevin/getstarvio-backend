package reminder

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/aalexanderkevin/getstarvio-backend/internal/models"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

type SchedulableServiceRow struct {
	BusinessID      string    `gorm:"column:business_id"`
	CustomerID      string    `gorm:"column:customer_id"`
	CustomerName    string    `gorm:"column:customer_name"`
	CustomerWA      string    `gorm:"column:customer_wa"`
	CategoryID      string    `gorm:"column:category_id"`
	ServiceName     string    `gorm:"column:service_name"`
	LastVisitAt     time.Time `gorm:"column:last_visit_at"`
	IntervalDays    int       `gorm:"column:interval_days"`
	TemplateID      string    `gorm:"column:template_id"`
	CategoryEnabled bool      `gorm:"column:category_enabled"`
}

type DispatchContext struct {
	Reminder models.Reminder
	Business models.Business
	Customer models.Customer
	Category *models.Category
}

func (r *Repo) FindBusinessByUser(userID string) (*models.Business, error) {
	var b models.Business
	err := r.db.Where("user_id = ?", userID).First(&b).Error
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *Repo) FindSettingsByBusiness(businessID string) (*models.BusinessSettings, error) {
	var s models.BusinessSettings
	err := r.db.Where("business_id = ?", businessID).First(&s).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *Repo) ListReminderLogs(businessID, status string, limit int) ([]models.Reminder, error) {
	if limit <= 0 {
		limit = 200
	}
	dbq := r.db.Where("business_id = ?", businessID)
	if status != "" {
		dbq = dbq.Where("status = ?", strings.ToLower(status))
	}
	var out []models.Reminder
	err := dbq.Order("scheduled_at desc").Limit(limit).Find(&out).Error
	return out, err
}

func (r *Repo) RetryReminder(businessID, reminderID string) error {
	now := time.Now().UTC()
	return r.db.Transaction(func(tx *gorm.DB) error {
		var rem models.Reminder
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("business_id = ? AND id = ?", businessID, reminderID).First(&rem).Error
		if err != nil {
			return err
		}
		if rem.Status == models.ReminderStatusSent {
			return fmt.Errorf("cannot retry sent reminder")
		}
		return tx.Model(&models.Reminder{}).Where("id = ?", rem.ID).Updates(map[string]interface{}{
			"status":       models.ReminderStatusPending,
			"scheduled_at": now,
			"error_reason": "",
			"updated_at":   now,
		}).Error
	})
}

func (r *Repo) CountCustomers(businessID string) (int64, error) {
	var n int64
	err := r.db.Model(&models.Customer{}).Where("business_id = ?", businessID).Count(&n).Error
	return n, err
}

func (r *Repo) CountRemindersByStatus(businessID, status string) (int64, error) {
	var n int64
	err := r.db.Model(&models.Reminder{}).Where("business_id = ? AND status = ?", businessID, status).Count(&n).Error
	return n, err
}

func (r *Repo) CountRemindersByStatusBetween(businessID, status string, fromUTC, toUTC time.Time) (int64, error) {
	var n int64
	err := r.db.Model(&models.Reminder{}).
		Where("business_id = ? AND status = ? AND scheduled_at >= ? AND scheduled_at <= ?", businessID, status, fromUTC, toUTC).
		Count(&n).Error
	return n, err
}

func (r *Repo) CountSentBetween(businessID string, fromUTC, toUTC time.Time) (int64, error) {
	var n int64
	err := r.db.Model(&models.Reminder{}).
		Where("business_id = ? AND status = ? AND sent_at >= ? AND sent_at <= ?", businessID, models.ReminderStatusSent, fromUTC, toUTC).
		Count(&n).Error
	return n, err
}

func (r *Repo) CountFailedBetween(businessID string, fromUTC, toUTC time.Time) (int64, error) {
	var n int64
	err := r.db.Model(&models.Reminder{}).
		Where("business_id = ? AND status = ? AND updated_at >= ? AND updated_at <= ?", businessID, models.ReminderStatusFailed, fromUTC, toUTC).
		Count(&n).Error
	return n, err
}

func (r *Repo) FindWalletByBusiness(businessID string) (*models.Wallet, error) {
	var w models.Wallet
	err := r.db.Where("business_id = ?", businessID).First(&w).Error
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *Repo) ListBusinesses() ([]models.Business, error) {
	var out []models.Business
	err := r.db.Order("created_at asc").Find(&out).Error
	return out, err
}

func (r *Repo) ListSchedulableServices(businessID string) ([]SchedulableServiceRow, error) {
	var rows []SchedulableServiceRow
	err := r.db.Table("customer_services cs").
		Select(`
			c.business_id,
			c.id AS customer_id,
			c.name AS customer_name,
			c.phone_number AS customer_wa,
			COALESCE(cs.category_id, '') AS category_id,
			COALESCE(cat.name, 'Layanan') AS service_name,
			cs.last_visit_at,
			cs.interval_days,
			COALESCE(cat.template_id, '') AS template_id,
			COALESCE(cat.is_enabled, true) AS category_enabled
		`).
		Joins("JOIN customers c ON c.id = cs.customer_id").
		Joins("LEFT JOIN categories cat ON cat.id = cs.category_id").
		Where("c.business_id = ?", businessID).
		Scan(&rows).Error
	return rows, err
}

func (r *Repo) ReminderExists(businessID, customerID, serviceName string, scheduledAt time.Time) (bool, error) {
	var n int64
	err := r.db.Model(&models.Reminder{}).
		Where("business_id = ? AND customer_id = ? AND svc_name = ? AND scheduled_at = ?", businessID, customerID, serviceName, scheduledAt).
		Count(&n).Error
	return n > 0, err
}

func (r *Repo) CreateReminder(rem models.Reminder) error {
	return r.db.Create(&rem).Error
}

func (r *Repo) ListDuePending(limit int) ([]models.Reminder, error) {
	if limit <= 0 {
		limit = 200
	}
	var out []models.Reminder
	err := r.db.Where("status = ? AND scheduled_at <= ?", models.ReminderStatusPending, time.Now().UTC()).
		Order("scheduled_at asc").
		Limit(limit).
		Find(&out).Error
	return out, err
}

func (r *Repo) GetDispatchContext(reminderID string) (*DispatchContext, error) {
	var rem models.Reminder
	if err := r.db.Where("id = ?", reminderID).First(&rem).Error; err != nil {
		return nil, err
	}

	var biz models.Business
	if err := r.db.Where("id = ?", rem.BusinessID).First(&biz).Error; err != nil {
		return nil, err
	}

	var cx models.Customer
	if err := r.db.Where("id = ? AND business_id = ?", rem.CustomerID, rem.BusinessID).First(&cx).Error; err != nil {
		return nil, err
	}

	var cat *models.Category
	if rem.CategoryID != "" {
		var c models.Category
		if err := r.db.Where("id = ? AND business_id = ?", rem.CategoryID, rem.BusinessID).First(&c).Error; err == nil {
			cat = &c
		}
	}

	return &DispatchContext{Reminder: rem, Business: biz, Customer: cx, Category: cat}, nil
}

func (r *Repo) MarkReminderFailed(reminderID, reason string) error {
	now := time.Now().UTC()
	return r.db.Model(&models.Reminder{}).Where("id = ?", reminderID).Updates(map[string]interface{}{
		"status":       models.ReminderStatusFailed,
		"error_reason": reason,
		"retry_count":  gorm.Expr("retry_count + 1"),
		"updated_at":   now,
	}).Error
}

func (r *Repo) MarkReminderSentAndDeduct(reminderID, metaMessageID string) error {
	now := time.Now().UTC()

	return r.db.Transaction(func(tx *gorm.DB) error {
		var rem models.Reminder
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", reminderID).First(&rem).Error; err != nil {
			return err
		}
		if rem.Status == models.ReminderStatusSent {
			return nil
		}
		if rem.Status != models.ReminderStatusPending {
			return nil
		}

		var w models.Wallet
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("business_id = ?", rem.BusinessID).First(&w).Error; err != nil {
			return err
		}

		if now.After(w.TrialEndsAt) && w.SubscriptionStatus != models.SubscriptionStatusActive {
			return updateFailedTx(tx, rem.ID, "trial ended, subscription required", now)
		}

		total := w.WelcomeCreditsLeft + w.SubCreditsLeft + w.TopupCreditsLeft
		if total <= 0 {
			return updateFailedTx(tx, rem.ID, "credits empty", now)
		}

		nw, ns, nt := deductCredits(w.WelcomeCreditsLeft, w.SubCreditsLeft, w.TopupCreditsLeft)
		if err := tx.Model(&models.Wallet{}).Where("id = ?", w.ID).Updates(map[string]interface{}{
			"welcome_credits_left": nw,
			"sub_credits_left":     ns,
			"topup_credits_left":   nt,
			"updated_at":           now,
		}).Error; err != nil {
			return err
		}

		if err := tx.Create(&models.BillingTransaction{
			ID:           uuid.NewString(),
			BusinessID:   rem.BusinessID,
			Type:         "usage",
			Label:        "Penggunaan Reminder",
			Delta:        -1,
			BalanceAfter: nw + ns + nt,
			Note:         fmt.Sprintf("Reminder %s (%s)", rem.CxName, rem.SvcName),
		}).Error; err != nil {
			return err
		}

		return tx.Model(&models.Reminder{}).Where("id = ?", rem.ID).Updates(map[string]interface{}{
			"status":          models.ReminderStatusSent,
			"sent_at":         now,
			"meta_message_id": metaMessageID,
			"error_reason":    "",
			"kredit":          1,
			"updated_at":      now,
		}).Error
	})
}

func (r *Repo) SetCategoryEnabledByMetaTemplateID(metaTemplateID string, status, category string) error {
	now := time.Now().UTC()
	if status != "" {
		return r.db.Model(&models.WATemplate{}).
			Where("meta_template_id = ?", metaTemplateID).
			Updates(map[string]any{
				"status":     status,
				"updated_at": now,
			}).Error
	} else {
		return r.db.Model(&models.WATemplate{}).
			Where("meta_template_id = ?", metaTemplateID).
			Updates(map[string]any{
				"category":   category,
				"updated_at": now,
			}).Error
	}
}

func (r *Repo) FindWATemplateByMetaTemplateName(templateName string) (*models.WATemplate, error) {
	var row models.WATemplate
	err := r.db.Where("meta_template_name = ?", templateName).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func updateFailedTx(tx *gorm.DB, reminderID, reason string, now time.Time) error {
	return tx.Model(&models.Reminder{}).Where("id = ?", reminderID).Updates(map[string]interface{}{
		"status":       models.ReminderStatusFailed,
		"error_reason": reason,
		"retry_count":  gorm.Expr("retry_count + 1"),
		"updated_at":   now,
	}).Error
}

func deductCredits(welcome, sub, topup int) (int, int, int) {
	if welcome > 0 {
		welcome--
		return welcome, sub, topup
	}
	if sub > 0 {
		sub--
		return welcome, sub, topup
	}
	topup--
	return welcome, sub, topup
}
