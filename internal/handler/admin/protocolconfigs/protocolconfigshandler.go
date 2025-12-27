package protocolconfigs

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	handlercommon "github.com/zero-net-panel/zero-net-panel/internal/handler/common"
	adminconfigs "github.com/zero-net-panel/zero-net-panel/internal/logic/admin/protocolconfigs"
	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// AdminListProtocolConfigsHandler lists protocol configs.
func AdminListProtocolConfigsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminListProtocolConfigsRequest
		if err := httpx.Parse(r, &req); err != nil {
			handlercommon.RespondError(w, r, repository.ErrInvalidArgument)
			return
		}

		logic := adminconfigs.NewListLogic(r.Context(), svcCtx)
		resp, err := logic.List(&req)
		if err != nil {
			handlercommon.RespondError(w, r, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

// AdminCreateProtocolConfigHandler creates protocol config.
func AdminCreateProtocolConfigHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminCreateProtocolConfigRequest
		if err := httpx.Parse(r, &req); err != nil {
			handlercommon.RespondError(w, r, repository.ErrInvalidArgument)
			return
		}

		logic := adminconfigs.NewCreateLogic(r.Context(), svcCtx)
		resp, err := logic.Create(&req)
		if err != nil {
			handlercommon.RespondError(w, r, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

// AdminUpdateProtocolConfigHandler updates protocol config.
func AdminUpdateProtocolConfigHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminUpdateProtocolConfigRequest
		if err := httpx.Parse(r, &req); err != nil {
			handlercommon.RespondError(w, r, repository.ErrInvalidArgument)
			return
		}

		logic := adminconfigs.NewUpdateLogic(r.Context(), svcCtx)
		resp, err := logic.Update(&req)
		if err != nil {
			handlercommon.RespondError(w, r, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

// AdminDeleteProtocolConfigHandler deletes protocol config.
func AdminDeleteProtocolConfigHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminDeleteProtocolConfigRequest
		if err := httpx.Parse(r, &req); err != nil {
			handlercommon.RespondError(w, r, repository.ErrInvalidArgument)
			return
		}

		logic := adminconfigs.NewDeleteLogic(r.Context(), svcCtx)
		if err := logic.Delete(&req); err != nil {
			handlercommon.RespondError(w, r, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
