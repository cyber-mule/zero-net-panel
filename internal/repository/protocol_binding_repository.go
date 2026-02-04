package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/zero-net-panel/zero-net-panel/internal/status"
)

// ProtocolBinding binds a protocol configuration to a node instance.
type ProtocolBinding struct {
	ID              uint64         `gorm:"primaryKey"`
	Name            string         `gorm:"size:255"`
	NodeID          uint64         `gorm:"index"`
	Protocol        string         `gorm:"size:32;index"`
	Role            string         `gorm:"size:32"`
	Listen          string         `gorm:"size:512"`
	Connect         string         `gorm:"size:512"`
	AccessPort      int            `gorm:"column:access_port"`
	Status          int            `gorm:"column:status"`
	KernelID        string         `gorm:"size:128;index"`
	SyncStatus      int            `gorm:"column:sync_status"`
	HealthStatus    int            `gorm:"column:health_status"`
	LastSyncedAt    time.Time      `gorm:"column:last_synced_at"`
	LastHeartbeatAt time.Time      `gorm:"column:last_heartbeat_at"`
	LastSyncError   string         `gorm:"type:text"`
	Tags            []string       `gorm:"serializer:json"`
	Description     string         `gorm:"type:text"`
	Profile         map[string]any `gorm:"serializer:json"`
	Metadata        map[string]any `gorm:"serializer:json"`
	UpdatedAt       time.Time
	CreatedAt       time.Time

	Node Node `gorm:"foreignKey:NodeID;references:ID"`
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
	Status    int
	Protocol  string
	NodeID    *uint64
}

// UpdateProtocolBindingInput defines mutable binding fields.
type UpdateProtocolBindingInput struct {
	Name            *string
	NodeID          *uint64
	Protocol        *string
	Role            *string
	Listen          *string
	Connect         *string
	AccessPort      *int
	Status          *int
	KernelID        *string
	SyncStatus      *int
	HealthStatus    *int
	LastSyncedAt    *time.Time
	LastHeartbeatAt *time.Time
	LastSyncError   *string
	Tags            *[]string
	Description     *string
	Profile         *map[string]any
	Metadata        *map[string]any
}

// ProtocolBindingRepository manages protocol binding persistence.
type ProtocolBindingRepository interface {
	List(ctx context.Context, opts ListProtocolBindingsOptions) ([]ProtocolBinding, int64, error)
	ListByIDs(ctx context.Context, ids []uint64) ([]ProtocolBinding, error)
	ListByNodeIDs(ctx context.Context, nodeIDs []uint64) ([]ProtocolBinding, error)
	ListAll(ctx context.Context) ([]ProtocolBinding, error)
	ListProtocols(ctx context.Context) ([]string, error)
	Get(ctx context.Context, id uint64) (ProtocolBinding, error)
	Create(ctx context.Context, binding ProtocolBinding) (ProtocolBinding, error)
	Update(ctx context.Context, id uint64, input UpdateProtocolBindingInput) (ProtocolBinding, error)
	UpdateSyncState(ctx context.Context, id uint64, input UpdateProtocolBindingInput) (ProtocolBinding, error)
	UpdateHealthByKernelID(ctx context.Context, kernelID string, statusCode int, observedAt time.Time, message string) (ProtocolBinding, error)
	UpdateHealthByKernelIDForNodes(ctx context.Context, kernelID string, nodeIDs []uint64, statusCode int, observedAt time.Time, message string) (ProtocolBinding, error)
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
	if opts.Status != 0 {
		base = base.Where("status = ?", opts.Status)
	}
	if opts.NodeID != nil {
		base = base.Where("node_id = ?", *opts.NodeID)
	}
	if protocol := strings.TrimSpace(strings.ToLower(opts.Protocol)); protocol != "" {
		base = base.Where("LOWER(protocol) = ?", protocol)
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
		Preload("Node")

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
		Preload("Node").
		Find(&bindings).Error; err != nil {
		return nil, err
	}
	return bindings, nil
}

func (r *protocolBindingRepository) ListByIDs(ctx context.Context, ids []uint64) ([]ProtocolBinding, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return []ProtocolBinding{}, nil
	}

	var bindings []ProtocolBinding
	if err := r.db.WithContext(ctx).
		Where("id IN ?", ids).
		Preload("Node").
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
		Find(&bindings).Error; err != nil {
		return nil, err
	}
	return bindings, nil
}

