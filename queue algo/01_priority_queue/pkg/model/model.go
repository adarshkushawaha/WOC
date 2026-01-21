package model

import (
	"fmt"
	"time"
)

// Job represents a task to be executed.
type Job struct {
	ID          string
	ArrivalTime time.Time
	Duration    time.Duration // Estimated execution time
	Size        int           // Abstract size for score calculation
}

// WaitingTime returns the duration the job has been waiting.
func (j *Job) WaitingTime() time.Duration {
	return time.Since(j.ArrivalTime)
}

// Score calculates the priority score. Higher score = Higher priority.
// Formula: (Waiting Time (seconds) * 1.5) + (Job Size * 0.5)
func (j *Job) Score() float64 {
	return (j.WaitingTime().Seconds() * 1.5) + (float64(j.Size) * 0.5)
}

func (j *Job) String() string {
	return fmt.Sprintf("Job{ID: %s, Size: %d, Score: %.2f}", j.ID, j.Size, j.Score())
}

// Worker represents a resource than can execute a job.
type Worker struct {
	ID            string
	AvailableTime time.Time
	Efficiency    float64 // Multiplier for score
}

// IdleTime returns the duration the worker has been idle.
func (w *Worker) IdleTime() time.Duration {
	return time.Since(w.AvailableTime)
}

// Score calculates the priority score. Higher score = Higher priority.
// Formula: (Idle Time (seconds) * 1.0) + (Efficiency * 10)
func (w *Worker) Score() float64 {
	return (w.IdleTime().Seconds() * 1.0) + (w.Efficiency * 10.0)
}

func (w *Worker) String() string {
	return fmt.Sprintf("Worker{ID: %s, Eff: %.2f, Score: %.2f}", w.ID, w.Efficiency, w.Score())
}
