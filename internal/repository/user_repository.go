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

// User 描述系统用户信息。
type User struct {
	ID                  uint64    `gorm:"primaryKey"`
	Email               string    `gorm:"size:255;uniqueIndex"`
	DisplayName         string    `gorm:"size:255"`
	PasswordHash        string    `gorm:"size:255"`
	Roles               []string  `gorm:"serializer:json"`
	Status              int       `gorm:"column:status"`
	EmailVerifiedAt     time.Time `gorm:"column:email_verified_at"`
	FailedLoginAttempts int       `gorm:"column:failed_login_attempts"`
	LockedUntil         time.Time `gorm:"column:locked_until"`
	TokenInvalidBefore  time.Time `gorm:"column:token_invalid_before"`
	PasswordUpdatedAt   time.Time `gorm:"column:password_updated_at"`
	PasswordResetAt     time.Time `gorm:"column:password_reset_at"`
	LastLoginAt         time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// TableName 自定义用户表名。
func (User) TableName() string { return "users" }

// UserRepository 定义用户仓储接口。
type UserRepository interface {
	Get(ctx context.Context, id uint64) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	List(ctx context.Context, opts ListUsersOptions) ([]User, int64, error)
	Create(ctx context.Context, user User) (User, error)
	UpdateProfile(ctx context.Context, id uint64, input UpdateUserProfileInput) (User, error)
	UpdateEmail(ctx context.Context, id uint64, email string, verifiedAt *time.Time) (User, error)
	UpdateStatus(ctx context.Context, id uint64, status int) (User, error)
	UpdateRoles(ctx context.Context, id uint64, roles []string) (User, error)
	UpdatePassword(ctx context.Context, id uint64, passwordHash string, resetAt *time.Time) error
	UpdateVerification(ctx context.Context, id uint64, verifiedAt time.Time, status int) (User, error)
	UpdateTokenInvalidBefore(ctx context.Context, id uint64, ts time.Time) error
	RecordLoginFailure(ctx context.Context, id uint64, maxAttempts int, lockDuration time.Duration) (User, error)
	UpdateLastLogin(ctx context.Context, id uint64, ts time.Time) error
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储，当前以内存实现模拟。
func NewUserRepository(db *gorm.DB) (UserRepository, error) {
	if db == nil {
		return nil, errors.New("repository: database connection is required")
	}
	return &userRepository{db: db}, nil
}

// ListUsersOptions controls filtering and pagination for user listing.
type ListUsersOptions struct {
	Page      int
	PerPage   int
	Sort      string
	Direction string
	Query     string
	Status    int
	Role      string
}

// UpdateUserProfileInput defines allowed profile updates.
type UpdateUserProfileInput struct {
	DisplayName *string
}

func (r *userRepository) Get(ctx context.Context, id uint64) (User, error) {
	if err := ctx.Err(); err != nil {
		return User{}, err
	}

	var user User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return User{}, translateError(err)
	}

	return user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (User, error) {
	if err := ctx.Err(); err != nil {
		return User{}, err
	}

	var user User
	if err := r.db.WithContext(ctx).
		Where("LOWER(email) = ?", normalizeEmail(email)).
		First(&user).Error; err != nil {
		return User{}, translateError(err)
	}

	return user, nil
}

func (r *userRepository) List(ctx context.Context, opts ListUsersOptions) ([]User, int64, error) {
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}

	opts = normalizeListUsersOptions(opts)

	base := r.db.WithContext(ctx).Model(&User{})
	if query := strings.TrimSpace(strings.ToLower(opts.Query)); query != "" {
		like := fmt.Sprintf("%%%s%%", query)
		base = base.Where("(LOWER(email) LIKE ? OR LOWER(display_name) LIKE ?)", like, like)
	}
	if opts.Status != 0 {
		base = base.Where("status = ?", opts.Status)
	}
	if role := strings.TrimSpace(strings.ToLower(opts.Role)); role != "" {
		like := fmt.Sprintf("%%\"%s\"%%", role)
		base = base.Where("LOWER(roles) LIKE ?", like)
	}

	countQuery := base.Session(&gorm.Session{})
	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []User{}, 0, nil
	}

	orderClause := buildUserOrderClause(opts.Sort, opts.Direction)
	offset := (opts.Page - 1) * opts.PerPage
	listQuery := base.Session(&gorm.Session{}).Order(orderClause).Limit(opts.PerPage).Offset(offset)

	var users []User
	if err := listQuery.Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) Create(ctx context.Context, user User) (User, error) {
	if err := ctx.Err(); err != nil {
		return User{}, err
	}

	user.Email = normalizeEmail(user.Email)
	if user.Email == "" || strings.TrimSpace(user.PasswordHash) == "" {
		return User{}, ErrInvalidArgument
	}

	now := time.Now().UTC()
	if user.CreatedAt.IsZero() {
		user.CreatedAt = now
	}
	user.UpdatedAt = now
	if user.PasswordUpdatedAt.IsZero() {
		user.PasswordUpdatedAt = now
	}
	if user.Status == 0 {
		user.Status = status.UserStatusActive
	}
	user.EmailVerifiedAt = NormalizeTime(user.EmailVerifiedAt)
	user.LockedUntil = NormalizeTime(user.LockedUntil)
	user.TokenInvalidBefore = NormalizeTime(user.TokenInvalidBefore)
	user.PasswordResetAt = NormalizeTime(user.PasswordResetAt)
	user.LastLoginAt = NormalizeTime(user.LastLoginAt)

	if err := r.db.WithContext(ctx).Create(&user).Error; err != nil {
		return User{}, translateError(err)
	}

	return user, nil
}

