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

// UpdateLogic handles billing option updates.
type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUpdateLogic constructs UpdateLogic.
func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Update patches a plan billing option.
func (l *UpdateLogic) Update(req *types.AdminUpdatePlanBillingOptionRequest) (*types.PlanBillingOptionSummary, error) {
	if req.PlanID == 0 || req.OptionID == 0 {
		return nil, repository.ErrInvalidArgument
	}

	plan, err := l.svcCtx.Repositories.Plan.Get(l.ctx, req.PlanID)
	if err != nil {
		return nil, err
	}

	option, err := l.svcCtx.Repositories.PlanBillingOption.Get(l.ctx, req.OptionID)
	if err != nil {
		return nil, err
	}
	if option.PlanID != req.PlanID {
		return nil, repository.ErrInvalidArgument
	}

	updates := option
	if req.Name != nil {
		updates.Name = strings.TrimSpace(*req.Name)
	}
	if req.DurationUnit != nil {
		unit := normalizeDurationUnit(*req.DurationUnit)
		if !isValidDurationUnit(unit) {
			return nil, repository.ErrInvalidArgument
		}
		updates.DurationUnit = unit
	}
	if req.DurationValue != nil {
		if *req.DurationValue <= 0 {
			return nil, repository.ErrInvalidArgument
		}
		updates.DurationValue = *req.DurationValue
	}
	if req.PriceCents != nil {
		if *req.PriceCents < 0 {
			return nil, repository.ErrInvalidArgument
		}
		updates.PriceCents = *req.PriceCents
	}
	if req.Currency != nil {
		currency := strings.TrimSpace(*req.Currency)
		if currency == "" {
			currency = strings.TrimSpace(plan.Currency)
		}
		if currency == "" {
			currency = "CNY"
		}
		updates.Currency = strings.ToUpper(currency)
	}
	if req.SortOrder != nil {
		updates.SortOrder = *req.SortOrder
	}
	if req.Status != nil {
		switch *req.Status {
		case status.PlanBillingOptionStatusDraft,
			status.PlanBillingOptionStatusActive,
			status.PlanBillingOptionStatusArchived:
			updates.Status = *req.Status
		default:
			return nil, repository.ErrInvalidArgument
		}
	}
	if req.Visible != nil {
		updates.Visible = *req.Visible
	}

	if updates.Name == "" {
		updates.Name = formatDurationLabel(updates.DurationValue, updates.DurationUnit)
	}

	updated, err := l.svcCtx.Repositories.PlanBillingOption.Update(l.ctx, option.ID, updates)
	if err != nil {
		return nil, err
	}

	summary := toPlanBillingOptionSummary(updated)
	return &summary, nil
}
