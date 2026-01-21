package core

// Reuse Phone struct
type Phone struct {
	ID        string
	FreeMemMB int
	Region    int
	Battery   int
}

// Reuse Tiered Mapping Logic
const (
	Tier1LimitMB  = 4096 // 4GB
	Tier1Interval = 64   // 64MB steps

	Tier2LimitMB  = 16384 // 16GB
	Tier2Interval = 192   // ~192MB steps
)

func MapSizeToClass(mb int) int {
	if mb <= 0 {
		return 0
	}

	if mb <= Tier1LimitMB {
		idx := mb / Tier1Interval
		if idx >= 64 {
			idx = 63
		}
		return idx
	} else {
		remaining := mb - Tier1LimitMB
		idx := remaining / Tier2Interval
		if idx >= 64 {
			idx = 63
		}
		return 64 + idx
	}
}
