package app

import (
	"github.com/aalexanderkevin/getstarvio-backend/internal/config"
	"github.com/aalexanderkevin/getstarvio-backend/internal/platform/meta"
	"github.com/aalexanderkevin/getstarvio-backend/internal/platform/xendit"
	"gorm.io/gorm"
)

type Container struct {
	Cfg    config.Config
	DB     *gorm.DB
	Meta   *meta.Client
	Xendit *xendit.Client
}

func NewContainer(cfg config.Config, db *gorm.DB) *Container {
	return &Container{
		Cfg:    cfg,
		DB:     db,
		Meta:   meta.NewClient(cfg.Meta, meta.NewGormFacebookLogStore(db)),
		Xendit: xendit.NewClient(cfg.Xendit),
	}
}
