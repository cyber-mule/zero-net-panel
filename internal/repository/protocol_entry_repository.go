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

// ProtocolEntry publishes a protocol binding to user-facing endpoints.
type ProtocolEntry struct {
	ID           uint64         `gorm:"primaryKey"`
	Name         string         `gorm:"size:255"`
	BindingID    uint64         `gorm:"index"`
	Protocol     string         `gorm:"size:32;index"`
	Status       int            `gorm:"column:status"`
	EntryAddress string         `gorm:"size:512"`
	EntryPort    int            `gorm:"column:entry_port"`
	Tags         []string       `gorm:"serializer:json"`
	Description  string         `gorm:"type:text"`
	Profile      map[string]any `gorm:"serializer:json"`
	CreatedAt    time.Time
	UpdatedAt    time.Time

	Binding ProtocolBinding `gorm:"foreignKey:BindingID;references:ID"`
}

// TableName binds the protocol entry table name.
func (ProtocolEntry) TableName() string { return "protocol_entries" }

// ListProtocolEntriesOptions controls filtering and pagination.
type ListProtocolEntriesOptions struct {
	Page      int
	PerPage   int
	Sort      string
	Direction string
	Query     string
	Status    int
	Protocol  string
	BindingID *uint64
}

// UpdateProtocolEntryInput defines mutable entry fields.
type UpdateProtocolEntryInput struct {
	Name         *string
	BindingID    *uint64
	Protocol     *string
	Status       *int
	EntryAddress *string
	EntryPort    *int
	Tags         *[]string
	Description  *string
	Profile      *map[string]any
}

// ProtocolEntryRepository manages protocol entry persistence.
type ProtocolEntryRepository interface {
	List(ctx context.Context, opts ListProtocolEntriesOptions) ([]ProtocolEntry, int64, error)
	ListByBindingIDs(ctx context.Context, bindingIDs []uint64) ([]ProtocolEntry, error)
	Get(ctx context.Context, id uint64) (ProtocolEntry, error)
	Create(ctx context.Context, entry ProtocolEntry) (ProtocolEntry, error)
	Update(ctx context.Context, id uint64, input UpdateProtocolEntryInput) (ProtocolEntry, error)
	Delete(ctx context.Context, id uint64) error
}

type protocolEntryRepository struct {
	db *gorm.DB
}

// NewProtocolEntryRepository constructs a protocol entry repository.
func NewProtocolEntryRepository(db *gorm.DB) (ProtocolEntryRepository, error) {
	if db == nil {
		return nil, errors.New("repository: database connection is required")
	}
	return &protocolEntryRepository{db: db}, nil
}

func (r *protocolEntryRepository) List(ctx context.Context, opts ListProtocolEntriesOptions) ([]ProtocolEntry, int64, error) {
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}

	opts = normalizeListProtocolEntriesOptions(opts)
	base := r.db.WithContext(ctx).Model(&ProtocolEntry{})

	if query := strings.TrimSpace(strings.ToLower(opts.Query)); query != "" {
		like := fmt.Sprintf("%%%s%%", query)
		base = base.Where("(LOWER(name) LIKE ? OR LOWER(description) LIKE ?)", like, like)
	}
	if opts.Status != 0 {
		base = base.Where("status = ?", opts.Status)
	}
	if protocol := strings.TrimSpace(strings.ToLower(opts.Protocol)); protocol != "" {
		base = base.Where("LOWER(protocol) = ?", protocol)
	}
	if opts.BindingID != nil {
		base = base.Where("binding_id = ?", *opts.BindingID)
	}

	countQuery := base.Session(&gorm.Session{})
	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []ProtocolEntry{}, 0, nil
	}

	orderClause := buildProtocolEntryOrderClause(opts.Sort, opts.Direction)
	offset := (opts.Page - 1) * opts.PerPage
	listQuery := base.Session(&gorm.Session{}).Order(orderClause).Limit(opts.PerPage).Offset(offset).
		Preload("Binding").Preload("Binding.Node")

	var entries []ProtocolEntry
	if err := listQuery.Find(&entries).Error; err != nil {
		return nil, 0, err
	}

	return entries, total, nil
}

