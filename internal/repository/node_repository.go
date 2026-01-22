package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/zero-net-panel/zero-net-panel/internal/status"
)

// Node 表示节点元信息。
type Node struct {
	ID                                        uint64         `gorm:"primaryKey"`
	Name                                      string         `gorm:"size:255;uniqueIndex"`
	Region                                    string         `gorm:"size:128"`
	Country                                   string         `gorm:"size:8"`
	ISP                                       string         `gorm:"size:128"`
	Status                                    int            `gorm:"column:status"`
	Tags                                      []string       `gorm:"serializer:json"`
	CapacityMbps                              int            `gorm:"column:capacity_mbps"`
	Description                               string         `gorm:"type:text"`
	AccessAddress                             string         `gorm:"size:512"`
	ControlEndpoint                           string         `gorm:"size:512"`
	ControlAccessKey                          string         `gorm:"size:255"`
	ControlSecretKey                          string         `gorm:"size:512"`
	ControlToken                              string         `gorm:"size:512"`
	KernelDefaultProtocol                     string         `gorm:"column:kernel_default_protocol;size:32;default:http"`
	KernelHTTPTimeoutSeconds                  int            `gorm:"column:kernel_http_timeout_seconds;default:5"`
	KernelStatusPollIntervalSeconds           int            `gorm:"column:kernel_status_poll_interval_seconds;default:30"`
	KernelStatusPollBackoffEnabled            bool           `gorm:"column:kernel_status_poll_backoff_enabled;default:true"`
	KernelStatusPollBackoffMaxIntervalSeconds int            `gorm:"column:kernel_status_poll_backoff_max_interval_seconds;default:300"`
	KernelStatusPollBackoffMultiplier         float64        `gorm:"column:kernel_status_poll_backoff_multiplier;default:2"`
	KernelStatusPollBackoffJitter             float64        `gorm:"column:kernel_status_poll_backoff_jitter;default:0.2"`
	KernelOfflineProbeMaxIntervalSeconds      int            `gorm:"column:kernel_offline_probe_max_interval_seconds;default:0"`
	StatusSyncEnabled                         bool           `gorm:"column:status_sync_enabled;default:true"`
	LastSyncedAt                              time.Time      `gorm:"column:last_synced_at"`
	DeletedAt                                 gorm.DeletedAt `gorm:"index"`
	UpdatedAt                                 time.Time
	CreatedAt                                 time.Time
}

// TableName 自定义节点表名。
func (Node) TableName() string { return "nodes" }

// NodeKernel 表示节点某一协议的配置摘要。
type NodeKernel struct {
	NodeID       uint64         `gorm:"primaryKey;autoIncrement:false"`
	Protocol     string         `gorm:"primaryKey;size:32"`
	Endpoint     string         `gorm:"size:512"`
	Revision     string         `gorm:"size:128"`
	Status       int            `gorm:"column:status"`
	Config       map[string]any `gorm:"serializer:json"`
	LastSyncedAt time.Time      `gorm:"column:last_synced_at"`
	UpdatedAt    time.Time
	CreatedAt    time.Time
}

// TableName 自定义节点内核表名。
func (NodeKernel) TableName() string { return "node_kernels" }

// ListNodesOptions 为节点列表提供过滤与排序选项。
type ListNodesOptions struct {
	Page      int
	PerPage   int
	Sort      string
	Direction string
	Query     string
	Status    int
	Protocol  string
	NodeIDs   []uint64
}

// NodeRepository 定义节点仓储接口。
type NodeRepository interface {
	List(ctx context.Context, opts ListNodesOptions) ([]Node, int64, error)
	ListAll(ctx context.Context) ([]Node, error)
	Get(ctx context.Context, nodeID uint64) (Node, error)
	Create(ctx context.Context, node Node) (Node, error)
	Update(ctx context.Context, nodeID uint64, input UpdateNodeInput) (Node, error)
	UpdateStatusByIDs(ctx context.Context, nodeIDs []uint64, status int) error
	Delete(ctx context.Context, nodeID uint64) error
	GetKernels(ctx context.Context, nodeID uint64) ([]NodeKernel, error)
	RecordKernelSync(ctx context.Context, nodeID uint64, kernel NodeKernel) (NodeKernel, error)
	UpsertKernel(ctx context.Context, nodeID uint64, input UpsertNodeKernelInput) (NodeKernel, error)
}

