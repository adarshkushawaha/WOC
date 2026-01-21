package manager

import (
	"container/heap"
	"fmt"

	"github.com/adarsh/woc1/queue_algo/01_priority_queue/pkg/model"
	"github.com/adarsh/woc1/queue_algo/01_priority_queue/pkg/queue"
)

type Dispatcher struct {
	// Job Queues
	TimeJobQueue  *queue.TimeJobQueue
	WaitJobQueue  *queue.WaitJobQueue
	ScoreJobQueue *queue.ScoreJobQueue

	// Worker Queues
	TimeWorkerQueue  *queue.TimeWorkerQueue
	WaitWorkerQueue  *queue.WaitWorkerQueue
	ScoreWorkerQueue *queue.ScoreWorkerQueue

	// Maps to track items across queues for safe removal (optional, but good for O(1) lookup if needed)
	// For this simulation, we will assume removing the popped item from other queues is key.
	// Since heap removal is O(n) without an index, and we are sharing pointers, we need to handle "already processed" items.
	// A simple way is to mark the Job/Worker as "Assigned" in the struct, or maintain a set of IDs.
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		TimeJobQueue:     &queue.TimeJobQueue{},
		WaitJobQueue:     &queue.WaitJobQueue{},
		ScoreJobQueue:    &queue.ScoreJobQueue{},
		TimeWorkerQueue:  &queue.TimeWorkerQueue{},
		WaitWorkerQueue:  &queue.WaitWorkerQueue{},
		ScoreWorkerQueue: &queue.ScoreWorkerQueue{},
	}
}

func (d *Dispatcher) AddJob(j *model.Job) {
	fmt.Printf("[Dispatcher] Adding Job: %s\n", j)
	d.TimeJobQueue.PushJob(j)
	d.WaitJobQueue.PushJob(j)
	d.ScoreJobQueue.PushJob(j)
}

func (d *Dispatcher) AddWorker(w *model.Worker) {
	fmt.Printf("[Dispatcher] Adding Worker: %s\n", w)
	d.TimeWorkerQueue.PushWorker(w)
	d.WaitWorkerQueue.PushWorker(w)
	d.ScoreWorkerQueue.PushWorker(w)
}

// removeJobFromAll removes the job from all queues using the Heap Remove if we tracked indices,
// or we rebuild/filter.
// For simplicity in this specific "container/heap" usage without a shared index map map:
// We will just Pop from the ScoreQueue (Master) and then lazily ignore it in others or do a linear scan remove.
// Given the prompt asks for "queues", synchronization is implied.
// Linear scan remove is O(N).
func (d *Dispatcher) removeJob(j *model.Job) {
	// Helper to remove a specific job from a specific heap interface
	remove := func(h interface{}, jobID string) {
		// This is tricky with standard heap interface without exposing the slice.
		// We will rely on unique IDs.
		// NOTE: In a real high-perf system, we'd use the map[ID]*Item to get the index.
		// Here we will just iterate and remove for correctness demonstration.
		// Since we defined the underlying types as slices in queue package, we can't easily access them without casting or adding Remove helpers there.
		// Let's rely on the fact that if we POP from one, we should remove from others.

		// To make this robust, let's implement a 'RemoveByID' in the queue package or just handle it here by
		// popping everything and pushing back except the target. (Inefficient but clear).
		// actually, let's just mark the job as "Processed" using a map in Dispatcher? No, the user wants queues.
	}
	_ = remove // placeholder
}

// Helper to remove by Ref (requires scanning the slice)
func removeJobFromHeap(h heap.Interface, job *model.Job) {
	// We need to access the underlying slice to find the index.
	// We can do this by repeated popping or better: modify the Queue implementation to support Remove(item).
	// But `heap.Remove` takes an index.
	// Let's defer this complexity. For this algorithm, we MUST pick from the Mathematical/Score queue as per instructions?
	// "final queue is based on something mathematical dependency... this will assign the job"
	// So we Pop from ScoreQueue.
}

// Match attempts to pair the highest priority Job with the highest priority Worker.
func (d *Dispatcher) Match() {
	if d.ScoreJobQueue.Len() == 0 || d.ScoreWorkerQueue.Len() == 0 {
		fmt.Println("[Dispatcher] Not enough jobs or workers to match.")
		return
	}

	// 1. Get Best Job (Mathematical Queue)
	// We use the ScoreQueue as the decision maker.
	jobItem := heap.Pop(d.ScoreJobQueue).(*queue.JobItem) // Pop from Score Queue
	job := jobItem.Job

	// 2. Get Best Worker (Mathematical Queue)
	workerItem := heap.Pop(d.ScoreWorkerQueue).(*queue.WorkerItem)
	worker := workerItem.Worker

	fmt.Printf("\n[MATCH] Assigned %s \n        To       %s\n\n", job, worker)

	// 3. Cleanup: Remove from other queues.
	// Since we popped from ScoreQueue, we need to remove this specific job/worker from Time and Wait queues.
	d.removeFromOtherJobQueues(job)
	d.removeFromOtherWorkerQueues(worker)
}

func (d *Dispatcher) removeFromOtherJobQueues(target *model.Job) {
	// Brute-force remove for correctness
	// TimeQueue
	for i := 0; i < d.TimeJobQueue.Len(); i++ {
		if d.TimeJobQueue.BaseJobQueue[i].Job.ID == target.ID {
			heap.Remove(d.TimeJobQueue, i)
			break
		}
	}
	// WaitQueue
	for i := 0; i < d.WaitJobQueue.Len(); i++ {
		if d.WaitJobQueue.BaseJobQueue[i].Job.ID == target.ID {
			heap.Remove(d.WaitJobQueue, i)
			break
		}
	}
}

func (d *Dispatcher) removeFromOtherWorkerQueues(target *model.Worker) {
	// TimeQueue
	for i := 0; i < d.TimeWorkerQueue.Len(); i++ {
		if d.TimeWorkerQueue.BaseWorkerQueue[i].Worker.ID == target.ID {
			heap.Remove(d.TimeWorkerQueue, i)
			break
		}
	}
	// WaitQueue
	for i := 0; i < d.WaitWorkerQueue.Len(); i++ {
		if d.WaitWorkerQueue.BaseWorkerQueue[i].Worker.ID == target.ID {
			heap.Remove(d.WaitWorkerQueue, i)
			break
		}
	}
}

func (d *Dispatcher) PrintStatus() {
	fmt.Println("--- Queue Status ---")
	fmt.Printf("Jobs Pending: %d\n", d.ScoreJobQueue.Len())
	fmt.Printf("Workers Available: %d\n", d.ScoreWorkerQueue.Len())
	fmt.Println("--------------------")
}
