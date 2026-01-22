package plan

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/status"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// ListLogic 用户套餐列表。
type ListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewListLogic 构造函数。
func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// List 返回用户可用套餐。
func (l *ListLogic) List(req *types.UserPlanListRequest) (*types.UserPlanListResponse, error) {
	visible := true
	opts := repository.ListPlansOptions{
		Page:    1,
		PerPage: 100,
		Query:   req.Query,
		Status:  status.PlanStatusActive,
		Visible: &visible,
	}

	plans, _, err := l.svcCtx.Repositories.Plan.List(l.ctx, opts)
	if err != nil {
		return nil, err
	}

	planIDs := make([]uint64, 0, len(plans))
	for _, plan := range plans {
		planIDs = append(planIDs, plan.ID)
	}

	optionMap := map[uint64][]repository.PlanBillingOption{}
	if len(planIDs) > 0 {
		visibleOnly := true
		options, err := l.svcCtx.Repositories.PlanBillingOption.List(l.ctx, repository.ListPlanBillingOptionsOptions{
			PlanIDs: planIDs,
			Status:  status.PlanBillingOptionStatusActive,
			Visible: &visibleOnly,
		})
		if err != nil {
			return nil, err
		}
		for _, option := range options {
			optionMap[option.PlanID] = append(optionMap[option.PlanID], option)
		}
	}

	result := make([]types.UserPlanSummary, 0, len(plans))
	for _, plan := range plans {
		result = append(result, toUserPlanSummary(plan, optionMap[plan.ID]))
	}

	return &types.UserPlanListResponse{Plans: result}, nil
}