type nodeRepository struct {
	db *gorm.DB
}

// NewNodeRepository 创建节点仓储。
func NewNodeRepository(db *gorm.DB) (NodeRepository, error) {
	if db == nil {
		return nil, errors.New("repository: database connection is required")
	}
	return &nodeRepository{db: db}, nil
}

func (r *nodeRepository) List(ctx context.Context, opts ListNodesOptions) ([]Node, int64, error) {
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}

	opts = normalizeListNodesOptions(opts)

	base := r.db.WithContext(ctx).Model(&Node{})

	if query := strings.TrimSpace(strings.ToLower(opts.Query)); query != "" {
		like := fmt.Sprintf("%%%s%%", query)
		base = base.Where("(LOWER(name) LIKE ? OR LOWER(region) LIKE ? OR LOWER(description) LIKE ?)", like, like, like)
	}
	if opts.Status != 0 {
		base = base.Where("status = ?", opts.Status)
	}
	if len(opts.NodeIDs) > 0 {
		base = base.Where("id IN ?", opts.NodeIDs)
	}
	if protocol := strings.TrimSpace(strings.ToLower(opts.Protocol)); protocol != "" {
		base = base.Joins("JOIN node_kernels nk ON nk.node_id = nodes.id AND LOWER(nk.protocol) = ?", protocol)
	}

	countQuery := base.Session(&gorm.Session{}).Distinct("nodes.id")
	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []Node{}, 0, nil
	}

	orderClause := buildNodeOrderClause(opts.Sort, opts.Direction)
	offset := (opts.Page - 1) * opts.PerPage
	listQuery := base.Session(&gorm.Session{}).Distinct().Order(orderClause).Limit(opts.PerPage).Offset(offset)

	var nodes []Node
	if err := listQuery.Find(&nodes).Error; err != nil {
		return nil, 0, err
	}

	return nodes, total, nil
}

func (r *nodeRepository) ListAll(ctx context.Context) ([]Node, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var nodes []Node
	if err := r.db.WithContext(ctx).
		Order("updated_at DESC").
		Find(&nodes).Error; err != nil {
		return nil, err
	}
	return nodes, nil
}

func (r *nodeRepository) Get(ctx context.Context, nodeID uint64) (Node, error) {
	if err := ctx.Err(); err != nil {
		return Node{}, err
	}

	var node Node
	if err := r.db.WithContext(ctx).First(&node, nodeID).Error; err != nil {
		return Node{}, translateError(err)
	}

	return node, nil
}

func (r *nodeRepository) Create(ctx context.Context, node Node) (Node, error) {
	if err := ctx.Err(); err != nil {
		return Node{}, err
	}

	node.Name = strings.TrimSpace(node.Name)
	if node.Name == "" {
		return Node{}, ErrInvalidArgument
	}
	node.Region = strings.TrimSpace(node.Region)
	node.Country = strings.TrimSpace(node.Country)
	node.ISP = strings.TrimSpace(node.ISP)
	node.Description = strings.TrimSpace(node.Description)
	node.AccessAddress = strings.TrimSpace(node.AccessAddress)
	node.ControlEndpoint = strings.TrimSpace(node.ControlEndpoint)
	node.ControlToken = strings.TrimSpace(node.ControlToken)
	node.KernelDefaultProtocol = strings.TrimSpace(node.KernelDefaultProtocol)
	if node.Tags == nil {
		node.Tags = []string{}
	}

	now := time.Now().UTC()
	if node.CreatedAt.IsZero() {
		node.CreatedAt = now
	}
	node.UpdatedAt = now
	node.LastSyncedAt = NormalizeTime(node.LastSyncedAt)

	if err := r.db.WithContext(ctx).Create(&node).Error; err != nil {
		return Node{}, translateError(err)
	}
	return node, nil
}

