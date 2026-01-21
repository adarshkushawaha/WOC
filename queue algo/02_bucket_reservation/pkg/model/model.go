package model

import (
	"fmt"
)

type Job struct {
	ID     string
	SizeMB int
}

func (j *Job) String() string {
	return fmt.Sprintf("Job{%s, %dMB}", j.ID, j.SizeMB)
}

type Worker struct {
	ID         string
	CapacityMB int
}

func (w *Worker) String() string {
	return fmt.Sprintf("Worker{%s, Cap:%dMB}", w.ID, w.CapacityMB)
}
