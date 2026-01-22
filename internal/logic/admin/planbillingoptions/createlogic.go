package planbillingoptions

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/status"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// CreateLogic handles billing option creation.
type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewCreateLogic constructs CreateLogic.
func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Create creates a plan billing option.
func (l *CreateLogic) Create(req *types.AdminCreatePlanBillingOptionRequest) (*types.PlanBillingOptionSummary, error) {
	if req.PlanID == 0 {
		return nil, repository.ErrInvalidArgument
	}

	plan, err := l.svcCtx.Repositories.Plan.Get(l.ctx, req.PlanID)
	if err != nil {
		return nil, err
	}

	unit := normalizeDurationUnit(req.DurationUnit)
	if !isValidDurationUnit(unit) || req.DurationValue <= 0 {
		return nil, repository.ErrInvalidArgument
	}
	if req.PriceCents < 0 {
		return nil, repository.ErrInvalidArgument
	}

	currency := strings.TrimSpace(req.Currency)
	if currency == "" {
		currency = strings.TrimSpace(plan.Currency)
	}
	if currency == "" {
		currency = "CNY"
	}

	statusCode := req.Status
	if statusCode == 0 {
		statusCode = status.PlanBillingOptionStatusDraft
	}
	switch statusCode {
	case status.PlanBillingOptionStatusDraft, status.PlanBillingOptionStatusActive, status.PlanBillingOptionStatusArchived:
	default:
		return nil, repository.ErrInvalidArgument
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		name = formatDurationLabel(req.DurationValue, unit)
	}

	option := repository.PlanBillingOption{
		PlanID:        plan.ID,
		Name:          name,
		DurationValue: req.DurationValue,
		DurationUnit:  unit,
		PriceCents:    req.PriceCents,
		Currency:      strings.ToUpper(currency),
		SortOrder:     req.SortOrder,
		Status:        statusCode,
		Visible:       req.Visible,
	}

	created, err := l.svcCtx.Repositories.PlanBillingOption.Create(l.ctx, option)
	if err != nil {
		return nil, err
	}

	summary := toPlanBillingOptionSummary(created)
	return &summary, nil
}
