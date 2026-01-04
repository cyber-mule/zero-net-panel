package subscriptionutil

import (
	"context"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
)

// LoadSubscriptionEntries returns protocol entries assigned to a subscription snapshot.
func LoadSubscriptionEntries(ctx context.Context, repos *repository.Repositories, sub repository.Subscription) ([]repository.ProtocolEntry, error) {
	if sub.ID == 0 {
		return nil, repository.ErrInvalidArgument
	}

	ids := ExtractBindingIDs(sub.PlanSnapshot)
	if len(ids) == 0 {
		var err error
		ids, err = repos.PlanProtocolBinding.ListBindingIDs(ctx, sub.PlanID)
		if err != nil {
			return nil, err
		}
	}
	if len(ids) == 0 {
		return []repository.ProtocolEntry{}, nil
	}

	unique := uniqueBindingIDs(ids)
	return repos.ProtocolEntry.ListByBindingIDs(ctx, unique)
}
