package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/adarsh/woc1/queue_algo/06_linked_queue/pkg/core"
	"github.com/adarsh/woc1/queue_algo/06_linked_queue/pkg/scheduler"
)

func main() {
	fmt.Println("=== ALGO 6: DYNAMIC LINKED-CHUNK SCHEDULER (INFINITE DENSITY) ===")

	sched := scheduler.NewScheduler()

	// 1. The Density Test
	// We want to simulate the "India Problem":
	// 500,000 phones that are EXACTLY the same (e.g. 4000MB, Region 1, Battery 1)
	// In Algo 5, this would overflow fixed RingBuffer[1024].
	// In Algo 6, it should work perfectly.

	count := 500000
	fmt.Printf("\n--- Phase 1: Registering %d Identical Phones ---\n", count)

	start := time.Now()
	for i := 0; i < count; i++ {
		sched.AddPhone(&core.Phone{
			ID:        fmt.Sprintf("IndPhone_%d", i),
			FreeMemMB: 4000,
			Region:    1,
			Battery:   1,
		})
	}
	fmt.Printf("Registered %d Phones in %s\n", count, time.Since(start))

	// 2. Memory Usage Check
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Heap Alloc: %d MB\n", m.Alloc/1024/1024)

	// 3. Retrieval Test
	// Retrieve all 500,000 phones
	fmt.Println("\n--- Phase 2: Scheduling 500,000 Jobs ---")
	start = time.Now()
	success := 0
	for i := 0; i < count; i++ {
		p := sched.GetBestPhone(3500, 1, 1)
		if p != nil {
			success++
		}
	}
	dur := time.Since(start)

	fmt.Printf("Scheduled %d Jobs in %s (%.2f ns/op)\n", success, dur, float64(dur.Nanoseconds())/float64(success))

	if success == count {
		fmt.Println("[PASS] All phones utilized. No drops.")
	} else {
		fmt.Printf("[FAIL] Dropped %d phones!\n", count-success)
	}
}