func (r *nodeRepository) Update(ctx context.Context, nodeID uint64, input UpdateNodeInput) (Node, error) {
	if err := ctx.Err(); err != nil {
		return Node{}, err
	}

	updates := map[string]any{}
	if input.Name != nil {
		updates["name"] = strings.TrimSpace(*input.Name)
	}
	if input.Region != nil {
		updates["region"] = strings.TrimSpace(*input.Region)
	}
	if input.Country != nil {
		updates["country"] = strings.TrimSpace(*input.Country)
	}
	if input.ISP != nil {
		updates["isp"] = strings.TrimSpace(*input.ISP)
	}
	if input.Status != nil {
		updates["status"] = *input.Status
	}
	if input.Tags != nil {
		serialized, err := serializeTags(*input.Tags)
		if err != nil {
			return Node{}, ErrInvalidArgument
		}
		updates["tags"] = serialized
	}
	if input.CapacityMbps != nil {
		updates["capacity_mbps"] = *input.CapacityMbps
	}
	if input.Description != nil {
		updates["description"] = strings.TrimSpace(*input.Description)
	}
	if input.AccessAddress != nil {
		updates["access_address"] = strings.TrimSpace(*input.AccessAddress)
	}
	if input.ControlEndpoint != nil {
		updates["control_endpoint"] = strings.TrimSpace(*input.ControlEndpoint)
	}
	if input.ControlAccessKey != nil {
		updates["control_access_key"] = strings.TrimSpace(*input.ControlAccessKey)
	}
	if input.ControlSecretKey != nil {
		updates["control_secret_key"] = strings.TrimSpace(*input.ControlSecretKey)
	}
	if input.ControlToken != nil {
		updates["control_token"] = strings.TrimSpace(*input.ControlToken)
	}
	if input.KernelDefaultProtocol != nil {
		updates["kernel_default_protocol"] = strings.TrimSpace(*input.KernelDefaultProtocol)
	}
	if input.KernelHTTPTimeoutSeconds != nil {
		updates["kernel_http_timeout_seconds"] = *input.KernelHTTPTimeoutSeconds
	}
	if input.KernelStatusPollIntervalSeconds != nil {
		updates["kernel_status_poll_interval_seconds"] = *input.KernelStatusPollIntervalSeconds
	}
	if input.KernelStatusPollBackoffEnabled != nil {
		updates["kernel_status_poll_backoff_enabled"] = *input.KernelStatusPollBackoffEnabled
	}
	if input.KernelStatusPollBackoffMaxIntervalSeconds != nil {
		updates["kernel_status_poll_backoff_max_interval_seconds"] = *input.KernelStatusPollBackoffMaxIntervalSeconds
	}
	if input.KernelStatusPollBackoffMultiplier != nil {
		updates["kernel_status_poll_backoff_multiplier"] = *input.KernelStatusPollBackoffMultiplier
	}
	if input.KernelStatusPollBackoffJitter != nil {
		updates["kernel_status_poll_backoff_jitter"] = *input.KernelStatusPollBackoffJitter
	}
	if input.KernelOfflineProbeMaxIntervalSeconds != nil {
		updates["kernel_offline_probe_max_interval_seconds"] = *input.KernelOfflineProbeMaxIntervalSeconds
	}
	if input.StatusSyncEnabled != nil {
		updates["status_sync_enabled"] = *input.StatusSyncEnabled
	}
	if input.LastSyncedAt != nil {
		updates["last_synced_at"] = input.LastSyncedAt.UTC()
	}

	if len(updates) == 0 {
		return Node{}, ErrInvalidArgument
	}
	updates["updated_at"] = time.Now().UTC()

	if err := r.db.WithContext(ctx).Model(&Node{}).Where("id = ?", nodeID).Updates(updates).Error; err != nil {
		return Node{}, translateError(err)
	}

	return r.Get(ctx, nodeID)
}

