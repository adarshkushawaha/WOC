package manager

import (
	"fmt"

	"github.com/adarsh/woc1/queue_algo/02_bucket_reservation/pkg/model"
)

// Constants for Reservation Logic
const (
	AvalancheThreshold = 3   // If > 3 jobs are waiting in Large buckets, trigger avalanche
	JobLargeThreshold  = 800 // Jobs > 800MB are considered "Large" (Heavy)
)

type Dispatcher struct {
	Buckets         *model.BucketManager
	Workers         []*model.Worker
	AvalancheActive bool
}

func NewDispatcher(maxMem int) *Dispatcher {
	return &Dispatcher{
		Buckets: model.NewBucketManager(maxMem),
		Workers: make([]*model.Worker, 0),
	}
}

func (d *Dispatcher) AddJob(j *model.Job) {
	d.Buckets.AddJob(j)
}

func (d *Dispatcher) AddWorker(w *model.Worker) {
	d.Workers = append(d.Workers, w)
}

// CheckAvalancheStatus determines if we are in a heavy-load state
func (d *Dispatcher) CheckAvalancheStatus() {
	heavyCount := d.Buckets.TotalHeavyJobs(JobLargeThreshold)
	if heavyCount >= AvalancheThreshold {
		if !d.AvalancheActive {
			fmt.Println("!!! AVALANCHE DETECTED !!! Triggering Reservation Mode.")
			d.AvalancheActive = true
		}
	} else {
		if d.AvalancheActive {
			fmt.Println("... Avalanche cleared. Resuming normal operation.")
			d.AvalancheActive = false
		}
	}
}

// Match performs a single round of matching
func (d *Dispatcher) Match() {
	// 1. Update State
	d.CheckAvalancheStatus()

	// 2. Iterate over Workers
	// Ideally we iterate available workers. For sim, we assume all in slice are "Available".
	for _, w := range d.Workers {
		// Determine the Minimum Bucket Index we are allowed to pick from.
		// Default: 0 (Any job)
		minBucketIdx := 0

		// RESERVATION LOGIC:
		// If Avalanche is Active AND Worker is capable of handling Heavy Jobs,
		// Then enforce Reservation: Worker can ONLY pick jobs >= JobLargeThreshold.
		if d.AvalancheActive && w.CapacityMB >= JobLargeThreshold {
			// Calculate bucket index for the Large Threshold
			minBucketIdx = d.Buckets.GetBucketIndex(JobLargeThreshold)
			// fmt.Printf(" [Reserved] Worker %s reserved for jobs > %d MB\n", w.ID, JobLargeThreshold)
		}

		// Try to find a job
		job := d.Buckets.GetHeaviestJobForCapacity(w.CapacityMB, minBucketIdx)

		if job != nil {
			fmt.Printf("[MATCH] Worker %s (Cap %d) -> Assigned %s\n", w.ID, w.CapacityMB, job)
			// In a real system, we'd mark worker as busy here.
			// d.RemoveWorker(w) or similar.
		} else {
			// If job is nil, it means:
			// 1. No jobs fit capacity.
			// 2. OR Reservation blocked taking small jobs.
			if d.AvalancheActive && w.CapacityMB >= JobLargeThreshold {
				// Verify if it was truly reservation blocking
				// Check if there ARE small jobs?
				// For now, just log that it's waiting
				// fmt.Printf(" [Wait] Worker %s is reserved and waiting for Heavy Job...\n", w.ID)
			}
		}
	}
}
