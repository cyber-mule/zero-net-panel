package account

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	handlercommon "github.com/zero-net-panel/zero-net-panel/internal/handler/common"
	useraccount "github.com/zero-net-panel/zero-net-panel/internal/logic/user/account"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// UserEmailChangeCodeHandler sends a verification code for email change.
func UserEmailChangeCodeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserEmailChangeCodeRequest
		if err := httpx.Parse(r, &req); err != nil {
			handlercommon.RespondInvalidRequest(w, r, err)
			return
		}

		logic := useraccount.NewEmailChangeCodeLogic(r.Context(), svcCtx)
		resp, err := logic.Send(&req)
		if err != nil {
			handlercommon.RespondError(w, r, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

// UserChangeEmailHandler applies the email change after verification.
func UserChangeEmailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserChangeEmailRequest
		if err := httpx.Parse(r, &req); err != nil {
			handlercommon.RespondInvalidRequest(w, r, err)
			return
		}

		logic := useraccount.NewChangeEmailLogic(r.Context(), svcCtx)
		resp, err := logic.Change(&req)
		if err != nil {
			handlercommon.RespondError(w, r, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
