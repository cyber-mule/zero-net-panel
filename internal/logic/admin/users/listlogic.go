package users

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// ListLogic handles admin user listing.
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

// List returns user list with pagination.
func (l *ListLogic) List(req *types.AdminListUsersRequest) (*types.AdminUserListResponse, error) {
	opts := repository.ListUsersOptions{
		Page:      req.Page,
		PerPage:   req.PerPage,
		Query:     req.Query,
		Status:    req.Status,
		Role:      req.Role,
		Sort:      "updated_at",
		Direction: "desc",
	}

	users, total, err := l.svcCtx.Repositories.User.List(l.ctx, opts)
	if err != nil {
		return nil, err
	}

	list := make([]types.AdminUserSummary, 0, len(users))
	for _, user := range users {
		list = append(list, toAdminUserSummary(user))
	}

	page, perPage := normalizePage(req.Page, req.PerPage)
	pagination := types.PaginationMeta{
		Page:       page,
		PerPage:    perPage,
		TotalCount: total,
		HasNext:    int64(page*perPage) < total,
		HasPrev:    page > 1,
	}

	return &types.AdminUserListResponse{
		Users:      list,
		Pagination: pagination,
	}, nil
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
