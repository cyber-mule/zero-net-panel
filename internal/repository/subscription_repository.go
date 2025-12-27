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
)

// Subscription 表示用户订阅信息。
type Subscription struct {
	ID                   uint64 `gorm:"primaryKey"`
	UserID               uint64 `gorm:"index"`
	Name                 string `gorm:"size:255"`
	PlanName             string `gorm:"size:255"`
	Status               string `gorm:"size:32"`
	TemplateID           uint64
	AvailableTemplateIDs []uint64 `gorm:"serializer:json"`
	Token                string   `gorm:"size:255"`
	ExpiresAt            time.Time
	TrafficTotalBytes    int64
	TrafficUsedBytes     int64
	DevicesLimit         int
	LastRefreshedAt      time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

// TableName 自定义订阅表名。
func (Subscription) TableName() string { return "subscriptions" }

// ListSubscriptionsOptions 控制订阅列表的分页与过滤。
type ListSubscriptionsOptions struct {
	Page       int
	PerPage    int
	Sort       string
	Direction  string
	Query      string
	Status     string
	UserID     *uint64
	PlanName   string
	TemplateID uint64
}

// SubscriptionRepository 提供订阅相关操作。
type SubscriptionRepository interface {
	List(ctx context.Context, opts ListSubscriptionsOptions) ([]Subscription, int64, error)
	ListByUser(ctx context.Context, userID uint64, opts ListSubscriptionsOptions) ([]Subscription, int64, error)
	Get(ctx context.Context, id uint64) (Subscription, error)
	GetActiveByUser(ctx context.Context, userID uint64) (Subscription, error)
	Create(ctx context.Context, sub Subscription) (Subscription, error)
	Update(ctx context.Context, id uint64, input UpdateSubscriptionInput) (Subscription, error)
	UpdateTemplate(ctx context.Context, subscriptionID uint64, templateID uint64, userID uint64) (Subscription, error)
	IncrementTrafficUsage(ctx context.Context, id uint64, delta int64) (Subscription, error)
}

type subscriptionRepository struct {
	db           *gorm.DB
	templateRepo SubscriptionTemplateRepository
}

// NewSubscriptionRepository 创建订阅仓储。
func NewSubscriptionRepository(db *gorm.DB, templateRepo SubscriptionTemplateRepository) (SubscriptionRepository, error) {
	if db == nil {
		return nil, errors.New("repository: database connection is required")
	}
	if templateRepo == nil {
		return nil, errors.New("repository: template repository is required")
	}
	return &subscriptionRepository{db: db, templateRepo: templateRepo}, nil
}

func (r *subscriptionRepository) ListByUser(ctx context.Context, userID uint64, opts ListSubscriptionsOptions) ([]Subscription, int64, error) {
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}

	opts = normalizeListSubscriptionsOptions(opts)
	opts.UserID = &userID
	return r.listWithOptions(ctx, opts)
}

func (r *subscriptionRepository) List(ctx context.Context, opts ListSubscriptionsOptions) ([]Subscription, int64, error) {
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}

	opts = normalizeListSubscriptionsOptions(opts)
	return r.listWithOptions(ctx, opts)
}

func (r *subscriptionRepository) Get(ctx context.Context, id uint64) (Subscription, error) {
	if err := ctx.Err(); err != nil {
		return Subscription{}, err
	}

	var subscription Subscription
	if err := r.db.WithContext(ctx).First(&subscription, id).Error; err != nil {
		return Subscription{}, translateError(err)
	}

	return subscription, nil
}

func (r *subscriptionRepository) GetActiveByUser(ctx context.Context, userID uint64) (Subscription, error) {
	if err := ctx.Err(); err != nil {
		return Subscription{}, err
	}
	if userID == 0 {
		return Subscription{}, ErrInvalidArgument
	}

	var subscription Subscription
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND status = ?", userID, "active").
		Where("expires_at > ?", time.Now().UTC()).
		Order("expires_at DESC").
		First(&subscription).Error; err != nil {
		return Subscription{}, translateError(err)
	}

	return subscription, nil
}

// UpdateSubscriptionInput defines mutable subscription fields.
type UpdateSubscriptionInput struct {
	Name                 *string
	PlanName             *string
	Status               *string
	TemplateID           *uint64
	AvailableTemplateIDs *[]uint64
	Token                *string
	ExpiresAt            *time.Time
	TrafficTotalBytes    *int64
	TrafficUsedBytes     *int64
	DevicesLimit         *int
	LastRefreshedAt      *time.Time
}

func (r *subscriptionRepository) Create(ctx context.Context, sub Subscription) (Subscription, error) {
	if err := ctx.Err(); err != nil {
		return Subscription{}, err
	}

	if sub.UserID == 0 || strings.TrimSpace(sub.Name) == "" || strings.TrimSpace(sub.PlanName) == "" {
		return Subscription{}, ErrInvalidArgument
	}
	now := time.Now().UTC()
	if sub.CreatedAt.IsZero() {
		sub.CreatedAt = now
	}
	sub.UpdatedAt = now
	if sub.Status == "" {
		sub.Status = "active"
	}
	if sub.LastRefreshedAt.IsZero() {
		sub.LastRefreshedAt = now
	}

	if err := r.db.WithContext(ctx).Create(&sub).Error; err != nil {
		return Subscription{}, translateError(err)
	}
	return sub, nil
}

