package internaladmin

import (
	"time"

	"gorm.io/gorm"

	"github.com/aalexanderkevin/getstarvio-backend/internal/models"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) FindAdminByEmail(email string) (*models.InternalAdmin, error) {
	var a models.InternalAdmin
	err := r.db.Where("LOWER(email) = LOWER(?)", email).First(&a).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *Repo) FindAdminByID(adminID string) (*models.InternalAdmin, error) {
	var a models.InternalAdmin
	err := r.db.Where("id = ?", adminID).First(&a).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *Repo) TouchAdminLastLogin(adminID string, at time.Time) error {
	return r.db.Model(&models.InternalAdmin{}).
		Where("id = ?", adminID).
		Updates(map[string]interface{}{
			"last_login_at": at,
			"updated_at":    at,
		}).Error
}

func (r *Repo) FindInternalRefreshToken(tokenHash string) (*models.InternalRefreshToken, error) {
	var t models.InternalRefreshToken
	err := r.db.Where("token_hash = ?", tokenHash).First(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *Repo) SaveInternalRefreshToken(t models.InternalRefreshToken) error {
	return r.db.Create(&t).Error
}

func (r *Repo) RevokeInternalRefreshToken(tokenHash string, at time.Time) error {
	return r.db.Model(&models.InternalRefreshToken{}).
		Where("token_hash = ?", tokenHash).
		Update("revoked_at", at).Error
}

func (r *Repo) ListDefaultCategories() ([]models.DefaultCategory, error) {
	var out []models.DefaultCategory
	err := r.db.Order("name asc").Find(&out).Error
	return out, err
}

func (r *Repo) CreateDefaultCategory(c models.DefaultCategory) error {
	return r.db.Create(&c).Error
}

func (r *Repo) GetPrimaryPlanConfig() (*models.PlanConfig, error) {
	var p models.PlanConfig
	err := r.db.Order("created_at asc").First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repo) UpdatePlanConfig(businessID string, payload map[string]interface{}) error {
	payload["updated_at"] = time.Now().UTC()
	return r.db.Model(&models.PlanConfig{}).Where("business_id = ?", businessID).Updates(payload).Error
}
