package protocolbindings

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/security"
	"github.com/zero-net-panel/zero-net-panel/internal/status"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
	"github.com/zero-net-panel/zero-net-panel/pkg/kernel"
)

// StatusSyncLogic handles manual protocol status synchronization.
type StatusSyncLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewStatusSyncLogic constructs StatusSyncLogic.
func NewStatusSyncLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StatusSyncLogic {
	return &StatusSyncLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Sync pulls protocol health status from kernel control plane.
func (l *StatusSyncLogic) Sync(req *types.AdminSyncProtocolBindingStatusRequest) (*types.AdminSyncProtocolBindingStatusResponse, error) {
	if req == nil {
		return nil, repository.ErrInvalidArgument
	}
	nodeIDs := uniqueNodeIDs(req.NodeIDs)
	if len(nodeIDs) == 0 {
		return nil, repository.ErrInvalidArgument
	}

	startedAt := time.Now().UTC()
	results := make([]types.ProtocolBindingStatusSyncResult, 0, len(nodeIDs))
	indexByID := make(map[uint64]int, len(nodeIDs))
	for _, nodeID := range nodeIDs {
		indexByID[nodeID] = len(results)
		results = append(results, types.ProtocolBindingStatusSyncResult{
			NodeID:   nodeID,
			Status:   status.SyncResultStatusError,
			SyncedAt: startedAt.Unix(),
		})
	}

	type controlKey struct {
		endpoint string
		token    string
		timeout  time.Duration
	}
	controlGroups := make(map[controlKey][]uint64)
	validNodeIDs := make([]uint64, 0, len(nodeIDs))

	for _, nodeID := range nodeIDs {
		res := &results[indexByID[nodeID]]
		node, err := l.svcCtx.Repositories.Node.Get(l.ctx, nodeID)
		if err != nil {
			res.Message = err.Error()
			continue
		}
		if node.Status == status.NodeStatusDisabled {
			res.Status = status.SyncResultStatusSkipped
			res.Message = "node disabled"
			continue
		}
		endpoint := strings.TrimSpace(node.ControlEndpoint)
		if endpoint == "" {
			res.Message = "node control endpoint not configured"
			continue
		}

		key := controlKey{
			endpoint: endpoint,
			token:    resolveControlToken(node),
			timeout:  resolveKernelHTTPTimeout(node),
		}
		controlGroups[key] = append(controlGroups[key], nodeID)
		validNodeIDs = append(validNodeIDs, nodeID)
	}

	bindingsByNode := make(map[uint64][]repository.ProtocolBinding)
	if len(validNodeIDs) > 0 {
		bindings, err := l.svcCtx.Repositories.ProtocolBinding.ListByNodeIDs(l.ctx, validNodeIDs)
		if err != nil {
			return nil, err
		}
		for _, binding := range bindings {
			bindingsByNode[binding.NodeID] = append(bindingsByNode[binding.NodeID], binding)
		}
	}

	for key, nodeGroup := range controlGroups {
		client, err := kernel.NewControlClient(kernel.HTTPOptions{
			BaseURL: key.endpoint,
			Token:   key.token,
			Timeout: key.timeout,
		})
		if err != nil {
			l.markResults(nodeGroup, status.SyncResultStatusError, err.Error(), nil, results, indexByID)
			continue
		}

		snapshot, err := client.GetStatus(l.ctx)
		if err != nil {
			l.markResults(nodeGroup, status.SyncResultStatusError, err.Error(), nil, results, indexByID)
			continue
		}

		healthByKernel := make(map[string]int)
		for _, node := range snapshot.Snapshot.Nodes {
			kernelID := strings.TrimSpace(node.ID)
			if kernelID == "" {
				continue
			}
			healthByKernel[kernelID] = mapKernelHealthStatus(node.Health.Status)
		}

		observedAt := time.Now().UTC()
		updatedCounts := make(map[uint64]int)
		var updateErr error
		for _, nodeID := range nodeGroup {
			for _, binding := range bindingsByNode[nodeID] {
				if binding.KernelID == "" || binding.Status != status.ProtocolBindingStatusActive {
					continue
				}
				health, ok := healthByKernel[binding.KernelID]
				if !ok {
					health = status.ProtocolBindingHealthStatusOffline
				}
				_, err := l.svcCtx.Repositories.ProtocolBinding.UpdateHealthByKernelIDForNodes(
					l.ctx,
					binding.KernelID,
					[]uint64{nodeID},
					health,
					observedAt,
					"",
				)
				if err != nil {
					if !errors.Is(err, repository.ErrNotFound) {
						updateErr = err
					}
					continue
				}
				updatedCounts[nodeID]++
			}
		}

		statusValue := status.SyncResultStatusSynced
		message := "ok"
		if updateErr != nil {
			statusValue = status.SyncResultStatusError
			message = updateErr.Error()
		}
		l.markResults(nodeGroup, statusValue, message, updatedCounts, results, indexByID)
	}

	if actor, ok := security.UserFromContext(l.ctx); ok {
		l.Infof("audit: protocol status sync by=%s nodes=%v", strings.TrimSpace(actor.Email), nodeIDs)
	} else {
		l.Infof("audit: protocol status sync by=unknown nodes=%v", nodeIDs)
	}

	return &types.AdminSyncProtocolBindingStatusResponse{Results: results}, nil
}

func (l *StatusSyncLogic) markResults(nodeIDs []uint64, statusCode int, message string, updated map[uint64]int, results []types.ProtocolBindingStatusSyncResult, indexByID map[uint64]int) {
	if len(nodeIDs) == 0 {
		return
	}
	ts := time.Now().UTC().Unix()
	for _, nodeID := range nodeIDs {
		idx, ok := indexByID[nodeID]
		if !ok {
			continue
		}
		results[idx].Status = statusCode
		results[idx].Message = message
		results[idx].SyncedAt = ts
		if updated != nil {
			results[idx].Updated = updated[nodeID]
		}
	}
}

func mapKernelHealthStatus(raw string) int {
	normalized := strings.ToLower(strings.TrimSpace(raw))
	switch normalized {
	case "healthy":
		return status.ProtocolBindingHealthStatusHealthy
	case "degraded":
		return status.ProtocolBindingHealthStatusDegraded
	case "unhealthy":
		return status.ProtocolBindingHealthStatusUnhealthy
	case "offline":
		return status.ProtocolBindingHealthStatusOffline
	default:
		return status.ProtocolBindingHealthStatusUnknown
	}
}

func uniqueNodeIDs(ids []uint64) []uint64 {
	if len(ids) == 0 {
		return nil
	}
	seen := make(map[uint64]struct{}, len(ids))
	result := make([]uint64, 0, len(ids))
	for _, id := range ids {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		result = append(result, id)
	}
	return result
}
