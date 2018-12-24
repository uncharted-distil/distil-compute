package description

// Provides an interface to assemble a D3M pipeline DAG as a protobuf PipelineDescription.  This created
// description can be passed to a TA2 system for execution and inference.

import (
	"fmt"
	"sync/atomic"

	"github.com/pkg/errors"
	"github.com/unchartedsoftware/distil-compute/pipeline"
)

const (
	stepKey             = "steps"
	pipelineInputsName  = "input"
	pipelineOutputsName = "outputs"
)

// if this needs to be thread safe make use a sync/atomic
var nextNodeID uint64

// PipelineNode creates a pipeline node that can be added to the pipeline DAG.
type PipelineNode struct {
	nodeID    uint64
	step      Step
	children  []*PipelineNode
	parents   []*PipelineNode
	outputIdx int
	visited   bool
}

// Add adds a child node.  If the node has alread been compiled into a pipeline
// this will action will fail.
func (s *PipelineNode) Add(outgoing *PipelineNode) error {
	if s.visited {
		return errors.New("cannot assign child to compiled pipeline element")
	}
	s.children = append(s.children, outgoing)
	outgoing.parents = append(outgoing.parents, s)

	return nil
}

// NewPipelineNode creates a new pipeline node that can be added to the pipeline
// DAG.  PipelneNode structs should only be instantiated through this function to
// ensure internal structure are properly initialized.
func NewPipelineNode(step Step) *PipelineNode {
	newStep := &PipelineNode{
		nodeID:   nextNodeID,
		step:     step,
		children: []*PipelineNode{},
		parents:  []*PipelineNode{},
	}
	atomic.AddUint64(&nextNodeID, 1)
	return newStep
}

func (s *PipelineNode) nextOutput() int {
	val := s.outputIdx
	s.outputIdx++
	return val
}

func (s *PipelineNode) isSink() bool {
	return len(s.children) == 0
}

// PipelineBuilder compiles a pipeline DAG into a protobuf pipeline description that can
// be passed to a downstream TA2 for inference (optional) and execution.
type PipelineBuilder struct {
	name        string
	description string
	sources     []*PipelineNode
	compiled    bool
	inferred    bool
}

// NewPipelineBuilder creates a new pipeline builder instance.  Source nodes need to be added in a subsequent call.
func NewPipelineBuilder(name string, description string) *PipelineBuilder {
	builder := &PipelineBuilder{
		sources:     []*PipelineNode{},
		name:        name,
		description: description,
	}
	return builder
}

// AddSource adds a root node to the pipeline builder.
func (p *PipelineBuilder) AddSource(child *PipelineNode) {
	p.sources = append(p.sources, child)
}

// Compile creates the protobuf pipeline description from the step graph.  It can only be
// called once.
func (p *PipelineBuilder) Compile() (*pipeline.PipelineDescription, error) {
	if p.compiled {
		return nil, errors.New("compile failed: pipeline already compiled")
	}

	if len(p.sources) == 0 {
		return nil, errors.New("compile failed: pipeline requires at least 1 step")
	}

	pipelineNodes := []*PipelineNode{}
	idToIndexMap := map[uint64]int{}
	traversalQueue := []*PipelineNode{}

	// ensure that there aren't any cycles in the graph - the traversal above assigns node IDs
	// so it needs to be deferred until after
	if checkCycles(p.sources) {
		return nil, errors.Errorf("compile error: detected cycle in graph")
	}

	// start processing from the roots
	for i, sourceNode := range p.sources {
		err := validate(sourceNode)
		if err != nil {
			return nil, err
		}

		// set the input to the dataset by default for each of the root steps
		key := fmt.Sprintf("%s.%d", pipelineInputsKey, i)
		args := sourceNode.step.GetArguments()
		_, ok := args[key]
		if ok {
			return nil, errors.Errorf("compile failed: argument `%s` is reserved for internal use", key)
		}
		sourceNode.step.UpdateArguments(pipelineInputsKey, key) // TODO: will just overwite until args can be a list

		// add to traversal queue
		traversalQueue = append(traversalQueue, sourceNode)
	}

	// perform a breadth first traversal of the DAG to establish connections between
	// steps
	for len(traversalQueue) > 0 {
		node := traversalQueue[0]
		traversalQueue = traversalQueue[1:]
		var err error
		pipelineNodes, err = p.processNode(node, pipelineNodes, idToIndexMap)
		if err != nil {
			return nil, err
		}

		// Process any children that haven't yet been visited
		for _, child := range node.children {
			traversalQueue = append(traversalQueue, child)
		}
	}

	// Set the outputs from the pipeline graph sinks
	pipelineOutputs := []*pipeline.PipelineDescriptionOutput{}
	for i, node := range pipelineNodes {
		if node.isSink() {
			for _, outputMethod := range node.step.GetOutputMethods() {
				output := &pipeline.PipelineDescriptionOutput{
					Name: fmt.Sprintf("%s %d", pipelineOutputsName, i),
					Data: fmt.Sprintf("%s.%d.%s", stepKey, idToIndexMap[node.nodeID], outputMethod),
				}
				pipelineOutputs = append(pipelineOutputs, output)
			}
		}
	}

	// Set the input to to the placeholder
	pipelineInputs := []*pipeline.PipelineDescriptionInput{}
	for i := range p.sources {
		pipelineInputs = append(pipelineInputs, &pipeline.PipelineDescriptionInput{
			Name: fmt.Sprintf("%s %d", pipelineInputsName, i),
		})
	}

	// Build the step descriptions now that all of the inputs/outputs defined
	compiledSteps := []*pipeline.PipelineDescriptionStep{}
	for _, node := range pipelineNodes {
		compiledStep, err := node.step.BuildDescriptionStep()
		compiledSteps = append(compiledSteps, compiledStep)
		if err != nil {
			return nil, err
		}
	}

	pipelineDesc := &pipeline.PipelineDescription{
		Name:        p.name,
		Description: p.description,
		Steps:       compiledSteps,
		Inputs:      pipelineInputs,
		Outputs:     pipelineOutputs,
		Context:     pipeline.PipelineContext_TESTING,
	}

	// mark the entire pipeline as compiled so it can't be compiled again
	p.compiled = true

	return pipelineDesc, nil
}

