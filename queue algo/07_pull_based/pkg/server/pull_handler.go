package server

import (
	"github.com/adarsh/woc1/queue_algo/07_pull_based/pkg/queue"
)

type Server struct {
	JobQueue *queue.ShardedQueue
}

func NewServer() *Server {
	return &Server{
		JobQueue: queue.NewShardedQueue(),
	}
}

// HandlePull is the RPC/HTTP handler called by Worker
func (s *Server) HandlePull(workerID string, ram, region, batt int) *queue.Job {
	// 1. Try Instant Pull
	j := s.JobQueue.PullJob(ram, region, batt)
	if j != nil {
		return j
	}

	// 2. Long Polling Simulation
	// In reality, we'd use a channel or condition variable.
	// For simulation, we just retry a few times fast or return nil (Exponential Backoff on client)

	// Return nil means "No Work, Sleep and Ask Later"
	return nil
}

// IngestJob simulates receiving a job from API gateway
func (s *Server) IngestJob(j *queue.Job) {
	s.JobQueue.AddJob(j)
}