func (r *nodeRepository) UpdateStatusByIDs(ctx context.Context, nodeIDs []uint64, statusCode int) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if len(nodeIDs) == 0 {
		return nil
	}
	if statusCode == 0 {
		return ErrInvalidArgument
	}
	return r.db.WithContext(ctx).
		Model(&Node{}).
		Where("id IN ?", nodeIDs).
		Where("status_sync_enabled = ?", true).
		Where("status != ?", status.NodeStatusDisabled).
		Updates(map[string]any{
			"status":     statusCode,
			"updated_at": time.Now().UTC(),
		}).Error
}

func serializeTags(tags []string) (string, error) {
	if tags == nil {
		tags = []string{}
	}
	raw, err := json.Marshal(tags)
	if err != nil {
		return "", err
	}
	return string(raw), nil
}

func (r *nodeRepository) Delete(ctx context.Context, nodeID uint64) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if nodeID == 0 {
		return ErrInvalidArgument
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("node_id = ?", nodeID).Delete(&NodeKernel{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&Node{}, nodeID).Error; err != nil {
			return translateError(err)
		}
		return nil
	})
}

// UpdateNodeInput defines mutable node fields.
type UpdateNodeInput struct {
	Name                                      *string
	Region                                    *string
	Country                                   *string
	ISP                                       *string
	Status                                    *int
	Tags                                      *[]string
	CapacityMbps                              *int
	Description                               *string
	AccessAddress                             *string
	ControlEndpoint                           *string
	ControlAccessKey                          *string
	ControlSecretKey                          *string
	ControlToken                              *string
	KernelDefaultProtocol                     *string
	KernelHTTPTimeoutSeconds                  *int
	KernelStatusPollIntervalSeconds           *int
	KernelStatusPollBackoffEnabled            *bool
	KernelStatusPollBackoffMaxIntervalSeconds *int
	KernelStatusPollBackoffMultiplier         *float64
	KernelStatusPollBackoffJitter             *float64
	KernelOfflineProbeMaxIntervalSeconds      *int
	StatusSyncEnabled                         *bool
	LastSyncedAt                              *time.Time
}

// UpsertNodeKernelInput defines node kernel configuration updates.
type UpsertNodeKernelInput struct {
	Protocol     string
	Endpoint     string
	Revision     *string
	Status       *int
	Config       map[string]any
	LastSyncedAt *time.Time
}

func (r *nodeRepository) GetKernels(ctx context.Context, nodeID uint64) ([]NodeKernel, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	var kernels []NodeKernel
	if err := r.db.WithContext(ctx).
		Where("node_id = ?", nodeID).
		Order("protocol ASC, last_synced_at DESC").
		Find(&kernels).Error; err != nil {
		return nil, err
	}

	return kernels, nil
}

func (r *nodeRepository) RecordKernelSync(ctx context.Context, nodeID uint64, kernel NodeKernel) (NodeKernel, error) {
	if err := ctx.Err(); err != nil {
		return NodeKernel{}, err
	}

	proto := normalizeProtocol(kernel.Protocol)
	now := time.Now().UTC()

	if kernel.LastSyncedAt.IsZero() {
		kernel.LastSyncedAt = now
	}
	if kernel.CreatedAt.IsZero() {
		kernel.CreatedAt = now
	}
	kernel.UpdatedAt = now
	kernel.Protocol = proto
	kernel.NodeID = nodeID
	if kernel.Status == 0 {
		kernel.Status = status.NodeKernelStatusSynced
	}
	if kernel.Config == nil {
		kernel.Config = map[string]any{}
	}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var node Node
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&node, nodeID).Error; err != nil {
			return err
		}

		var existing NodeKernel
		err := tx.Where("node_id = ? AND LOWER(protocol) = ?", nodeID, proto).First(&existing).Error
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			if err := tx.Create(&kernel).Error; err != nil {
				return err
			}
		case err != nil:
			return err
		default:
			existing.Endpoint = kernel.Endpoint
			existing.Revision = kernel.Revision
			existing.Status = kernel.Status
			existing.Config = kernel.Config
			existing.LastSyncedAt = kernel.LastSyncedAt
			existing.UpdatedAt = kernel.UpdatedAt
			if err := tx.Save(&existing).Error; err != nil {
				return err
			}
			kernel = existing
		}

		if kernel.LastSyncedAt.After(node.LastSyncedAt) {
			node.LastSyncedAt = kernel.LastSyncedAt
		}
		if kernel.UpdatedAt.After(node.UpdatedAt) {
			node.UpdatedAt = kernel.UpdatedAt
		}
		node.Status = status.NodeStatusOnline

		return tx.Save(&node).Error
	})

	if err != nil {
		return NodeKernel{}, translateError(err)
	}

	return kernel, nil
}

