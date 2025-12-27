package auditlogs

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/security"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// ExportLogic handles audit log export.
type ExportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewExportLogic constructs ExportLogic.
func NewExportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExportLogic {
	return &ExportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Export returns audit logs for export.
func (l *ExportLogic) Export(req *types.AdminAuditLogExportRequest) (*types.AdminAuditLogExportResponse, error) {
	actor, ok := security.UserFromContext(l.ctx)
	if !ok {
		return nil, repository.ErrUnauthorized
	}
	if !security.HasRole(actor, "admin") {
		return nil, repository.ErrForbidden
	}

	listReq := &types.AdminAuditLogListRequest{
		Page:         req.Page,
		PerPage:      req.PerPage,
		ActorID:      req.ActorID,
		Action:       req.Action,
		ResourceType: req.ResourceType,
		ResourceID:   req.ResourceID,
		Since:        req.Since,
		Until:        req.Until,
	}

	opts, err := buildAuditLogOptions(listReq)
	if err != nil {
		return nil, err
	}

	page, perPage := normalizeExportPage(req.Page, req.PerPage)
	opts.Page = page
	opts.PerPage = perPage

	logs, total, err := l.svcCtx.Repositories.AuditLog.List(l.ctx, opts)
	if err != nil {
		return nil, err
	}

	result := make([]types.AuditLogSummary, 0, len(logs))
	for _, entry := range logs {
		result = append(result, mapAuditLogSummary(entry))
	}

	return &types.AdminAuditLogExportResponse{
		Logs:       result,
		TotalCount: total,
		ExportedAt: time.Now().UTC().Unix(),
	}, nil
}
