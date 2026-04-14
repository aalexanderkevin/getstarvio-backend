package business

import (
	"time"

	"gorm.io/gorm"

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

func (r *Repo) FindSettings(businessID string) (*models.BusinessSettings, error) {
	var s models.BusinessSettings
	err := r.db.Where("business_id = ?", businessID).First(&s).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *Repo) FindWallet(businessID string) (*models.Wallet, error) {
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

func (r *Repo) ListCategories(businessID string) ([]models.Category, error) {
	var out []models.Category
	err := r.db.Where("business_id = ?", businessID).Order("created_at asc").Find(&out).Error
	return out, err
}

func (r *Repo) ListCustomers(businessID string) ([]models.Customer, error) {
	var out []models.Customer
	err := r.db.Where("business_id = ?", businessID).Order("created_at asc").Find(&out).Error
	return out, err
}

func (r *Repo) ListCustomerServices(customerIDs []string) ([]models.CustomerService, error) {
	if len(customerIDs) == 0 {
		return []models.CustomerService{}, nil
	}
	var out []models.CustomerService
	err := r.db.Table("customer_services cs").
		Select(`
			cs.id,
			cs.customer_id,
			cs.category_id,
			cs.last_visit_at,
			cs.interval_days,
			cs.created_at,
			cs.updated_at,
			COALESCE(cat.name, 'Layanan') AS service_name,
			COALESCE(cat.icon, '✨') AS service_icon
		`).
		Joins("LEFT JOIN categories cat ON cat.id = cs.category_id").
		Where("cs.customer_id IN ?", customerIDs).
		Find(&out).Error
	return out, err
}

func (r *Repo) ListReminders(businessID string, limit int) ([]models.Reminder, error) {
	if limit <= 0 {
		limit = 500
	}
	var out []models.Reminder
	err := r.db.Where("business_id = ?", businessID).Order("scheduled_at desc").Limit(limit).Find(&out).Error
	return out, err
}

func (r *Repo) UpdateProfile(businessID string, payload map[string]interface{}) error {
	payload["updated_at"] = time.Now().UTC()
	return r.db.Model(&models.Business{}).Where("id = ?", businessID).Updates(payload).Error
}

func (r *Repo) UpdateWhatsApp(businessID string, ownerWA string, waNum string, metaWabaID string, metaAccessToken string) error {
	return r.db.Model(&models.Business{}).Where("id = ?", businessID).Updates(map[string]interface{}{
		"owner_wa":          ownerWA,
		"wa_num":            waNum,
		"meta_waba_id":      metaWabaID,
		"meta_access_token": metaAccessToken,
		"updated_at":        time.Now().UTC(),
	}).Error
}

func (r *Repo) UpdateSettings(businessID string, payload map[string]interface{}) error {
	payload["updated_at"] = time.Now().UTC()
	return r.db.Model(&models.BusinessSettings{}).Where("business_id = ?", businessID).Updates(payload).Error
}
