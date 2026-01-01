package plans

import (
	"context"
	"strings"

	"github.com/zero-net-panel/zero-net-panel/internal/repository"
)

func normalizeTrafficMultipliers(input map[string]float64) map[string]float64 {
	if input == nil {
		return map[string]float64{}
	}
	result := make(map[string]float64, len(input))
	for key, value := range input {
		key = strings.ToLower(strings.TrimSpace(key))
		if key == "" || value <= 0 {
			continue
		}
		result[key] = value
	}
	return result
}

func normalizeBindingIDs(input []uint64) ([]uint64, error) {
	if len(input) == 0 {
		return []uint64{}, nil
	}
	seen := make(map[uint64]struct{}, len(input))
	result := make([]uint64, 0, len(input))
	for _, id := range input {
		if id == 0 {
			return nil, repository.ErrInvalidArgument
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		result = append(result, id)
	}
	return result, nil
}

func ensureBindingsExist(ctx context.Context, repo repository.ProtocolBindingRepository, bindingIDs []uint64) error {
	if len(bindingIDs) == 0 {
		return nil
	}
	bindings, err := repo.ListByIDs(ctx, bindingIDs)
	if err != nil {
		return err
	}
	if len(bindings) != len(bindingIDs) {
		return repository.ErrInvalidArgument
	}
	return nil
}
