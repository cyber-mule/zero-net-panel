package repository

import "encoding/json"

func serializeStringSlice(values []string) (string, error) {
	if values == nil {
		values = []string{}
	}
	payload, err := json.Marshal(values)
	if err != nil {
		return "", err
	}
	return string(payload), nil
}

func serializeAnyMap(values map[string]any) (string, error) {
	if values == nil {
		values = map[string]any{}
	}
	payload, err := json.Marshal(values)
	if err != nil {
		return "", err
	}
	return string(payload), nil
}