func (r *subscriptionRepository) Update(ctx context.Context, id uint64, input UpdateSubscriptionInput) (Subscription, error) {
	if err := ctx.Err(); err != nil {
		return Subscription{}, err
	}

	updates := map[string]any{}
	if input.Name != nil {
		updates["name"] = strings.TrimSpace(*input.Name)
	}
	if input.PlanName != nil {
		updates["plan_name"] = strings.TrimSpace(*input.PlanName)
	}
	if input.Status != nil {
		updates["status"] = strings.TrimSpace(*input.Status)
	}
	if input.TemplateID != nil {
		updates["template_id"] = *input.TemplateID
	}
	if input.AvailableTemplateIDs != nil {
		payload, err := json.Marshal(*input.AvailableTemplateIDs)
		if err != nil {
			return Subscription{}, err
		}
		updates["available_template_ids"] = string(payload)
	}
	if input.Token != nil {
		updates["token"] = strings.TrimSpace(*input.Token)
	}
	if input.ExpiresAt != nil {
		updates["expires_at"] = *input.ExpiresAt
	}
	if input.TrafficTotalBytes != nil {
		updates["traffic_total_bytes"] = *input.TrafficTotalBytes
	}
	if input.TrafficUsedBytes != nil {
		updates["traffic_used_bytes"] = *input.TrafficUsedBytes
	}
	if input.DevicesLimit != nil {
		updates["devices_limit"] = *input.DevicesLimit
	}
	if input.LastRefreshedAt != nil {
		updates["last_refreshed_at"] = *input.LastRefreshedAt
	}
	if len(updates) == 0 {
		return Subscription{}, ErrInvalidArgument
	}

	updates["updated_at"] = time.Now().UTC()

	if err := r.db.WithContext(ctx).Model(&Subscription{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return Subscription{}, translateError(err)
	}

	return r.Get(ctx, id)
}

func (r *subscriptionRepository) IncrementTrafficUsage(ctx context.Context, id uint64, delta int64) (Subscription, error) {
	if err := ctx.Err(); err != nil {
		return Subscription{}, err
	}
	if id == 0 {
		return Subscription{}, ErrInvalidArgument
	}

	updates := map[string]any{
		"traffic_used_bytes": gorm.Expr("traffic_used_bytes + ?", delta),
		"updated_at":         time.Now().UTC(),
	}
	if err := r.db.WithContext(ctx).Model(&Subscription{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return Subscription{}, translateError(err)
	}
	return r.Get(ctx, id)
}

func (r *subscriptionRepository) UpdateTemplate(ctx context.Context, subscriptionID uint64, templateID uint64, userID uint64) (Subscription, error) {
	if err := ctx.Err(); err != nil {
		return Subscription{}, err
	}

	var subscription Subscription

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&subscription, subscriptionID).Error; err != nil {
			return err
		}

		if subscription.UserID != userID {
			return ErrForbidden
		}

		targetTemplate := templateID
		if targetTemplate == 0 {
			targetTemplate = subscription.TemplateID
		}

		if len(subscription.AvailableTemplateIDs) > 0 {
			allowed := false
			for _, id := range subscription.AvailableTemplateIDs {
				if id == targetTemplate {
					allowed = true
					break
				}
			}
			if !allowed {
				return ErrForbidden
			}
		}

		var tpl SubscriptionTemplate
		if err := tx.First(&tpl, targetTemplate).Error; err != nil {
			return err
		}

		subscription.TemplateID = tpl.ID
		now := time.Now().UTC()
		subscription.LastRefreshedAt = now
		subscription.UpdatedAt = now

		return tx.Save(&subscription).Error
	})

	if err != nil {
		switch {
		case errors.Is(err, ErrForbidden):
			return Subscription{}, err
		default:
			return Subscription{}, translateError(err)
		}
	}

	return subscription, nil
}

func (r *subscriptionRepository) listWithOptions(ctx context.Context, opts ListSubscriptionsOptions) ([]Subscription, int64, error) {
	base := r.db.WithContext(ctx).Model(&Subscription{})
	if opts.UserID != nil {
		base = base.Where("user_id = ?", *opts.UserID)
	}

	if query := strings.TrimSpace(strings.ToLower(opts.Query)); query != "" {
		like := fmt.Sprintf("%%%s%%", query)
		base = base.Where("(LOWER(name) LIKE ? OR LOWER(plan_name) LIKE ?)", like, like)
	}
	if status := strings.TrimSpace(strings.ToLower(opts.Status)); status != "" {
		base = base.Where("LOWER(status) = ?", status)
	}
	if planName := strings.TrimSpace(strings.ToLower(opts.PlanName)); planName != "" {
		base = base.Where("LOWER(plan_name) = ?", planName)
	}
	if opts.TemplateID != 0 {
		base = base.Where("template_id = ?", opts.TemplateID)
	}

	countQuery := base.Session(&gorm.Session{})
	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []Subscription{}, 0, nil
	}

	orderClause := buildSubscriptionOrderClause(opts.Sort, opts.Direction)
	offset := (opts.Page - 1) * opts.PerPage
	listQuery := base.Session(&gorm.Session{}).Order(orderClause).Limit(opts.PerPage).Offset(offset)

	var subscriptions []Subscription
	if err := listQuery.Find(&subscriptions).Error; err != nil {
		return nil, 0, err
	}

	return subscriptions, total, nil
}

func buildSubscriptionOrderClause(field, direction string) string {
	column := "updated_at"
	switch strings.ToLower(field) {
	case "name":
		column = "name"
	case "plan_name":
		column = "plan_name"
	case "status":
		column = "status"
	case "expires_at":
		column = "expires_at"
	case "created_at":
		column = "created_at"
	}

	dir := "ASC"
	if strings.EqualFold(direction, "desc") {
		dir = "DESC"
	}

	return fmt.Sprintf("%s %s", column, dir)
}

func normalizeListSubscriptionsOptions(opts ListSubscriptionsOptions) ListSubscriptionsOptions {
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
