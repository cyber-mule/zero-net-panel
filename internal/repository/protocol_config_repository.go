package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// ProtocolConfig stores protocol-specific configuration profiles.
type ProtocolConfig struct {
	ID          uint64         `gorm:"primaryKey"`
	Name        string         `gorm:"size:255;uniqueIndex"`
	Protocol    string         `gorm:"size:32;index"`
	Status      string         `gorm:"size:32"`
	Tags        []string       `gorm:"serializer:json"`
	Description string         `gorm:"type:text"`
	Profile     map[string]any `gorm:"serializer:json"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TableName binds the protocol config table name.
func (ProtocolConfig) TableName() string { return "protocol_configs" }

// ListProtocolConfigsOptions controls filtering and pagination.
type ListProtocolConfigsOptions struct {
	Page      int
	PerPage   int
	Sort      string
	Direction string
	Query     string
	Protocol  string
	Status    string
}

// UpdateProtocolConfigInput defines mutable protocol config fields.
type UpdateProtocolConfigInput struct {
	Name        *string
	Protocol    *string
	Status      *string
	Tags        *[]string
	Description *string
	Profile     *map[string]any
}

// ProtocolConfigRepository manages protocol configuration persistence.
type ProtocolConfigRepository interface {
	List(ctx context.Context, opts ListProtocolConfigsOptions) ([]ProtocolConfig, int64, error)
	Get(ctx context.Context, id uint64) (ProtocolConfig, error)
	Create(ctx context.Context, cfg ProtocolConfig) (ProtocolConfig, error)
	Update(ctx context.Context, id uint64, input UpdateProtocolConfigInput) (ProtocolConfig, error)
	Delete(ctx context.Context, id uint64) error
}

type protocolConfigRepository struct {
	db *gorm.DB
}

// NewProtocolConfigRepository constructs a protocol config repository.
func NewProtocolConfigRepository(db *gorm.DB) (ProtocolConfigRepository, error) {
	if db == nil {
		return nil, errors.New("repository: database connection is required")
	}
	return &protocolConfigRepository{db: db}, nil
}

func (r *protocolConfigRepository) List(ctx context.Context, opts ListProtocolConfigsOptions) ([]ProtocolConfig, int64, error) {
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}

	opts = normalizeListProtocolConfigsOptions(opts)
	base := r.db.WithContext(ctx).Model(&ProtocolConfig{})

	if query := strings.TrimSpace(strings.ToLower(opts.Query)); query != "" {
		like := fmt.Sprintf("%%%s%%", query)
		base = base.Where("(LOWER(name) LIKE ? OR LOWER(description) LIKE ?)", like, like)
	}
	if protocol := strings.TrimSpace(strings.ToLower(opts.Protocol)); protocol != "" {
		base = base.Where("LOWER(protocol) = ?", protocol)
	}
	if status := strings.TrimSpace(strings.ToLower(opts.Status)); status != "" {
		base = base.Where("LOWER(status) = ?", status)
	}

	countQuery := base.Session(&gorm.Session{})
	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []ProtocolConfig{}, 0, nil
	}

	orderClause := buildProtocolConfigOrderClause(opts.Sort, opts.Direction)
	offset := (opts.Page - 1) * opts.PerPage
	listQuery := base.Session(&gorm.Session{}).Order(orderClause).Limit(opts.PerPage).Offset(offset)

	var configs []ProtocolConfig
	if err := listQuery.Find(&configs).Error; err != nil {
		return nil, 0, err
	}

	return configs, total, nil
}

func (r *protocolConfigRepository) Get(ctx context.Context, id uint64) (ProtocolConfig, error) {
	if err := ctx.Err(); err != nil {
		return ProtocolConfig{}, err
	}

	var cfg ProtocolConfig
	if err := r.db.WithContext(ctx).First(&cfg, id).Error; err != nil {
		return ProtocolConfig{}, translateError(err)
	}
	return cfg, nil
}

func (r *protocolConfigRepository) Create(ctx context.Context, cfg ProtocolConfig) (ProtocolConfig, error) {
	if err := ctx.Err(); err != nil {
		return ProtocolConfig{}, err
	}

	cfg.Name = strings.TrimSpace(cfg.Name)
	cfg.Protocol = strings.ToLower(strings.TrimSpace(cfg.Protocol))
	cfg.Status = strings.TrimSpace(cfg.Status)
	cfg.Description = strings.TrimSpace(cfg.Description)
	if cfg.Name == "" || cfg.Protocol == "" {
		return ProtocolConfig{}, ErrInvalidArgument
	}
	if cfg.Status == "" {
		cfg.Status = "active"
	}
	if cfg.Tags == nil {
		cfg.Tags = []string{}
	}
	if cfg.Profile == nil {
		cfg.Profile = map[string]any{}
	}

	now := time.Now().UTC()
	if cfg.CreatedAt.IsZero() {
		cfg.CreatedAt = now
	}
	cfg.UpdatedAt = now

	if err := r.db.WithContext(ctx).Create(&cfg).Error; err != nil {
		return ProtocolConfig{}, translateError(err)
	}
	return cfg, nil
}

func (r *protocolConfigRepository) Update(ctx context.Context, id uint64, input UpdateProtocolConfigInput) (ProtocolConfig, error) {
	if err := ctx.Err(); err != nil {
		return ProtocolConfig{}, err
	}

	updates := map[string]any{}
	if input.Name != nil {
		updates["name"] = strings.TrimSpace(*input.Name)
	}
	if input.Protocol != nil {
		updates["protocol"] = strings.ToLower(strings.TrimSpace(*input.Protocol))
	}
	if input.Status != nil {
		updates["status"] = strings.TrimSpace(*input.Status)
	}
	if input.Tags != nil {
		updates["tags"] = append([]string(nil), (*input.Tags)...)
	}
	if input.Description != nil {
		updates["description"] = strings.TrimSpace(*input.Description)
	}
	if input.Profile != nil {
		updates["profile"] = *input.Profile
	}
	if len(updates) == 0 {
		return ProtocolConfig{}, ErrInvalidArgument
	}
	updates["updated_at"] = time.Now().UTC()

	if err := r.db.WithContext(ctx).Model(&ProtocolConfig{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return ProtocolConfig{}, translateError(err)
	}

	return r.Get(ctx, id)
}

func (r *protocolConfigRepository) Delete(ctx context.Context, id uint64) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).Delete(&ProtocolConfig{}, id).Error; err != nil {
		return translateError(err)
	}
	return nil
}

func normalizeListProtocolConfigsOptions(opts ListProtocolConfigsOptions) ListProtocolConfigsOptions {
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
		opts.Sort = "updated_at"
	}
	if opts.Direction == "" {
		opts.Direction = "desc"
	}
	return opts
}

func buildProtocolConfigOrderClause(field, direction string) string {
	column := "updated_at"
	switch strings.ToLower(field) {
	case "name":
		column = "name"
	case "protocol":
		column = "protocol"
	case "status":
		column = "status"
	case "created_at":
		column = "created_at"
	}

	dir := "ASC"
	if strings.EqualFold(direction, "desc") {
		dir = "DESC"
	}

	return fmt.Sprintf("%s %s", column, dir)
}
