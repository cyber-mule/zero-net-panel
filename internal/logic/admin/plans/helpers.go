package plans

import "strings"

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
