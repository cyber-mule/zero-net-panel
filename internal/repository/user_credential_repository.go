package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/zero-net-panel/zero-net-panel/internal/status"
)

// UserCredential stores per-user authentication material (encrypted).
type UserCredential struct {
	ID               uint64    `gorm:"primaryKey"`
	UserID           uint64    `gorm:"index"`
	Version          int       `gorm:"index"`
	Status           int       `gorm:"column:status;index"`
	MasterKeyID      string    `gorm:"size:64"`
	SecretCiphertext string    `gorm:"type:text"`
	SecretNonce      string    `gorm:"size:64"`
	Fingerprint      string    `gorm:"size:128;index"`
	IssuedAt         time.Time `gorm:"column:issued_at"`
	DeprecatedAt     time.Time `gorm:"column:deprecated_at"`
	RevokedAt        time.Time `gorm:"column:revoked_at"`
	RotatedFromID    *uint64   `gorm:"column:rotated_from_id"`
	LastSeenAt       time.Time `gorm:"column:last_seen_at"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// TableName binds the credential table name.
func (UserCredential) TableName() string { return "user_credentials" }

// UpdateUserCredentialInput defines mutable credential fields.
type UpdateUserCredentialInput struct {
	Status       *int
	DeprecatedAt *time.Time
	RevokedAt    *time.Time
	LastSeenAt   *time.Time
}

// UserCredentialRepository manages credential persistence.
type UserCredentialRepository interface {
	GetActiveByUser(ctx context.Context, userID uint64) (UserCredential, error)
	Create(ctx context.Context, credential UserCredential) (UserCredential, error)
	Update(ctx context.Context, id uint64, input UpdateUserCredentialInput) (UserCredential, error)
}

type userCredentialRepository struct {
	db *gorm.DB
}

// NewUserCredentialRepository constructs a credential repository.
func NewUserCredentialRepository(db *gorm.DB) (UserCredentialRepository, error) {
	if db == nil {
		return nil, errors.New("repository: database connection is required")
	}
	return &userCredentialRepository{db: db}, nil
}

func (r *userCredentialRepository) GetActiveByUser(ctx context.Context, userID uint64) (UserCredential, error) {
	if err := ctx.Err(); err != nil {
		return UserCredential{}, err
	}
	if userID == 0 {
		return UserCredential{}, ErrInvalidArgument
	}

	var credential UserCredential
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND status = ?", userID, status.UserCredentialStatusActive).
		Order("version DESC").
		First(&credential).Error; err != nil {
		return UserCredential{}, translateError(err)
	}
	return credential, nil
}

func (r *userCredentialRepository) Create(ctx context.Context, credential UserCredential) (UserCredential, error) {
	if err := ctx.Err(); err != nil {
		return UserCredential{}, err
	}

	credential.MasterKeyID = strings.TrimSpace(credential.MasterKeyID)
	credential.SecretCiphertext = strings.TrimSpace(credential.SecretCiphertext)
	credential.SecretNonce = strings.TrimSpace(credential.SecretNonce)
	credential.Fingerprint = strings.TrimSpace(credential.Fingerprint)
	if credential.UserID == 0 || credential.Version <= 0 || credential.Status == 0 {
		return UserCredential{}, ErrInvalidArgument
	}
	if credential.SecretCiphertext == "" || credential.SecretNonce == "" || credential.Fingerprint == "" {
		return UserCredential{}, ErrInvalidArgument
	}

	now := time.Now().UTC()
	if credential.IssuedAt.IsZero() {
		credential.IssuedAt = now
	}
	if credential.CreatedAt.IsZero() {
		credential.CreatedAt = now
	}
	credential.UpdatedAt = now

	if err := r.db.WithContext(ctx).Create(&credential).Error; err != nil {
		return UserCredential{}, translateError(err)
	}
	return credential, nil
}

func (r *userCredentialRepository) Update(ctx context.Context, id uint64, input UpdateUserCredentialInput) (UserCredential, error) {
	if err := ctx.Err(); err != nil {
		return UserCredential{}, err
	}
	if id == 0 {
		return UserCredential{}, ErrInvalidArgument
	}

	updates := map[string]any{}
	if input.Status != nil {
		updates["status"] = *input.Status
	}
	if input.DeprecatedAt != nil {
		updates["deprecated_at"] = input.DeprecatedAt.UTC()
	}
	if input.RevokedAt != nil {
		updates["revoked_at"] = input.RevokedAt.UTC()
	}
	if input.LastSeenAt != nil {
		updates["last_seen_at"] = input.LastSeenAt.UTC()
	}
	if len(updates) == 0 {
		return UserCredential{}, ErrInvalidArgument
	}
	updates["updated_at"] = time.Now().UTC()

	if err := r.db.WithContext(ctx).Model(&UserCredential{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return UserCredential{}, translateError(err)
	}

	var credential UserCredential
	if err := r.db.WithContext(ctx).First(&credential, id).Error; err != nil {
		return UserCredential{}, translateError(err)
	}
	return credential, nil
}
