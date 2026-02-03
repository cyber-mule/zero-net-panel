package templates

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zero-net-panel/zero-net-panel/internal/svc"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
	subscriptionclient "github.com/zero-net-panel/zero-net-panel/pkg/subscription/client"
)

// ClientListLogic returns supported subscription template clients.
type ClientListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewClientListLogic constructs ClientListLogic.
func NewClientListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClientListLogic {
	return &ClientListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// List returns client types supported by subscription templates.
func (l *ClientListLogic) List() (*types.AdminSubscriptionTemplateClientListResponse, error) {
	rules := subscriptionclient.Rules()
	clients := make([]types.SubscriptionTemplateClient, 0, len(rules))
	for _, rule := range rules {
		clients = append(clients, types.SubscriptionTemplateClient{
			ClientType:      rule.Type,
			DisplayName:     rule.Label,
			UserAgentTokens: append([]string(nil), rule.Tokens...),
			Source:          rule.Source,
		})
	}

	return &types.AdminSubscriptionTemplateClientListResponse{
		Clients: clients,
	}, nil
}
