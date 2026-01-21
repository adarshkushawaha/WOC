package queue

import (
	"container/heap"

	"github.com/adarsh/woc1/queue_algo/01_priority_queue/pkg/model"
)

// Item wraps the actual data to manage heap indices
type JobItem struct {
	Job   *model.Job
	Index int
}

type WorkerItem struct {
	Worker *model.Worker
	Index  int
}

// --- Job Queues ---

// JobQueue interface
type JobQueue interface {
	heap.Interface
	PushJob(*model.Job)
	PopJob() *model.Job
	PeekJob() *model.Job
}

// BaseJobQueue containing the slice
type BaseJobQueue []*JobItem

func (pq BaseJobQueue) Len() int { return len(pq) }

func (pq BaseJobQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *BaseJobQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*JobItem)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *BaseJobQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // Avoid memory leak
	item.Index = -1
	*pq = old[0 : n-1]
	return item
}

// --- Specific Job Queues ---

// TimeJobQueue: Min-heap based on ArrivalTime
type TimeJobQueue struct {
	BaseJobQueue
}

func (pq TimeJobQueue) Less(i, j int) bool {
	return pq.BaseJobQueue[i].Job.ArrivalTime.Before(pq.BaseJobQueue[j].Job.ArrivalTime)
}
func (pq *TimeJobQueue) PushJob(j *model.Job) { heap.Push(pq, &JobItem{Job: j}) }
func (pq *TimeJobQueue) PopJob() *model.Job {
	if pq.Len() == 0 {
		return nil
	}
	return heap.Pop(pq).(*JobItem).Job
}
func (pq *TimeJobQueue) PeekJob() *model.Job {
	if pq.Len() == 0 {
		return nil
	}
	return pq.BaseJobQueue[0].Job
}

// WaitJobQueue: Max-heap based on WaitingTime
type WaitJobQueue struct {
	BaseJobQueue
}

func (pq WaitJobQueue) Less(i, j int) bool {
	// Earlier arrival = Longer wait. We want Longest wait first (Max-heap logic).
	// So we want 'Before' to be TRUE if i has EARLIER arrival than j.
	// Logic: If i arriving at 10:00 (Wait 20m) and j at 10:10 (Wait 10m).
	// We want i < j to be true if i should pop first?
	// Heap Min-Heap pops the "smallest".
	// If we want Max-Wait (Earliest Arrival) to pop first, then strictly speaking:
	// A standard TimeQueue (Min Arrival) IS a Max-Wait Queue.
	// So Less logic is same: I.Arrival < J.Arrival.
	return pq.BaseJobQueue[i].Job.ArrivalTime.Before(pq.BaseJobQueue[j].Job.ArrivalTime)
}
func (pq *WaitJobQueue) PushJob(j *model.Job) { heap.Push(pq, &JobItem{Job: j}) }
func (pq *WaitJobQueue) PopJob() *model.Job {
	if pq.Len() == 0 {
		return nil
	}
	return heap.Pop(pq).(*JobItem).Job
}
func (pq *WaitJobQueue) PeekJob() *model.Job {
	if pq.Len() == 0 {
		return nil
	}
	return pq.BaseJobQueue[0].Job
}

// ScoreJobQueue: Max-heap based on Score
type ScoreJobQueue struct {
	BaseJobQueue
}

func (pq ScoreJobQueue) Less(i, j int) bool {
	// Max Heap: We want Higher Score first.
	// Less(i, j) should return true if i > j? No, heap is Min-Heap.
	// To make it Max-Heap, Less must return true if i > j.
	return pq.BaseJobQueue[i].Job.Score() > pq.BaseJobQueue[j].Job.Score()
}
func (pq *ScoreJobQueue) PushJob(j *model.Job) { heap.Push(pq, &JobItem{Job: j}) }
func (pq *ScoreJobQueue) PopJob() *model.Job {
	if pq.Len() == 0 {
		return nil
	}
	return heap.Pop(pq).(*JobItem).Job
}
func (pq *ScoreJobQueue) PeekJob() *model.Job {
	if pq.Len() == 0 {
		return nil
	}
	return pq.BaseJobQueue[0].Job
}

// --- Worker Queues ---

type BaseWorkerQueue []*WorkerItem

func (pq BaseWorkerQueue) Len() int { return len(pq) }

func (pq BaseWorkerQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *BaseWorkerQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*WorkerItem)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *BaseWorkerQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.Index = -1
	*pq = old[0 : n-1]
	return item
}

// TimeWorkerQueue
type TimeWorkerQueue struct {
	BaseWorkerQueue
}

func (pq TimeWorkerQueue) Less(i, j int) bool {
	return pq.BaseWorkerQueue[i].Worker.AvailableTime.Before(pq.BaseWorkerQueue[j].Worker.AvailableTime)
}
func (pq *TimeWorkerQueue) PushWorker(w *model.Worker) { heap.Push(pq, &WorkerItem{Worker: w}) }
func (pq *TimeWorkerQueue) PopWorker() *model.Worker {
	if pq.Len() == 0 {
		return nil
	}
	return heap.Pop(pq).(*WorkerItem).Worker
}
func (pq *TimeWorkerQueue) PeekWorker() *model.Worker {
	if pq.Len() == 0 {
		return nil
	}
	return pq.BaseWorkerQueue[0].Worker
}

// WaitWorkerQueue
type WaitWorkerQueue struct {
	BaseWorkerQueue
}

func (pq WaitWorkerQueue) Less(i, j int) bool {
	return pq.BaseWorkerQueue[i].Worker.AvailableTime.Before(pq.BaseWorkerQueue[j].Worker.AvailableTime)
}
func (pq *WaitWorkerQueue) PushWorker(w *model.Worker) { heap.Push(pq, &WorkerItem{Worker: w}) }
func (pq *WaitWorkerQueue) PopWorker() *model.Worker {
	if pq.Len() == 0 {
		return nil
	}
	return heap.Pop(pq).(*WorkerItem).Worker
}
func (pq *WaitWorkerQueue) PeekWorker() *model.Worker {
	if pq.Len() == 0 {
		return nil
	}
	return pq.BaseWorkerQueue[0].Worker
}

// ScoreWorkerQueue
type ScoreWorkerQueue struct {
	BaseWorkerQueue
}

func (pq ScoreWorkerQueue) Less(i, j int) bool {
	return pq.BaseWorkerQueue[i].Worker.Score() > pq.BaseWorkerQueue[j].Worker.Score()
}
func (pq *ScoreWorkerQueue) PushWorker(w *model.Worker) { heap.Push(pq, &WorkerItem{Worker: w}) }
func (pq *ScoreWorkerQueue) PopWorker() *model.Worker {
	if pq.Len() == 0 {
		return nil
	}
	return heap.Pop(pq).(*WorkerItem).Worker
}
func (pq *ScoreWorkerQueue) PeekWorker() *model.Worker {
	if pq.Len() == 0 {
		return nil
	}
	return pq.BaseWorkerQueue[0].Worker
}
