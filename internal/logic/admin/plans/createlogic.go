package plans

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/status"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// CreateLogic 处理套餐创建。
type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewCreateLogic 构造函数。
func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Create 创建套餐。
func (l *CreateLogic) Create(req *types.AdminCreatePlanRequest) (*types.PlanSummary, error) {
	bindingIDs, err := normalizeBindingIDs(req.BindingIDs)
	if err != nil {
		return nil, err
	}
	if err := ensureBindingsExist(l.ctx, l.svcCtx.Repositories.ProtocolBinding, bindingIDs); err != nil {
		return nil, err
	}

	currency := strings.TrimSpace(req.Currency)
	if currency == "" {
		currency = "CNY"
	}
	statusCode := status.PlanStatusDraft
	if req.Status != 0 {
		normalized, err := normalizePlanStatus(req.Status)
		if err != nil {
			return nil, err
		}
		statusCode = normalized
	}

	plan := repository.Plan{
		Name:               strings.TrimSpace(req.Name),
		Slug:               strings.TrimSpace(req.Slug),
		Description:        strings.TrimSpace(req.Description),
		Tags:               append([]string(nil), req.Tags...),
		Features:           append([]string(nil), req.Features...),
		PriceCents:         req.PriceCents,
		Currency:           strings.ToUpper(currency),
		DurationDays:       req.DurationDays,
		TrafficLimitBytes:  req.TrafficLimitBytes,
		TrafficMultipliers: normalizeTrafficMultipliers(req.TrafficMultipliers),
		DevicesLimit:       req.DevicesLimit,
		SortOrder:          req.SortOrder,
		Status:             statusCode,
		Visible:            req.Visible,
	}

	var created repository.Plan
	if err := l.svcCtx.Repositories.Transaction(l.ctx, func(txRepos *repository.Repositories) error {
		createdPlan, err := txRepos.Plan.Create(l.ctx, plan)
		if err != nil {
			return err
		}
		created = createdPlan
		return txRepos.PlanProtocolBinding.Replace(l.ctx, created.ID, bindingIDs)
	}); err != nil {
		return nil, err
	}

	summary := toPlanSummary(created, nil, bindingIDs)
	return &summary, nil
}
