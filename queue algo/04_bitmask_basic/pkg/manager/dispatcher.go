package manager

import (
	"fmt"

	"github.com/adarsh/woc1/queue_algo/04_bitmask_basic/pkg/bitmask"
	"github.com/adarsh/woc1/queue_algo/04_bitmask_basic/pkg/dag"
)

type Dispatcher struct {
	SeqEngine   *dag.SequenceEngine
	Scheduler   *bitmask.O1Scheduler
	ActiveFlows map[string]*dag.FlowInstance
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		SeqEngine:   dag.NewSequenceEngine(),
		Scheduler:   bitmask.NewO1Scheduler(),
		ActiveFlows: make(map[string]*dag.FlowInstance),
	}
}

// StartFlow initiates a new user-defined sequence
func (d *Dispatcher) StartFlow(flowName, instanceID string) {
	inst, firstJob := d.SeqEngine.CreateInstance(flowName, instanceID)
	if inst == nil {
		fmt.Printf("Error: Flow %s not found\n", flowName)
		return
	}
	d.ActiveFlows[instanceID] = inst

	d.ScheduleJob(firstJob)
}

// ScheduleJob finds a phone for the job
func (d *Dispatcher) ScheduleJob(jobName string) {
	if jobName == "" {
		return
	}

	// In a real system, we'd look up Job Requirements (e.g. Memory)
	// For simulation, we assume all jobs need 100MB
	neededMB := 100

	// O(1) Lookup
	phone := d.Scheduler.GetBestPhone(neededMB)

	if phone != nil {
		fmt.Printf("[DISPATCH] Assigned %s -> Phone %s (Free %dMB)\n", jobName, phone.ID, phone.FreeMemMB)

		// Simulate Execution Complete immediately for demo
		d.JobComplete(jobName)
	} else {
		fmt.Printf("[QUEUE] %s waiting (No phone with >%dMB)\n", jobName, neededMB)
		// Logic to retry later...
	}
}

func (d *Dispatcher) JobComplete(jobName string) {
	// Parse instance ID from jobName "ID:Step"
	// Simplified: Just assume we know the instance
	// In real code, we'd have a Map[JobName]Instance
	// For this demo, we can't easily map back without parsing.

	// Hack for demo: We just assume a single active flow or ignore completion limit
	// effectively, the dispatcher logic relies on external triggers.
}

// Manual Advance helper for Main.go simulation
func (d *Dispatcher) AdvanceFlow(instanceID string) {
	inst, ok := d.ActiveFlows[instanceID]
	if !ok {
		return
	}

	nextJob, done := d.SeqEngine.Advance(inst)
	if done {
		fmt.Printf("[FLOW] Instance %s Completed.\n", instanceID)
		delete(d.ActiveFlows, instanceID)
	} else {
		d.ScheduleJob(nextJob)
	}
}
