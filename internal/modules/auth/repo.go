package auth

import (
	"time"

	"github.com/aalexanderkevin/getstarvio-backend/internal/models"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) FindUserByGoogleSub(sub string) (*models.User, error) {
	var u models.User
	err := r.db.Where("google_sub = ?", sub).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repo) FindUserByID(id string) (*models.User, error) {
	var u models.User
	err := r.db.Where("id = ?", id).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repo) FindRefreshToken(tokenHash string) (*models.RefreshToken, error) {
	var t models.RefreshToken
	err := r.db.Where("token_hash = ?", tokenHash).First(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *Repo) RevokeRefreshToken(tokenHash string, at time.Time) error {
	return r.db.Model(&models.RefreshToken{}).Where("token_hash = ?", tokenHash).Update("revoked_at", at).Error
}

func (r *Repo) SaveRefreshToken(t models.RefreshToken) error {
	return r.db.Create(&t).Error
}

func (r *Repo) BootstrapNewUser(user models.User, biz models.Business, settings models.BusinessSettings, wallet models.Wallet, plan models.PlanConfig, tx models.BillingTransaction) error {
	return r.db.Transaction(func(db *gorm.DB) error {
		if err := db.Create(&user).Error; err != nil {
			return err
		}
		if err := db.Create(&biz).Error; err != nil {
			return err
		}
		if err := db.Create(&settings).Error; err != nil {
			return err
		}
		if err := db.Create(&wallet).Error; err != nil {
			return err
		}
		if err := db.Create(&plan).Error; err != nil {
			return err
		}
		if err := db.Create(&tx).Error; err != nil {
			return err
		}
		return nil
	})
}