func (r *nodeRepository) UpsertKernel(ctx context.Context, nodeID uint64, input UpsertNodeKernelInput) (NodeKernel, error) {
	if err := ctx.Err(); err != nil {
		return NodeKernel{}, err
	}

	proto := normalizeProtocol(input.Protocol)
	endpoint := strings.TrimSpace(input.Endpoint)
	if endpoint == "" || proto == "" {
		return NodeKernel{}, ErrInvalidArgument
	}

	now := time.Now().UTC()
	kernel := NodeKernel{
		NodeID:   nodeID,
		Protocol: proto,
		Endpoint: endpoint,
	}
	if input.Revision != nil {
		kernel.Revision = strings.TrimSpace(*input.Revision)
	}
	if input.Status != nil {
		kernel.Status = *input.Status
	}
	if input.Config != nil {
		kernel.Config = input.Config
	} else {
		kernel.Config = map[string]any{}
	}
	if input.LastSyncedAt != nil {
		kernel.LastSyncedAt = input.LastSyncedAt.UTC()
	}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var node Node
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&node, nodeID).Error; err != nil {
			return err
		}

		var existing NodeKernel
		err := tx.Where("node_id = ? AND LOWER(protocol) = ?", nodeID, proto).First(&existing).Error
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			if kernel.CreatedAt.IsZero() {
				kernel.CreatedAt = now
			}
			kernel.UpdatedAt = now
			kernel.LastSyncedAt = NormalizeTime(kernel.LastSyncedAt)
			if kernel.Status == 0 {
				kernel.Status = status.NodeKernelStatusConfigured
			}
			if err := tx.Create(&kernel).Error; err != nil {
				return err
			}
		case err != nil:
			return err
		default:
			existing.Endpoint = kernel.Endpoint
			if input.Revision != nil {
				existing.Revision = kernel.Revision
			}
			if input.Status != nil {
				existing.Status = kernel.Status
			}
			if input.Config != nil {
				existing.Config = kernel.Config
			}
			if input.LastSyncedAt != nil {
				existing.LastSyncedAt = kernel.LastSyncedAt
			}
			existing.UpdatedAt = now
			if err := tx.Save(&existing).Error; err != nil {
				return err
			}
			kernel = existing
		}

		if input.LastSyncedAt != nil && kernel.LastSyncedAt.After(node.LastSyncedAt) {
			node.LastSyncedAt = kernel.LastSyncedAt
		}
		node.UpdatedAt = now

		return tx.Save(&node).Error
	})
	if err != nil {
		return NodeKernel{}, translateError(err)
	}

	return kernel, nil
}

func normalizeListNodesOptions(opts ListNodesOptions) ListNodesOptions {
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
	opts.Sort = strings.ToLower(opts.Sort)
	if opts.Direction == "" {
		opts.Direction = "desc"
	}
	return opts
}

func buildNodeOrderClause(field, direction string) string {
	column := "nodes.updated_at"
	switch strings.ToLower(field) {
	case "name":
		column = "nodes.name"
	case "region":
		column = "nodes.region"
	case "last_synced_at":
		column = "nodes.last_synced_at"
	case "capacity_mbps":
		column = "nodes.capacity_mbps"
	}

	dir := "ASC"
	if strings.EqualFold(direction, "desc") {
		dir = "DESC"
	}

	return fmt.Sprintf("%s %s", column, dir)
}

func normalizeProtocol(protocol string) string {
	proto := strings.TrimSpace(strings.ToLower(protocol))
	if proto == "" {
		return "http"
	}
	return proto
}
