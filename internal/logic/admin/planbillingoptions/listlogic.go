package planbillingoptions

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// ListLogic handles plan billing option listing.
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

// List returns billing options for a plan.
func (l *ListLogic) List(req *types.AdminListPlanBillingOptionsRequest) (*types.AdminPlanBillingOptionListResponse, error) {
	if req.PlanID == 0 {
		return nil, repository.ErrInvalidArgument
	}

	if _, err := l.svcCtx.Repositories.Plan.Get(l.ctx, req.PlanID); err != nil {
		return nil, err
	}

	options, err := l.svcCtx.Repositories.PlanBillingOption.List(l.ctx, repository.ListPlanBillingOptionsOptions{
		PlanID:  req.PlanID,
		Status:  req.Status,
		Visible: req.Visible,
	})
	if err != nil {
		return nil, err
	}

	result := make([]types.PlanBillingOptionSummary, 0, len(options))
	for _, option := range options {
		result = append(result, toPlanBillingOptionSummary(option))
	}

	return &types.AdminPlanBillingOptionListResponse{Options: result}, nil
}
