package main

import (
	"fmt"
	"time"

	"github.com/adarsh/woc1/queue_algo/04_bitmask_basic/pkg/bitmask"
	"github.com/adarsh/woc1/queue_algo/04_bitmask_basic/pkg/manager"
)

func main() {
	fmt.Println("=== ALGO 4: ULTIMATE SERVER SCHEDULER ===")

	disp := manager.NewDispatcher()

	// 1. Register a BluePrint
	disp.SeqEngine.RegisterFlow("ImageProcess", []string{"Download", "Resize", "Upload"})

	// 2. Register Phones (Workers)
	// We add 1000 phones to simulate scale
	start := time.Now()
	for i := 0; i < 1000; i++ {
		// Phones have random memory between 50MB and 2000MB
		mem := (i%20)*100 + 50
		disp.Scheduler.AddPhone(&bitmask.Phone{
			ID:        fmt.Sprintf("Phone_%d", i),
			FreeMemMB: mem,
		})
	}
	fmt.Printf("Registered 1000 Phones in %s\n", time.Since(start))

	// 3. Start a Flow
	fmt.Println("\n--- Starting Flow A ---")
	disp.StartFlow("ImageProcess", "FlowA")

	// 4. Simulate Async Completions
	// FlowA: Download -> Done -> Resize -> Done -> Upload
	disp.AdvanceFlow("FlowA") // Finish Download, Start Resize
	disp.AdvanceFlow("FlowA") // Finish Resize, Start Upload
	disp.AdvanceFlow("FlowA") // Finish Upload, Flow Complete

	// 5. O(1) Performance Test
	fmt.Println("\n--- Performance Check ---")
	// Try to schedule 10,000 jobs instantly
	start = time.Now()
	success := 0
	for i := 0; i < 10000; i++ {
		// Direct scheduler access to measure raw throughput
		// Ask for 500MB phone
		p := disp.Scheduler.GetBestPhone(500)
		if p != nil {
			success++
			// In real code, we'd decrement memory here.
			// Re-add to simulate "Returning to pool"
			disp.Scheduler.AddPhone(p)
		}
	}
	dur := time.Since(start)
	fmt.Printf("Scheduled 10,000 Jobs in %s (%.2f ns/op)\n", dur, float64(dur.Nanoseconds())/10000.0)
}
