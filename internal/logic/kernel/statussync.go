package kernel

import (
	"context"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/svc"
)

// SyncStatus pulls kernel status snapshot and updates binding health.
func SyncStatus(ctx context.Context, svcCtx *svc.ServiceContext) error {
	if svcCtx == nil || svcCtx.KernelControl == nil {
		return nil
	}

	status, err := svcCtx.KernelControl.GetStatus(ctx)
	if err != nil {
		return err
	}

	observedAt := time.Now().UTC()
	seen := make(map[string]struct{})
	for _, node := range status.Snapshot.Nodes {
		kernelID := strings.TrimSpace(node.ID)
		if kernelID == "" {
			continue
		}
		health := strings.ToLower(strings.TrimSpace(node.Health.Status))
		if health == "" {
			health = "unknown"
		}
		_, err := svcCtx.Repositories.ProtocolBinding.UpdateHealthByKernelID(ctx, kernelID, health, observedAt, "")
		if err != nil {
			logx.WithContext(ctx).Errorf("kernel status update failed for %s: %v", kernelID, err)
		}
		seen[kernelID] = struct{}{}
	}

	bindings, err := svcCtx.Repositories.ProtocolBinding.ListAll(ctx)
	if err != nil {
		return err
	}

	for _, binding := range bindings {
		if binding.KernelID == "" || strings.ToLower(binding.Status) != "active" {
			continue
		}
		if _, ok := seen[binding.KernelID]; ok {
			continue
		}
		_, err := svcCtx.Repositories.ProtocolBinding.UpdateHealthByKernelID(ctx, binding.KernelID, "offline", observedAt, "")
		if err != nil {
			logx.WithContext(ctx).Errorf("kernel offline mark failed for %s: %v", binding.KernelID, err)
		}
	}

	return nil
}
