package credentialutil

import (
	"context"
	"errors"
	"time"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/security"
	"github.com/zero-net-panel/zero-net-panel/internal/status"
)

// EnsureActiveCredential returns the current active credential or creates one.
func EnsureActiveCredential(ctx context.Context, repos *repository.Repositories, manager *security.CredentialManager, userID uint64) (repository.UserCredential, error) {
	if userID == 0 {
		return repository.UserCredential{}, repository.ErrInvalidArgument
	}
	if manager == nil {
		return repository.UserCredential{}, repository.ErrInvalidState
	}

	credential, err := repos.UserCredential.GetActiveByUser(ctx, userID)
	if err == nil {
		return credential, nil
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return repository.UserCredential{}, err
	}

	credential, err = createCredential(manager, userID, nil)
	if err != nil {
		return repository.UserCredential{}, err
	}
	return repos.UserCredential.Create(ctx, credential)
}

// RotateCredential creates a new active credential and deprecates the old one.
func RotateCredential(ctx context.Context, repos *repository.Repositories, manager *security.CredentialManager, userID uint64) (repository.UserCredential, error) {
	if userID == 0 {
		return repository.UserCredential{}, repository.ErrInvalidArgument
	}
	if manager == nil {
		return repository.UserCredential{}, repository.ErrInvalidState
	}

	var created repository.UserCredential
	err := repos.Transaction(ctx, func(txRepos *repository.Repositories) error {
		credential, err := txRepos.UserCredential.GetActiveByUser(ctx, userID)
		if err != nil && !errors.Is(err, repository.ErrNotFound) {
			return err
		}

		var rotatedFrom *uint64
		version := 1
		if err == nil {
			version = credential.Version + 1
			rotatedFrom = &credential.ID
			now := time.Now().UTC()
			statusCode := status.UserCredentialStatusDeprecated
			_, err = txRepos.UserCredential.Update(ctx, credential.ID, repository.UpdateUserCredentialInput{
				Status:       &statusCode,
				DeprecatedAt: &now,
			})
			if err != nil {
				return err
			}
		}

		newCredential, err := createCredential(manager, userID, rotatedFrom)
		if err != nil {
			return err
		}
		newCredential.Version = version
		newCredential.IssuedAt = time.Now().UTC()
		newCredential.Status = status.UserCredentialStatusActive
		newCredential.RotatedFromID = rotatedFrom

		created, err = txRepos.UserCredential.Create(ctx, newCredential)
		return err
	})
	if err != nil {
		return repository.UserCredential{}, err
	}
	return created, nil
}

// BuildIdentity decrypts and derives identity fields.
func BuildIdentity(manager *security.CredentialManager, userID uint64, credential repository.UserCredential) (security.DerivedIdentity, error) {
	secret, err := manager.DecryptForUser(userID, credential.SecretCiphertext, credential.SecretNonce)
	if err != nil {
		return security.DerivedIdentity{}, err
	}
	return manager.DeriveIdentity(userID, credential.Version, secret)
}

func createCredential(manager *security.CredentialManager, userID uint64, rotatedFrom *uint64) (repository.UserCredential, error) {
	secret, err := manager.GenerateSecret()
	if err != nil {
		return repository.UserCredential{}, err
	}
	ciphertext, nonce, fingerprint, err := manager.EncryptForUser(userID, secret)
	if err != nil {
		return repository.UserCredential{}, err
	}

	return repository.UserCredential{
		UserID:           userID,
		Version:          1,
		Status:           status.UserCredentialStatusActive,
		MasterKeyID:      manager.KeyID(),
		SecretCiphertext: ciphertext,
		SecretNonce:      nonce,
		Fingerprint:      fingerprint,
		IssuedAt:         time.Now().UTC(),
		RotatedFromID:    rotatedFrom,
	}, nil
}
