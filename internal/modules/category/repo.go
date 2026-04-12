package category

import (
	"time"

	"github.com/aalexanderkevin/getstarvio-backend/internal/models"
	"gorm.io/gorm"
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

func (r *Repo) List(businessID string) ([]models.Category, error) {
	var out []models.Category
	err := r.db.Where("business_id = ?", businessID).Order("created_at asc").Find(&out).Error
	return out, err
}

func (r *Repo) Create(c models.Category) error {
	return r.db.Create(&c).Error
}

func (r *Repo) FindByID(businessID, id string) (*models.Category, error) {
	var c models.Category
	err := r.db.Where("business_id = ? AND id = ?", businessID, id).First(&c).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *Repo) Update(businessID, id string, payload map[string]interface{}) error {
	payload["updated_at"] = time.Now().UTC()
	return r.db.Model(&models.Category{}).Where("business_id = ? AND id = ?", businessID, id).Updates(payload).Error
}

func (r *Repo) Delete(businessID, id string) error {
	return r.db.Where("business_id = ? AND id = ?", businessID, id).Delete(&models.Category{}).Error
}
