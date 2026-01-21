package manager

import (
	"fmt"

	"github.com/adarsh/woc1/queue_algo/03_segment_tree/pkg/model"
)

const (
	AvalancheThreshold = 3
	JobLargeThreshold  = 800
)

type Dispatcher struct {
	Tree            *model.SegmentTree
	Workers         []*model.Worker
	AvalancheActive bool
}

func NewDispatcher(maxMem int) *Dispatcher {
	return &Dispatcher{
		Tree:    model.NewSegmentTree(maxMem),
		Workers: make([]*model.Worker, 0),
	}
}

func (d *Dispatcher) AddJob(j *model.Job) {
	d.Tree.AddJob(j)
	fmt.Printf("[Dispatcher] Job Added: %s\n", j)
}

func (d *Dispatcher) AddWorker(w *model.Worker) {
	d.Workers = append(d.Workers, w)
}

func (d *Dispatcher) CheckAvalancheStatus() {
	// Query tree for total jobs in the "Large" range
	heavyCount := d.Tree.TotalJobsInRange(JobLargeThreshold, d.Tree.MaxSize)

	if heavyCount >= AvalancheThreshold {
		if !d.AvalancheActive {
			fmt.Println("\n!!! AVALANCHE DETECTED !!! Tree switching to Reservation Mode.")
			d.AvalancheActive = true
		}
	} else {
		if d.AvalancheActive {
			fmt.Println("\n... Avalanche cleared. Returning to normal tree search.")
			d.AvalancheActive = false
		}
	}
}

func (d *Dispatcher) Match() {
	d.CheckAvalancheStatus()

	for _, w := range d.Workers {
		// Keep trying to fill the worker until no more jobs fit or memory is full
		for {
			avail := w.AvailableMemory()
			if avail <= 0 {
				break
			}

			minSize := 0
			// If Avalanche is active and this worker is a "Heavy" resource,
			// we reserve its FIRST slot for a heavy job.
			// Once it has at least one heavy job (or if avalanche is off),
			// it can fill its remaining space with anything.
			if d.AvalancheActive && w.CapacityMB >= JobLargeThreshold && len(w.CurrentJobs) == 0 {
				minSize = JobLargeThreshold
			}

			// Optimized O(log B) search
			job := d.Tree.FindHeaviest(avail, minSize)

			if job != nil {
				w.UsedMemory += job.SizeMB
				w.CurrentJobs = append(w.CurrentJobs, job)
				fmt.Printf("[BIN-PACKING] Worker %s -> Added %s (Remaining: %dMB)\n", w.ID, job, w.AvailableMemory())
			} else {
				// No more jobs fit in this worker's remaining space
				break
			}
		}
	}
}
