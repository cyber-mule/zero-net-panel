package subscriptionutil

import (
	"context"
	"sort"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
)

// LoadPlanBindings returns protocol bindings assigned to a plan.
func LoadPlanBindings(ctx context.Context, repos *repository.Repositories, planID uint64) ([]repository.ProtocolBinding, error) {
	if planID == 0 {
		return nil, repository.ErrInvalidArgument
	}

	ids, err := repos.PlanProtocolBinding.ListBindingIDs(ctx, planID)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return []repository.ProtocolBinding{}, nil
	}

	unique := uniqueBindingIDs(ids)
	bindings, err := repos.ProtocolBinding.ListByIDs(ctx, unique)
	if err != nil {
		return nil, err
	}

	sort.Slice(bindings, func(i, j int) bool {
		return bindings[i].ID < bindings[j].ID
	})
	return bindings, nil
}

// LoadSubscriptionBindings returns protocol bindings assigned to a subscription snapshot.
func LoadSubscriptionBindings(ctx context.Context, repos *repository.Repositories, sub repository.Subscription) ([]repository.ProtocolBinding, error) {
	if sub.ID == 0 {
		return nil, repository.ErrInvalidArgument
	}

	ids := ExtractBindingIDs(sub.PlanSnapshot)
	if len(ids) == 0 {
		return LoadPlanBindings(ctx, repos, sub.PlanID)
	}

	unique := uniqueBindingIDs(ids)
	bindings, err := repos.ProtocolBinding.ListByIDs(ctx, unique)
	if err != nil {
		return nil, err
	}

	sort.Slice(bindings, func(i, j int) bool {
		return bindings[i].ID < bindings[j].ID
	})
	return bindings, nil
}

func uniqueBindingIDs(input []uint64) []uint64 {
	seen := make(map[uint64]struct{}, len(input))
	result := make([]uint64, 0, len(input))
	for _, id := range input {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		result = append(result, id)
	}
	return result
}
