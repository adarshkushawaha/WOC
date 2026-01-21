package core

// Constants for Bucket Mapping
const (
	Tier1LimitMB  = 4096 // 4GB
	Tier1Interval = 64   // 64MB steps ( 4096 / 64 = 64 buckets )

	Tier2LimitMB  = 16384 // 16GB
	Tier2Interval = 192   // ~192MB steps ( (16384-4096)/192 ~= 64 buckets )

	TotalClasses = 128 // 64 + 64
)

func MapSizeToClass(mb int) int {
	if mb <= 0 {
		return 0
	}

	if mb <= Tier1LimitMB {
		// Tier 1: 0 - 4096 MB
		idx := mb / Tier1Interval
		if idx >= 64 {
			idx = 63
		}
		return idx
	} else {
		// Tier 2: 4096 - 16384 MB
		remaining := mb - Tier1LimitMB
		idx := remaining / Tier2Interval
		if idx >= 64 {
			idx = 63
		}
		return 64 + idx // Offest by 64
	}
}

// MapClassToMinSize returns the minimum MB for a given class
func MapClassToMinSize(class int) int {
	if class < 64 {
		return class * Tier1Interval
	}
	// Tier 2
	idx := class - 64
	return Tier1LimitMB + (idx * Tier2Interval)
}
