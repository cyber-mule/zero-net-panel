package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest"

	adminAnnouncements "github.com/zero-net-panel/zero-net-panel/internal/handler/admin/announcements"
	adminAuditLogs "github.com/zero-net-panel/zero-net-panel/internal/handler/admin/auditlogs"
	adminCoupons "github.com/zero-net-panel/zero-net-panel/internal/handler/admin/coupons"
	adminDashboard "github.com/zero-net-panel/zero-net-panel/internal/handler/admin/dashboard"
	adminNodes "github.com/zero-net-panel/zero-net-panel/internal/handler/admin/nodes"
	adminOrders "github.com/zero-net-panel/zero-net-panel/internal/handler/admin/orders"
	adminPaymentChannels "github.com/zero-net-panel/zero-net-panel/internal/handler/admin/paymentchannels"
	adminPlanBillingOptions "github.com/zero-net-panel/zero-net-panel/internal/handler/admin/planbillingoptions"
	adminPlans "github.com/zero-net-panel/zero-net-panel/internal/handler/admin/plans"
	adminProtocolBindings "github.com/zero-net-panel/zero-net-panel/internal/handler/admin/protocolbindings"
	adminProtocolConfigs "github.com/zero-net-panel/zero-net-panel/internal/handler/admin/protocolconfigs"
	adminSecurity "github.com/zero-net-panel/zero-net-panel/internal/handler/admin/security"
	adminSite "github.com/zero-net-panel/zero-net-panel/internal/handler/admin/site"
	adminSubscriptions "github.com/zero-net-panel/zero-net-panel/internal/handler/admin/subscriptions"
	adminTemplates "github.com/zero-net-panel/zero-net-panel/internal/handler/admin/templates"
	adminUsers "github.com/zero-net-panel/zero-net-panel/internal/handler/admin/users"
	authhandlers "github.com/zero-net-panel/zero-net-panel/internal/handler/auth"
	kernelhandlers "github.com/zero-net-panel/zero-net-panel/internal/handler/kernel"
	sharedhandlers "github.com/zero-net-panel/zero-net-panel/internal/handler/shared"
	userAccount "github.com/zero-net-panel/zero-net-panel/internal/handler/user/account"
	userAnnouncements "github.com/zero-net-panel/zero-net-panel/internal/handler/user/announcements"
	userNodes "github.com/zero-net-panel/zero-net-panel/internal/handler/user/nodes"
	userOrders "github.com/zero-net-panel/zero-net-panel/internal/handler/user/orders"
	userPaymentChannels "github.com/zero-net-panel/zero-net-panel/internal/handler/user/paymentchannels"
	userPlans "github.com/zero-net-panel/zero-net-panel/internal/handler/user/plans"
	userSubscriptions "github.com/zero-net-panel/zero-net-panel/internal/handler/user/subscriptions"
	"github.com/zero-net-panel/zero-net-panel/internal/middleware"
	"github.com/zero-net-panel/zero-net-panel/internal/svc"
)

