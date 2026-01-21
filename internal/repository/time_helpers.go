package repository

import "time"

var zeroTimeSentinel = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

// ZeroTime returns a non-zero sentinel used to represent unset timestamps.
func ZeroTime() time.Time {
	return zeroTimeSentinel
}

// IsZeroTime treats the sentinel and zero value as unset.
func IsZeroTime(ts time.Time) bool {
	return ts.IsZero() || ts.Equal(zeroTimeSentinel)
}

// NormalizeTime replaces zero values with the sentinel to avoid invalid dates.
func NormalizeTime(ts time.Time) time.Time {
	if ts.IsZero() {
		return zeroTimeSentinel
	}
	return ts
}
