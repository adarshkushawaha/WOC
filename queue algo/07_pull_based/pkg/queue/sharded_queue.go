package queue

import (
	"sync"
	"time"
)

// Job definition reusing core concepts
type Job struct {
	ID         string
	MinRam     int
	MinBatt    int
	Region     int
	DurationMs int // Added for Simulator
	CreatedAt  time.Time
}

// ShardedQueue manages jobs partitioned by Region and Requirements
// We use a chart: Queues[Region][Battery][RamTier] -> []Job
const RamTiers = 3 // 0: <1GB, 1: 1-4GB, 2: >4GB

type ShardedQueue struct {
	// [Region][Battery][RamTier] -> List of Jobs
	Queues [5][3][3][]*Job

	Locks [5]sync.Mutex
}

func getRamTier(mb int) int {
	if mb < 1000 {
		return 0
	}
	if mb < 4000 {
		return 1
	}
	return 2
}

func NewShardedQueue() *ShardedQueue {
	return &ShardedQueue{}
}

func (sq *ShardedQueue) AddJob(j *Job) {
	if j.Region >= 5 {
		j.Region = 0
	}
	tier := getRamTier(j.MinRam)

	sq.Locks[j.Region].Lock()
	defer sq.Locks[j.Region].Unlock()

	sq.Queues[j.Region][j.MinBatt][tier] = append(sq.Queues[j.Region][j.MinBatt][tier], j)
}

// PullJob attempts to find a job for a worker with specific specs
func (sq *ShardedQueue) PullJob(workerRam, workerRegion, workerBatt int) *Job {
	// Worker in Region X can pull jobs from Region X or Region 0 (Global)
	// We check specific region first

	j := sq.tryPullFromRegion(workerRegion, workerRam, workerBatt)
	if j != nil {
		return j
	}

	// Try Global
	if workerRegion != 0 {
		j = sq.tryPullFromRegion(0, workerRam, workerBatt)
	}

	return j
}

func (sq *ShardedQueue) tryPullFromRegion(reg, workerRam, batt int) *Job {
	sq.Locks[reg].Lock()
	defer sq.Locks[reg].Unlock()

	workerTier := getRamTier(workerRam)

	// 1. Iterate Battery Levels (High -> Low)
	for b := batt; b >= 0; b-- {
		// 2. Iterate RAM Tiers (Worker Tier -> 0)
		// Crucial: Heavy Worker checks Heavy Jobs first.
		for t := workerTier; t >= 0; t-- {
			q := sq.Queues[reg][b][t]
			if len(q) > 0 {
				// Found a job!
				job := q[0]
				sq.Queues[reg][b][t] = q[1:]
				return job
			}
		}
	}
	return nil
}
