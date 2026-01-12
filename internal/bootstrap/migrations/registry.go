package migrations

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"gorm.io/gorm"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
)

// SchemaMigration stores executed migration metadata.
type SchemaMigration struct {
	Version   uint64    `gorm:"primaryKey"`
	Name      string    `gorm:"size:128"`
	AppliedAt time.Time `gorm:"column:applied_at"`
}

// TableName ensures deterministic naming.
func (SchemaMigration) TableName() string { return "schema_migrations" }

// Migration describes an idempotent schema evolution step.
type Migration struct {
	Version uint64
	Name    string
	Up      func(ctx context.Context, db *gorm.DB) error
	Down    func(ctx context.Context, db *gorm.DB) error
}

// Result captures the state transition produced by Apply.
type Result struct {
	BeforeVersion      uint64
	AfterVersion       uint64
	TargetVersion      uint64
	AppliedVersions    []uint64
	RolledBackVersions []uint64
}

// ApplyResult captures details about a migration attempt.
type ApplyResult struct {
	// PreviousVersion reflects the version recorded before applying
	// any pending migrations.
	PreviousVersion uint64
	// CurrentVersion reflects the schema version after migrations have
	// been executed.
	CurrentVersion uint64
	// Applied enumerates the migrations that were newly executed within
	// the current invocation.
	Applied []SchemaMigration
	// Seeded indicates whether demo seed data was populated.
	Seeded bool
}

var migrationRegistry = []Migration{
	{
		Version: 2026010501,
		Name:    "release-init",
		Up: func(ctx context.Context, db *gorm.DB) error {
			return db.WithContext(ctx).AutoMigrate(
				&repository.AdminModule{},
				&repository.User{},
				&repository.UserCredential{},
				&repository.AuditLog{},
				&repository.SiteSetting{},
				&repository.SecuritySetting{},
				&repository.Announcement{},
				&repository.Node{},
				&repository.NodeKernel{},
				&repository.ProtocolBinding{},
				&repository.ProtocolEntry{},
				&repository.Plan{},
				&repository.PlanBillingOption{},
				&repository.PlanProtocolBinding{},
				&repository.SubscriptionTemplate{},
				&repository.SubscriptionTemplateHistory{},
				&repository.Subscription{},
				&repository.TrafficUsageRecord{},
				&repository.PaymentChannel{},
				&repository.Coupon{},
				&repository.CouponRedemption{},
				&repository.UserBalance{},
				&repository.BalanceTransaction{},
				&repository.Order{},
				&repository.OrderItem{},
				&repository.OrderRefund{},
				&repository.OrderPayment{},
			)
		},
		Down: func(ctx context.Context, db *gorm.DB) error {
			migrator := db.WithContext(ctx).Migrator()
			tables := []any{
				&repository.OrderPayment{},
				&repository.OrderRefund{},
				&repository.OrderItem{},
				&repository.Order{},
				&repository.BalanceTransaction{},
				&repository.UserBalance{},
				&repository.CouponRedemption{},
				&repository.Coupon{},
				&repository.PaymentChannel{},
				&repository.TrafficUsageRecord{},
				&repository.Subscription{},
				&repository.SubscriptionTemplateHistory{},
				&repository.SubscriptionTemplate{},
				&repository.PlanProtocolBinding{},
				&repository.PlanBillingOption{},
				&repository.Plan{},
				&repository.ProtocolEntry{},
				&repository.ProtocolBinding{},
				&repository.NodeKernel{},
				&repository.Node{},
				&repository.Announcement{},
				&repository.SecuritySetting{},
				&repository.SiteSetting{},
				&repository.AuditLog{},
				&repository.UserCredential{},
				&repository.User{},
				&repository.AdminModule{},
			}

			for _, table := range tables {
				if err := migrator.DropTable(table); err != nil {
					return err
				}
			}

			return nil
		},
	},
	{
		Version: 2026010502,
		Name:    "add-site-setting-offline-probe",
		Up: func(ctx context.Context, db *gorm.DB) error {
			return db.WithContext(ctx).AutoMigrate(&repository.SiteSetting{})
		},
		Down: func(ctx context.Context, db *gorm.DB) error {
			return db.WithContext(ctx).Migrator().
				DropColumn(&repository.SiteSetting{}, "kernel_offline_probe_max_interval_seconds")
		},
	},
}

func init() {
	sort.Slice(migrationRegistry, func(i, j int) bool {
		return migrationRegistry[i].Version < migrationRegistry[j].Version
	})
}