func (r *userRepository) UpdateProfile(ctx context.Context, id uint64, input UpdateUserProfileInput) (User, error) {
	if err := ctx.Err(); err != nil {
		return User{}, err
	}

	if input.DisplayName == nil {
		return User{}, ErrInvalidArgument
	}
	displayName := strings.TrimSpace(*input.DisplayName)
	if displayName == "" {
		return User{}, ErrInvalidArgument
	}

	updates := map[string]any{
		"display_name": displayName,
		"updated_at":   time.Now().UTC(),
	}
	if err := r.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return User{}, translateError(err)
	}

	return r.Get(ctx, id)
}

func (r *userRepository) UpdateEmail(ctx context.Context, id uint64, email string, verifiedAt *time.Time) (User, error) {
	if err := ctx.Err(); err != nil {
		return User{}, err
	}

	email = normalizeEmail(email)
	if email == "" {
		return User{}, ErrInvalidArgument
	}

	now := time.Now().UTC()
	updates := map[string]any{
		"email":      email,
		"updated_at": now,
	}
	if verifiedAt != nil {
		updates["email_verified_at"] = verifiedAt.UTC()
	}

	if err := r.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return User{}, translateError(err)
	}

	return r.Get(ctx, id)
}

func (r *userRepository) UpdateStatus(ctx context.Context, id uint64, statusCode int) (User, error) {
	if err := ctx.Err(); err != nil {
		return User{}, err
	}

	if statusCode == 0 {
		return User{}, ErrInvalidArgument
	}

	updates := map[string]any{
		"status":     statusCode,
		"updated_at": time.Now().UTC(),
	}
	if err := r.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return User{}, translateError(err)
	}

	return r.Get(ctx, id)
}

func (r *userRepository) UpdateRoles(ctx context.Context, id uint64, roles []string) (User, error) {
	if err := ctx.Err(); err != nil {
		return User{}, err
	}

	roles = normalizeRoles(roles)
	if len(roles) == 0 {
		return User{}, ErrInvalidArgument
	}

	updates := map[string]any{
		"roles":      roles,
		"updated_at": time.Now().UTC(),
	}

	if err := r.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return User{}, translateError(err)
	}

	return r.Get(ctx, id)
}

