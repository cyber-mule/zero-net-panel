package kernel

import (
	"bytes"
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// ServiceEventLogic handles kernel service event callbacks.
type ServiceEventLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewServiceEventLogic constructs ServiceEventLogic.
func NewServiceEventLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ServiceEventLogic {
	return &ServiceEventLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Handle processes service event callbacks.
func (l *ServiceEventLogic) Handle(req *types.KernelServiceEventRequest) (*types.KernelServiceEventResponse, error) {
	if req == nil {
		return &types.KernelServiceEventResponse{Status: "ignored"}, nil
	}

	event := normalizeKernelEventName(req.Event)
	switch event {
	case "user_traffic_reported":
		return l.handleUserTrafficReported(req)
	default:
		return &types.KernelServiceEventResponse{Status: "ignored"}, nil
	}
}

func (l *ServiceEventLogic) handleUserTrafficReported(req *types.KernelServiceEventRequest) (*types.KernelServiceEventResponse, error) {
	payload, err := decodeServiceEventPayload(req.Payload)
	if err != nil {
		return nil, err
	}

	userID, _ := parseUint64FromAny(payload["user_id"])
	subscriptionID, _ := parseUint64FromAny(payload["subscription_id"])
	used, ok := resolveTrafficUsed(payload)
	if !ok {
		return nil, repository.ErrInvalidArgument
	}

	if used < 0 {
		used = 0
	}

	var sub repository.Subscription
	switch {
	case subscriptionID != 0:
		sub, err = l.svcCtx.Repositories.Subscription.Get(l.ctx, subscriptionID)
	case userID != 0:
		sub, err = l.svcCtx.Repositories.Subscription.GetActiveByUser(l.ctx, userID)
	default:
		return &types.KernelServiceEventResponse{Status: "ignored"}, nil
	}
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.Repositories.Subscription.Update(l.ctx, sub.ID, repository.UpdateSubscriptionInput{
		TrafficUsedBytes: &used,
	})
	if err != nil {
		return nil, err
	}

	return &types.KernelServiceEventResponse{
		Status:   "accepted",
		Accepted: 1,
	}, nil
}

func normalizeKernelEventName(event string) string {
	name := strings.ToLower(strings.TrimSpace(event))
	name = strings.ReplaceAll(name, ".", "_")
	name = strings.ReplaceAll(name, "-", "_")
	return name
}

func decodeServiceEventPayload(raw json.RawMessage) (map[string]any, error) {
	if len(raw) == 0 {
		return nil, repository.ErrInvalidArgument
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	var payload map[string]any
	if err := decoder.Decode(&payload); err != nil {
		return nil, err
	}
	return payload, nil
}

func resolveTrafficUsed(payload map[string]any) (int64, bool) {
	if current, ok := payload["current"].(map[string]any); ok {
		if used, ok := parseInt64FromAny(current["used"]); ok {
			return used, true
		}
	}
	if traffic, ok := payload["traffic"].(map[string]any); ok {
		if used, ok := parseInt64FromAny(traffic["used"]); ok {
			return used, true
		}
	}
	if used, ok := parseInt64FromAny(payload["used"]); ok {
		return used, true
	}
	return 0, false
}

func parseUint64FromAny(value any) (uint64, bool) {
	switch v := value.(type) {
	case json.Number:
		if parsed, err := strconv.ParseUint(strings.TrimSpace(v.String()), 10, 64); err == nil {
			return parsed, true
		}
	case float64:
		if v >= 0 {
			return uint64(v), true
		}
	case int:
		if v >= 0 {
			return uint64(v), true
		}
	case int64:
		if v >= 0 {
			return uint64(v), true
		}
	case uint64:
		return v, true
	case uint:
		return uint64(v), true
	case string:
		if parsed, err := strconv.ParseUint(strings.TrimSpace(v), 10, 64); err == nil {
			return parsed, true
		}
	}
	return 0, false
}

func parseInt64FromAny(value any) (int64, bool) {
	switch v := value.(type) {
	case json.Number:
		if parsed, err := v.Int64(); err == nil {
			return parsed, true
		}
		if parsed, err := v.Float64(); err == nil {
			return int64(parsed), true
		}
	case float64:
		return int64(v), true
	case int:
		return int64(v), true
	case int64:
		return v, true
	case uint64:
		return int64(v), true
	case uint:
		return int64(v), true
	case string:
		if parsed, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64); err == nil {
			return parsed, true
		}
	}
	return 0, false
}
