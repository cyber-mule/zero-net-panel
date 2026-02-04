package users

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	handlercommon "github.com/zero-net-panel/zero-net-panel/internal/handler/common"
	adminusers "github.com/zero-net-panel/zero-net-panel/internal/logic/admin/users"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// AdminListUsersHandler lists users.
func AdminListUsersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminListUsersRequest
		if err := httpx.Parse(r, &req); err != nil {
			handlercommon.RespondInvalidRequest(w, r, err)
			return
		}

		logic := adminusers.NewListLogic(r.Context(), svcCtx)
		resp, err := logic.List(&req)
		if err != nil {
			handlercommon.RespondError(w, r, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

// AdminCreateUserHandler creates a new user.
func AdminCreateUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminCreateUserRequest
		if err := httpx.Parse(r, &req); err != nil {
			handlercommon.RespondInvalidRequest(w, r, err)
			return
		}

		logic := adminusers.NewCreateLogic(r.Context(), svcCtx)
		resp, err := logic.Create(&req)
		if err != nil {
			handlercommon.RespondError(w, r, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

// AdminUpdateUserStatusHandler updates user status.
func AdminUpdateUserStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminUpdateUserStatusRequest
		if err := httpx.Parse(r, &req); err != nil {
			handlercommon.RespondInvalidRequest(w, r, err)
			return
		}

		logic := adminusers.NewUpdateStatusLogic(r.Context(), svcCtx)
		resp, err := logic.Update(&req)
		if err != nil {
			handlercommon.RespondError(w, r, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

// AdminUpdateUserRolesHandler updates user roles.
func AdminUpdateUserRolesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminUpdateUserRolesRequest
		if err := httpx.Parse(r, &req); err != nil {
			handlercommon.RespondInvalidRequest(w, r, err)
			return
		}

		logic := adminusers.NewUpdateRolesLogic(r.Context(), svcCtx)
		resp, err := logic.Update(&req)
		if err != nil {
			handlercommon.RespondError(w, r, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

// AdminResetUserPasswordHandler resets user password.
func AdminResetUserPasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminResetUserPasswordRequest
		if err := httpx.Parse(r, &req); err != nil {
			handlercommon.RespondInvalidRequest(w, r, err)
			return
		}

		logic := adminusers.NewResetPasswordLogic(r.Context(), svcCtx)
		resp, err := logic.Reset(&req)
		if err != nil {
			handlercommon.RespondError(w, r, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

// AdminForceLogoutHandler invalidates tokens.
func AdminForceLogoutHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminForceLogoutRequest
		if err := httpx.Parse(r, &req); err != nil {
			handlercommon.RespondInvalidRequest(w, r, err)
			return
		}

		logic := adminusers.NewForceLogoutLogic(r.Context(), svcCtx)
		resp, err := logic.Force(&req)
		if err != nil {
			handlercommon.RespondError(w, r, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

// AdminRotateUserCredentialHandler rotates a user's credential.
func AdminRotateUserCredentialHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminRotateUserCredentialRequest
		if err := httpx.Parse(r, &req); err != nil {
			handlercommon.RespondInvalidRequest(w, r, err)
			return
		}

		logic := adminusers.NewRotateCredentialLogic(r.Context(), svcCtx)
		resp, err := logic.Rotate(&req)
		if err != nil {
			handlercommon.RespondError(w, r, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
