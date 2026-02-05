package subscriptions

import (
	"net/http"
	"strings"

	"github.com/zeromicro/go-zero/rest/httpx"

	handlercommon "github.com/zero-net-panel/zero-net-panel/internal/handler/common"
	usersub "github.com/zero-net-panel/zero-net-panel/internal/logic/user/subscription"
	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

// UserListSubscriptionsHandler returns the authenticated user's subscription list.
func UserListSubscriptionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserListSubscriptionsRequest
		if err := httpx.Parse(r, &req); err != nil {
			handlercommon.RespondInvalidRequest(w, r, err)
			return
		}

		logic := usersub.NewListLogic(r.Context(), svcCtx)
		subscriptionBase := resolveSubscriptionBaseURL(r, svcCtx)
		resp, err := logic.List(&req, subscriptionBase)
		if err != nil {
			handlercommon.RespondError(w, r, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

// UserUpdateSubscriptionTemplateHandler switches the template bound to the subscription.
func UserUpdateSubscriptionTemplateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserUpdateSubscriptionTemplateRequest
		if err := httpx.Parse(r, &req); err != nil {
			handlercommon.RespondInvalidRequest(w, r, err)
			return
		}

		logic := usersub.NewUpdateTemplateLogic(r.Context(), svcCtx)
		resp, err := logic.UpdateTemplate(&req)
		if err != nil {
			handlercommon.RespondError(w, r, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

// UserSubscriptionTrafficHandler returns traffic usage details.
func UserSubscriptionTrafficHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserSubscriptionTrafficRequest
		if err := httpx.Parse(r, &req); err != nil {
			handlercommon.RespondInvalidRequest(w, r, err)
			return
		}

		logic := usersub.NewTrafficLogic(r.Context(), svcCtx)
		resp, err := logic.Traffic(&req)
		if err != nil {
			handlercommon.RespondError(w, r, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

func resolveSubscriptionBaseURL(r *http.Request, svcCtx *svc.ServiceContext) string {
	if svcCtx != nil && svcCtx.Repositories != nil && svcCtx.Repositories.Site != nil {
		defaults := repository.SiteSettingDefaults{
			Name:    svcCtx.Config.Site.Name,
			LogoURL: svcCtx.Config.Site.LogoURL,
		}
		if setting, err := svcCtx.Repositories.Site.GetSiteSetting(r.Context(), defaults); err == nil {
			if base := normalizeConfiguredBase(setting.SubscriptionDomain, r); base != "" {
				return base
			}
			if base := normalizeConfiguredBase(setting.ServiceDomain, r); base != "" {
				return base
			}
		}
	}
	return buildRequestBaseURL(r)
}

func normalizeConfiguredBase(raw string, r *http.Request) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	if strings.Contains(raw, "://") {
		return strings.TrimRight(raw, "/")
	}
	scheme := requestScheme(r)
	if scheme == "" {
		return strings.TrimRight(raw, "/")
	}
	return strings.TrimRight(scheme+"://"+raw, "/")
}

func buildRequestBaseURL(r *http.Request) string {
	host := requestHost(r)
	if host == "" {
		return ""
	}
	scheme := requestScheme(r)
	if scheme == "" {
		return host
	}
	return scheme + "://" + host
}

func requestScheme(r *http.Request) string {
	if forwarded := r.Header.Get("Forwarded"); forwarded != "" {
		if proto := parseForwardedParam(forwarded, "proto"); proto != "" {
			return strings.ToLower(proto)
		}
	}
	if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
		return strings.ToLower(strings.TrimSpace(strings.Split(proto, ",")[0]))
	}
	if r.TLS != nil {
		return "https"
	}
	return "http"
}

func requestHost(r *http.Request) string {
	if forwarded := r.Header.Get("Forwarded"); forwarded != "" {
		if host := parseForwardedParam(forwarded, "host"); host != "" {
			return host
		}
	}
	if host := r.Header.Get("X-Forwarded-Host"); host != "" {
		return strings.TrimSpace(strings.Split(host, ",")[0])
	}
	return r.Host
}

func parseForwardedParam(value, key string) string {
	parts := strings.Split(value, ",")
	if len(parts) == 0 {
		return ""
	}
	segment := strings.TrimSpace(parts[0])
	for _, item := range strings.Split(segment, ";") {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		pair := strings.SplitN(item, "=", 2)
		if len(pair) != 2 {
			continue
		}
		if strings.EqualFold(strings.TrimSpace(pair[0]), key) {
			return strings.Trim(strings.TrimSpace(pair[1]), "\"")
		}
	}
	return ""
}
