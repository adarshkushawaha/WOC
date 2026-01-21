package dag

import (
	"fmt"
)

// StepID represents a unique step in a sequence
type StepID int

// FlowDef defines a sequence of steps (Blueprints)
type FlowDef struct {
	Name  string
	Steps []string // e.g. ["Fetch", "Process", "Save"]
}

// FlowInstance is a running instance of a Flow
type FlowInstance struct {
	ID           string
	Def          *FlowDef
	CurrentStep  int
	Dependencies []string // Just for external ref
}

// SequenceEngine manages flows
type SequenceEngine struct {
	Flows map[string]*FlowDef
}

func NewSequenceEngine() *SequenceEngine {
	return &SequenceEngine{
		Flows: make(map[string]*FlowDef),
	}
}

// RegisterFlow creates a new blueprint.
func (se *SequenceEngine) RegisterFlow(name string, steps []string) {
	se.Flows[name] = &FlowDef{Name: name, Steps: steps}
}

// CreateInstance starts a flow. Returns the first Job Name.
func (se *SequenceEngine) CreateInstance(flowName string, instanceID string) (*FlowInstance, string) {
	def, ok := se.Flows[flowName]
	if !ok {
		return nil, ""
	}

	inst := &FlowInstance{
		ID:          instanceID,
		Def:         def,
		CurrentStep: 0,
	}

	// Return first task
	return inst, se.GetJobName(inst)
}

func (se *SequenceEngine) GetJobName(inst *FlowInstance) string {
	if inst.CurrentStep >= len(inst.Def.Steps) {
		return ""
	}
	// Format: "InstanceID:StepName"
	return fmt.Sprintf("%s:%s", inst.ID, inst.Def.Steps[inst.CurrentStep])
}

// Advance moves the flow to the next step.
// Returns: (NextJobName, IsComplete)
func (se *SequenceEngine) Advance(inst *FlowInstance) (string, bool) {
	inst.CurrentStep++
	if inst.CurrentStep >= len(inst.Def.Steps) {
		return "", true
	}
	return se.GetJobName(inst), false
}
