package client

import (
	"time"

	"github.com/adarsh/woc1/queue_algo/07_pull_based/pkg/server"
)

type Worker struct {
	ID        string
	Region    int
	FreeMemMB int
	Battery   int
	Server    *server.Server

	JobsProcessed int
}

func (w *Worker) Start(stopChan chan bool) {
	// Simple Loop: Pull -> Work -> Sleep -> Repeat
	ticker := time.NewTicker(time.Millisecond * 10) // Poll every 10ms
	defer ticker.Stop()

	for {
		select {
		case <-stopChan:
			return
		case <-ticker.C:
			// "Pull" from server
			job := w.Server.HandlePull(w.ID, w.FreeMemMB, w.Region, w.Battery)

			if job != nil {
				// Process Job
				// fmt.Printf("Worker %s got job %s\n", w.ID, job.ID)
				w.JobsProcessed++

				// Simulate work duration?
				// For high throughput test, we assume instant lambda
			} else {
				// Exponential Backoff would go here
			}
		}
	}
}
