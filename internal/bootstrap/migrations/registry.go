package migrations

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/status"
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
		Version: 2026030101,
		Name:    "status-code-migration",
		Up: func(ctx context.Context, db *gorm.DB) error {
			columns := []statusColumn{
				{table: "users", column: "status", mapping: map[string]int{
					"active":   status.UserStatusActive,
					"pending":  status.UserStatusPending,
					"disabled": status.UserStatusDisabled,
				}},
				{table: "user_credentials", column: "status", mapping: map[string]int{
					"active":     status.UserCredentialStatusActive,
					"deprecated": status.UserCredentialStatusDeprecated,
					"revoked":    status.UserCredentialStatusRevoked,
				}},
				{table: "announcements", column: "status", mapping: map[string]int{
					"draft":     status.AnnouncementStatusDraft,
					"published": status.AnnouncementStatusPublished,
					"archived":  status.AnnouncementStatusArchived,
				}},
				{table: "coupons", column: "status", mapping: map[string]int{
					"active":   status.CouponStatusActive,
					"disabled": status.CouponStatusDisabled,
				}},
				{table: "coupon_redemptions", column: "status", mapping: map[string]int{
					"reserved": status.CouponRedemptionStatusReserved,
					"applied":  status.CouponRedemptionStatusApplied,
					"released": status.CouponRedemptionStatusReleased,
				}},
				{table: "plans", column: "status", mapping: map[string]int{
					"draft":    status.PlanStatusDraft,
					"active":   status.PlanStatusActive,
					"archived": status.PlanStatusArchived,
				}},
				{table: "plan_billing_options", column: "status", mapping: map[string]int{
					"draft":    status.PlanBillingOptionStatusDraft,
					"active":   status.PlanBillingOptionStatusActive,
					"archived": status.PlanBillingOptionStatusArchived,
				}},
				{table: "subscriptions", column: "status", mapping: map[string]int{
					"active":   status.SubscriptionStatusActive,
					"disabled": status.SubscriptionStatusDisabled,
					"expired":  status.SubscriptionStatusExpired,
					"pending":  status.SubscriptionStatusUnknown,
				}},
				{table: "nodes", column: "status", mapping: map[string]int{
					"online":      status.NodeStatusOnline,
					"offline":     status.NodeStatusOffline,
					"maintenance": status.NodeStatusMaintenance,
					"disabled":    status.NodeStatusDisabled,
				}},
				{table: "node_kernels", column: "status", mapping: map[string]int{
					"configured": status.NodeKernelStatusConfigured,
					"synced":     status.NodeKernelStatusSynced,
				}},
				{table: "protocol_bindings", column: "status", mapping: map[string]int{
					"active":   status.ProtocolBindingStatusActive,
					"disabled": status.ProtocolBindingStatusDisabled,
				}},
				{table: "protocol_bindings", column: "sync_status", mapping: map[string]int{
					"pending": status.ProtocolBindingSyncStatusPending,
					"synced":  status.ProtocolBindingSyncStatusSynced,
					"error":   status.ProtocolBindingSyncStatusError,
				}},
				{table: "protocol_bindings", column: "health_status", mapping: map[string]int{
					"unknown":   status.ProtocolBindingHealthStatusUnknown,
					"healthy":   status.ProtocolBindingHealthStatusHealthy,
					"degraded":  status.ProtocolBindingHealthStatusDegraded,
					"unhealthy": status.ProtocolBindingHealthStatusUnhealthy,
					"offline":   status.ProtocolBindingHealthStatusOffline,
				}},
				{table: "protocol_entries", column: "status", mapping: map[string]int{
					"active":   status.ProtocolEntryStatusActive,
					"disabled": status.ProtocolEntryStatusDisabled,
				}},
				{table: "orders", column: "status", mapping: map[string]int{
					"pending_payment":    status.OrderStatusPendingPayment,
					"pending":            status.OrderStatusPendingPayment,
					"paid":               status.OrderStatusPaid,
					"payment_failed":     status.OrderStatusPaymentFailed,
					"cancelled":          status.OrderStatusCancelled,
					"canceled":           status.OrderStatusCancelled,
					"partially_refunded": status.OrderStatusPartiallyRefunded,
					"refunded":           status.OrderStatusRefunded,
				}},
				{table: "orders", column: "payment_status", mapping: map[string]int{
					"pending":   status.OrderPaymentStatusPending,
					"succeeded": status.OrderPaymentStatusSucceeded,
					"failed":    status.OrderPaymentStatusFailed,
				}},
				{table: "order_payments", column: "status", mapping: map[string]int{
					"pending":   status.OrderPaymentStatusPending,
					"succeeded": status.OrderPaymentStatusSucceeded,
					"failed":    status.OrderPaymentStatusFailed,
				}},
			}

			for _, column := range columns {
				if err := updateStatusColumn(ctx, db, column); err != nil {
					return err
				}
			}

			switch db.Dialector.Name() {
			case "postgres", "mysql":
				for _, column := range columns {
					if err := alterStatusColumn(ctx, db, column); err != nil {
						return err
					}
				}
			}

			return nil
		},
		Down: func(ctx context.Context, db *gorm.DB) error {
			return nil
		},
	},
}

type statusColumn struct {
	table   string
	column  string
	mapping map[string]int
}

func updateStatusColumn(ctx context.Context, db *gorm.DB, column statusColumn) error {
	if db == nil {
		return fmt.Errorf("migrations: database connection is required")
	}
	if column.table == "" || column.column == "" {
		return fmt.Errorf("migrations: invalid status column")
	}

	dialect := db.Dialector.Name()
	valueExpr := column.column
	switch dialect {
	case "mysql":
		valueExpr = fmt.Sprintf("LOWER(CAST(%s AS CHAR))", column.column)
	default:
		valueExpr = fmt.Sprintf("LOWER(CAST(%s AS TEXT))", column.column)
	}

	codes := collectStatusCodes(column.mapping)
	clauses := make([]string, 0, len(column.mapping)+len(codes))
	for key, value := range column.mapping {
		clauses = append(clauses, fmt.Sprintf("WHEN %s = '%s' THEN %d", valueExpr, key, value))
	}
	for _, code := range codes {
		clauses = append(clauses, fmt.Sprintf("WHEN %s = '%d' THEN %d", valueExpr, code, code))
	}

	statement := fmt.Sprintf(
		"UPDATE %s SET %s = CASE %s ELSE 0 END",
		column.table,
		column.column,
		strings.Join(clauses, " "),
	)

	return db.WithContext(ctx).Exec(statement).Error
}

func alterStatusColumn(ctx context.Context, db *gorm.DB, column statusColumn) error {
	if db == nil {
		return fmt.Errorf("migrations: database connection is required")
	}

	switch db.Dialector.Name() {
	case "postgres":
		statement := fmt.Sprintf(
			"ALTER TABLE %s ALTER COLUMN %s TYPE INTEGER USING (%s::integer)",
			column.table,
			column.column,
			column.column,
		)
		return db.WithContext(ctx).Exec(statement).Error
	case "mysql":
		statement := fmt.Sprintf("ALTER TABLE %s MODIFY COLUMN %s INT", column.table, column.column)
		return db.WithContext(ctx).Exec(statement).Error
	default:
		return nil
	}
}

func collectStatusCodes(mapping map[string]int) []int {
	seen := make(map[int]struct{}, len(mapping)+1)
	seen[0] = struct{}{}
	for _, value := range mapping {
		seen[value] = struct{}{}
	}
	codes := make([]int, 0, len(seen))
	for value := range seen {
		codes = append(codes, value)
	}
	sort.Ints(codes)
	return codes
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
