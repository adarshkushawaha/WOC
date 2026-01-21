package main

import (
	"fmt"

	"github.com/adarsh/woc1/queue_algo/03_segment_tree/pkg/manager"
	"github.com/adarsh/woc1/queue_algo/03_segment_tree/pkg/model"
)

func main() {
	fmt.Println("=== ALGO 3: BIN-PACKING SEGMENT TREE DISPATCHER ===")

	disp := manager.NewDispatcher(4096)

	// 1. Add many small jobs
	fmt.Println("\n--- Phase 1: Bin-Packing Demonstration ---")
	disp.AddJob(&model.Job{ID: "Small_1", SizeMB: 100})
	disp.AddJob(&model.Job{ID: "Small_2", SizeMB: 100})
	disp.AddJob(&model.Job{ID: "Small_3", SizeMB: 100})
	disp.AddJob(&model.Job{ID: "Small_4", SizeMB: 100})

	// One Big Worker (4000MB)
	disp.AddWorker(&model.Worker{ID: "GiantWorker", CapacityMB: 4000})

	// Match: GiantWorker should take ALL 4 small jobs because it has space
	disp.Match()

	// 2. Avalanche Phase with Bin-Packing
	fmt.Println("\n--- Phase 2: Heavy Avalanche (Reservation + Packing) ---")
	disp.AddJob(&model.Job{ID: "Heavy_1", SizeMB: 1200})
	disp.AddJob(&model.Job{ID: "Heavy_2", SizeMB: 1200})
	disp.AddJob(&model.Job{ID: "Heavy_3", SizeMB: 1200})
	disp.AddJob(&model.Job{ID: "Temptation_1", SizeMB: 50})
	disp.AddJob(&model.Job{ID: "Temptation_2", SizeMB: 50})

	// GiantWorker (now empty from previous phase if we reset or just new worker)
	// Let's assume a fresh match round
	disp.Match()

	fmt.Println("\nExecution Complete.")
}
