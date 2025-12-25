package paymentchannels

import (
	"github.com/zero-net-panel/zero-net-panel/internal/repository"
	"github.com/zero-net-panel/zero-net-panel/internal/types"
)

func toPaymentChannelSummary(channel repository.PaymentChannel) types.PaymentChannelSummary {
	return types.PaymentChannelSummary{
		ID:        channel.ID,
		Name:      channel.Name,
		Code:      channel.Code,
		Provider:  channel.Provider,
		Enabled:   channel.Enabled,
		SortOrder: channel.SortOrder,
		Config:    channel.Config,
		CreatedAt: channel.CreatedAt.Unix(),
		UpdatedAt: channel.UpdatedAt.Unix(),
	}
}
