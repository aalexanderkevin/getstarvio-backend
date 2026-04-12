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
