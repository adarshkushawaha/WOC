package main

import (
	"fmt"
	"time"

	"github.com/adarsh/woc1/queue_algo/01_priority_queue/pkg/manager"
	"github.com/adarsh/woc1/queue_algo/01_priority_queue/pkg/model"
)

func main() {
	fmt.Println("Starting Job-Worker Queue Simulation...")

	disp := manager.NewDispatcher()

	// 1. Simulate finding some Jobs
	jobs := []*model.Job{
		{ID: "Job1", ArrivalTime: time.Now().Add(-10 * time.Minute), Size: 5}, // Late arrival, small size
		{ID: "Job2", ArrivalTime: time.Now().Add(-5 * time.Minute), Size: 50}, // Mid arrival, big size
		{ID: "Job3", ArrivalTime: time.Now().Add(-2 * time.Minute), Size: 10}, // Recent arrival, small size
		{ID: "Job4", ArrivalTime: time.Now().Add(-60 * time.Minute), Size: 2}, // Very old arrival! Should have high wait score
	}

	for _, j := range jobs {
		disp.AddJob(j)
	}

	// 2. Simulate workers becoming available
	workers := []*model.Worker{
		{ID: "WorkerA", AvailableTime: time.Now().Add(-1 * time.Minute), Efficiency: 1.0},
		{ID: "WorkerB", AvailableTime: time.Now().Add(-30 * time.Minute), Efficiency: 0.8}, // Very idle, low efficiency
		{ID: "WorkerC", AvailableTime: time.Now(), Efficiency: 2.0},                        // Just now, high efficiency
	}

	for _, w := range workers {
		disp.AddWorker(w)
	}

	disp.PrintStatus()

	// 3. Perform Matches
	// We expect Job4 (Long wait) to be matched first if weight is high enough.
	// We expect WorkerB (Long idle) or WorkerC (High efficiency) depending on formula.

	fmt.Println(">>> Round 1 Matching")
	disp.Match()
	disp.PrintStatus()

	fmt.Println(">>> Round 2 Matching")
	disp.Match()
	disp.PrintStatus()

	fmt.Println(">>> Round 3 Matching")
	disp.Match()
	disp.PrintStatus()

	fmt.Println("Simulation Complete.")
}
