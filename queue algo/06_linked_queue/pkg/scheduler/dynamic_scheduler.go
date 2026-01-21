package scheduler

import (
	"math/bits"
	"sync"

	"github.com/adarsh/woc1/queue_algo/06_linked_queue/pkg/core"
)

// Reuse Dimension Constants
const (
	NumRegions   = 4
	NumBatteries = 3
)

// DynamicScheduler uses LinkedQueue for infinite density
type DynamicScheduler struct {
	Masks [NumRegions][NumBatteries][2]uint64

	// Dynamic Map instead of fixed array?
	// To minimize memory, we can use a map[int]*Queue or a fixed array of pointers where nil = not initialized.
	// Since 64 pointers is small (512 bytes), fixed array is better for O(1) access than map.
	// We just Lazily Allocate the *LinkedQueue inside.

	Queues [2][64]*core.LinkedQueue

	// Mutex for Lazy Allocation (Only hit once per class)
	AllocLock sync.RWMutex
}

func NewScheduler() *DynamicScheduler {
	return &DynamicScheduler{}
}

func (s *DynamicScheduler) getQueue(tier, cls int) *core.LinkedQueue {
	// Fast path: Check matches
	q := s.Queues[tier][cls]
	if q != nil {
		return q
	}

	// Slow path: Allocate
	s.AllocLock.Lock()
	defer s.AllocLock.Unlock()

	// Double check
	if s.Queues[tier][cls] == nil {
		s.Queues[tier][cls] = core.NewLinkedQueue()
	}
	return s.Queues[tier][cls]
}

func (s *DynamicScheduler) AddPhone(p *core.Phone) {
	absClass := core.MapSizeToClass(p.FreeMemMB)
	tier := 0
	relClass := absClass
	if absClass >= 64 {
		tier = 1
		relClass = absClass - 64
	}

	// 1. Update Bitmasks (Atomic OR or simplified)
	s.Masks[p.Region][p.Battery][tier] |= (1 << relClass)
	s.Masks[0][p.Battery][tier] |= (1 << relClass)

	// 2. Get (or Create) Queue & Enqueue
	q := s.getQueue(tier, relClass)
	q.Enqueue(p)
}

func (s *DynamicScheduler) GetBestPhone(neededMB, region, battery int) *core.Phone {
	// Try Exact Battery first
	p := s.findInTiers(neededMB, region, battery)
	if p != nil {
		return p
	}

	// Fallback/Borrow logic
	for b := battery - 1; b >= 0; b-- {
		p := s.findInTiers(neededMB, region, b)
		if p != nil {
			return p
		}
	}
	return nil
}

func (s *DynamicScheduler) findInTiers(neededMB, region, battery int) *core.Phone {
	absClass := core.MapSizeToClass(neededMB)
	startTier := 0
	if absClass >= 64 {
		startTier = 1
	}

	for t := startTier; t < 2; t++ {
		mask := s.Masks[region][battery][t]

		targetMask := uint64(0)
		if t == startTier {
			relClass := absClass
			if t == 1 {
				relClass -= 64
			}
			targetMask = ^(uint64(1)<<relClass - 1)
		} else {
			targetMask = ^uint64(0)
		}

		validOptions := mask & targetMask

		// Loop through options until we find a phone or run out
		for validOptions != 0 {
			bestClass := bits.TrailingZeros64(validOptions)

			q := s.Queues[t][bestClass]
			if q != nil {
				p := q.Dequeue()
				if p != nil {
					return p
				}
			}

			// Queue empty or nil? Clear bit and continue searching this mask
			validOptions &^= (1 << bestClass)
			// Also update global mask? Ideally yes.
		}
	}
	return nil
}
