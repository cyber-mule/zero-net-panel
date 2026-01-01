package subscriptions

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// ListLogic handles admin subscription listing.
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

// List returns subscriptions with filters.
func (l *ListLogic) List(req *types.AdminListSubscriptionsRequest) (*types.AdminSubscriptionListResponse, error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	perPage := req.PerPage
	if perPage <= 0 || perPage > 100 {
		perPage = 20
	}

	opts := repository.ListSubscriptionsOptions{
		Page:       page,
		PerPage:    perPage,
		Sort:       "updated_at",
		Direction:  "desc",
		Query:      req.Query,
		Status:     req.Status,
		PlanName:   req.PlanName,
		PlanID:     req.PlanID,
		TemplateID: req.TemplateID,
	}
	if req.UserID != 0 {
		userID := req.UserID
		opts.UserID = &userID
	}

	subs, total, err := l.svcCtx.Repositories.Subscription.List(l.ctx, opts)
	if err != nil {
		return nil, err
	}

	userCache := make(map[uint64]repository.User)
	items := make([]types.AdminSubscriptionSummary, 0, len(subs))
	for _, sub := range subs {
		user, ok := userCache[sub.UserID]
		if !ok {
			u, err := l.svcCtx.Repositories.User.Get(l.ctx, sub.UserID)
			if err != nil {
				return nil, err
			}
			user = u
			userCache[sub.UserID] = u
		}
		items = append(items, toAdminSubscriptionSummary(sub, user))
	}

	pagination := types.PaginationMeta{
		Page:       page,
		PerPage:    perPage,
		TotalCount: total,
		HasNext:    int64(page*perPage) < total,
		HasPrev:    page > 1,
	}

	return &types.AdminSubscriptionListResponse{
		Subscriptions: items,
		Pagination:    pagination,
	}, nil
}