func (r *userRepository) UpdatePassword(ctx context.Context, id uint64, passwordHash string, resetAt *time.Time) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	passwordHash = strings.TrimSpace(passwordHash)
	if passwordHash == "" {
		return ErrInvalidArgument
	}

	now := time.Now().UTC()
	updates := map[string]any{
		"password_hash":         passwordHash,
		"password_updated_at":   now,
		"failed_login_attempts": 0,
		"locked_until":          ZeroTime(),
		"updated_at":            now,
	}
	if resetAt != nil {
		updates["password_reset_at"] = *resetAt
	}

	if err := r.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return translateError(err)
	}

	return nil
}

func (r *userRepository) UpdateVerification(ctx context.Context, id uint64, verifiedAt time.Time, statusCode int) (User, error) {
	if err := ctx.Err(); err != nil {
		return User{}, err
	}

	if verifiedAt.IsZero() {
		return User{}, ErrInvalidArgument
	}

	updates := map[string]any{
		"email_verified_at": verifiedAt,
		"updated_at":        time.Now().UTC(),
	}
	if statusCode != 0 {
		updates["status"] = statusCode
	}

	if err := r.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return User{}, translateError(err)
	}

	return r.Get(ctx, id)
}

func (r *userRepository) UpdateTokenInvalidBefore(ctx context.Context, id uint64, ts time.Time) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if ts.IsZero() {
		return ErrInvalidArgument
	}

	updates := map[string]any{
		"token_invalid_before": ts,
		"updated_at":           time.Now().UTC(),
	}

	if err := r.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return translateError(err)
	}

	return nil
}

func (r *userRepository) RecordLoginFailure(ctx context.Context, id uint64, maxAttempts int, lockDuration time.Duration) (User, error) {
	if err := ctx.Err(); err != nil {
		return User{}, err
	}

	now := time.Now().UTC()
	var user User

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&user, id).Error; err != nil {
			return err
		}

		user.FailedLoginAttempts++
		if maxAttempts > 0 && user.FailedLoginAttempts >= maxAttempts {
			if lockDuration > 0 {
				user.LockedUntil = now.Add(lockDuration)
			} else {
				user.LockedUntil = now
			}
		}
		user.UpdatedAt = now

		return tx.Model(&User{}).Where("id = ?", id).Updates(map[string]any{
			"failed_login_attempts": user.FailedLoginAttempts,
			"locked_until":          user.LockedUntil,
			"updated_at":            user.UpdatedAt,
		}).Error
	})
	if err != nil {
		return User{}, translateError(err)
	}

	return user, nil
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, id uint64, ts time.Time) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	updates := map[string]any{
		"last_login_at":         ts,
		"updated_at":            ts,
		"failed_login_attempts": 0,
		"locked_until":          ZeroTime(),
	}

	if err := r.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return translateError(err)
	}

	return nil
}

func normalizeListUsersOptions(opts ListUsersOptions) ListUsersOptions {
	if opts.Page <= 0 {
		opts.Page = 1
	}
	if opts.PerPage <= 0 || opts.PerPage > 100 {
		opts.PerPage = 20
	}
	if opts.Sort == "" {
		opts.Sort = "updated_at"
	}
	opts.Sort = strings.ToLower(strings.TrimSpace(opts.Sort))
	if opts.Direction == "" {
		opts.Direction = "desc"
	}
	return opts
}

func buildUserOrderClause(sort, direction string) string {
	column := "updated_at"
	switch strings.ToLower(strings.TrimSpace(sort)) {
	case "email":
		column = "email"
	case "created_at":
		column = "created_at"
	case "last_login_at":
		column = "last_login_at"
	case "status":
		column = "status"
	}

	dir := "ASC"
	if strings.EqualFold(direction, "desc") {
		dir = "DESC"
	}

	return fmt.Sprintf("%s %s, id ASC", column, dir)
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func normalizeRoles(roles []string) []string {
	seen := make(map[string]struct{}, len(roles))
	var normalized []string
	for _, role := range roles {
		role = strings.ToLower(strings.TrimSpace(role))
		if role == "" {
			continue
		}
		if _, ok := seen[role]; ok {
			continue
		}
		seen[role] = struct{}{}
		normalized = append(normalized, role)
	}
	return normalized
}
