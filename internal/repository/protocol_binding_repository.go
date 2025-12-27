package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// ProtocolBinding binds a protocol configuration to a node instance.
type ProtocolBinding struct {
	ID               uint64         `gorm:"primaryKey"`
	Name             string         `gorm:"size:255"`
	NodeID           uint64         `gorm:"index"`
	ProtocolConfigID uint64         `gorm:"index"`
	Role             string         `gorm:"size:32"`
	Listen           string         `gorm:"size:512"`
	Connect          string         `gorm:"size:512"`
	Status           string         `gorm:"size:32"`
	KernelID         string         `gorm:"size:128;index"`
	SyncStatus       string         `gorm:"size:32"`
	HealthStatus     string         `gorm:"size:32"`
	LastSyncedAt     time.Time      `gorm:"column:last_synced_at"`
	LastHeartbeatAt  time.Time      `gorm:"column:last_heartbeat_at"`
	LastSyncError    string         `gorm:"type:text"`
	Tags             []string       `gorm:"serializer:json"`
	Description      string         `gorm:"type:text"`
	Metadata         map[string]any `gorm:"serializer:json"`
	UpdatedAt        time.Time
	CreatedAt        time.Time

	Node           Node           `gorm:"foreignKey:NodeID;references:ID"`
	ProtocolConfig ProtocolConfig `gorm:"foreignKey:ProtocolConfigID;references:ID"`
}

// TableName binds the protocol binding table name.
func (ProtocolBinding) TableName() string { return "protocol_bindings" }

// ListProtocolBindingsOptions controls filtering and pagination.
type ListProtocolBindingsOptions struct {
	Page      int
	PerPage   int
	Sort      string
	Direction string
	Query     string
	Status    string
	Protocol  string
	NodeID    *uint64
	ConfigID  *uint64
}

// UpdateProtocolBindingInput defines mutable binding fields.
type UpdateProtocolBindingInput struct {
	Name             *string
	NodeID           *uint64
	ProtocolConfigID *uint64
	Role             *string
	Listen           *string
	Connect          *string
	Status           *string
	KernelID         *string
	SyncStatus       *string
	HealthStatus     *string
	LastSyncedAt     *time.Time
	LastHeartbeatAt  *time.Time
	LastSyncError    *string
	Tags             *[]string
	Description      *string
	Metadata         *map[string]any
}

// ProtocolBindingRepository manages protocol binding persistence.
type ProtocolBindingRepository interface {
	List(ctx context.Context, opts ListProtocolBindingsOptions) ([]ProtocolBinding, int64, error)
	ListByNodeIDs(ctx context.Context, nodeIDs []uint64) ([]ProtocolBinding, error)
	ListAll(ctx context.Context) ([]ProtocolBinding, error)
	Get(ctx context.Context, id uint64) (ProtocolBinding, error)
	Create(ctx context.Context, binding ProtocolBinding) (ProtocolBinding, error)
	Update(ctx context.Context, id uint64, input UpdateProtocolBindingInput) (ProtocolBinding, error)
	UpdateSyncState(ctx context.Context, id uint64, input UpdateProtocolBindingInput) (ProtocolBinding, error)
	UpdateHealthByKernelID(ctx context.Context, kernelID string, status string, observedAt time.Time, message string) (ProtocolBinding, error)
	Delete(ctx context.Context, id uint64) error
}

type protocolBindingRepository struct {
	db *gorm.DB
}

// NewProtocolBindingRepository constructs a protocol binding repository.
func NewProtocolBindingRepository(db *gorm.DB) (ProtocolBindingRepository, error) {
	if db == nil {
		return nil, errors.New("repository: database connection is required")
	}
	return &protocolBindingRepository{db: db}, nil
}

func (r *protocolBindingRepository) List(ctx context.Context, opts ListProtocolBindingsOptions) ([]ProtocolBinding, int64, error) {
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}

	opts = normalizeListProtocolBindingsOptions(opts)
	base := r.db.WithContext(ctx).Model(&ProtocolBinding{})

	if query := strings.TrimSpace(strings.ToLower(opts.Query)); query != "" {
		like := fmt.Sprintf("%%%s%%", query)
		base = base.Where("(LOWER(name) LIKE ? OR LOWER(description) LIKE ?)", like, like)
	}
	if status := strings.TrimSpace(strings.ToLower(opts.Status)); status != "" {
		base = base.Where("LOWER(status) = ?", status)
	}
	if opts.NodeID != nil {
		base = base.Where("node_id = ?", *opts.NodeID)
	}
	if opts.ConfigID != nil {
		base = base.Where("protocol_config_id = ?", *opts.ConfigID)
	}
	if protocol := strings.TrimSpace(strings.ToLower(opts.Protocol)); protocol != "" {
		base = base.Joins("JOIN protocol_configs pc ON pc.id = protocol_bindings.protocol_config_id AND LOWER(pc.protocol) = ?", protocol)
	}

	countQuery := base.Session(&gorm.Session{})
	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []ProtocolBinding{}, 0, nil
	}

	orderClause := buildProtocolBindingOrderClause(opts.Sort, opts.Direction)
	offset := (opts.Page - 1) * opts.PerPage
	listQuery := base.Session(&gorm.Session{}).Order(orderClause).Limit(opts.PerPage).Offset(offset).
		Preload("Node").Preload("ProtocolConfig")

	var bindings []ProtocolBinding
	if err := listQuery.Find(&bindings).Error; err != nil {
		return nil, 0, err
	}

	return bindings, total, nil
}

