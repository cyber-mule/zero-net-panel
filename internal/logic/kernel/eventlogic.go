package kernel

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// EventLogic handles kernel event callbacks.
type EventLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewEventLogic constructs EventLogic.
func NewEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EventLogic {
	return &EventLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Handle processes kernel node events.
func (l *EventLogic) Handle(req *types.KernelNodeEventRequest) (*types.KernelNodeEventResponse, error) {
	kernelID := strings.TrimSpace(req.ID)
	if kernelID == "" {
		kernelID = strings.TrimSpace(req.NodeID)
	}
	if kernelID == "" {
		return &types.KernelNodeEventResponse{Status: "ignored"}, nil
	}

	if msg := strings.TrimSpace(req.Message); msg != "" {
		l.Logger.Infof("kernel event %s for %s: %s", req.Event, kernelID, msg)
	}

	return &types.KernelNodeEventResponse{Status: "ignored"}, nil
}
