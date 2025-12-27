package account

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	handlercommon "github.com/zero-net-panel/zero-net-panel/internal/handler/common"
	useraccount "github.com/zero-net-panel/zero-net-panel/internal/logic/user/account"
	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// UserChangePasswordHandler allows users to update their password.
func UserChangePasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserChangePasswordRequest
		if err := httpx.Parse(r, &req); err != nil {
			handlercommon.RespondError(w, r, repository.ErrInvalidArgument)
			return
		}

		logic := useraccount.NewChangePasswordLogic(r.Context(), svcCtx)
		resp, err := logic.Change(&req)
		if err != nil {
			handlercommon.RespondError(w, r, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