// Apply executes migrations up to targetVersion (0 denotes latest).
func Apply(ctx context.Context, db *gorm.DB, targetVersion uint64, allowRollback bool) (Result, error) {
	var result Result

	if db == nil {
		return result, fmt.Errorf("migrations: database connection is required")
	}

	if err := db.WithContext(ctx).AutoMigrate(&SchemaMigration{}); err != nil {
		return result, fmt.Errorf("migrations: prepare metadata table: %w", err)
	}

	var applied []SchemaMigration
	if err := db.WithContext(ctx).Order("version ASC").Find(&applied).Error; err != nil {
		return result, fmt.Errorf("migrations: load applied versions: %w", err)
	}

	appliedSet := make(map[uint64]SchemaMigration, len(applied))
	var currentVersion uint64
	for _, record := range applied {
		appliedSet[record.Version] = record
		if record.Version > currentVersion {
			currentVersion = record.Version
		}
	}

	registryMap := make(map[uint64]Migration, len(migrationRegistry))
	var latestVersion uint64
	for _, migration := range migrationRegistry {
		if _, exists := registryMap[migration.Version]; exists {
			return result, fmt.Errorf("migrations: duplicate migration version %d registered", migration.Version)
		}
		registryMap[migration.Version] = migration
		if migration.Version > latestVersion {
			latestVersion = migration.Version
		}
	}

	for version := range appliedSet {
		if _, ok := registryMap[version]; !ok {
			return result, fmt.Errorf("migrations: applied version %d is not registered", version)
		}
	}

	result.BeforeVersion = currentVersion

	effectiveTarget := targetVersion
	if targetVersion == 0 {
		effectiveTarget = latestVersion
	} else {
		if len(migrationRegistry) == 0 {
			return result, fmt.Errorf("migrations: no migrations registered, cannot reach target version %d", targetVersion)
		}
		if targetVersion > latestVersion {
			return result, fmt.Errorf("migrations: target version %d is newer than latest registered version %d", targetVersion, latestVersion)
		}
	}
	result.TargetVersion = effectiveTarget

	if effectiveTarget > currentVersion {
		for _, migration := range migrationRegistry {
			if migration.Version <= currentVersion {
				continue
			}
			if migration.Version > effectiveTarget {
				break
			}

			entry, err := applyMigration(ctx, db, migration)
			if err != nil {
				return result, err
			}
			result.AppliedVersions = append(result.AppliedVersions, migration.Version)
			appliedSet[migration.Version] = entry
			currentVersion = migration.Version
		}
	} else if effectiveTarget < currentVersion {
		if !allowRollback {
			return result, fmt.Errorf("migrations: target version %d is older than current version %d", effectiveTarget, currentVersion)
		}

		sort.Slice(applied, func(i, j int) bool {
			return applied[i].Version > applied[j].Version
		})

		for _, record := range applied {
			if record.Version <= effectiveTarget {
				break
			}
			migration, ok := registryMap[record.Version]
			if !ok {
				return result, fmt.Errorf("migrations: applied version %d is not registered", record.Version)
			}
			if err := rollbackMigration(ctx, db, migration); err != nil {
				return result, err
			}
			delete(appliedSet, migration.Version)
			result.RolledBackVersions = append(result.RolledBackVersions, migration.Version)
		}

		currentVersion = 0
		for version := range appliedSet {
			if version > currentVersion {
				currentVersion = version
			}
		}
	}

	result.AfterVersion = currentVersion

	return result, nil
}

func applyMigration(ctx context.Context, db *gorm.DB, migration Migration) (SchemaMigration, error) {
	var entry SchemaMigration

	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if migration.Up == nil {
			return fmt.Errorf("migrations: up migration %d (%s) is nil", migration.Version, migration.Name)
		}
		if err := migration.Up(ctx, tx); err != nil {
			return fmt.Errorf("up: %w", err)
		}

		entry = SchemaMigration{
			Version:   migration.Version,
			Name:      migration.Name,
			AppliedAt: time.Now().UTC(),
		}
		if result := tx.Create(&entry); result.Error != nil {
			return fmt.Errorf("record: %w", result.Error)
		} else if result.RowsAffected != 1 {
			return fmt.Errorf("record: affected %d rows", result.RowsAffected)
		}
		return nil
	})
	if err != nil {
		var cleanupErr error
		if migration.Down != nil {
			cleanupErr = migration.Down(ctx, db.WithContext(ctx))
		}
		if cleanupErr != nil {
			return SchemaMigration{}, fmt.Errorf("migrations: apply %d (%s): %w", migration.Version, migration.Name, errors.Join(err, cleanupErr))
		}
		return SchemaMigration{}, fmt.Errorf("migrations: apply %d (%s): %w", migration.Version, migration.Name, err)
	}

	return entry, nil
}

func rollbackMigration(ctx context.Context, db *gorm.DB, migration Migration) error {
	if migration.Down == nil {
		return fmt.Errorf("migrations: rollback %d (%s): down migration not defined", migration.Version, migration.Name)
	}

	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := migration.Down(ctx, tx); err != nil {
			return fmt.Errorf("down: %w", err)
		}

		result := tx.Where("version = ?", migration.Version).Delete(&SchemaMigration{})
		if result.Error != nil {
			return fmt.Errorf("delete metadata: %w", result.Error)
		}
		if result.RowsAffected != 1 {
			return fmt.Errorf("delete metadata: affected %d rows", result.RowsAffected)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("migrations: rollback %d (%s): %w", migration.Version, migration.Name, err)
	}

	return nil
}
