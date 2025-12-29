package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

const (
	DurationUnitHour  = "hour"
	DurationUnitDay   = "day"
	DurationUnitMonth = "month"
	DurationUnitYear  = "year"
)

// PlanBillingOption defines a billing cycle and price for a plan.
type PlanBillingOption struct {
	ID            uint64 `gorm:"primaryKey"`
	PlanID        uint64 `gorm:"column:plan_id;index"`
	Name          string `gorm:"size:255"`
	DurationValue int    `gorm:"column:duration_value"`
	DurationUnit  string `gorm:"size:16;column:duration_unit"`
	PriceCents    int64  `gorm:"column:price_cents"`
	Currency      string `gorm:"size:16"`
	SortOrder     int    `gorm:"column:sort_order"`
	Status        string `gorm:"size:32"`
	Visible       bool   `gorm:"column:is_visible"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// TableName binds to the billing options table.
func (PlanBillingOption) TableName() string { return "plan_billing_options" }

// ListPlanBillingOptionsOptions controls filters for billing option listing.
type ListPlanBillingOptionsOptions struct {
	PlanID  uint64
	PlanIDs []uint64
	Status  string
	Visible *bool
}

// PlanBillingOptionRepository exposes persistence helpers for plan billing options.
type PlanBillingOptionRepository interface {
	List(ctx context.Context, opts ListPlanBillingOptionsOptions) ([]PlanBillingOption, error)
	Create(ctx context.Context, option PlanBillingOption) (PlanBillingOption, error)
	Update(ctx context.Context, id uint64, updates PlanBillingOption) (PlanBillingOption, error)
	Get(ctx context.Context, id uint64) (PlanBillingOption, error)
}

type planBillingOptionRepository struct {
	db *gorm.DB
}

// NewPlanBillingOptionRepository constructs the repository using a gorm DB.
func NewPlanBillingOptionRepository(db *gorm.DB) (PlanBillingOptionRepository, error) {
	if db == nil {
		return nil, errors.New("repository: database connection is required")
	}
	return &planBillingOptionRepository{db: db}, nil
}

func (r *planBillingOptionRepository) List(ctx context.Context, opts ListPlanBillingOptionsOptions) ([]PlanBillingOption, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	base := r.db.WithContext(ctx).Model(&PlanBillingOption{})
	if opts.PlanID > 0 {
		base = base.Where("plan_id = ?", opts.PlanID)
	}
	if len(opts.PlanIDs) > 0 {
		base = base.Where("plan_id IN ?", opts.PlanIDs)
	}
	if status := strings.TrimSpace(strings.ToLower(opts.Status)); status != "" {
		base = base.Where("LOWER(status) = ?", status)
	}
	if opts.Visible != nil {
		base = base.Where("is_visible = ?", *opts.Visible)
	}

	var options []PlanBillingOption
	if err := base.Order(buildPlanBillingOptionOrderClause()).Find(&options).Error; err != nil {
		return nil, err
	}
	return options, nil
}

func (r *planBillingOptionRepository) Create(ctx context.Context, option PlanBillingOption) (PlanBillingOption, error) {
	if err := ctx.Err(); err != nil {
		return PlanBillingOption{}, err
	}

	now := time.Now().UTC()
	if option.CreatedAt.IsZero() {
		option.CreatedAt = now
	}
	option.UpdatedAt = now
	if option.Status == "" {
		option.Status = "draft"
	}

	if err := r.db.WithContext(ctx).Create(&option).Error; err != nil {
		return PlanBillingOption{}, translateError(err)
	}
	return option, nil
}

func (r *planBillingOptionRepository) Update(ctx context.Context, id uint64, updates PlanBillingOption) (PlanBillingOption, error) {
	if err := ctx.Err(); err != nil {
		return PlanBillingOption{}, err
	}

	updates.UpdatedAt = time.Now().UTC()

	if err := r.db.WithContext(ctx).Model(&PlanBillingOption{}).Where("id = ?", id).Updates(map[string]any{
		"name":           updates.Name,
		"duration_value": updates.DurationValue,
		"duration_unit":  updates.DurationUnit,
		"price_cents":    updates.PriceCents,
		"currency":       updates.Currency,
		"sort_order":     updates.SortOrder,
		"status":         updates.Status,
		"is_visible":     updates.Visible,
		"updated_at":     updates.UpdatedAt,
	}).Error; err != nil {
		return PlanBillingOption{}, translateError(err)
	}

	return r.Get(ctx, id)
}

func (r *planBillingOptionRepository) Get(ctx context.Context, id uint64) (PlanBillingOption, error) {
	if err := ctx.Err(); err != nil {
		return PlanBillingOption{}, err
	}

	var option PlanBillingOption
	if err := r.db.WithContext(ctx).First(&option, id).Error; err != nil {
		return PlanBillingOption{}, translateError(err)
	}
	return option, nil
}

func buildPlanBillingOptionOrderClause() string {
	return fmt.Sprintf("%s %s, %s %s, id ASC", "sort_order", "ASC", "duration_value", "ASC")
}