func (r *protocolBindingRepository) ListProtocols(ctx context.Context) ([]string, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var protocols []string
	if err := r.db.WithContext(ctx).
		Model(&ProtocolBinding{}).
		Where("protocol <> ''").
		Distinct("LOWER(protocol)").
		Order("LOWER(protocol) ASC").
		Pluck("LOWER(protocol)", &protocols).Error; err != nil {
		return nil, err
	}

	if protocols == nil {
		protocols = []string{}
	}
	return protocols, nil
}

func (r *protocolBindingRepository) Get(ctx context.Context, id uint64) (ProtocolBinding, error) {
	if err := ctx.Err(); err != nil {
		return ProtocolBinding{}, err
	}

	var binding ProtocolBinding
	if err := r.db.WithContext(ctx).Preload("Node").First(&binding, id).Error; err != nil {
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
	binding.KernelID = strings.TrimSpace(binding.KernelID)
	binding.Protocol = strings.ToLower(strings.TrimSpace(binding.Protocol))
	binding.Description = strings.TrimSpace(binding.Description)
	if binding.NodeID == 0 || binding.Protocol == "" || binding.Role == "" {
		return ProtocolBinding{}, ErrInvalidArgument
	}
	if binding.Status == 0 {
		binding.Status = status.ProtocolBindingStatusActive
	}
	if binding.SyncStatus == 0 {
		binding.SyncStatus = status.ProtocolBindingSyncStatusPending
	}
	if binding.HealthStatus == 0 {
		binding.HealthStatus = status.ProtocolBindingHealthStatusUnknown
	}
	if binding.Tags == nil {
		binding.Tags = []string{}
	}
	if binding.Profile == nil {
		binding.Profile = map[string]any{}
	}
	if binding.Metadata == nil {
		binding.Metadata = map[string]any{}
	}

	now := time.Now().UTC()
	if binding.CreatedAt.IsZero() {
		binding.CreatedAt = now
	}
	binding.UpdatedAt = now
	binding.LastSyncedAt = NormalizeTime(binding.LastSyncedAt)
	binding.LastHeartbeatAt = NormalizeTime(binding.LastHeartbeatAt)

	if err := r.db.WithContext(ctx).Create(&binding).Error; err != nil {
		return ProtocolBinding{}, translateError(err)
	}

	return r.Get(ctx, binding.ID)
}

func (r *protocolBindingRepository) Update(ctx context.Context, id uint64, input UpdateProtocolBindingInput) (ProtocolBinding, error) {
	if err := ctx.Err(); err != nil {
		return ProtocolBinding{}, err
	}

	updates, err := r.buildBindingUpdates(input)
	if err != nil {
		return ProtocolBinding{}, ErrInvalidArgument
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

func (r *protocolBindingRepository) UpdateSyncState(ctx context.Context, id uint64, input UpdateProtocolBindingInput) (ProtocolBinding, error) {
	if err := ctx.Err(); err != nil {
		return ProtocolBinding{}, err
	}

	updates := map[string]any{}
	if input.SyncStatus != nil {
		updates["sync_status"] = *input.SyncStatus
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

func (r *protocolBindingRepository) UpdateHealthByKernelID(ctx context.Context, kernelID string, statusCode int, observedAt time.Time, message string) (ProtocolBinding, error) {
	if err := ctx.Err(); err != nil {
		return ProtocolBinding{}, err
	}
	kernelID = strings.TrimSpace(kernelID)
	if kernelID == "" {
		return ProtocolBinding{}, ErrInvalidArgument
	}
	if statusCode == 0 {
		statusCode = status.ProtocolBindingHealthStatusUnknown
	}

	updates := map[string]any{
		"health_status": statusCode,
		"updated_at":    time.Now().UTC(),
	}
	if !observedAt.IsZero() {
		updates["last_heartbeat_at"] = observedAt.UTC()
	}
	if message != "" {
		updates["last_sync_error"] = strings.TrimSpace(message)
	}

	nodeFilter := r.db.WithContext(ctx).Model(&Node{}).Select("id").Where("status_sync_enabled = ?", true)
	result := r.db.WithContext(ctx).
		Model(&ProtocolBinding{}).
		Where("kernel_id = ?", kernelID).
		Where("node_id IN (?)", nodeFilter).
		Updates(updates)
	if result.Error != nil {
		return ProtocolBinding{}, translateError(result.Error)
	}
	if result.RowsAffected == 0 {
		return ProtocolBinding{}, ErrNotFound
	}

	var binding ProtocolBinding
	if err := r.db.WithContext(ctx).Preload("Node").
		Where("kernel_id = ?", kernelID).
		Where("node_id IN (?)", nodeFilter).
		First(&binding).Error; err != nil {
		return ProtocolBinding{}, translateError(err)
	}
	return binding, nil
}

func (r *protocolBindingRepository) UpdateHealthByKernelIDForNodes(ctx context.Context, kernelID string, nodeIDs []uint64, statusCode int, observedAt time.Time, message string) (ProtocolBinding, error) {
	if err := ctx.Err(); err != nil {
		return ProtocolBinding{}, err
	}
	kernelID = strings.TrimSpace(kernelID)
	if kernelID == "" || len(nodeIDs) == 0 {
		return ProtocolBinding{}, ErrInvalidArgument
	}
	if statusCode == 0 {
		statusCode = status.ProtocolBindingHealthStatusUnknown
	}

	updates := map[string]any{
		"health_status": statusCode,
		"updated_at":    time.Now().UTC(),
	}
	if !observedAt.IsZero() {
		updates["last_heartbeat_at"] = observedAt.UTC()
	}
	if message != "" {
		updates["last_sync_error"] = strings.TrimSpace(message)
	}

	result := r.db.WithContext(ctx).
		Model(&ProtocolBinding{}).
		Where("kernel_id = ?", kernelID).
		Where("node_id IN ?", nodeIDs).
		Updates(updates)
	if result.Error != nil {
		return ProtocolBinding{}, translateError(result.Error)
	}
	if result.RowsAffected == 0 {
		return ProtocolBinding{}, ErrNotFound
	}

	var binding ProtocolBinding
	if err := r.db.WithContext(ctx).Preload("Node").
		Where("kernel_id = ?", kernelID).
		Where("node_id IN ?", nodeIDs).
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

func (r *protocolBindingRepository) buildBindingUpdates(input UpdateProtocolBindingInput) (map[string]any, error) {
	updates := map[string]any{}
	if input.Name != nil {
		updates["name"] = strings.TrimSpace(*input.Name)
	}
	if input.NodeID != nil {
		updates["node_id"] = *input.NodeID
	}
	if input.Protocol != nil {
		updates["protocol"] = strings.ToLower(strings.TrimSpace(*input.Protocol))
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
	if input.AccessPort != nil {
		updates["access_port"] = *input.AccessPort
	}
	if input.Status != nil {
		updates["status"] = *input.Status
	}
	if input.KernelID != nil {
		updates["kernel_id"] = strings.TrimSpace(*input.KernelID)
	}
	if input.SyncStatus != nil {
		updates["sync_status"] = *input.SyncStatus
	}
	if input.HealthStatus != nil {
		updates["health_status"] = *input.HealthStatus
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
		serialized, err := serializeStringSlice(*input.Tags)
		if err != nil {
			return nil, err
		}
		updates["tags"] = serialized
	}
	if input.Description != nil {
		updates["description"] = strings.TrimSpace(*input.Description)
	}
	if input.Profile != nil {
		serialized, err := serializeAnyMap(*input.Profile)
		if err != nil {
			return nil, err
		}
		updates["profile"] = serialized
	}
	if input.Metadata != nil {
		serialized, err := serializeAnyMap(*input.Metadata)
		if err != nil {
			return nil, err
		}
		updates["metadata"] = serialized
	}
	return updates, nil
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
