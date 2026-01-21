package model

import "fmt"

type Job struct {
	ID     string
	SizeMB int
}

func (j *Job) String() string {
	return fmt.Sprintf("Job{%s, %dMB}", j.ID, j.SizeMB)
}

type Worker struct {
	ID          string
	CapacityMB  int
	UsedMemory  int
	CurrentJobs []*Job
}

func (w *Worker) AvailableMemory() int {
	return w.CapacityMB - w.UsedMemory
}

func (w *Worker) String() string {
	return fmt.Sprintf("Worker{%s, Cap:%dMB, Used:%dMB, Jobs:%d}", w.ID, w.CapacityMB, w.UsedMemory, len(w.CurrentJobs))
}
