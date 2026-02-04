package auth

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	handlercommon "github.com/zero-net-panel/zero-net-panel/internal/handler/common"
	authlogic "github.com/zero-net-panel/zero-net-panel/internal/logic/auth"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// AuthLoginHandler handles user sign-in requests and returns issued tokens.
func AuthLoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AuthLoginRequest
		if err := httpx.Parse(r, &req); err != nil {
			handlercommon.RespondInvalidRequest(w, r, err)
			return
		}

		logic := authlogic.NewLoginLogic(r.Context(), svcCtx)
		resp, err := logic.Login(&req)
		if err != nil {
			handlercommon.RespondError(w, r, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
