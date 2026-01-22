package nodes

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/nodecfg"
	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/security"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// UpdateLogic handles node updates.
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

// Update patches node fields.
func (l *UpdateLogic) Update(req *types.AdminUpdateNodeRequest) (*types.AdminNodeResponse, error) {
	actor, ok := security.UserFromContext(l.ctx)
	if !ok {
		return nil, repository.ErrUnauthorized
	}
	if !security.HasRole(actor, "admin") {
		return nil, repository.ErrForbidden
	}

	input := repository.UpdateNodeInput{}
	metadata := map[string]any{}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, repository.ErrInvalidArgument
		}
		input.Name = &name
		metadata["name"] = name
	}
	if req.Region != nil {
		region := strings.TrimSpace(*req.Region)
		input.Region = &region
		metadata["region"] = region
	}
	if req.Country != nil {
		country := strings.TrimSpace(*req.Country)
		input.Country = &country
		metadata["country"] = country
	}
	if req.ISP != nil {
		isp := strings.TrimSpace(*req.ISP)
		input.ISP = &isp
		metadata["isp"] = isp
	}
	if req.Status != nil {
		statusCode, err := normalizeNodeStatus(*req.Status)
		if err != nil {
			return nil, err
		}
		input.Status = &statusCode
		metadata["status"] = statusCode
	}
	if req.Tags != nil {
		tags := normalizeTags(req.Tags)
		input.Tags = &tags
		metadata["tags"] = tags
	}
	if req.CapacityMbps != nil {
		if *req.CapacityMbps < 0 {
			return nil, repository.ErrInvalidArgument
		}
		input.CapacityMbps = req.CapacityMbps
		metadata["capacity_mbps"] = *req.CapacityMbps
	}
	if req.Description != nil {
		desc := strings.TrimSpace(*req.Description)
		input.Description = &desc
		metadata["description"] = desc
	}
	if req.AccessAddress != nil {
		address := strings.TrimSpace(*req.AccessAddress)
		input.AccessAddress = &address
		metadata["access_address"] = address
	}
	if req.ControlEndpoint != nil {
		endpoint := strings.TrimSpace(*req.ControlEndpoint)
		if endpoint == "" {
			return nil, fmt.Errorf("%w: node control endpoint required", repository.ErrInvalidArgument)
		}
		input.ControlEndpoint = &endpoint
		metadata["control_endpoint"] = endpoint
	}
	if req.ControlAccessKey != nil {
		accessKey := strings.TrimSpace(*req.ControlAccessKey)
		input.ControlAccessKey = &accessKey
	} else if req.AK != nil {
		accessKey := strings.TrimSpace(*req.AK)
		input.ControlAccessKey = &accessKey
	}
	if req.ControlSecretKey != nil {
		secretKey := strings.TrimSpace(*req.ControlSecretKey)
		input.ControlSecretKey = &secretKey
	} else if req.SK != nil {
		secretKey := strings.TrimSpace(*req.SK)
		input.ControlSecretKey = &secretKey
	}
	if req.ControlToken != nil {
		token := strings.TrimSpace(*req.ControlToken)
		input.ControlToken = &token
	}
	kernelUpdate := req.KernelDefaultProtocol != nil ||
		req.KernelHTTPTimeoutSeconds != nil ||
		req.KernelStatusPollIntervalSeconds != nil ||
		req.KernelStatusPollBackoffEnabled != nil ||
		req.KernelStatusPollBackoffMaxIntervalSeconds != nil ||
		req.KernelStatusPollBackoffMultiplier != nil ||
		req.KernelStatusPollBackoffJitter != nil ||
		req.KernelOfflineProbeMaxIntervalSeconds != nil
	if kernelUpdate {
		current, err := l.svcCtx.Repositories.Node.Get(l.ctx, req.NodeID)
		if err != nil {
			return nil, err
		}

		kernelDefaultProtocol := strings.TrimSpace(current.KernelDefaultProtocol)
		if kernelDefaultProtocol == "" {
			kernelDefaultProtocol = nodecfg.DefaultKernelProtocol
		}
		kernelHTTPTimeoutSeconds := current.KernelHTTPTimeoutSeconds
		if kernelHTTPTimeoutSeconds <= 0 {
			kernelHTTPTimeoutSeconds = nodecfg.DefaultKernelHTTPTimeoutSeconds
		}
		kernelStatusPollIntervalSeconds := current.KernelStatusPollIntervalSeconds
		if kernelStatusPollIntervalSeconds < 0 {
			kernelStatusPollIntervalSeconds = nodecfg.DefaultKernelStatusPollIntervalSeconds
		}
		kernelStatusPollBackoffEnabled := current.KernelStatusPollBackoffEnabled
		kernelStatusPollBackoffMaxIntervalSeconds := current.KernelStatusPollBackoffMaxIntervalSeconds
		if kernelStatusPollBackoffMaxIntervalSeconds < 0 {
			kernelStatusPollBackoffMaxIntervalSeconds = nodecfg.DefaultKernelStatusPollBackoffMaxIntervalSeconds
		}
		kernelStatusPollBackoffMultiplier := current.KernelStatusPollBackoffMultiplier
		if kernelStatusPollBackoffMultiplier <= 1 {
			kernelStatusPollBackoffMultiplier = nodecfg.DefaultKernelStatusPollBackoffMultiplier
		}
		kernelStatusPollBackoffJitter := current.KernelStatusPollBackoffJitter
		if kernelStatusPollBackoffJitter < 0 || kernelStatusPollBackoffJitter > 1 {
			kernelStatusPollBackoffJitter = nodecfg.DefaultKernelStatusPollBackoffJitter
		}
		kernelOfflineProbeMaxIntervalSeconds := current.KernelOfflineProbeMaxIntervalSeconds
		if kernelOfflineProbeMaxIntervalSeconds < 0 {
			kernelOfflineProbeMaxIntervalSeconds = nodecfg.DefaultKernelOfflineProbeMaxIntervalSeconds
		}

		if req.KernelDefaultProtocol != nil {
			kernelDefaultProtocol = strings.TrimSpace(*req.KernelDefaultProtocol)
			if kernelDefaultProtocol == "" {
				return nil, repository.ErrInvalidArgument
			}
			input.KernelDefaultProtocol = &kernelDefaultProtocol
			metadata["kernel_default_protocol"] = kernelDefaultProtocol
		}
		if req.KernelHTTPTimeoutSeconds != nil {
			kernelHTTPTimeoutSeconds = *req.KernelHTTPTimeoutSeconds
			input.KernelHTTPTimeoutSeconds = req.KernelHTTPTimeoutSeconds
			metadata["kernel_http_timeout_seconds"] = kernelHTTPTimeoutSeconds
		}
		if req.KernelStatusPollIntervalSeconds != nil {
			kernelStatusPollIntervalSeconds = *req.KernelStatusPollIntervalSeconds
			input.KernelStatusPollIntervalSeconds = req.KernelStatusPollIntervalSeconds
			metadata["kernel_status_poll_interval_seconds"] = kernelStatusPollIntervalSeconds
		}
		if req.KernelStatusPollBackoffEnabled != nil {
			kernelStatusPollBackoffEnabled = *req.KernelStatusPollBackoffEnabled
			input.KernelStatusPollBackoffEnabled = req.KernelStatusPollBackoffEnabled
			metadata["kernel_status_poll_backoff_enabled"] = kernelStatusPollBackoffEnabled
		}
		if req.KernelStatusPollBackoffMaxIntervalSeconds != nil {
			kernelStatusPollBackoffMaxIntervalSeconds = *req.KernelStatusPollBackoffMaxIntervalSeconds
			input.KernelStatusPollBackoffMaxIntervalSeconds = req.KernelStatusPollBackoffMaxIntervalSeconds
			metadata["kernel_status_poll_backoff_max_interval_seconds"] = kernelStatusPollBackoffMaxIntervalSeconds
		}
		if req.KernelStatusPollBackoffMultiplier != nil {
			kernelStatusPollBackoffMultiplier = *req.KernelStatusPollBackoffMultiplier
			input.KernelStatusPollBackoffMultiplier = req.KernelStatusPollBackoffMultiplier
			metadata["kernel_status_poll_backoff_multiplier"] = kernelStatusPollBackoffMultiplier
		}
		if req.KernelStatusPollBackoffJitter != nil {
			kernelStatusPollBackoffJitter = *req.KernelStatusPollBackoffJitter
			input.KernelStatusPollBackoffJitter = req.KernelStatusPollBackoffJitter
			metadata["kernel_status_poll_backoff_jitter"] = kernelStatusPollBackoffJitter
		}
		if req.KernelOfflineProbeMaxIntervalSeconds != nil {
			kernelOfflineProbeMaxIntervalSeconds = *req.KernelOfflineProbeMaxIntervalSeconds
			input.KernelOfflineProbeMaxIntervalSeconds = req.KernelOfflineProbeMaxIntervalSeconds
			metadata["kernel_offline_probe_max_interval_seconds"] = kernelOfflineProbeMaxIntervalSeconds
		}

		if kernelDefaultProtocol == "" {
			return nil, repository.ErrInvalidArgument
		}
		if kernelHTTPTimeoutSeconds <= 0 {
			return nil, repository.ErrInvalidArgument
		}
		if kernelStatusPollIntervalSeconds < 0 {
			return nil, repository.ErrInvalidArgument
		}
		if kernelStatusPollBackoffMaxIntervalSeconds < 0 {
			return nil, repository.ErrInvalidArgument
		}
		if kernelStatusPollBackoffMultiplier <= 1 {
			return nil, repository.ErrInvalidArgument
		}
		if kernelStatusPollBackoffJitter < 0 || kernelStatusPollBackoffJitter > 1 {
			return nil, repository.ErrInvalidArgument
		}
		if kernelOfflineProbeMaxIntervalSeconds < 0 {
			return nil, repository.ErrInvalidArgument
		}
		if kernelStatusPollBackoffEnabled && kernelStatusPollIntervalSeconds <= 0 {
			return nil, repository.ErrInvalidArgument
		}
		if kernelStatusPollBackoffEnabled && kernelStatusPollBackoffMaxIntervalSeconds > 0 &&
			kernelStatusPollBackoffMaxIntervalSeconds < kernelStatusPollIntervalSeconds {
			return nil, repository.ErrInvalidArgument
		}
	}
	if req.StatusSyncEnabled != nil {
		input.StatusSyncEnabled = req.StatusSyncEnabled
		metadata["status_sync_enabled"] = *req.StatusSyncEnabled
	}

	var updated repository.Node
	if err := l.svcCtx.Repositories.Transaction(l.ctx, func(txRepos *repository.Repositories) error {
		node, err := txRepos.Node.Update(l.ctx, req.NodeID, input)
		if err != nil {
			return err
		}
		updated = node

		if len(metadata) == 0 {
			return nil
		}
		_, err = txRepos.AuditLog.Create(l.ctx, repository.AuditLog{
			ActorID:      &actor.ID,
			ActorEmail:   actor.Email,
			ActorRoles:   actor.Roles,
			Action:       "admin.node.update",
			ResourceType: "node",
			ResourceID:   fmt.Sprintf("%d", updated.ID),
			Metadata:     metadata,
		})
		return err
	}); err != nil {
		return nil, err
	}

	return &types.AdminNodeResponse{
		Node: mapNodeSummary(updated),
	}, nil
}
