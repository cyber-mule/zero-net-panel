package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

// AuditLog stores security-relevant actions for audit trails.
type AuditLog struct {
	ID           uint64         `gorm:"primaryKey"`
	ActorID      *uint64        `gorm:"column:actor_id"`
	ActorEmail   string         `gorm:"size:255"`
	ActorRoles   []string       `gorm:"serializer:json"`
	Action       string         `gorm:"size:64"`
	ResourceType string         `gorm:"size:64"`
	ResourceID   string         `gorm:"size:64"`
	SourceIP     string         `gorm:"size:64"`
	Metadata     map[string]any `gorm:"serializer:json"`
	CreatedAt    time.Time
}

// TableName binds to audit_logs.
func (AuditLog) TableName() string { return "audit_logs" }

// AuditLogListOptions controls list filters.
type AuditLogListOptions struct {
	Page         int
	PerPage      int
	ActorID      *uint64
	Action       string
	ResourceType string
	ResourceID   string
	Since        *time.Time
	Until        *time.Time
}

// AuditLogRepository exposes audit log persistence.
type AuditLogRepository interface {
	Create(ctx context.Context, entry AuditLog) (AuditLog, error)
	List(ctx context.Context, opts AuditLogListOptions) ([]AuditLog, int64, error)
}

type auditLogRepository struct {
	db *gorm.DB
}

// NewAuditLogRepository constructs audit log repository.
func NewAuditLogRepository(db *gorm.DB) (AuditLogRepository, error) {
	if db == nil {
		return nil, errors.New("repository: database connection is required")
	}
	return &auditLogRepository{db: db}, nil
}

func (r *auditLogRepository) Create(ctx context.Context, entry AuditLog) (AuditLog, error) {
	if err := ctx.Err(); err != nil {
		return AuditLog{}, err
	}

	entry.Action = strings.TrimSpace(entry.Action)
	entry.ResourceType = strings.TrimSpace(entry.ResourceType)
	entry.ResourceID = strings.TrimSpace(entry.ResourceID)
	entry.ActorEmail = strings.TrimSpace(entry.ActorEmail)
	entry.SourceIP = strings.TrimSpace(entry.SourceIP)
	if entry.Action == "" {
		return AuditLog{}, ErrInvalidArgument
	}
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = time.Now().UTC()
	}
	if entry.Metadata == nil {
		entry.Metadata = map[string]any{}
	}

	if err := r.db.WithContext(ctx).Create(&entry).Error; err != nil {
		return AuditLog{}, translateError(err)
	}

	return entry, nil
}

func (r *auditLogRepository) List(ctx context.Context, opts AuditLogListOptions) ([]AuditLog, int64, error) {
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}

	if opts.Page <= 0 {
		opts.Page = 1
	}
	if opts.PerPage <= 0 || opts.PerPage > 100 {
		opts.PerPage = 20
	}

	base := r.db.WithContext(ctx).Model(&AuditLog{})
	if opts.ActorID != nil {
		base = base.Where("actor_id = ?", *opts.ActorID)
	}
	if action := strings.TrimSpace(strings.ToLower(opts.Action)); action != "" {
		base = base.Where("LOWER(action) = ?", action)
	}
	if resourceType := strings.TrimSpace(strings.ToLower(opts.ResourceType)); resourceType != "" {
		base = base.Where("LOWER(resource_type) = ?", resourceType)
	}
	if resourceID := strings.TrimSpace(opts.ResourceID); resourceID != "" {
		base = base.Where("resource_id = ?", resourceID)
	}
	if opts.Since != nil && !opts.Since.IsZero() {
		base = base.Where("created_at >= ?", opts.Since.UTC())
	}
	if opts.Until != nil && !opts.Until.IsZero() {
		base = base.Where("created_at <= ?", opts.Until.UTC())
	}

	countQuery := base.Session(&gorm.Session{})
	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []AuditLog{}, 0, nil
	}

	offset := (opts.Page - 1) * opts.PerPage
	listQuery := base.Session(&gorm.Session{}).Order("created_at DESC, id DESC").Limit(opts.PerPage).Offset(offset)

	var logs []AuditLog
	if err := listQuery.Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
