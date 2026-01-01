package subscriptions

import (
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/rest/httpx"

	handlercommon "github.com/zero-net-panel/zero-net-panel/internal/handler/common"
	publicsub "github.com/zero-net-panel/zero-net-panel/internal/logic/public/subscription"
	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
)

type publicSubscriptionRequest struct {
	Token string `path:"token"`
}

// PublicSubscriptionDownloadHandler renders subscription content by token.
func PublicSubscriptionDownloadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req publicSubscriptionRequest
		if err := httpx.Parse(r, &req); err != nil {
			handlercommon.RespondError(w, r, repository.ErrInvalidArgument)
			return
		}

		logic := publicsub.NewDownloadLogic(r.Context(), svcCtx)
		result, err := logic.Download(req.Token, r.Header.Get("User-Agent"))
		if err != nil {
			handlercommon.RespondError(w, r, err)
			return
		}

		if match := strings.TrimSpace(r.Header.Get("If-None-Match")); match != "" {
			candidate := strings.TrimPrefix(match, "W/")
			candidate = strings.Trim(candidate, "\"")
			if candidate == result.ETag {
				w.Header().Set("ETag", result.ETag)
				w.WriteHeader(http.StatusNotModified)
				return
			}
		}

		w.Header().Set("Content-Type", result.ContentType)
		w.Header().Set("ETag", result.ETag)
		w.Header().Set("Cache-Control", "no-cache")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(result.Content))
	}
}
