package model

import (
	"fmt"
)

const BucketInterval = 50

// JobBucket represents a single queue for a specific memory range
type JobBucket struct {
	MinSize int
	MaxSize int
	Jobs    []*Job
}

func (b *JobBucket) Push(j *Job) {
	b.Jobs = append(b.Jobs, j)
}

func (b *JobBucket) Pop() *Job {
	if len(b.Jobs) == 0 {
		return nil
	}
	j := b.Jobs[0]
	b.Jobs = b.Jobs[1:]
	return j
}

func (b *JobBucket) IsEmpty() bool {
	return len(b.Jobs) == 0
}

func (b *JobBucket) Len() int {
	return len(b.Jobs)
}

// BucketManager holds all buckets
type BucketManager struct {
	Buckets []*JobBucket
}

func NewBucketManager(maxSizeMB int) *BucketManager {
	// Create buckets up to maxSizeMB
	// e.g. 0-50, 51-100, ...
	numBuckets := (maxSizeMB / BucketInterval) + 1
	bm := &BucketManager{
		Buckets: make([]*JobBucket, numBuckets),
	}

	for i := 0; i < numBuckets; i++ {
		min := i * BucketInterval
		max := (i + 1) * BucketInterval
		if i == 0 {
			min = 0
		} // Handle 0-50
		bm.Buckets[i] = &JobBucket{
			MinSize: min,
			MaxSize: max,
			Jobs:    make([]*Job, 0),
		}
	}
	return bm
}

// GetBucketIndex returns the index for a given size
func (bm *BucketManager) GetBucketIndex(sizeMB int) int {
	// if size is 0-49 -> index 0. 50-99 -> index 1.
	// Actually logic: size / 50.
	idx := sizeMB / BucketInterval
	if idx >= len(bm.Buckets) {
		// Extend or cap? Ideally we'd extend dynamicall logic, but for now lets cap at max bucket or panic
		// Returning last bucket for overflow
		return len(bm.Buckets) - 1
	}
	return idx
}

func (bm *BucketManager) AddJob(j *Job) {
	idx := bm.GetBucketIndex(j.SizeMB)
	bm.Buckets[idx].Push(j)
	fmt.Printf("[BucketManager] Added %s to Bucket[%d] (%d-%d MB)\n", j, idx, bm.Buckets[idx].MinSize, bm.Buckets[idx].MaxSize)
}

// GetHeaviestJobForCapacity finds the heaviest available job that fits within workerCapacity.
// It searches from the largest possible bucket downwards.
// excludedBuckets: useful if we want to forbid "Small" buckets (Reservation Mode)
func (bm *BucketManager) GetHeaviestJobForCapacity(capacityMB int, minBucketIndex int) *Job {
	// Start from the bucket corresponding to capacityMB
	startIdx := bm.GetBucketIndex(capacityMB)

	// Scan downwards
	for i := startIdx; i >= minBucketIndex; i-- {
		if !bm.Buckets[i].IsEmpty() {
			return bm.Buckets[i].Pop()
		}
	}
	return nil
}

// TotalHeavyJobs returns count of jobs above a certain size threshold
func (bm *BucketManager) TotalHeavyJobs(thresholdMB int) int {
	startIdx := bm.GetBucketIndex(thresholdMB)
	count := 0
	for i := startIdx; i < len(bm.Buckets); i++ {
		count += bm.Buckets[i].Len()
	}
	return count
}
