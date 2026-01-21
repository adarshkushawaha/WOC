package scheduler

import (
	"math/bits"

	"github.com/adarsh/woc1/queue_algo/05_tiered_bitmask/pkg/core"
)

// Dimension Constants
const (
	NumRegions   = 4 // 0=Global, 1=US, 2=EU, 3=APAC
	NumBatteries = 3 // 0=Low, 1=Med, 2=High
)

// TieredScheduler manages phones using multi-dimensional bitmasks
type TieredScheduler struct {
	// Masks[Region][Battery][Tier] => uint64
	// Tier 0 = HighRes (0-4GB), Tier 1 = MedRes (4-16GB)
	Masks [NumRegions][NumBatteries][2]uint64

	// Two tiers of RingBuffers: Queues[Tier][Class]
	// Tier 0 has 64 classes, Tier 1 has 64 classes
	Queues [2][64]*core.LockFreeQueue
}

func NewScheduler() *TieredScheduler {
	s := &TieredScheduler{}

	// Initialize Queues
	// Capacity: 1024 phones per bucket (Power of 2 for RingBuffer)
	queueCap := uint64(1024)

	for tier := 0; tier < 2; tier++ {
		for cls := 0; cls < 64; cls++ {
			s.Queues[tier][cls] = core.NewLockFreeQueue(queueCap)
		}
	}
	return s
}

func (s *TieredScheduler) AddPhone(p *core.Phone) {
	// 1. Map RAM to Class & Tier
	absClass := core.MapSizeToClass(p.FreeMemMB)
	tier := 0
	relClass := absClass
	if absClass >= 64 {
		tier = 1
		relClass = absClass - 64
	}

	// 2. Add to Mask (Atomic logic would go here, simplified using simple OR)
	// We update the mask for specific Region and Battery
	s.Masks[p.Region][p.Battery][tier] |= (1 << relClass)

	// Also update Global Region (0)
	s.Masks[0][p.Battery][tier] |= (1 << relClass)

	// 3. Add to Lock-Free Queue
	ok := s.Queues[tier][relClass].Enqueue(p)
	if !ok {
		// fmt.Printf("Warning: Queue full for class %d\n", absClass)
	}
}

// GetBestPhone attempts to find a phone.
// It tries Exact Match first, then "Smart Borrows" by relaxing Battery constraints.
func (s *TieredScheduler) GetBestPhone(neededMB, region, minBattery int) *core.Phone {
	// 1. Try Exact Criteria
	p := s.findInTiers(neededMB, region, minBattery)
	if p != nil {
		return p
	}

	// 2. Smart Borrowing (Relax Battery)
	// If we asked for High Battery (2) and failed, try Med (1)...
	for b := minBattery - 1; b >= 0; b-- {
		p := s.findInTiers(neededMB, region, b)
		if p != nil {
			// fmt.Printf(" [Borrowing] Found lower battery phone (Level %d)\n", b)
			return p
		}
	}

	return nil
}

// findInTiers checks Tier 0 then Tier 1
func (s *TieredScheduler) findInTiers(neededMB, region, battery int) *core.Phone {
	absClass := core.MapSizeToClass(neededMB)

	// Start checking from the Tier where neededMB falls
	startTier := 0
	if absClass >= 64 {
		startTier = 1
	}

	for t := startTier; t < 2; t++ {
		mask := s.Masks[region][battery][t]

		// If we are in the starting tier, we need to mask out bits below us
		targetMask := uint64(0)
		if t == startTier {
			relClass := absClass
			if t == 1 {
				relClass -= 64
			}
			// Create mask of all classes >= relClass
			targetMask = ^(uint64(1)<<relClass - 1)
		} else {
			// If we moved up a tier, we can take ANY class (it's all bigger than needed)
			targetMask = ^uint64(0)
		}

		validOptions := mask & targetMask

		if validOptions != 0 {
			// Found a candidate bucket!
			bestClass := bits.TrailingZeros64(validOptions)

			// Pop from queue
			p, ok := s.Queues[t][bestClass].Dequeue()
			if ok {
				return p
			} else {
				// Race: Mask said yes, Queue said empty.
				// Clear bit
				s.Masks[region][battery][t] &^= (1 << bestClass)
				// Retry this tier? For Speed, we just continue loop or fail this tick.
			}
		}
	}
	return nil
}
