package meta

import (
	"context"

	"gorm.io/gorm"

	"github.com/aalexanderkevin/getstarvio-backend/internal/models"
)

type FacebookLogEntry struct {
	Operation    string
	URL          string
	RequestBody  string
	ResponseBody string
	ResponseCode int
	RefID        string
}

type FacebookLogStore interface {
	CreateFacebookLog(ctx context.Context, entry FacebookLogEntry) error
}

type GormFacebookLogStore struct {
	db *gorm.DB
}

func NewGormFacebookLogStore(db *gorm.DB) *GormFacebookLogStore {
	return &GormFacebookLogStore{db: db}
}

func (s *GormFacebookLogStore) CreateFacebookLog(ctx context.Context, entry FacebookLogEntry) error {
	row := models.FacebookLog{
		Operation:    entry.Operation,
		URL:          entry.URL,
		RequestBody:  entry.RequestBody,
		ResponseBody: entry.ResponseBody,
		ResponseCode: entry.ResponseCode,
		RefID:        entry.RefID,
	}
	return s.db.WithContext(ctx).Create(&row).Error
}
