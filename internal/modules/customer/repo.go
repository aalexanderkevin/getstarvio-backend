package customer

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

func (r *Repo) ListCustomers(businessID string, q string) ([]models.Customer, error) {
	var out []models.Customer
	db := r.db.Where("business_id = ?", businessID)
	if q != "" {
		like := "%" + q + "%"
		db = db.Where("LOWER(name) LIKE LOWER(?) OR wa LIKE ?", like, like)
	}
	err := db.Order("created_at desc").Find(&out).Error
	return out, err
}

func (r *Repo) FindCustomer(businessID, id string) (*models.Customer, error) {
	var c models.Customer
	err := r.db.Where("business_id = ? AND id = ?", businessID, id).First(&c).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *Repo) FindCustomerByWA(businessID, wa string) (*models.Customer, error) {
	var c models.Customer
	err := r.db.Where("business_id = ? AND wa = ?", businessID, wa).First(&c).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *Repo) ListServices(customerIDs []string) ([]models.CustomerService, error) {
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

func (r *Repo) ListCategories(businessID string) ([]models.Category, error) {
	var out []models.Category
	err := r.db.Where("business_id = ?", businessID).Find(&out).Error
	return out, err
}

func (r *Repo) CreateCustomerWithServices(customer models.Customer, services []models.CustomerService) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&customer).Error; err != nil {
			return err
		}
		for _, s := range services {
			if err := tx.Create(&s).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *Repo) UpdateCustomerAndServices(businessID, customerID string, customerPayload map[string]interface{}, services []models.CustomerService) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if len(customerPayload) > 0 {
			customerPayload["updated_at"] = time.Now().UTC()
			if err := tx.Model(&models.Customer{}).Where("business_id = ? AND id = ?", businessID, customerID).Updates(customerPayload).Error; err != nil {
				return err
			}
		}
		for _, s := range services {
			var existing models.CustomerService
			err := tx.Where("customer_id = ? AND category_id = ?", customerID, s.CategoryID).First(&existing).Error
			if err == nil {
				if err := tx.Model(&models.CustomerService{}).Where("id = ?", existing.ID).Updates(map[string]interface{}{
					"last_visit_at": s.LastVisitAt,
					"interval_days": s.IntervalDays,
					"updated_at":    time.Now().UTC(),
				}).Error; err != nil {
					return err
				}
			} else if err == gorm.ErrRecordNotFound {
				if err := tx.Create(&s).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
		return nil
	})
}

func (r *Repo) DeleteCustomer(businessID, customerID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("customer_id = ?", customerID).Delete(&models.CustomerService{}).Error; err != nil {
			return err
		}
		return tx.Where("business_id = ? AND id = ?", businessID, customerID).Delete(&models.Customer{}).Error
	})
}