func RegisterHandlers(server *rest.Server, svcCtx *svc.ServiceContext) {
	authMiddleware := middleware.NewAuthMiddleware(svcCtx.Auth, svcCtx.Repositories.User)
	thirdPartyMiddleware := middleware.NewThirdPartyMiddleware(svcCtx.Repositories.Security)
	accessMiddleware := middleware.NewAccessMiddleware(svcCtx.Config.Admin.Access)
	webhookMiddleware := middleware.NewWebhookMiddleware(svcCtx.Config.Webhook)

	server.Use(middleware.HTTPMetricsMiddleware{}.Handler)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/ping",
				Handler: sharedhandlers.PingHandler(svcCtx),
			},
		},
		rest.WithPrefix("/api/v1"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/login",
				Handler: authhandlers.AuthLoginHandler(svcCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/refresh",
				Handler: authhandlers.AuthRefreshHandler(svcCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/register",
				Handler: authhandlers.AuthRegisterHandler(svcCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/verify",
				Handler: authhandlers.AuthVerifyHandler(svcCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/forgot",
				Handler: authhandlers.AuthForgotPasswordHandler(svcCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/reset",
				Handler: authhandlers.AuthResetPasswordHandler(svcCtx),
			},
		},
		rest.WithPrefix("/api/v1/auth"),
	)

	adminRoutes := []rest.Route{
		{
			Method:  http.MethodGet,
			Path:    "/dashboard",
			Handler: adminDashboard.AdminDashboardHandler(svcCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/users",
			Handler: adminUsers.AdminListUsersHandler(svcCtx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: adminUsers.AdminCreateUserHandler(svcCtx),
		},
		{
			Method:  http.MethodPatch,
			Path:    "/users/:id/status",
			Handler: adminUsers.AdminUpdateUserStatusHandler(svcCtx),
		},
		{
			Method:  http.MethodPatch,
			Path:    "/users/:id/roles",
			Handler: adminUsers.AdminUpdateUserRolesHandler(svcCtx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/users/:id/reset-password",
			Handler: adminUsers.AdminResetUserPasswordHandler(svcCtx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/users/:id/force-logout",
			Handler: adminUsers.AdminForceLogoutHandler(svcCtx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/users/:id/credentials/rotate",
			Handler: adminUsers.AdminRotateUserCredentialHandler(svcCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/nodes",
			Handler: adminNodes.AdminListNodesHandler(svcCtx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/nodes",
			Handler: adminNodes.AdminCreateNodeHandler(svcCtx),
		},
		{
			Method:  http.MethodPatch,
			Path:    "/nodes/:id",
			Handler: adminNodes.AdminUpdateNodeHandler(svcCtx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/nodes/:id/disable",
			Handler: adminNodes.AdminDisableNodeHandler(svcCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/nodes/:id/kernels",
			Handler: adminNodes.AdminNodeKernelsHandler(svcCtx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/nodes/:id/kernels",
			Handler: adminNodes.AdminUpsertNodeKernelHandler(svcCtx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/nodes/:id/kernels/sync",
			Handler: adminNodes.AdminSyncNodeKernelHandler(svcCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/protocol-configs",
			Handler: adminProtocolConfigs.AdminListProtocolConfigsHandler(svcCtx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/protocol-configs",
			Handler: adminProtocolConfigs.AdminCreateProtocolConfigHandler(svcCtx),
		},
		{
			Method:  http.MethodPatch,
			Path:    "/protocol-configs/:id",
			Handler: adminProtocolConfigs.AdminUpdateProtocolConfigHandler(svcCtx),
		},
		{
			Method:  http.MethodDelete,
			Path:    "/protocol-configs/:id",
			Handler: adminProtocolConfigs.AdminDeleteProtocolConfigHandler(svcCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/protocol-bindings",
			Handler: adminProtocolBindings.AdminListProtocolBindingsHandler(svcCtx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/protocol-bindings",
			Handler: adminProtocolBindings.AdminCreateProtocolBindingHandler(svcCtx),
		},
		{
			Method:  http.MethodPatch,
			Path:    "/protocol-bindings/:id",
			Handler: adminProtocolBindings.AdminUpdateProtocolBindingHandler(svcCtx),
		},
		{
			Method:  http.MethodDelete,
			Path:    "/protocol-bindings/:id",
			Handler: adminProtocolBindings.AdminDeleteProtocolBindingHandler(svcCtx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/protocol-bindings/:id/sync",
			Handler: adminProtocolBindings.AdminSyncProtocolBindingHandler(svcCtx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/protocol-bindings/sync",
			Handler: adminProtocolBindings.AdminSyncProtocolBindingsHandler(svcCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/subscription-templates",
			Handler: adminTemplates.AdminListSubscriptionTemplatesHandler(svcCtx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/subscription-templates",
			Handler: adminTemplates.AdminCreateSubscriptionTemplateHandler(svcCtx),
		},
		{
			Method:  http.MethodPatch,
			Path:    "/subscription-templates/:id",
			Handler: adminTemplates.AdminUpdateSubscriptionTemplateHandler(svcCtx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/subscription-templates/:id/publish",
			Handler: adminTemplates.AdminPublishSubscriptionTemplateHandler(svcCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/subscription-templates/:id/history",
			Handler: adminTemplates.AdminSubscriptionTemplateHistoryHandler(svcCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/subscriptions",
			Handler: adminSubscriptions.AdminListSubscriptionsHandler(svcCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/subscriptions/:id",
			Handler: adminSubscriptions.AdminGetSubscriptionHandler(svcCtx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/subscriptions",
			Handler: adminSubscriptions.AdminCreateSubscriptionHandler(svcCtx),
		},
		{
			Method:  http.MethodPatch,
			Path:    "/subscriptions/:id",
			Handler: adminSubscriptions.AdminUpdateSubscriptionHandler(svcCtx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/subscriptions/:id/disable",
			Handler: adminSubscriptions.AdminDisableSubscriptionHandler(svcCtx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/subscriptions/:id/extend",
			Handler: adminSubscriptions.AdminExtendSubscriptionHandler(svcCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/plans",
			Handler: adminPlans.AdminListPlansHandler(svcCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/plans/:plan_id/billing-options",
			Handler: adminPlanBillingOptions.AdminListPlanBillingOptionsHandler(svcCtx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/plans/:plan_id/billing-options",
			Handler: adminPlanBillingOptions.AdminCreatePlanBillingOptionHandler(svcCtx),
		},
		{
			Method:  http.MethodPatch,
			Path:    "/plans/:plan_id/billing-options/:id",
			Handler: adminPlanBillingOptions.AdminUpdatePlanBillingOptionHandler(svcCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/audit-logs",
			Handler: adminAuditLogs.AdminAuditLogListHandler(svcCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/audit-logs/export",
			Handler: adminAuditLogs.AdminAuditLogExportHandler(svcCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/payment-channels",
			Handler: adminPaymentChannels.AdminListPaymentChannelsHandler(svcCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/payment-channels/:id",
			Handler: adminPaymentChannels.AdminGetPaymentChannelHandler(svcCtx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/payment-channels",
			Handler: adminPaymentChannels.AdminCreatePaymentChannelHandler(svcCtx),
		},
		{
			Method:  http.MethodPatch,
			Path:    "/payment-channels/:id",
			Handler: adminPaymentChannels.AdminUpdatePaymentChannelHandler(svcCtx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/plans",
			Handler: adminPlans.AdminCreatePlanHandler(svcCtx),
		},
		{
			Method:  http.MethodPatch,
			Path:    "/plans/:id",
			Handler: adminPlans.AdminUpdatePlanHandler(svcCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/announcements",
			Handler: adminAnnouncements.AdminListAnnouncementsHandler(svcCtx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/announcements",
			Handler: adminAnnouncements.AdminCreateAnnouncementHandler(svcCtx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/announcements/:id/publish",
			Handler: adminAnnouncements.AdminPublishAnnouncementHandler(svcCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/coupons",
			Handler: adminCoupons.AdminListCouponsHandler(svcCtx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/coupons",
			Handler: adminCoupons.AdminCreateCouponHandler(svcCtx),
		},
		{
			Method:  http.MethodPatch,
			Path:    "/coupons/:id",
			Handler: adminCoupons.AdminUpdateCouponHandler(svcCtx),
		},
		{
			Method:  http.MethodDelete,
			Path:    "/coupons/:id",
			Handler: adminCoupons.AdminDeleteCouponHandler(svcCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/site-settings",
			Handler: adminSite.AdminGetSiteSettingHandler(svcCtx),
		},
		{
			Method:  http.MethodPatch,
			Path:    "/site-settings",
			Handler: adminSite.AdminUpdateSiteSettingHandler(svcCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/security-settings",
			Handler: adminSecurity.AdminGetSecuritySettingHandler(svcCtx),
		},
		{
			Method:  http.MethodPatch,
			Path:    "/security-settings",
			Handler: adminSecurity.AdminUpdateSecuritySettingHandler(svcCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/orders",
			Handler: adminOrders.AdminListOrdersHandler(svcCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/orders/:id",
			Handler: adminOrders.AdminGetOrderHandler(svcCtx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/orders/:id/pay",
			Handler: adminOrders.AdminMarkOrderPaidHandler(svcCtx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/orders/:id/cancel",
			Handler: adminOrders.AdminCancelOrderHandler(svcCtx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/orders/:id/refund",
			Handler: adminOrders.AdminRefundOrderHandler(svcCtx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/orders/payments/reconcile",
			Handler: adminOrders.AdminReconcilePaymentHandler(svcCtx),
		},
	}
	adminRoutes = rest.WithMiddlewares([]rest.Middleware{accessMiddleware.Handler, authMiddleware.RequireRoles("admin")}, adminRoutes...)
	adminPrefix := svcCtx.Config.Admin.RoutePrefix
	adminBase := "/api/v1"
	if adminPrefix != "" {
		adminBase += "/" + adminPrefix
	}
	server.AddRoutes(adminRoutes, rest.WithPrefix(adminBase))

	webhookRoutes := []rest.Route{
		{
			Method:  http.MethodPost,
			Path:    "/orders/payments/callback",
			Handler: adminOrders.AdminPaymentCallbackHandler(svcCtx),
		},
	}
	webhookRoutes = rest.WithMiddlewares([]rest.Middleware{webhookMiddleware.Handler}, webhookRoutes...)
	server.AddRoutes(webhookRoutes, rest.WithPrefix(adminBase))

	publicWebhookRoutes := []rest.Route{
		{
			Method:  http.MethodPost,
			Path:    "/payments/callback",
			Handler: adminOrders.AdminPaymentCallbackPublicHandler(svcCtx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/kernel/traffic",
			Handler: kernelhandlers.KernelTrafficHandler(svcCtx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/kernel/events",
			Handler: kernelhandlers.KernelEventHandler(svcCtx),
		},
	}
	publicWebhookRoutes = rest.WithMiddlewares([]rest.Middleware{webhookMiddleware.Handler}, publicWebhookRoutes...)
	server.AddRoutes(publicWebhookRoutes, rest.WithPrefix("/api/v1"))

	userRoutes := []rest.Route{
		{
			Method:  http.MethodGet,
			Path:    "/subscriptions",
			Handler: userSubscriptions.UserListSubscriptionsHandler(svcCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/subscriptions/:id/preview",
			Handler: userSubscriptions.UserSubscriptionPreviewHandler(svcCtx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/subscriptions/:id/template",
			Handler: userSubscriptions.UserUpdateSubscriptionTemplateHandler(svcCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/subscriptions/:id/traffic",
			Handler: userSubscriptions.UserSubscriptionTrafficHandler(svcCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/plans",
			Handler: userPlans.UserListPlansHandler(svcCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/nodes",
			Handler: userNodes.UserListNodesHandler(svcCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/payment-channels",
			Handler: userPaymentChannels.UserListPaymentChannelsHandler(svcCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/announcements",
			Handler: userAnnouncements.UserListAnnouncementsHandler(svcCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/account/balance",
			Handler: userAccount.UserBalanceHandler(svcCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/account/profile",
			Handler: userAccount.UserProfileHandler(svcCtx),
		},
		{
			Method:  http.MethodPatch,
			Path:    "/account/profile",
			Handler: userAccount.UserUpdateProfileHandler(svcCtx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/account/password",
			Handler: userAccount.UserChangePasswordHandler(svcCtx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/account/credentials/rotate",
			Handler: userAccount.UserRotateCredentialHandler(svcCtx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/account/email/code",
			Handler: userAccount.UserEmailChangeCodeHandler(svcCtx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/account/email",
			Handler: userAccount.UserChangeEmailHandler(svcCtx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/orders",
			Handler: userOrders.UserCreateOrderHandler(svcCtx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/orders/:id/cancel",
			Handler: userOrders.UserCancelOrderHandler(svcCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/orders",
			Handler: userOrders.UserListOrdersHandler(svcCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/orders/:id",
			Handler: userOrders.UserGetOrderHandler(svcCtx),
		},
		{
			Method:  http.MethodGet,
			Path:    "/orders/:id/payment-status",
			Handler: userOrders.UserGetOrderPaymentStatusHandler(svcCtx),
		},
	}
	userRoutes = rest.WithMiddlewares([]rest.Middleware{authMiddleware.RequireRoles("user"), thirdPartyMiddleware.Handler}, userRoutes...)
	server.AddRoutes(userRoutes, rest.WithPrefix("/api/v1/user"))
}
