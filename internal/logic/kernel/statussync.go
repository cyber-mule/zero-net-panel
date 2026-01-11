package kernel

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	adminprotocolbindings "github.com/zero-net-panel/zero-net-panel/internal/logic/admin/protocolbindings"
	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
	"github.com/zero-net-panel/zero-net-panel/pkg/kernel"
)

// SyncStatus pulls kernel status snapshot and updates node availability.
func SyncStatus(ctx context.Context, svcCtx *svc.ServiceContext) error {
	if svcCtx == nil {
		return nil
	}

	nodes, err := svcCtx.Repositories.Node.ListAll(ctx)
	if err != nil {
		return err
	}

	statusByID := make(map[uint64]string, len(nodes))
	type controlKey struct {
		endpoint string
		token    string
	}
	pairs := make(map[controlKey][]uint64)
	metaByKey := make(map[controlKey]authDebug)
	for _, node := range nodes {
		statusByID[node.ID] = strings.ToLower(strings.TrimSpace(node.Status))
		if !node.StatusSyncEnabled {
			continue
		}
		endpoint := strings.TrimSpace(node.ControlEndpoint)
		if endpoint == "" {
			continue
		}
		token := resolveControlToken(node)
		key := controlKey{endpoint: endpoint, token: token}
		pairs[key] = append(pairs[key], node.ID)
		if _, ok := metaByKey[key]; !ok {
			metaByKey[key] = buildAuthDebug(node)
		}
	}

	var lastErr error
	hadSuccess := false

	for key, nodeIDs := range pairs {
		client, err := kernel.NewControlClient(kernel.HTTPOptions{
			BaseURL: key.endpoint,
			Token:   key.token,
			Timeout: svcCtx.Config.Kernel.HTTP.Timeout,
		})
		if err != nil {
			lastErr = err
			logx.WithContext(ctx).Errorf("kernel control client init failed for %s: %v", key.endpoint, err)
			markNodeStatus(ctx, svcCtx, nodeIDs, "offline")
			continue
		}

		_, err = client.GetStatus(ctx)
		if err != nil {
			lastErr = err
			logx.WithContext(ctx).Errorf("kernel status poll failed for %s: %v", key.endpoint, err)
			if isUnauthorized(err) {
				meta := metaByKey[key]
				logx.WithContext(ctx).Errorf(
					"kernel status auth debug endpoint=%s auth=%s ak=%s sk_fp=%s token_fp=%s nodes=%v",
					key.endpoint,
					meta.AuthType,
					meta.AccessKeyMasked,
					meta.SecretFingerprint,
					meta.TokenFingerprint,
					nodeIDs,
				)
			}
			markNodeStatus(ctx, svcCtx, nodeIDs, "offline")
			continue
		}

		hadSuccess = true
		markNodeStatus(ctx, svcCtx, nodeIDs, "online")
		recovered := resolveRecoveredNodes(nodeIDs, statusByID)
		if len(recovered) > 0 {
			triggerKernelRecovery(ctx, svcCtx, recovered)
		}
	}

	if !hadSuccess {
		return lastErr
	}
	return nil
}

func resolveRecoveredNodes(nodeIDs []uint64, statusByID map[uint64]string) []uint64 {
	recovered := make([]uint64, 0, len(nodeIDs))
	for _, id := range nodeIDs {
		status := statusByID[id]
		if status == "online" || status == "disabled" {
			continue
		}
		recovered = append(recovered, id)
	}
	return recovered
}

func triggerKernelRecovery(ctx context.Context, svcCtx *svc.ServiceContext, nodeIDs []uint64) {
	if svcCtx == nil || len(nodeIDs) == 0 {
		return
	}
	logic := adminprotocolbindings.NewSyncLogic(ctx, svcCtx)
	_, err := logic.SyncBatch(&types.AdminSyncProtocolBindingsRequest{
		NodeIDs: nodeIDs,
	})
	if err != nil && !errors.Is(err, repository.ErrInvalidArgument) {
		logx.WithContext(ctx).Errorf("kernel recovery sync failed nodes=%v: %v", nodeIDs, err)
	}
}

func resolveControlToken(node repository.Node) string {
	accessKey := strings.TrimSpace(node.ControlAccessKey)
	secretKey := strings.TrimSpace(node.ControlSecretKey)
	if accessKey != "" && secretKey != "" {
		encoded := base64.StdEncoding.EncodeToString([]byte(accessKey + ":" + secretKey))
		return "Basic " + encoded
	}
	token := strings.TrimSpace(node.ControlToken)
	if token != "" {
		return token
	}
	return ""
}

type authDebug struct {
	AuthType          string
	AccessKeyMasked   string
	SecretFingerprint string
	TokenFingerprint  string
}

func buildAuthDebug(node repository.Node) authDebug {
	accessKey := strings.TrimSpace(node.ControlAccessKey)
	secretKey := strings.TrimSpace(node.ControlSecretKey)
	if accessKey != "" && secretKey != "" {
		return authDebug{
			AuthType:          "basic",
			AccessKeyMasked:   maskKey(accessKey),
			SecretFingerprint: fingerprint(secretKey),
		}
	}
	token := strings.TrimSpace(node.ControlToken)
	if token != "" {
		return authDebug{
			AuthType:         "token",
			TokenFingerprint: fingerprint(token),
		}
	}
	return authDebug{AuthType: "none"}
}

func maskKey(value string) string {
	if value == "" {
		return ""
	}
	if len(value) <= 4 {
		return "****"
	}
	return value[:2] + "***" + value[len(value)-2:]
}

func fingerprint(value string) string {
	if value == "" {
		return ""
	}
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:4])
}

func isUnauthorized(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "401") || strings.Contains(msg, "unauthorized")
}

func markNodeStatus(ctx context.Context, svcCtx *svc.ServiceContext, nodeIDs []uint64, status string) {
	if svcCtx == nil || len(nodeIDs) == 0 {
		return
	}
	if err := svcCtx.Repositories.Node.UpdateStatusByIDs(ctx, nodeIDs, status); err != nil && !errors.Is(err, repository.ErrNotFound) {
		logx.WithContext(ctx).Errorf("node status update failed (%s): %v", status, err)
	}
}
