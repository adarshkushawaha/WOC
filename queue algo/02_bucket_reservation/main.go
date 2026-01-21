package main

import (
	"fmt"

	"github.com/adarsh/woc1/queue_algo/02_bucket_reservation/pkg/manager"
	"github.com/adarsh/woc1/queue_algo/02_bucket_reservation/pkg/model"
)

func main() {
	fmt.Println("--- starting job-worker bucket simulation (algo 2) ---")

	// system max mem 2000 mb
	disp := manager.NewDispatcher(2000)

	// 1. add small jobs (under 200 mb)
	for i := 1; i <= 5; i++ {
		disp.AddJob(&model.Job{ID: fmt.Sprintf("SmallJob_%d", i), SizeMB: 100})
	}

	// 2. add workers
	// worker a: 2000mb capacity (Very Capable)
	// worker b: 200mb capacity (Small Capable)
	disp.AddWorker(&model.Worker{ID: "HeavyWorker", CapacityMB: 2000})
	disp.AddWorker(&model.Worker{ID: "SmallWorker", CapacityMB: 200})

	fmt.Println("\n>>> Round 1: Normal Mode (No Avalanche)")
	// HeavyWorker should pick a SmallJob because no Heavy jobs exist, so why be idle?
	disp.Match()

	fmt.Println("\n>>> Triggering Avalanche: Adding 5 Heavy Jobs (>800MB)")
	for i := 1; i <= 5; i++ {
		disp.AddJob(&model.Job{ID: fmt.Sprintf("HeavyJob_%d", i), SizeMB: 1500})
	}

	fmt.Println("\n>>> Round 2: Avalanche Detection")
	// The next Match calls should detect avalanche.
	// We expect HeavyWorker to switching to RESERVED mode.
	// If HeavyWorker is free, it should pick HeavyJob ONLY.
	// If HeavyWorker was busy in Round 1 (simulated), it would finish and then check again.
	// Here, we simulate new match round.
	disp.Match()

	fmt.Println("\n>>> Round 3: Continued Matching")
	disp.Match()

	// Add a new small job to see if HeavyWorker ignores it
	disp.AddJob(&model.Job{ID: "SmallJob_New", SizeMB: 50})
	fmt.Println("\n>>> Round 4: Temptation Test")
	// HeavyWorker should distinct HeavyJob_2 over SmallJob_New
	disp.Match()
}