func (r *protocolBindingRepository) ListByNodeIDs(ctx context.Context, nodeIDs []uint64) ([]ProtocolBinding, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if len(nodeIDs) == 0 {
		return []ProtocolBinding{}, nil
	}

	var bindings []ProtocolBinding
	if err := r.db.WithContext(ctx).
		Where("node_id IN ?", nodeIDs).
		Order("node_id ASC, updated_at DESC").
		Preload("ProtocolConfig").
		Find(&bindings).Error; err != nil {
		return nil, err
	}
	return bindings, nil
}

func (r *protocolBindingRepository) ListAll(ctx context.Context) ([]ProtocolBinding, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var bindings []ProtocolBinding
	if err := r.db.WithContext(ctx).
		Order("updated_at DESC").
		Preload("Node").
		Preload("ProtocolConfig").
		Find(&bindings).Error; err != nil {
		return nil, err
	}
	return bindings, nil
}

func (r *protocolBindingRepository) Get(ctx context.Context, id uint64) (ProtocolBinding, error) {
	if err := ctx.Err(); err != nil {
		return ProtocolBinding{}, err
	}

	var binding ProtocolBinding
	if err := r.db.WithContext(ctx).Preload("Node").Preload("ProtocolConfig").First(&binding, id).Error; err != nil {
		return ProtocolBinding{}, translateError(err)
	}
	return binding, nil
}

func (r *protocolBindingRepository) Create(ctx context.Context, binding ProtocolBinding) (ProtocolBinding, error) {
	if err := ctx.Err(); err != nil {
		return ProtocolBinding{}, err
	}

	binding.Name = strings.TrimSpace(binding.Name)
	binding.Role = strings.TrimSpace(binding.Role)
	binding.Listen = strings.TrimSpace(binding.Listen)
	binding.Connect = strings.TrimSpace(binding.Connect)
	binding.Status = strings.TrimSpace(binding.Status)
	binding.KernelID = strings.TrimSpace(binding.KernelID)
	binding.SyncStatus = strings.TrimSpace(binding.SyncStatus)
	binding.HealthStatus = strings.TrimSpace(binding.HealthStatus)
	binding.Description = strings.TrimSpace(binding.Description)
	if binding.NodeID == 0 || binding.ProtocolConfigID == 0 || binding.Role == "" {
		return ProtocolBinding{}, ErrInvalidArgument
	}
	if binding.Status == "" {
		binding.Status = "active"
	}
	if binding.SyncStatus == "" {
		binding.SyncStatus = "pending"
	}
	if binding.HealthStatus == "" {
		binding.HealthStatus = "unknown"
	}
	if binding.Tags == nil {
		binding.Tags = []string{}
	}
	if binding.Metadata == nil {
		binding.Metadata = map[string]any{}
	}

	now := time.Now().UTC()
	if binding.CreatedAt.IsZero() {
		binding.CreatedAt = now
	}
	binding.UpdatedAt = now

	if err := r.db.WithContext(ctx).Create(&binding).Error; err != nil {
		return ProtocolBinding{}, translateError(err)
	}

	if binding.KernelID == "" {
		kernelID := fmt.Sprintf("binding-%d", binding.ID)
		if err := r.db.WithContext(ctx).
			Model(&ProtocolBinding{}).
			Where("id = ?", binding.ID).
			Updates(map[string]any{"kernel_id": kernelID, "updated_at": time.Now().UTC()}).
			Error; err != nil {
			return ProtocolBinding{}, translateError(err)
		}
		binding.KernelID = kernelID
	}

	return r.Get(ctx, binding.ID)
}

