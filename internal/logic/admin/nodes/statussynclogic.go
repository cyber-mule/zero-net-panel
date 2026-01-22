package nodes

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

// SyncStatusLogic handles manual node status synchronization.
type SyncStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewSyncStatusLogic constructs SyncStatusLogic.
func NewSyncStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncStatusLogic {
	return &SyncStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Sync triggers status sync for a subset of nodes.
func (l *SyncStatusLogic) Sync(req *types.AdminSyncNodeStatusRequest) (*types.AdminSyncNodeStatusResponse, error) {
	if req == nil {
		return nil, repository.ErrInvalidArgument
	}
	nodeIDs := uniqueNodeIDs(req.NodeIDs)
	if len(nodeIDs) == 0 {
		return nil, repository.ErrInvalidArgument
	}

	startedAt := time.Now().UTC()
	results := make([]types.NodeStatusSyncResult, 0, len(nodeIDs))
	indexByID := make(map[uint64]int, len(nodeIDs))
	for _, nodeID := range nodeIDs {
		indexByID[nodeID] = len(results)
		results = append(results, types.NodeStatusSyncResult{
			NodeID:   nodeID,
			Status:   status.NodeSyncResultStatusError,
			SyncedAt: startedAt.Unix(),
		})
	}

	type controlKey struct {
		endpoint string
		token    string
		timeout  time.Duration
	}
	controlGroups := make(map[controlKey][]uint64)

	for _, nodeID := range nodeIDs {
		res := &results[indexByID[nodeID]]
		node, err := l.svcCtx.Repositories.Node.Get(l.ctx, nodeID)
		if err != nil {
			res.Message = err.Error()
			continue
		}
		if node.Status == status.NodeStatusDisabled {
			res.Status = status.NodeSyncResultStatusSkipped
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
			token:    resolveNodeControlToken(node),
			timeout:  resolveKernelHTTPTimeout(node),
		}
		controlGroups[key] = append(controlGroups[key], nodeID)
	}

	for key, nodeGroup := range controlGroups {
		client, err := kernel.NewControlClient(kernel.HTTPOptions{
			BaseURL: key.endpoint,
			Token:   key.token,
			Timeout: key.timeout,
		})
		if err != nil {
			l.markNodeGroup(nodeGroup, status.NodeStatusOffline, status.NodeSyncResultStatusOffline, err.Error(), results, indexByID)
			continue
		}

		_, err = client.GetStatus(l.ctx)
		if err != nil {
			l.markNodeGroup(nodeGroup, status.NodeStatusOffline, status.NodeSyncResultStatusOffline, err.Error(), results, indexByID)
			continue
		}

		l.markNodeGroup(nodeGroup, status.NodeStatusOnline, status.NodeSyncResultStatusOnline, "ok", results, indexByID)
	}

	if actor, ok := security.UserFromContext(l.ctx); ok {
		l.Infof("audit: node status sync by=%s nodes=%v", strings.TrimSpace(actor.Email), nodeIDs)
	} else {
		l.Infof("audit: node status sync by=unknown nodes=%v", nodeIDs)
	}

	return &types.AdminSyncNodeStatusResponse{Results: results}, nil
}

func (l *SyncStatusLogic) markNodeGroup(nodeIDs []uint64, nodeStatus int, resultStatus int, message string, results []types.NodeStatusSyncResult, indexByID map[uint64]int) {
	if len(nodeIDs) == 0 {
		return
	}

	if err := l.svcCtx.Repositories.Node.UpdateStatusByIDs(l.ctx, nodeIDs, nodeStatus); err != nil && !errors.Is(err, repository.ErrNotFound) {
		l.Errorf("node status update failed (%d): %v", nodeStatus, err)
	}

	ts := time.Now().UTC().Unix()
	for _, nodeID := range nodeIDs {
		idx, ok := indexByID[nodeID]
		if !ok {
			continue
		}
		results[idx].Status = resultStatus
		results[idx].Message = message
		results[idx].SyncedAt = ts
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
