//
//   Copyright © 2019 Uncharted Software Inc.
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package description

// Provides an interface to assemble a D3M pipeline DAG as a protobuf PipelineDescription.  This created
// description can be passed to a TA2 system for execution and inference.  The pipeline description is
// covered in detail at https://gitlab.com/datadrivendiscovery/metalearning#pipeline with example JSON
// pipeline definitions found in that same repository.

import (
	"fmt"
	"sync/atomic"

	"github.com/pkg/errors"
	"github.com/uncharted-distil/distil-compute/pipeline"
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

func (s *PipelineNode) isSource() bool {
	return len(s.parents) == 0
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
func NewPipelineBuilder(name string, description string, sources ...*PipelineNode) *PipelineBuilder {
	builder := &PipelineBuilder{
		sources:     sources,
		name:        name,
		description: description,
	}
	return builder
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

	// ensure that there aren't any cycles in the graph
	if checkCycles(p.sources) {
		return nil, errors.Errorf("compile error: detected cycle in graph")
	}

	// start processing from the roots
	refCount := 0
	for _, sourceNode := range p.sources {
		err := validate(sourceNode)
		if err != nil {
			return nil, err
		}

		// set the primitive inputs to the pipeline inputs in a 1:1 fashion in order
		args := sourceNode.step.GetArguments()
		for _, arg := range args {
			key := fmt.Sprintf("%s.%d", pipelineInputsKey, refCount)
			sourceNode.step.UpdateArguments(arg.Name, key)
			refCount++
		}

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
	for _, arg := range node.step.GetArguments() {
		foundFreeOutput := false
		for _, parent := range node.parents {
			nextOutputIdx := parent.nextOutput()
			numParentOutputs := len(parent.step.GetOutputMethods())
			if nextOutputIdx < numParentOutputs {
				parentOutput := parent.step.GetOutputMethods()[nextOutputIdx]
				inputsRef := fmt.Sprintf("%s.%d.%s", stepKey, idToIndexMap[parent.nodeID], parentOutput)
				node.step.UpdateArguments(arg.Name, inputsRef)

				foundFreeOutput = true
				break
			}
		}

		if !foundFreeOutput && !node.isSource() {
			return nil, errors.Errorf(
				"compile failed: can't find output for step \"%s\" arg \"%s\"",
				node.step.GetPrimitive().GetName(),
				arg.Name)
		}
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
