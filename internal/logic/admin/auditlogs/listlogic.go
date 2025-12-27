package auditlogs

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/security"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// ListLogic handles audit log listing.
type ListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewListLogic constructs ListLogic.
func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// List returns audit log records.
func (l *ListLogic) List(req *types.AdminAuditLogListRequest) (*types.AdminAuditLogListResponse, error) {
	actor, ok := security.UserFromContext(l.ctx)
	if !ok {
		return nil, repository.ErrUnauthorized
	}
	if !security.HasRole(actor, "admin") {
		return nil, repository.ErrForbidden
	}

	opts, err := buildAuditLogOptions(req)
	if err != nil {
		return nil, err
	}

	logs, total, err := l.svcCtx.Repositories.AuditLog.List(l.ctx, opts)
	if err != nil {
		return nil, err
	}

	result := make([]types.AuditLogSummary, 0, len(logs))
	for _, entry := range logs {
		result = append(result, mapAuditLogSummary(entry))
	}

	page, perPage := normalizePage(req.Page, req.PerPage)
	pagination := types.PaginationMeta{
		Page:       page,
		PerPage:    perPage,
		TotalCount: total,
		HasNext:    int64(page*perPage) < total,
		HasPrev:    page > 1,
	}

	return &types.AdminAuditLogListResponse{
		Logs:       result,
		Pagination: pagination,
	}, nil
}
