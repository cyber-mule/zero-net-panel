package auditlogs

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/zeromicro/go-zero/rest/httpx"

	handlercommon "github.com/zero-net-panel/zero-net-panel/internal/handler/common"
	adminaudit "github.com/zero-net-panel/zero-net-panel/internal/logic/admin/auditlogs"
	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// AdminAuditLogListHandler returns audit log records.
func AdminAuditLogListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminAuditLogListRequest
		if err := httpx.Parse(r, &req); err != nil {
			handlercommon.RespondError(w, r, repository.ErrInvalidArgument)
			return
		}

		logic := adminaudit.NewListLogic(r.Context(), svcCtx)
		resp, err := logic.List(&req)
		if err != nil {
			handlercommon.RespondError(w, r, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

// AdminAuditLogExportHandler exports audit log records as JSON or CSV.
func AdminAuditLogExportHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminAuditLogExportRequest
		if err := httpx.Parse(r, &req); err != nil {
			handlercommon.RespondError(w, r, repository.ErrInvalidArgument)
			return
		}

		format := strings.ToLower(strings.TrimSpace(req.Format))
		if format == "" {
			format = "json"
		}

		logic := adminaudit.NewExportLogic(r.Context(), svcCtx)
		resp, err := logic.Export(&req)
		if err != nil {
			handlercommon.RespondError(w, r, err)
			return
		}

		switch format {
		case "json":
			httpx.OkJsonCtx(r.Context(), w, resp)
			return
		case "csv":
			writeAuditLogCSV(w, resp.Logs)
			return
		default:
			handlercommon.RespondError(w, r, repository.ErrInvalidArgument)
			return
		}
	}
}

func writeAuditLogCSV(w http.ResponseWriter, logs []types.AuditLogSummary) {
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=\"audit-logs.csv\"")

	writer := csv.NewWriter(w)
	_ = writer.Write([]string{
		"id",
		"actor_id",
		"actor_email",
		"actor_roles",
		"action",
		"resource_type",
		"resource_id",
		"source_ip",
		"metadata",
		"created_at",
	})

	for _, entry := range logs {
		actorID := ""
		if entry.ActorID != nil {
			actorID = strconv.FormatUint(*entry.ActorID, 10)
		}
		metadata := ""
		if entry.Metadata != nil {
			if payload, err := json.Marshal(entry.Metadata); err == nil {
				metadata = string(payload)
			}
		}
		_ = writer.Write([]string{
			strconv.FormatUint(entry.ID, 10),
			actorID,
			entry.ActorEmail,
			strings.Join(entry.ActorRoles, "|"),
			entry.Action,
			entry.ResourceType,
			entry.ResourceID,
			entry.SourceIP,
			metadata,
			strconv.FormatInt(entry.CreatedAt, 10),
		})
	}
	writer.Flush()
}