func (r *protocolEntryRepository) ListByBindingIDs(ctx context.Context, bindingIDs []uint64) ([]ProtocolEntry, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if len(bindingIDs) == 0 {
		return []ProtocolEntry{}, nil
	}

	var entries []ProtocolEntry
	if err := r.db.WithContext(ctx).
		Where("binding_id IN ?", bindingIDs).
		Order("binding_id ASC, updated_at DESC").
		Preload("Binding").
		Preload("Binding.Node").
		Find(&entries).Error; err != nil {
		return nil, err
	}
	return entries, nil
}

func (r *protocolEntryRepository) Get(ctx context.Context, id uint64) (ProtocolEntry, error) {
	if err := ctx.Err(); err != nil {
		return ProtocolEntry{}, err
	}

	var entry ProtocolEntry
	if err := r.db.WithContext(ctx).Preload("Binding").Preload("Binding.Node").First(&entry, id).Error; err != nil {
		return ProtocolEntry{}, translateError(err)
	}
	return entry, nil
}

func (r *protocolEntryRepository) Create(ctx context.Context, entry ProtocolEntry) (ProtocolEntry, error) {
	if err := ctx.Err(); err != nil {
		return ProtocolEntry{}, err
	}

	entry.Name = strings.TrimSpace(entry.Name)
	entry.Protocol = strings.ToLower(strings.TrimSpace(entry.Protocol))
	entry.EntryAddress = strings.TrimSpace(entry.EntryAddress)
	entry.Description = strings.TrimSpace(entry.Description)

	if entry.Name == "" || entry.BindingID == 0 || entry.Protocol == "" {
		return ProtocolEntry{}, ErrInvalidArgument
	}
	if entry.EntryAddress == "" || entry.EntryPort <= 0 {
		return ProtocolEntry{}, ErrInvalidArgument
	}
	if entry.Status == 0 {
		entry.Status = status.ProtocolEntryStatusActive
	}
	if entry.EntryPort < 0 {
		return ProtocolEntry{}, ErrInvalidArgument
	}
	if entry.Tags == nil {
		entry.Tags = []string{}
	}
	if entry.Profile == nil {
		entry.Profile = map[string]any{}
	}

	now := time.Now().UTC()
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = now
	}
	entry.UpdatedAt = now

	if err := r.db.WithContext(ctx).Create(&entry).Error; err != nil {
		return ProtocolEntry{}, translateError(err)
	}
	return r.Get(ctx, entry.ID)
}

func (r *protocolEntryRepository) Update(ctx context.Context, id uint64, input UpdateProtocolEntryInput) (ProtocolEntry, error) {
	if err := ctx.Err(); err != nil {
		return ProtocolEntry{}, err
	}

	updates := r.buildEntryUpdates(input)
	if len(updates) == 0 {
		return ProtocolEntry{}, ErrInvalidArgument
	}
	updates["updated_at"] = time.Now().UTC()

	if err := r.db.WithContext(ctx).Model(&ProtocolEntry{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return ProtocolEntry{}, translateError(err)
	}
	return r.Get(ctx, id)
}

func (r *protocolEntryRepository) Delete(ctx context.Context, id uint64) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).Delete(&ProtocolEntry{}, id).Error; err != nil {
		return translateError(err)
	}
	return nil
}

func (r *protocolEntryRepository) buildEntryUpdates(input UpdateProtocolEntryInput) map[string]any {
	updates := map[string]any{}
	if input.Name != nil {
		updates["name"] = strings.TrimSpace(*input.Name)
	}
	if input.BindingID != nil {
		updates["binding_id"] = *input.BindingID
	}
	if input.Protocol != nil {
		updates["protocol"] = strings.ToLower(strings.TrimSpace(*input.Protocol))
	}
	if input.Status != nil {
		updates["status"] = *input.Status
	}
	if input.EntryAddress != nil {
		updates["entry_address"] = strings.TrimSpace(*input.EntryAddress)
	}
	if input.EntryPort != nil {
		updates["entry_port"] = *input.EntryPort
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
	return updates
}

func normalizeListProtocolEntriesOptions(opts ListProtocolEntriesOptions) ListProtocolEntriesOptions {
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

func buildProtocolEntryOrderClause(field, direction string) string {
	column := "updated_at"
	switch strings.ToLower(field) {
	case "name":
		column = "name"
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
