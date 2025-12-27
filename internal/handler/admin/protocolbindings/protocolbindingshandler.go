package protocolbindings

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	handlercommon "github.com/zero-net-panel/zero-net-panel/internal/handler/common"
	adminbindings "github.com/zero-net-panel/zero-net-panel/internal/logic/admin/protocolbindings"
	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// AdminListProtocolBindingsHandler lists protocol bindings.
func AdminListProtocolBindingsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminListProtocolBindingsRequest
		if err := httpx.Parse(r, &req); err != nil {
			handlercommon.RespondError(w, r, repository.ErrInvalidArgument)
			return
		}

		logic := adminbindings.NewListLogic(r.Context(), svcCtx)
		resp, err := logic.List(&req)
		if err != nil {
			handlercommon.RespondError(w, r, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

// AdminCreateProtocolBindingHandler creates a new binding.
func AdminCreateProtocolBindingHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminCreateProtocolBindingRequest
		if err := httpx.Parse(r, &req); err != nil {
			handlercommon.RespondError(w, r, repository.ErrInvalidArgument)
			return
		}

		logic := adminbindings.NewCreateLogic(r.Context(), svcCtx)
		resp, err := logic.Create(&req)
		if err != nil {
			handlercommon.RespondError(w, r, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

// AdminUpdateProtocolBindingHandler updates a binding.
func AdminUpdateProtocolBindingHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminUpdateProtocolBindingRequest
		if err := httpx.Parse(r, &req); err != nil {
			handlercommon.RespondError(w, r, repository.ErrInvalidArgument)
			return
		}

		logic := adminbindings.NewUpdateLogic(r.Context(), svcCtx)
		resp, err := logic.Update(&req)
		if err != nil {
			handlercommon.RespondError(w, r, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

// AdminDeleteProtocolBindingHandler deletes a binding.
func AdminDeleteProtocolBindingHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminDeleteProtocolBindingRequest
		if err := httpx.Parse(r, &req); err != nil {
			handlercommon.RespondError(w, r, repository.ErrInvalidArgument)
			return
		}

		logic := adminbindings.NewDeleteLogic(r.Context(), svcCtx)
		if err := logic.Delete(&req); err != nil {
			handlercommon.RespondError(w, r, err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// AdminSyncProtocolBindingHandler triggers sync for a binding.
func AdminSyncProtocolBindingHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminSyncProtocolBindingRequest
		if err := httpx.Parse(r, &req); err != nil {
			handlercommon.RespondError(w, r, repository.ErrInvalidArgument)
			return
		}

		logic := adminbindings.NewSyncLogic(r.Context(), svcCtx)
		resp, err := logic.SyncSingle(&req)
		if err != nil {
			handlercommon.RespondError(w, r, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

// AdminSyncProtocolBindingsHandler triggers sync for multiple bindings.
func AdminSyncProtocolBindingsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminSyncProtocolBindingsRequest
		if err := httpx.Parse(r, &req); err != nil {
			handlercommon.RespondError(w, r, repository.ErrInvalidArgument)
			return
		}

		logic := adminbindings.NewSyncLogic(r.Context(), svcCtx)
		resp, err := logic.SyncBatch(&req)
		if err != nil {
			handlercommon.RespondError(w, r, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
