package subscriptions

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	handlercommon "github.com/zero-net-panel/zero-net-panel/internal/handler/common"
	adminsubs "github.com/zero-net-panel/zero-net-panel/internal/logic/admin/subscriptions"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// AdminListSubscriptionsHandler lists subscriptions.
func AdminListSubscriptionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminListSubscriptionsRequest
		if err := httpx.Parse(r, &req); err != nil {
			handlercommon.RespondInvalidRequest(w, r, err)
			return
		}

		logic := adminsubs.NewListLogic(r.Context(), svcCtx)
		resp, err := logic.List(&req)
		if err != nil {
			handlercommon.RespondError(w, r, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

// AdminGetSubscriptionHandler gets a subscription detail.
func AdminGetSubscriptionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminGetSubscriptionRequest
		if err := httpx.Parse(r, &req); err != nil {
			handlercommon.RespondInvalidRequest(w, r, err)
			return
		}

		logic := adminsubs.NewGetLogic(r.Context(), svcCtx)
		resp, err := logic.Get(&req)
		if err != nil {
			handlercommon.RespondError(w, r, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

// AdminCreateSubscriptionHandler creates a subscription.
func AdminCreateSubscriptionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminCreateSubscriptionRequest
		if err := httpx.Parse(r, &req); err != nil {
			handlercommon.RespondInvalidRequest(w, r, err)
			return
		}

		logic := adminsubs.NewCreateLogic(r.Context(), svcCtx)
		resp, err := logic.Create(&req)
		if err != nil {
			handlercommon.RespondError(w, r, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

// AdminUpdateSubscriptionHandler updates a subscription.
func AdminUpdateSubscriptionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminUpdateSubscriptionRequest
		if err := httpx.Parse(r, &req); err != nil {
			handlercommon.RespondInvalidRequest(w, r, err)
			return
		}

		logic := adminsubs.NewUpdateLogic(r.Context(), svcCtx)
		resp, err := logic.Update(&req)
		if err != nil {
			handlercommon.RespondError(w, r, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

// AdminDisableSubscriptionHandler disables a subscription.
func AdminDisableSubscriptionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminDisableSubscriptionRequest
		if err := httpx.Parse(r, &req); err != nil {
			handlercommon.RespondInvalidRequest(w, r, err)
			return
		}

		logic := adminsubs.NewDisableLogic(r.Context(), svcCtx)
		resp, err := logic.Disable(&req)
		if err != nil {
			handlercommon.RespondError(w, r, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

// AdminExtendSubscriptionHandler extends subscription expiry.
func AdminExtendSubscriptionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminExtendSubscriptionRequest
		if err := httpx.Parse(r, &req); err != nil {
			handlercommon.RespondInvalidRequest(w, r, err)
			return
		}

		logic := adminsubs.NewExtendLogic(r.Context(), svcCtx)
		resp, err := logic.Extend(&req)
		if err != nil {
			handlercommon.RespondError(w, r, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