func validate(node *PipelineNode) error {
	// Validate step parameters.  This is currently pretty surface level, but we could
	// go in validate the struct hierarchy to catch more potential caller errors during
	// the compile step.
	//
	// NOTE: Hyperparameters and Primitive are optional so there is no included check at this time.
	args := node.step.GetArguments()
	if args == nil {
		return errors.Errorf("compile failed: step \"%s\" missing argument list", node.step.GetPrimitive().GetName())
	}

	outputs := node.step.GetOutputMethods()
	if len(outputs) == 0 {
		return errors.Errorf("compile failed: expected at least 1 output for step \"%s\"", node.step.GetPrimitive().GetName())
	}

	// If this is an inference step, make sure it has no children.
	if _, ok := node.step.(*InferenceStepData); ok && !node.isSink() {
		return errors.Errorf("compile failed: inference step cannot have children")
	}
	return nil
}

func (p *PipelineBuilder) processNode(node *PipelineNode, pipelineNodes []*PipelineNode, idToIndexMap map[uint64]int) ([]*PipelineNode, error) {
	// don't re-process
	if node.visited {
		return pipelineNodes, nil
	}

	// validate node args
	if err := validate(node); err != nil {
		return nil, err
	}

	// Enforce a single inferred node.
	if _, ok := node.step.(*InferenceStepData); ok {
		if !p.inferred {
			p.inferred = true
		} else {
			return nil, errors.Errorf("compile failed: attempted to define more than one inference step")
		}
	}

	// Connect each input to the next unattached parent output
	for _, parent := range node.parents {
		numParentOutputs := len(parent.step.GetOutputMethods())
		nextOutputIdx := parent.nextOutput()
		if nextOutputIdx >= numParentOutputs {
			return nil, errors.Errorf("compile failed: step \"%s\" can't find parent input", node.step.GetPrimitive().GetName())
		}

		parentOutput := parent.step.GetOutputMethods()[nextOutputIdx]
		inputsRef := fmt.Sprintf("%s.%d.%s", stepKey, idToIndexMap[parent.nodeID], parentOutput)
		node.step.UpdateArguments(stepInputsKey, inputsRef)
	}

	// add to the node list
	pipelineNodes = append(pipelineNodes, node)
	idToIndexMap[node.nodeID] = len(pipelineNodes) - 1

	// Mark as visited so we don't reprocess
	node.visited = true

	return pipelineNodes, nil
}

func checkCycles(sources []*PipelineNode) bool {
	ids := map[uint64]bool{}
	for _, sourceNode := range sources {
		cycle := checkNode(sourceNode, ids)
		if cycle {
			return true
		}
	}
	return false
}

func checkNode(node *PipelineNode, ids map[uint64]bool) bool {

	if _, ok := ids[node.nodeID]; ok {
		return true
	}

	ids[node.nodeID] = true
	for _, childNode := range node.children {
		cycle := checkNode(childNode, ids)
		if cycle {
			return true
		}
	}
	delete(ids, node.nodeID)
	return false
}
