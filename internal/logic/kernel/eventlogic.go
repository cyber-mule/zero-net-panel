package kernel

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
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

	status := strings.ToLower(strings.TrimSpace(req.Status))
	if status == "" {
		status = mapEventToStatus(req.Event)
	}

	observedAt := time.Now().UTC()
	if req.ObservedAt > 0 {
		observedAt = time.Unix(req.ObservedAt, 0).UTC()
	}

	if msg := strings.TrimSpace(req.Message); msg != "" {
		l.Logger.Infof("kernel event %s for %s: %s", req.Event, kernelID, msg)
	}

	_, err := l.svcCtx.Repositories.ProtocolBinding.UpdateHealthByKernelID(l.ctx, kernelID, status, observedAt, "")
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return &types.KernelNodeEventResponse{Status: "ignored"}, nil
		}
		return nil, err
	}

	return &types.KernelNodeEventResponse{Status: "ok"}, nil
}

func mapEventToStatus(event string) string {
	switch strings.ToLower(strings.TrimSpace(event)) {
	case "node_healthy":
		return "healthy"
	case "node_degraded":
		return "degraded"
	case "node_unhealthy":
		return "unhealthy"
	case "node_removed":
		return "offline"
	case "node_added":
		return "unknown"
	default:
		return "unknown"
	}
}
