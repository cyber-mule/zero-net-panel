package protocols

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	handlercommon "github.com/zero-net-panel/zero-net-panel/internal/handler/common"
	adminprotocols "github.com/zero-net-panel/zero-net-panel/internal/logic/admin/protocols"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
)

// AdminListProtocolsHandler returns available protocols for the admin console.
func AdminListProtocolsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logic := adminprotocols.NewListLogic(r.Context(), svcCtx)
		resp, err := logic.List()
		if err != nil {
			handlercommon.RespondError(w, r, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
