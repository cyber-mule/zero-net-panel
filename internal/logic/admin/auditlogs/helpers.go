package auditlogs

import (
	"strings"
	"time"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

func mapAuditLogSummary(entry repository.AuditLog) types.AuditLogSummary {
	return types.AuditLogSummary{
		ID:           entry.ID,
		ActorID:      entry.ActorID,
		ActorEmail:   entry.ActorEmail,
		ActorRoles:   append([]string(nil), entry.ActorRoles...),
		Action:       entry.Action,
		ResourceType: entry.ResourceType,
		ResourceID:   entry.ResourceID,
		SourceIP:     entry.SourceIP,
		Metadata:     entry.Metadata,
		CreatedAt:    toUnixOrZero(entry.CreatedAt),
	}
}

func toUnixOrZero(ts time.Time) int64 {
	if ts.IsZero() {
		return 0
	}
	return ts.Unix()
}

func buildAuditLogOptions(req *types.AdminAuditLogListRequest) (repository.AuditLogListOptions, error) {
	opts := repository.AuditLogListOptions{
		Page:         req.Page,
		PerPage:      req.PerPage,
		ActorID:      req.ActorID,
		Action:       strings.TrimSpace(req.Action),
		ResourceType: strings.TrimSpace(req.ResourceType),
		ResourceID:   strings.TrimSpace(req.ResourceID),
	}

	if req.Since > 0 {
		since := time.Unix(req.Since, 0).UTC()
		opts.Since = &since
	}
	if req.Until > 0 {
		until := time.Unix(req.Until, 0).UTC()
		opts.Until = &until
	}
	if opts.Since != nil && opts.Until != nil && opts.Since.After(*opts.Until) {
		return repository.AuditLogListOptions{}, repository.ErrInvalidArgument
	}

	return opts, nil
}

func normalizePage(page, perPage int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}
	return page, perPage
}

func normalizeExportPage(page, perPage int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 1000
	}
	if perPage > 5000 {
		perPage = 5000
	}
	return page, perPage
}
