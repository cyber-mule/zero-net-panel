package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// PaymentChannel stores configurable external payment gateways.
type PaymentChannel struct {
	ID        uint64         `gorm:"primaryKey"`
	Name      string         `gorm:"size:128"`
	Code      string         `gorm:"size:64;uniqueIndex"`
	Provider  string         `gorm:"size:64"`
	Enabled   bool           `gorm:"column:is_enabled"`
	SortOrder int            `gorm:"column:sort_order"`
	Config    map[string]any `gorm:"serializer:json"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName provides explicit table binding.
func (PaymentChannel) TableName() string { return "payment_channels" }

// ListPaymentChannelsOptions controls filtering and pagination.
type ListPaymentChannelsOptions struct {
	Page      int
	PerPage   int
	Query     string
	Provider  string
	Enabled   *bool
	Sort      string
	Direction string
}

// PaymentChannelRepository exposes persistence helpers for payment channels.
type PaymentChannelRepository interface {
	List(ctx context.Context, opts ListPaymentChannelsOptions) ([]PaymentChannel, int64, error)
	Create(ctx context.Context, channel PaymentChannel) (PaymentChannel, error)
	Update(ctx context.Context, id uint64, updates PaymentChannel) (PaymentChannel, error)
	Get(ctx context.Context, id uint64) (PaymentChannel, error)
	GetByCode(ctx context.Context, code string) (PaymentChannel, error)
}

type paymentChannelRepository struct {
	db *gorm.DB
}

// NewPaymentChannelRepository constructs the repository using a gorm DB.
func NewPaymentChannelRepository(db *gorm.DB) (PaymentChannelRepository, error) {
	if db == nil {
		return nil, errors.New("repository: database connection is required")
	}
	return &paymentChannelRepository{db: db}, nil
}

func (r *paymentChannelRepository) List(ctx context.Context, opts ListPaymentChannelsOptions) ([]PaymentChannel, int64, error) {
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}

	opts = normalizeListPaymentChannelsOptions(opts)

	base := r.db.WithContext(ctx).Model(&PaymentChannel{})

	if query := strings.TrimSpace(strings.ToLower(opts.Query)); query != "" {
		like := fmt.Sprintf("%%%s%%", query)
		base = base.Where("(LOWER(name) LIKE ? OR LOWER(code) LIKE ?)", like, like)
	}
	if provider := strings.TrimSpace(strings.ToLower(opts.Provider)); provider != "" {
		base = base.Where("LOWER(provider) = ?", provider)
	}
	if opts.Enabled != nil {
		base = base.Where("is_enabled = ?", *opts.Enabled)
	}

	countQuery := base.Session(&gorm.Session{})
	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []PaymentChannel{}, 0, nil
	}

	orderClause := buildPaymentChannelOrderClause(opts.Sort, opts.Direction)
	offset := (opts.Page - 1) * opts.PerPage
	listQuery := base.Session(&gorm.Session{}).Order(orderClause).Limit(opts.PerPage).Offset(offset)

	var channels []PaymentChannel
	if err := listQuery.Find(&channels).Error; err != nil {
		return nil, 0, err
	}

	return channels, total, nil
}

func (r *paymentChannelRepository) Create(ctx context.Context, channel PaymentChannel) (PaymentChannel, error) {
	if err := ctx.Err(); err != nil {
		return PaymentChannel{}, err
	}

	now := time.Now().UTC()
	if channel.CreatedAt.IsZero() {
		channel.CreatedAt = now
	}
	channel.UpdatedAt = now

	normalizePaymentChannel(&channel)

	if err := r.db.WithContext(ctx).Create(&channel).Error; err != nil {
		return PaymentChannel{}, translateError(err)
	}

	return channel, nil
}

func (r *paymentChannelRepository) Update(ctx context.Context, id uint64, updates PaymentChannel) (PaymentChannel, error) {
	if err := ctx.Err(); err != nil {
		return PaymentChannel{}, err
	}

	updates.UpdatedAt = time.Now().UTC()
	normalizePaymentChannel(&updates)

	if err := r.db.WithContext(ctx).Model(&PaymentChannel{}).Where("id = ?", id).Updates(map[string]any{
		"name":       updates.Name,
		"code":       updates.Code,
		"provider":   updates.Provider,
		"is_enabled": updates.Enabled,
		"sort_order": updates.SortOrder,
		"config":     updates.Config,
		"updated_at": updates.UpdatedAt,
	}).Error; err != nil {
		return PaymentChannel{}, translateError(err)
	}

	return r.Get(ctx, id)
}

func (r *paymentChannelRepository) Get(ctx context.Context, id uint64) (PaymentChannel, error) {
	if err := ctx.Err(); err != nil {
		return PaymentChannel{}, err
	}

	var channel PaymentChannel
	if err := r.db.WithContext(ctx).First(&channel, id).Error; err != nil {
		return PaymentChannel{}, translateError(err)
	}

	return channel, nil
}

func (r *paymentChannelRepository) GetByCode(ctx context.Context, code string) (PaymentChannel, error) {
	if err := ctx.Err(); err != nil {
		return PaymentChannel{}, err
	}

	code = strings.TrimSpace(strings.ToLower(code))
	if code == "" {
		return PaymentChannel{}, ErrInvalidArgument
	}

	var channel PaymentChannel
	if err := r.db.WithContext(ctx).Where("LOWER(code) = ?", code).First(&channel).Error; err != nil {
		return PaymentChannel{}, translateError(err)
	}

	return channel, nil
}

func normalizePaymentChannel(channel *PaymentChannel) {
	channel.Name = strings.TrimSpace(channel.Name)
	channel.Code = strings.ToLower(strings.TrimSpace(channel.Code))
	channel.Provider = strings.ToLower(strings.TrimSpace(channel.Provider))
	if channel.Provider == "" {
		channel.Provider = channel.Code
	}
}

func normalizeListPaymentChannelsOptions(opts ListPaymentChannelsOptions) ListPaymentChannelsOptions {
	if opts.Page <= 0 {
		opts.Page = 1
	}
	if opts.PerPage <= 0 || opts.PerPage > 100 {
		opts.PerPage = 20
	}
	return opts
}

func buildPaymentChannelOrderClause(sort, direction string) string {
	column := "sort_order"
	dir := "ASC"

	switch strings.ToLower(strings.TrimSpace(sort)) {
	case "name":
		column = "name"
	case "created":
		column = "created_at"
	case "updated":
		column = "updated_at"
	}

	if strings.EqualFold(direction, "desc") {
		dir = "DESC"
	}

	return fmt.Sprintf("%s %s, id ASC", column, dir)
}
