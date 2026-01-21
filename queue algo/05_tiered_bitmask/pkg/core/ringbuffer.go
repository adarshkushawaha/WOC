package core

import (
	"sync/atomic"
)

type Phone struct {
	ID        string
	FreeMemMB int
	Region    int
	Battery   int
}

// Node wraps the item for the ring buffer
type Node struct {
	Value *Phone
	Seq   uint64
}

// LockFreeQueue is a RingBuffer implementation using CAS
type LockFreeQueue struct {
	buffer []Node
	mask   uint64
	head   uint64
	tail   uint64
}

func NewLockFreeQueue(size uint64) *LockFreeQueue {
	// Size must be power of 2
	if size&(size-1) != 0 {
		panic("size must be power of 2")
	}

	q := &LockFreeQueue{
		buffer: make([]Node, size),
		mask:   size - 1,
	}

	for i := range q.buffer {
		q.buffer[i].Seq = uint64(i)
	}

	return q
}

// Enqueue adds an item to the buffer
func (q *LockFreeQueue) Enqueue(p *Phone) bool {
	for {
		tail := atomic.LoadUint64(&q.tail)
		idx := tail & q.mask
		node := &q.buffer[idx]
		seq := atomic.LoadUint64(&node.Seq)

		dif := int64(seq) - int64(tail)
		if dif == 0 {
			if atomic.CompareAndSwapUint64(&q.tail, tail, tail+1) {
				node.Value = p
				atomic.StoreUint64(&node.Seq, tail+1)
				return true
			}
		} else if dif < 0 {
			// Buffer full
			return false
		} else {
			// Tail lagging, help move it
			atomic.CompareAndSwapUint64(&q.tail, tail, tail+1)
		}
	}
}

// Dequeue removes an item
func (q *LockFreeQueue) Dequeue() (*Phone, bool) {
	for {
		head := atomic.LoadUint64(&q.head)
		idx := head & q.mask
		node := &q.buffer[idx]
		seq := atomic.LoadUint64(&node.Seq)

		dif := int64(seq) - int64(head+1)
		if dif == 0 {
			if atomic.CompareAndSwapUint64(&q.head, head, head+1) {
				val := node.Value
				// atomic.StoreUint64(&node.Seq, head + q.mask + 1) // Correct for Ring wrap
				// Actually for single-producer single-consumer Seq logic is complex.
				// For simplified MPMC, we just need to mark it safe to write.
				node.Value = nil // GC safety

				// Reset sequence to allow overwriting in next cycle
				// Formula: Current Seq + Size
				atomic.StoreUint64(&node.Seq, head+q.mask+1)
				return val, true
			}
		} else if dif < 0 {
			// Buffer empty
			return nil, false
		} else {
			// Head lagging
			atomic.CompareAndSwapUint64(&q.head, head, head+1)
		}
	}
}

func (q *LockFreeQueue) Len() int {
	// Approximate length
	h := atomic.LoadUint64(&q.head)
	t := atomic.LoadUint64(&q.tail)
	if t < h {
		return 0
	}
	return int(t - h)
}
