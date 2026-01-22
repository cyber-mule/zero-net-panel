package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/zero-net-panel/zero-net-panel/internal/status"
)

const (
	CouponStatusActive   = status.CouponStatusActive
	CouponStatusDisabled = status.CouponStatusDisabled

	CouponTypePercent = "percent"
	CouponTypeFixed   = "fixed"

	CouponRedemptionReserved = status.CouponRedemptionStatusReserved
	CouponRedemptionApplied  = status.CouponRedemptionStatusApplied
	CouponRedemptionReleased = status.CouponRedemptionStatusReleased
)

// Coupon defines a discount rule.
type Coupon struct {
	ID                    uint64    `gorm:"primaryKey"`
	Code                  string    `gorm:"size:64;uniqueIndex"`
	Name                  string    `gorm:"size:255"`
	Description           string    `gorm:"type:text"`
	Status                int       `gorm:"column:status;index"`
	DiscountType          string    `gorm:"size:32"`
	DiscountValue         int64     `gorm:"column:discount_value"`
	Currency              string    `gorm:"size:16"`
	MaxRedemptions        int       `gorm:"column:max_redemptions"`
	MaxRedemptionsPerUser int       `gorm:"column:max_redemptions_per_user"`
	MinOrderCents         int64     `gorm:"column:min_order_cents"`
	StartsAt              time.Time `gorm:"column:starts_at"`
	EndsAt                time.Time `gorm:"column:ends_at"`
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

// TableName binds the coupon table name.
func (Coupon) TableName() string { return "coupons" }

// CouponRedemption records coupon usage for an order.
type CouponRedemption struct {
	ID          uint64 `gorm:"primaryKey"`
	CouponID    uint64 `gorm:"index"`
	UserID      uint64 `gorm:"index"`
	OrderID     uint64 `gorm:"index"`
	Status      int    `gorm:"column:status;index"`
	AmountCents int64  `gorm:"column:amount_cents"`
	Currency    string `gorm:"size:16"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TableName binds the redemption table name.
func (CouponRedemption) TableName() string { return "coupon_redemptions" }

// ListCouponsOptions controls filters and pagination.
type ListCouponsOptions struct {
	Page      int
	PerPage   int
	Sort      string
	Direction string
	Query     string
	Status    int
}

// UpdateCouponInput defines mutable coupon fields.
type UpdateCouponInput struct {
	Name                  *string
	Description           *string
	Status                *int
	DiscountType          *string
	DiscountValue         *int64
	Currency              *string
	MaxRedemptions        *int
	MaxRedemptionsPerUser *int
	MinOrderCents         *int64
	StartsAt              *time.Time
	EndsAt                *time.Time
}

// CouponRepository manages coupons and redemptions.
type CouponRepository interface {
	List(ctx context.Context, opts ListCouponsOptions) ([]Coupon, int64, error)
	Get(ctx context.Context, id uint64) (Coupon, error)
	GetByCode(ctx context.Context, code string) (Coupon, error)
	GetByCodeForUpdate(ctx context.Context, code string) (Coupon, error)
	Create(ctx context.Context, coupon Coupon) (Coupon, error)
	Update(ctx context.Context, id uint64, input UpdateCouponInput) (Coupon, error)
	Delete(ctx context.Context, id uint64) error

	CountRedemptions(ctx context.Context, couponID uint64) (int64, error)
	CountRedemptionsByUser(ctx context.Context, couponID, userID uint64) (int64, error)
	CreateRedemption(ctx context.Context, redemption CouponRedemption) (CouponRedemption, error)
	UpdateRedemptionStatusByOrder(ctx context.Context, orderID uint64, status int) error
}

type couponRepository struct {
	db *gorm.DB
}

// NewCouponRepository constructs a coupon repository.
func NewCouponRepository(db *gorm.DB) (CouponRepository, error) {
	if db == nil {
		return nil, errors.New("repository: database connection is required")
	}
	return &couponRepository{db: db}, nil
}

func (r *couponRepository) List(ctx context.Context, opts ListCouponsOptions) ([]Coupon, int64, error) {
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}

	opts = normalizeListCouponsOptions(opts)
	base := r.db.WithContext(ctx).Model(&Coupon{})

	if query := strings.TrimSpace(strings.ToLower(opts.Query)); query != "" {
		like := fmt.Sprintf("%%%s%%", query)
		base = base.Where("(LOWER(code) LIKE ? OR LOWER(name) LIKE ?)", like, like)
	}
	if opts.Status != 0 {
		base = base.Where("status = ?", opts.Status)
	}

	countQuery := base.Session(&gorm.Session{})
	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []Coupon{}, 0, nil
	}

	orderClause := buildCouponOrderClause(opts.Sort, opts.Direction)
	offset := (opts.Page - 1) * opts.PerPage
	listQuery := base.Session(&gorm.Session{}).Order(orderClause).Limit(opts.PerPage).Offset(offset)

	var coupons []Coupon
	if err := listQuery.Find(&coupons).Error; err != nil {
		return nil, 0, err
	}
	return coupons, total, nil
}

func (r *couponRepository) Get(ctx context.Context, id uint64) (Coupon, error) {
	if err := ctx.Err(); err != nil {
		return Coupon{}, err
	}

	var coupon Coupon
	if err := r.db.WithContext(ctx).First(&coupon, id).Error; err != nil {
		return Coupon{}, translateError(err)
	}
	return coupon, nil
}

func (r *couponRepository) GetByCode(ctx context.Context, code string) (Coupon, error) {
	if err := ctx.Err(); err != nil {
		return Coupon{}, err
	}

	var coupon Coupon
	if err := r.db.WithContext(ctx).Where("code = ?", normalizeCouponCode(code)).First(&coupon).Error; err != nil {
		return Coupon{}, translateError(err)
	}
	return coupon, nil
}

func (r *couponRepository) GetByCodeForUpdate(ctx context.Context, code string) (Coupon, error) {
	if err := ctx.Err(); err != nil {
		return Coupon{}, err
	}

	var coupon Coupon
	if err := r.db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("code = ?", normalizeCouponCode(code)).
		First(&coupon).Error; err != nil {
		return Coupon{}, translateError(err)
	}
	return coupon, nil
}

func (r *couponRepository) Create(ctx context.Context, coupon Coupon) (Coupon, error) {
	if err := ctx.Err(); err != nil {
		return Coupon{}, err
	}

	coupon.Code = normalizeCouponCode(coupon.Code)
	coupon.Name = strings.TrimSpace(coupon.Name)
	coupon.Description = strings.TrimSpace(coupon.Description)
	coupon.DiscountType = strings.ToLower(strings.TrimSpace(coupon.DiscountType))
	coupon.Currency = strings.ToUpper(strings.TrimSpace(coupon.Currency))

	if coupon.Code == "" || coupon.DiscountType == "" || coupon.DiscountValue <= 0 {
		return Coupon{}, ErrInvalidArgument
	}
	if coupon.Status == 0 {
		coupon.Status = CouponStatusActive
	}
	if coupon.DiscountType != CouponTypePercent && coupon.DiscountType != CouponTypeFixed {
		return Coupon{}, ErrInvalidArgument
	}
	if coupon.DiscountType == CouponTypePercent && coupon.DiscountValue > 10000 {
		return Coupon{}, ErrInvalidArgument
	}
	if coupon.DiscountType == CouponTypeFixed && coupon.Currency == "" {
		return Coupon{}, ErrInvalidArgument
	}

	now := time.Now().UTC()
	if coupon.CreatedAt.IsZero() {
		coupon.CreatedAt = now
	}
	coupon.UpdatedAt = now

	if err := r.db.WithContext(ctx).Create(&coupon).Error; err != nil {
		return Coupon{}, translateError(err)
	}
	return coupon, nil
}

func (r *couponRepository) Update(ctx context.Context, id uint64, input UpdateCouponInput) (Coupon, error) {
	if err := ctx.Err(); err != nil {
		return Coupon{}, err
	}

	updates := map[string]any{}
	if input.Name != nil {
		updates["name"] = strings.TrimSpace(*input.Name)
	}
	if input.Description != nil {
		updates["description"] = strings.TrimSpace(*input.Description)
	}
	if input.Status != nil {
		updates["status"] = *input.Status
	}
	if input.DiscountType != nil {
		updates["discount_type"] = strings.ToLower(strings.TrimSpace(*input.DiscountType))
	}
	if input.DiscountValue != nil {
		updates["discount_value"] = *input.DiscountValue
	}
	if input.Currency != nil {
		updates["currency"] = strings.ToUpper(strings.TrimSpace(*input.Currency))
	}
	if input.MaxRedemptions != nil {
		updates["max_redemptions"] = *input.MaxRedemptions
	}
	if input.MaxRedemptionsPerUser != nil {
		updates["max_redemptions_per_user"] = *input.MaxRedemptionsPerUser
	}
	if input.MinOrderCents != nil {
		updates["min_order_cents"] = *input.MinOrderCents
	}
	if input.StartsAt != nil {
		updates["starts_at"] = input.StartsAt.UTC()
	}
	if input.EndsAt != nil {
		updates["ends_at"] = input.EndsAt.UTC()
	}
	if len(updates) == 0 {
		return Coupon{}, ErrInvalidArgument
	}
	updates["updated_at"] = time.Now().UTC()

	if err := r.db.WithContext(ctx).Model(&Coupon{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return Coupon{}, translateError(err)
	}
	return r.Get(ctx, id)
}

func (r *couponRepository) Delete(ctx context.Context, id uint64) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).Delete(&Coupon{}, id).Error; err != nil {
		return translateError(err)
	}
	return nil
}

func (r *couponRepository) CountRedemptions(ctx context.Context, couponID uint64) (int64, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	if couponID == 0 {
		return 0, ErrInvalidArgument
	}

	var count int64
	err := r.db.WithContext(ctx).Model(&CouponRedemption{}).
		Where("coupon_id = ? AND status != ?", couponID, CouponRedemptionReleased).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *couponRepository) CountRedemptionsByUser(ctx context.Context, couponID, userID uint64) (int64, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	if couponID == 0 || userID == 0 {
		return 0, ErrInvalidArgument
	}

	var count int64
	err := r.db.WithContext(ctx).Model(&CouponRedemption{}).
		Where("coupon_id = ? AND user_id = ? AND status != ?", couponID, userID, CouponRedemptionReleased).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *couponRepository) CreateRedemption(ctx context.Context, redemption CouponRedemption) (CouponRedemption, error) {
	if err := ctx.Err(); err != nil {
		return CouponRedemption{}, err
	}

	redemption.Currency = strings.ToUpper(strings.TrimSpace(redemption.Currency))
	if redemption.CouponID == 0 || redemption.UserID == 0 || redemption.OrderID == 0 {
		return CouponRedemption{}, ErrInvalidArgument
	}
	if redemption.Status == 0 {
		redemption.Status = CouponRedemptionReserved
	}
	if redemption.AmountCents <= 0 {
		return CouponRedemption{}, ErrInvalidArgument
	}

	now := time.Now().UTC()
	if redemption.CreatedAt.IsZero() {
		redemption.CreatedAt = now
	}
	if redemption.UpdatedAt.IsZero() {
		redemption.UpdatedAt = now
	}

	if err := r.db.WithContext(ctx).Create(&redemption).Error; err != nil {
		return CouponRedemption{}, translateError(err)
	}
	return redemption, nil
}

func (r *couponRepository) UpdateRedemptionStatusByOrder(ctx context.Context, orderID uint64, status int) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if orderID == 0 {
		return ErrInvalidArgument
	}
	if status == 0 {
		return ErrInvalidArgument
	}

	updates := map[string]any{
		"status":     status,
		"updated_at": time.Now().UTC(),
	}
	if err := r.db.WithContext(ctx).Model(&CouponRedemption{}).Where("order_id = ?", orderID).Updates(updates).Error; err != nil {
		return translateError(err)
	}
	return nil
}

func normalizeCouponCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}

func normalizeListCouponsOptions(opts ListCouponsOptions) ListCouponsOptions {
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

func buildCouponOrderClause(field, direction string) string {
	column := "updated_at"
	switch strings.ToLower(field) {
	case "code":
		column = "code"
	case "status":
		column = "status"
	case "created_at":
		column = "created_at"
	case "starts_at":
		column = "starts_at"
	case "ends_at":
		column = "ends_at"
	}

	dir := "ASC"
	if strings.EqualFold(direction, "desc") {
		dir = "DESC"
	}

	return fmt.Sprintf("%s %s", column, dir)
}
