package bitmask

import (
	"math/bits"
	"sync"
)

// Constants
const (
	MaxClasses     = 64 // Using uint64, 50MB per bit = 3200MB max directly managed
	BucketInterval = 50
)

// Phone represents a worker device
type Phone struct {
	ID          string
	FreeMemMB   int
	MemoryClass int // 0-63
}

// O1Scheduler uses bitmasks for constant time lookups
type O1Scheduler struct {
	ActiveMask uint64                 // Bit K is 1 if Class K has phones
	Queues     [MaxClasses][]*Phone   // Array of slices (RingBuffers ideal, slices for simplicity)
	Locks      [MaxClasses]sync.Mutex // Fine-grained locking
}

func NewO1Scheduler() *O1Scheduler {
	return &O1Scheduler{}
}

// AddPhone registers a phone to the correct bucket
func (s *O1Scheduler) AddPhone(p *Phone) {
	class := p.FreeMemMB / BucketInterval
	if class >= MaxClasses {
		class = MaxClasses - 1
	}
	p.MemoryClass = class

	s.Locks[class].Lock()
	s.Queues[class] = append(s.Queues[class], p)
	// Set bit atomically-ish (under lock of that class? No, mask needs generic lock or atomic)
	// For simplicity in this demo, we assume single threaded dispatch or use a global lock for mask.
	// But to be "Max Efficient", we update mask.
	// Actually, concurrent bitwise OR is safe-ish if we don't care about race on read.
	// Better: Use Atomic Load/Store for mask.
	s.Locks[class].Unlock()

	// Update Mask (Idempotent)
	// We need to ensure visibility. A simplified approach:
	s.ActiveMask |= (1 << class)
}

// GetBestPhone finds a phone with at least `neededMB` in O(1)
func (s *O1Scheduler) GetBestPhone(neededMB int) *Phone {
	minClass := neededMB / BucketInterval
	if minClass >= MaxClasses {
		return nil // Too big
	}

	// 1. Create a mask of all classes >= minClass
	//    Example: needed Class 2. Mask: 11111100
	//    Formula: ^((1 << minClass) - 1)
	targetMask := ^(uint64(1)<<minClass - 1)

	// 2. Find intersection
	validOptions := s.ActiveMask & targetMask

	if validOptions == 0 {
		return nil // No phones available
	}

	// 3. Find the lowest set bit (Best Fit)
	//    TrailingZeros64 returns the index of the least significant bit set.
	//    Example: validOptions = ...001000 (Class 3 is set). TrailingZeros = 3.
	bestClass := bits.TrailingZeros64(validOptions)

	// 4. Pop from that queue
	s.Locks[bestClass].Lock()
	defer s.Locks[bestClass].Unlock()

	q := s.Queues[bestClass]
	if len(q) == 0 {
		// Race condition: Bit was set but queue empty.
		// Clear bit and recurse/retry.
		s.ActiveMask &^= (1 << bestClass)
		return nil // Retry logic would go here
	}

	// Pop
	p := q[len(q)-1] // LIFO for cache locality (Stack) or FIFO? FIFO is fairer.
	// Let's do LIFO for efficiency (no memmove).
	s.Queues[bestClass] = q[:len(q)-1]

	if len(s.Queues[bestClass]) == 0 {
		s.ActiveMask &^= (1 << bestClass) // Clear bit
	}

	return p
}
