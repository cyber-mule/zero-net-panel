package account

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	handlercommon "github.com/zero-net-panel/zero-net-panel/internal/handler/common"
	useraccount "github.com/zero-net-panel/zero-net-panel/internal/logic/user/account"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// UserRotateCredentialHandler allows users to rotate their credential.
func UserRotateCredentialHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserRotateCredentialRequest
		if err := httpx.Parse(r, &req); err != nil {
			handlercommon.RespondInvalidRequest(w, r, err)
			return
		}

		logic := useraccount.NewRotateCredentialLogic(r.Context(), svcCtx)
		resp, err := logic.Rotate(&req)
		if err != nil {
			handlercommon.RespondError(w, r, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
