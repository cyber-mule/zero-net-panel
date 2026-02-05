package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

// SiteSetting stores branding configuration.
type SiteSetting struct {
	ID                 uint64 `gorm:"primaryKey"`
	Name               string `gorm:"size:128"`
	LogoURL            string `gorm:"size:512"`
	ServiceDomain      string `gorm:"column:access_domain;size:512"`
	SubscriptionDomain string `gorm:"size:512"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// TableName custom binding.
func (SiteSetting) TableName() string { return "site_settings" }

// SiteSettingDefaults contains fallback values when initializing settings.
type SiteSettingDefaults struct {
	Name               string
	LogoURL            string
	ServiceDomain      string
	SubscriptionDomain string
}

// SiteRepository exposes accessors for site settings.
type SiteRepository interface {
	GetSiteSetting(ctx context.Context, defaults SiteSettingDefaults) (SiteSetting, error)
	UpsertSiteSetting(ctx context.Context, setting SiteSetting) (SiteSetting, error)
}

type siteRepository struct {
	db *gorm.DB
}

// NewSiteRepository constructs repo.
func NewSiteRepository(db *gorm.DB) (SiteRepository, error) {
	if db == nil {
		return nil, errors.New("repository: database connection is required")
	}
	return &siteRepository{db: db}, nil
}

func (r *siteRepository) GetSiteSetting(ctx context.Context, defaults SiteSettingDefaults) (SiteSetting, error) {
	if err := ctx.Err(); err != nil {
		return SiteSetting{}, err
	}

	var setting SiteSetting
	if err := r.db.WithContext(ctx).Limit(1).First(&setting).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			now := time.Now().UTC()
			setting = SiteSetting{
				Name:               strings.TrimSpace(defaults.Name),
				LogoURL:            strings.TrimSpace(defaults.LogoURL),
				ServiceDomain:      strings.TrimSpace(defaults.ServiceDomain),
				SubscriptionDomain: strings.TrimSpace(defaults.SubscriptionDomain),
				CreatedAt:          now,
				UpdatedAt:          now,
			}
			if err := r.db.WithContext(ctx).Create(&setting).Error; err != nil {
				return SiteSetting{}, err
			}
			return setting, nil
		}
		return SiteSetting{}, err
	}

	return setting, nil
}

func (r *siteRepository) UpsertSiteSetting(ctx context.Context, setting SiteSetting) (SiteSetting, error) {
	if err := ctx.Err(); err != nil {
		return SiteSetting{}, err
	}

	setting.Name = strings.TrimSpace(setting.Name)
	setting.LogoURL = strings.TrimSpace(setting.LogoURL)
	setting.ServiceDomain = strings.TrimSpace(setting.ServiceDomain)
	setting.SubscriptionDomain = strings.TrimSpace(setting.SubscriptionDomain)

	now := time.Now().UTC()
	setting.UpdatedAt = now

	if setting.ID == 0 {
		setting.CreatedAt = now
		if err := r.db.WithContext(ctx).Create(&setting).Error; err != nil {
			return SiteSetting{}, err
		}
		return setting, nil
	}

	if err := r.db.WithContext(ctx).Model(&SiteSetting{}).
		Where("id = ?", setting.ID).
		Updates(map[string]any{
			"name":                setting.Name,
			"logo_url":            setting.LogoURL,
			"access_domain":       setting.ServiceDomain,
			"subscription_domain": setting.SubscriptionDomain,
			"updated_at":          setting.UpdatedAt,
		}).Error; err != nil {
		return SiteSetting{}, err
	}

	return r.GetSiteSetting(ctx, SiteSettingDefaults{})
}
