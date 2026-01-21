package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/adarsh/woc1/queue_algo/07_pull_based/pkg/client"
	"github.com/adarsh/woc1/queue_algo/07_pull_based/pkg/queue"
	"github.com/adarsh/woc1/queue_algo/07_pull_based/pkg/server"
)

func main() {
	fmt.Println("=== ALGO 7: YGGDRASIL PULL-BASED SCHEDULER ===")

	srv := server.NewServer()

	// 1. Ingest Jobs (Simulate API Gateway)
	// Add 500,000 Jobs
	count := 500000
	fmt.Printf("\n--- Phase 1: Ingesting %d Jobs ---\n", count)

	start := time.Now()
	for i := 0; i < count; i++ {
		srv.IngestJob(&queue.Job{
			ID:      fmt.Sprintf("Job-%d", i),
			Region:  rand.Intn(4),          // 0-3
			MinRam:  100 + rand.Intn(2000), // 100-2100MB
			MinBatt: rand.Intn(3),          // 0-2
		})
	}
	fmt.Printf("Ingested Jobs in %s\n", time.Since(start))

	// 2. Start Workers (Simulate Phones coming online)
	// 10,000 Conurrent Phones pulling
	workerCount := 10000
	fmt.Printf("\n--- Phase 2: Starting %d Workers ---\n", workerCount)

	var wg sync.WaitGroup
	workers := make([]*client.Worker, workerCount)
	stopChan := make(chan bool)

	// Launch Workers
	for i := 0; i < workerCount; i++ {
		workers[i] = &client.Worker{
			ID:        fmt.Sprintf("W-%d", i),
			Region:    rand.Intn(4),
			FreeMemMB: 2000 + rand.Intn(4000), // High capacity phones
			Battery:   2,                      // High Battery
			Server:    srv,
		}

		wg.Add(1)
		go func(w *client.Worker) {
			defer wg.Done()
			w.Start(stopChan)
		}(workers[i])
	}

	// Let them run for a bit?
	// Or wait until all jobs drained?
	// Since workers loop forever, we need to check Server Queue.
	// But ShardedQueue doesn't expose "TotalCount".
	// We'll just run for a fixed time (e.g. 2 seocnds) and check throughput.

	fmt.Println("Workers running...")
	start = time.Now()
	time.Sleep(time.Millisecond * 2000)
	close(stopChan)
	wg.Wait() // Wait for workers to stop

	// 3. Calculate Stats
	totalProcessed := 0
	for _, w := range workers {
		totalProcessed += w.JobsProcessed
	}

	dur := time.Since(start)
	// Correct duration is 2000ms roughly, plus cleanup
	// Lets use 2.0s for math approximation or use actual start-stop delta

	opsPerSec := float64(totalProcessed) / dur.Seconds()

	fmt.Printf("\nProcessed %d jobs in %s\n", totalProcessed, dur)
	fmt.Printf("Throughput: %.2f Jobs/Sec\n", opsPerSec)

	if totalProcessed > 0 {
		fmt.Printf("Latency per Job: %.2f ns\n", float64(dur.Nanoseconds())/float64(totalProcessed))
	}
}
