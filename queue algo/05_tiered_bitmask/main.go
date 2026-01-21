package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/adarsh/woc1/queue_algo/05_tiered_bitmask/pkg/core"
	"github.com/adarsh/woc1/queue_algo/05_tiered_bitmask/pkg/scheduler"
)

func main() {
	fmt.Println("=== ALGO 5: ULTRA SCALAR TIERED SCHEDULER (16GB) ===")

	sched := scheduler.NewScheduler()

	// 1. Populate with Diverse Phones
	// We add 10,000 phones with varied specs
	// Region: 1=US, 2=EU
	// Battery: 0=Low, 2=High
	// Memory: 50MB to 16000MB
	start := time.Now()
	for i := 0; i < 10000; i++ {
		mem := rand.Intn(16000) + 50 // Random up to 16GB
		region := rand.Intn(3) + 1   // 1-3
		batt := rand.Intn(3)         // 0-2

		sched.AddPhone(&core.Phone{
			ID:        fmt.Sprintf("P-%d", i),
			FreeMemMB: mem,
			Region:    region,
			Battery:   batt,
		})
	}
	fmt.Printf("Registered 10,000 Phones (Dynamic 16GB Tiering) in %s\n", time.Since(start))

	// 2. Test Case A: Heavy 16GB Job
	// Needs 12000 MB, Region 1 (US), High Battery
	fmt.Println("\n--- Test A: Scheduling 12GB AI Model Task ---")
	p := sched.GetBestPhone(12000, 1, 2)
	if p != nil {
		fmt.Printf("[SUCCESS] Assigned to %s (Mem: %dMB, Reg: %d, Bat: %d)\n", p.ID, p.FreeMemMB, p.Region, p.Battery)
	} else {
		fmt.Println("[FAIL] No phone found")
	}

	// 3. Test Case B: Smart Borrowing
	// Needs 500 MB, Region 2 (EU), High Battery (2)
	// IF no High Battery exists, it should find Med Battery automatically
	fmt.Println("\n--- Test B: Borrowing Logic (Fallback) ---")
	// Force a borrow scenario: Add a specific EU phone with LOW battery
	sched.AddPhone(&core.Phone{ID: "LowBatEU", FreeMemMB: 1000, Region: 2, Battery: 0})

	// Ask for High Battery
	p2 := sched.GetBestPhone(500, 2, 2)
	if p2 != nil {
		fmt.Printf("[SUCCESS] Asked for Battery 2, Got Battery %d -> %s\n", p2.Battery, p2.ID)
	} else {
		fmt.Println("[FAIL] Borrowing logic failed")
	}

	// 4. Concurrency Test
	fmt.Println("\n--- Test C: High Concurrency (Lock-Free) ---")
	// Launch 100 goroutines trying to schedule simultaneously
	done := make(chan bool)
	ops := 100000
	start = time.Now()
	for i := 0; i < 10; i++ { // 10 routines doing 10k ops each
		go func() {
			for j := 0; j < 10000; j++ {
				sched.GetBestPhone(100, 1, 0)
				// Re-add occasionally to prevent empty
				if j%2 == 0 {
					sched.AddPhone(&core.Phone{FreeMemMB: 100, Region: 1, Battery: 1})
				}
			}
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
	dur := time.Since(start)
	fmt.Printf("Processed %d Ops in %s (%.2f ns/op)\n", ops, dur, float64(dur.Nanoseconds())/float64(ops))
}
