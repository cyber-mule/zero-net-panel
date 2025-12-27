package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// TrafficUsageRecord stores raw and charged traffic usage entries.
type TrafficUsageRecord struct {
	ID                uint64    `gorm:"primaryKey"`
	UserID            uint64    `gorm:"index"`
	SubscriptionID    uint64    `gorm:"index"`
	ProtocolBindingID uint64    `gorm:"index"`
	NodeID            uint64    `gorm:"index"`
	Protocol          string    `gorm:"size:32;index"`
	BytesUp           int64     `gorm:"column:bytes_up"`
	BytesDown         int64     `gorm:"column:bytes_down"`
	RawBytes          int64     `gorm:"column:raw_bytes"`
	ChargedBytes      int64     `gorm:"column:charged_bytes"`
	Multiplier        float64   `gorm:"column:multiplier"`
	ObservedAt        time.Time `gorm:"index;column:observed_at"`
	CreatedAt         time.Time
}

// TableName binds the traffic usage table name.
func (TrafficUsageRecord) TableName() string { return "traffic_usage_records" }

// ListTrafficUsageOptions controls usage listing.
type ListTrafficUsageOptions struct {
	Page              int
	PerPage           int
	Sort              string
	Direction         string
	Protocol          string
	NodeID            *uint64
	ProtocolBindingID *uint64
	From              *time.Time
	To                *time.Time
}

// TrafficUsageRepository manages usage records.
type TrafficUsageRepository interface {
	ListBySubscription(ctx context.Context, subscriptionID uint64, opts ListTrafficUsageOptions) ([]TrafficUsageRecord, int64, error)
	Create(ctx context.Context, record TrafficUsageRecord) (TrafficUsageRecord, error)
	SumBySubscription(ctx context.Context, subscriptionID uint64) (int64, int64, error)
}

type trafficUsageRepository struct {
	db *gorm.DB
}

// NewTrafficUsageRepository constructs a usage repository.
func NewTrafficUsageRepository(db *gorm.DB) (TrafficUsageRepository, error) {
	if db == nil {
		return nil, errors.New("repository: database connection is required")
	}
	return &trafficUsageRepository{db: db}, nil
}

func (r *trafficUsageRepository) ListBySubscription(ctx context.Context, subscriptionID uint64, opts ListTrafficUsageOptions) ([]TrafficUsageRecord, int64, error) {
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}
	if subscriptionID == 0 {
		return nil, 0, ErrInvalidArgument
	}

	opts = normalizeListTrafficUsageOptions(opts)
	base := r.db.WithContext(ctx).Model(&TrafficUsageRecord{}).Where("subscription_id = ?", subscriptionID)

	if protocol := strings.TrimSpace(strings.ToLower(opts.Protocol)); protocol != "" {
		base = base.Where("LOWER(protocol) = ?", protocol)
	}
	if opts.NodeID != nil {
		base = base.Where("node_id = ?", *opts.NodeID)
	}
	if opts.ProtocolBindingID != nil {
		base = base.Where("protocol_binding_id = ?", *opts.ProtocolBindingID)
	}
	if opts.From != nil && !opts.From.IsZero() {
		base = base.Where("observed_at >= ?", opts.From.UTC())
	}
	if opts.To != nil && !opts.To.IsZero() {
		base = base.Where("observed_at <= ?", opts.To.UTC())
	}

	countQuery := base.Session(&gorm.Session{})
	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []TrafficUsageRecord{}, 0, nil
	}

	orderClause := buildTrafficUsageOrderClause(opts.Sort, opts.Direction)
	offset := (opts.Page - 1) * opts.PerPage
	listQuery := base.Session(&gorm.Session{}).Order(orderClause).Limit(opts.PerPage).Offset(offset)

	var records []TrafficUsageRecord
	if err := listQuery.Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

func (r *trafficUsageRepository) Create(ctx context.Context, record TrafficUsageRecord) (TrafficUsageRecord, error) {
	if err := ctx.Err(); err != nil {
		return TrafficUsageRecord{}, err
	}
	if record.UserID == 0 || record.SubscriptionID == 0 {
		return TrafficUsageRecord{}, ErrInvalidArgument
	}

	record.Protocol = strings.ToLower(strings.TrimSpace(record.Protocol))
	now := time.Now().UTC()
	if record.ObservedAt.IsZero() {
		record.ObservedAt = now
	}
	if record.CreatedAt.IsZero() {
		record.CreatedAt = now
	}

	if err := r.db.WithContext(ctx).Create(&record).Error; err != nil {
		return TrafficUsageRecord{}, translateError(err)
	}
	return record, nil
}

func (r *trafficUsageRepository) SumBySubscription(ctx context.Context, subscriptionID uint64) (int64, int64, error) {
	if err := ctx.Err(); err != nil {
		return 0, 0, err
	}
	if subscriptionID == 0 {
		return 0, 0, ErrInvalidArgument
	}

	var result struct {
		RawBytes     int64 `gorm:"column:raw_bytes"`
		ChargedBytes int64 `gorm:"column:charged_bytes"`
	}
	err := r.db.WithContext(ctx).
		Model(&TrafficUsageRecord{}).
		Select("COALESCE(SUM(raw_bytes), 0) AS raw_bytes, COALESCE(SUM(charged_bytes), 0) AS charged_bytes").
		Where("subscription_id = ?", subscriptionID).
		Scan(&result).Error
	if err != nil {
		return 0, 0, err
	}
	return result.RawBytes, result.ChargedBytes, nil
}

func normalizeListTrafficUsageOptions(opts ListTrafficUsageOptions) ListTrafficUsageOptions {
	if opts.Page <= 0 {
		opts.Page = 1
	}
	if opts.PerPage <= 0 {
		opts.PerPage = 20
	}
	if opts.PerPage > 100 {
		opts.PerPage = 100
	}
	if opts.Sort == "" {
		opts.Sort = "observed_at"
	}
	if opts.Direction == "" {
		opts.Direction = "desc"
	}
	return opts
}

func buildTrafficUsageOrderClause(field, direction string) string {
	column := "observed_at"
	switch strings.ToLower(field) {
	case "created_at":
		column = "created_at"
	case "charged_bytes":
		column = "charged_bytes"
	case "raw_bytes":
		column = "raw_bytes"
	}

	dir := "ASC"
	if strings.EqualFold(direction, "desc") {
		dir = "DESC"
	}
	return fmt.Sprintf("%s %s", column, dir)
}
