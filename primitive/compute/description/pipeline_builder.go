package description

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/unchartedsoftware/distil-compute/pipeline"
)

var nextNodeID = 0

// PipelineNode creates a pipeline node that can be added to the pipeline DAG.
type PipelineNode struct {
	nodeID    int
	step      Step
	children  []*PipelineNode
	parents   []*PipelineNode
	outputIdx int
	visited   bool
}

// Add adds an outoing reference to a node.
func (s *PipelineNode) Add(outgoing *PipelineNode) {
	s.children = append(s.children, outgoing)
	outgoing.parents = append(outgoing.parents, s)
}

func (s *PipelineNode) nextOutput() int {
	val := s.outputIdx
	s.outputIdx++
	return val
}

// NewPipelineNode creates a new pipeline node that can be added to the pipeline
// DAG.
func NewPipelineNode(step Step) *PipelineNode {
	newStep := &PipelineNode{
		nodeID:   nextNodeID,
		step:     step,
		children: []*PipelineNode{},
		parents:  []*PipelineNode{},
	}
	nextNodeID++
	return newStep
}

// PipelineBuilder compiles a pipeline DAG into a protobuf pipeline description that can
// be passed to a downstream TA2 for inference (optional) and execution.
type PipelineBuilder struct {
	name        string
	description string
	sources     []*PipelineNode
	sinks       []*PipelineNode
	compiled    bool
	inferred    bool
}

// NewPipelineBuilder creates a new pipeline builder instance.  Source nodes need to be added in a subsequent call.
func NewPipelineBuilder(name string, description string) *PipelineBuilder {
	builder := &PipelineBuilder{
		sources:     []*PipelineNode{},
		sinks:       []*PipelineNode{},
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
	// TODO: Add a check for cycles.  There is no enforcement of acyclic constraint currently.

	if p.compiled {
		return nil, errors.New("compile failed: pipeline already compiled")
	}

	if len(p.sources) == 0 {
		return nil, errors.New("compile failed: pipeline requires at least 1 step")
	}

	// start processing from the roots
	pipelineNodes := []*PipelineNode{}
	for _, sourceNode := range p.sources {
		err := validate(sourceNode)
		if err != nil {
			return nil, err
		}

		// set the input to the dataset by default for each of the root steps
		key := fmt.Sprintf("%s.0", pipelineInputsKey)
		args := sourceNode.step.GetArguments()
		_, ok := args[key]
		if ok {
			return nil, errors.Errorf("compile failed: argument `%s` is reserved for internal use", key)
		}
		sourceNode.step.UpdateArguments(pipelineInputsKey, key) // TODO: will just overwite until args can be a list

		// compile strating at the roots
		if len(sourceNode.children) == 0 {
			p.sinks = append(p.sinks, sourceNode)
		}
		pipelineNodes, err = p.processNode(sourceNode, pipelineNodes)
		if err != nil {
			return nil, err
		}
	}

	// Set the outputs from the pipeline graph sinks
	pipelineOutputs := []*pipeline.PipelineDescriptionOutput{}
	for _, sink := range p.sinks {
		for _, outputMethod := range sink.step.GetOutputMethods() {
			output := &pipeline.PipelineDescriptionOutput{
				Name: "outputs",
				Data: fmt.Sprintf("steps.%d.%s", sink.nodeID, outputMethod),
			}
			pipelineOutputs = append(pipelineOutputs, output)
		}
	}

	// Set the input to to the placeholder
	pipelineInputs := []*pipeline.PipelineDescriptionInput{
		{
			Name: "input",
		},
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
	return nil
}

func (p *PipelineBuilder) processNode(node *PipelineNode, pipelineNodes []*PipelineNode) ([]*PipelineNode, error) {
	// Connect each input to the next unattached parent output
	for _, parent := range node.parents {
		parentOutput := parent.step.GetOutputMethods()[parent.nextOutput()]
		inputsRef := fmt.Sprintf("steps.%d.%s", parent.nodeID, parentOutput)
		node.step.UpdateArguments(stepInputsKey, inputsRef)
	}

	// add to the node list
	pipelineNodes = append(pipelineNodes, node)

	// Mark as visited so we don't reprocess
	node.visited = true

	// Mark as a leaf for end-of-pipeline output processing if necessary
	if len(node.children) == 0 {
		p.sinks = append(p.sinks, node)
	}

	// Process any children that haven't yet been visited
	var err error
	for _, child := range node.children {
		if !child.visited {
			pipelineNodes, err = p.processNode(child, pipelineNodes)
			if err != nil {
				return nil, err
			}
		}
	}
	return pipelineNodes, nil
}
