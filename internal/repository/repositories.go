package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// Repositories 聚合各领域仓储，方便在 ServiceContext 中注入。
type Repositories struct {
	db *gorm.DB

	AdminModule          AdminModuleRepository
	Node                 NodeRepository
	SubscriptionTemplate SubscriptionTemplateRepository
	Subscription         SubscriptionRepository
	User                 UserRepository
	UserCredential       UserCredentialRepository
	Plan                 PlanRepository
	Announcement         AnnouncementRepository
	Balance              BalanceRepository
	Security             SecurityRepository
	Site                 SiteRepository
	PaymentChannel       PaymentChannelRepository
	Order                OrderRepository
	AuditLog             AuditLogRepository
	ProtocolConfig       ProtocolConfigRepository
	ProtocolBinding      ProtocolBindingRepository
	TrafficUsage         TrafficUsageRepository
}

// NewRepositories 根据数据库实例创建仓储集合。
func NewRepositories(db *gorm.DB) (*Repositories, error) {
	if db == nil {
		return nil, errors.New("repository: database connection is required")
	}

	adminModuleRepo, err := NewAdminModuleRepository(db)
	if err != nil {
		return nil, err
	}

	templateRepo, err := NewSubscriptionTemplateRepository(db)
	if err != nil {
		return nil, err
	}

	userRepo, err := NewUserRepository(db)
	if err != nil {
		return nil, err
	}

	credentialRepo, err := NewUserCredentialRepository(db)
	if err != nil {
		return nil, err
	}

	nodeRepo, err := NewNodeRepository(db)
	if err != nil {
		return nil, err
	}

	subscriptionRepo, err := NewSubscriptionRepository(db, templateRepo)
	if err != nil {
		return nil, err
	}

	planRepo, err := NewPlanRepository(db)
	if err != nil {
		return nil, err
	}

	announcementRepo, err := NewAnnouncementRepository(db)
	if err != nil {
		return nil, err
	}

	balanceRepo, err := NewBalanceRepository(db)
	if err != nil {
		return nil, err
	}

	securityRepo, err := NewSecurityRepository(db)
	if err != nil {
		return nil, err
	}

	siteRepo, err := NewSiteRepository(db)
	if err != nil {
		return nil, err
	}

	channelRepo, err := NewPaymentChannelRepository(db)
	if err != nil {
		return nil, err
	}

	orderRepo, err := NewOrderRepository(db)
	if err != nil {
		return nil, err
	}

	auditRepo, err := NewAuditLogRepository(db)
	if err != nil {
		return nil, err
	}

	protocolConfigRepo, err := NewProtocolConfigRepository(db)
	if err != nil {
		return nil, err
	}

	protocolBindingRepo, err := NewProtocolBindingRepository(db)
	if err != nil {
		return nil, err
	}

	trafficRepo, err := NewTrafficUsageRepository(db)
	if err != nil {
		return nil, err
	}

	return &Repositories{
		db:                   db,
		AdminModule:          adminModuleRepo,
		Node:                 nodeRepo,
		SubscriptionTemplate: templateRepo,
		Subscription:         subscriptionRepo,
		User:                 userRepo,
		UserCredential:       credentialRepo,
		Plan:                 planRepo,
		Announcement:         announcementRepo,
		Balance:              balanceRepo,
		Security:             securityRepo,
		Site:                 siteRepo,
		PaymentChannel:       channelRepo,
		Order:                orderRepo,
		AuditLog:             auditRepo,
		ProtocolConfig:       protocolConfigRepo,
		ProtocolBinding:      protocolBindingRepo,
		TrafficUsage:         trafficRepo,
	}, nil
}

// Transaction executes the callback within a DB transaction with repositories bound to the same transaction.
func (r *Repositories) Transaction(ctx context.Context, fn func(txRepos *Repositories) error) error {
	if r == nil || r.db == nil {
		return errors.New("repository: database connection is required")
	}
	if fn == nil {
		return errors.New("repository: transaction callback is required")
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepos, err := NewRepositories(tx)
		if err != nil {
			return err
		}
		return fn(txRepos)
	})
}
