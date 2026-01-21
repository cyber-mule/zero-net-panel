package shared

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/zero-net-panel/zero-net-panel/internal/svc"
)

// RootHandler responds with a lightweight banner message.
func RootHandler(_ *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httpx.OkJsonCtx(r.Context(), w, map[string]string{
			"message": "Hello, Network. by zeronet",
		})
	}
}