func (r *protocolBindingRepository) Update(ctx context.Context, id uint64, input UpdateProtocolBindingInput) (ProtocolBinding, error) {
	if err := ctx.Err(); err != nil {
		return ProtocolBinding{}, err
	}

	updates := r.buildBindingUpdates(input)
	if len(updates) == 0 {
		return ProtocolBinding{}, ErrInvalidArgument
	}
	updates["updated_at"] = time.Now().UTC()

	if err := r.db.WithContext(ctx).Model(&ProtocolBinding{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return ProtocolBinding{}, translateError(err)
	}
	return r.Get(ctx, id)
}

func (r *protocolBindingRepository) UpdateSyncState(ctx context.Context, id uint64, input UpdateProtocolBindingInput) (ProtocolBinding, error) {
	if err := ctx.Err(); err != nil {
		return ProtocolBinding{}, err
	}

	updates := map[string]any{}
	if input.SyncStatus != nil {
		updates["sync_status"] = strings.TrimSpace(*input.SyncStatus)
	}
	if input.LastSyncedAt != nil {
		updates["last_synced_at"] = input.LastSyncedAt.UTC()
	}
	if input.LastSyncError != nil {
		updates["last_sync_error"] = strings.TrimSpace(*input.LastSyncError)
	}
	if input.KernelID != nil {
		updates["kernel_id"] = strings.TrimSpace(*input.KernelID)
	}
	if len(updates) == 0 {
		return ProtocolBinding{}, ErrInvalidArgument
	}
	updates["updated_at"] = time.Now().UTC()

	if err := r.db.WithContext(ctx).Model(&ProtocolBinding{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return ProtocolBinding{}, translateError(err)
	}
	return r.Get(ctx, id)
}

func (r *protocolBindingRepository) UpdateHealthByKernelID(ctx context.Context, kernelID string, status string, observedAt time.Time, message string) (ProtocolBinding, error) {
	if err := ctx.Err(); err != nil {
		return ProtocolBinding{}, err
	}
	kernelID = strings.TrimSpace(kernelID)
	if kernelID == "" {
		return ProtocolBinding{}, ErrInvalidArgument
	}
	status = strings.TrimSpace(status)
	if status == "" {
		status = "unknown"
	}

	updates := map[string]any{
		"health_status": status,
		"updated_at":    time.Now().UTC(),
	}
	if !observedAt.IsZero() {
		updates["last_heartbeat_at"] = observedAt.UTC()
	}
	if message != "" {
		updates["last_sync_error"] = strings.TrimSpace(message)
	}

	if err := r.db.WithContext(ctx).
		Model(&ProtocolBinding{}).
		Where("kernel_id = ?", kernelID).
		Updates(updates).Error; err != nil {
		return ProtocolBinding{}, translateError(err)
	}

	var binding ProtocolBinding
	if err := r.db.WithContext(ctx).Preload("Node").Preload("ProtocolConfig").
		Where("kernel_id = ?", kernelID).
		First(&binding).Error; err != nil {
		return ProtocolBinding{}, translateError(err)
	}
	return binding, nil
}

func (r *protocolBindingRepository) Delete(ctx context.Context, id uint64) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).Delete(&ProtocolBinding{}, id).Error; err != nil {
		return translateError(err)
	}
	return nil
}

func (r *protocolBindingRepository) buildBindingUpdates(input UpdateProtocolBindingInput) map[string]any {
	updates := map[string]any{}
	if input.Name != nil {
		updates["name"] = strings.TrimSpace(*input.Name)
	}
	if input.NodeID != nil {
		updates["node_id"] = *input.NodeID
	}
	if input.ProtocolConfigID != nil {
		updates["protocol_config_id"] = *input.ProtocolConfigID
	}
	if input.Role != nil {
		updates["role"] = strings.TrimSpace(*input.Role)
	}
	if input.Listen != nil {
		updates["listen"] = strings.TrimSpace(*input.Listen)
	}
	if input.Connect != nil {
		updates["connect"] = strings.TrimSpace(*input.Connect)
	}
	if input.Status != nil {
		updates["status"] = strings.TrimSpace(*input.Status)
	}
	if input.KernelID != nil {
		updates["kernel_id"] = strings.TrimSpace(*input.KernelID)
	}
	if input.SyncStatus != nil {
		updates["sync_status"] = strings.TrimSpace(*input.SyncStatus)
	}
	if input.HealthStatus != nil {
		updates["health_status"] = strings.TrimSpace(*input.HealthStatus)
	}
	if input.LastSyncedAt != nil {
		updates["last_synced_at"] = input.LastSyncedAt.UTC()
	}
	if input.LastHeartbeatAt != nil {
		updates["last_heartbeat_at"] = input.LastHeartbeatAt.UTC()
	}
	if input.LastSyncError != nil {
		updates["last_sync_error"] = strings.TrimSpace(*input.LastSyncError)
	}
	if input.Tags != nil {
		updates["tags"] = append([]string(nil), (*input.Tags)...)
	}
	if input.Description != nil {
		updates["description"] = strings.TrimSpace(*input.Description)
	}
	if input.Metadata != nil {
		updates["metadata"] = *input.Metadata
	}
	return updates
}

func normalizeListProtocolBindingsOptions(opts ListProtocolBindingsOptions) ListProtocolBindingsOptions {
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

func buildProtocolBindingOrderClause(field, direction string) string {
	column := "updated_at"
	switch strings.ToLower(field) {
	case "name":
		column = "name"
	case "status":
		column = "status"
	case "last_synced_at":
		column = "last_synced_at"
	case "last_heartbeat_at":
		column = "last_heartbeat_at"
	case "created_at":
		column = "created_at"
	}

	dir := "ASC"
	if strings.EqualFold(direction, "desc") {
		dir = "DESC"
	}

	return fmt.Sprintf("%s %s", column, dir)
}
