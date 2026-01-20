package svc

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/zero-net-panel/zero-net-panel/internal/config"
	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/security"
	"github.com/zero-net-panel/zero-net-panel/pkg/auth"
	"github.com/zero-net-panel/zero-net-panel/pkg/cache"
	"github.com/zero-net-panel/zero-net-panel/pkg/database"
)

type ServiceContext struct {
	Config        config.Config
	DB            *gorm.DB
	Cache         cache.Cache
	Repositories  *repository.Repositories
	Auth          *auth.Generator
	Credentials   *security.CredentialManager

	Ctx    context.Context
	cancel context.CancelFunc

	cleanup func()
}

func NewServiceContext(c config.Config) (*ServiceContext, error) {
	c.Normalize()

	db, dbClose, err := database.NewGorm(c.Database)
	if err != nil {
		return nil, fmt.Errorf("init database: %w", err)
	}
	if db == nil {
		dbClose()
		return nil, fmt.Errorf("init database: %w", errors.New("database configuration is required"))
	}

	cacheProvider, err := cache.New(c.Cache)
	if err != nil {
		dbClose()
		return nil, fmt.Errorf("init cache: %w", err)
	}

	repos, err := repository.NewRepositories(db)
	if err != nil {
		_ = cacheProvider.Close()
		dbClose()
		return nil, fmt.Errorf("init repositories: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	authGenerator := auth.NewGenerator(
		c.Auth.AccessSecret,
		c.Auth.RefreshSecret,
		c.Auth.AccessExpire,
		c.Auth.RefreshExpire,
	)

	credentialManager, err := security.NewCredentialManager(c.Credentials.MasterKey)
	if err != nil {
		_ = cacheProvider.Close()
		dbClose()
		return nil, fmt.Errorf("init credential manager: %w", err)
	}

	svcCtx := &ServiceContext{
		Config:        c,
		DB:            db,
		Cache:         cacheProvider,
		Repositories:  repos,
		Auth:          authGenerator,
		Credentials:   credentialManager,
		Ctx:           ctx,
		cancel:        cancel,
	}

	svcCtx.cleanup = func() {
		if svcCtx.cancel != nil {
			svcCtx.cancel()
		}
		if cacheProvider != nil {
			_ = cacheProvider.Close()
		}
		dbClose()
	}

	return svcCtx, nil
}

func (s *ServiceContext) Cleanup() {
	if s.cleanup != nil {
		s.cleanup()
	}
}

func (s *ServiceContext) Context() context.Context {
	return s.Ctx
}

func (s *ServiceContext) Cancel() {
	if s.cancel != nil {
		s.cancel()
	}
}
