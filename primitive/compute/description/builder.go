//
//   Copyright © 2020 Uncharted Software Inc.
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

	"github.com/pkg/errors"
	"github.com/uncharted-distil/distil-compute/pipeline"
)

// PipelineBuilder compiles a pipeline DAG into a protobuf pipeline description that can
// be passed to a downstream TA2 for inference (optional) and execution.
type PipelineBuilder struct {
	name        string
	description string
	inputs      []string
	outputs     []DataRef
	steps       []Step
	compiled    bool
}

// NewPipelineBuilder creates a new pipeline builder instance.  All of the source nodes in the pipeline
// DAG need to be passed in to the builder via the sources argument, which is variadic.
func NewPipelineBuilder(name string, description string, inputs []string, outputs []DataRef, steps []Step) *PipelineBuilder {
	builder := &PipelineBuilder{
		name:        name,
		description: description,
		inputs:      inputs,
		outputs:     outputs,
		steps:       steps,
	}
	return builder
}

// GetSteps returns compiled steps.
func (p *PipelineBuilder) GetSteps() []Step {
	return p.steps
}

// Compile creates the protobuf pipeline description from the step graph.  It can only be
// called once.
func (p *PipelineBuilder) Compile() (*pipeline.PipelineDescription, error) {
	if p.compiled {
		return nil, errors.New("compile failed: pipeline already compiled")
	}

	if len(p.steps) == 0 {
		return nil, errors.New("compile failed: pipeline requires at least 1 step")
	}

	// Set the inputs
	pipelineInputs := []*pipeline.PipelineDescriptionInput{}
	for _, input := range p.inputs {
		pipelineInputs = append(pipelineInputs, &pipeline.PipelineDescriptionInput{
			Name: input,
		})
	}

	// Set the outputs
	pipelineOutputs := []*pipeline.PipelineDescriptionOutput{}
	for i, output := range p.outputs {
		output := &pipeline.PipelineDescriptionOutput{
			Name: fmt.Sprintf("%s %d", pipelineOutputsName, i),
			Data: output.RefString(),
		}
		pipelineOutputs = append(pipelineOutputs, output)
	}

	// Compile the build steps
	compileResults := []*pipeline.PipelineDescriptionStep{}
	for _, step := range p.steps {
		compileResult, err := step.BuildDescriptionStep()
		if err != nil {
			return nil, err
		}
		compileResults = append(compileResults, compileResult)
	}

	pipelineDesc := &pipeline.PipelineDescription{
		Name:        p.name,
		Description: p.description,
		Steps:       compileResults,
		Inputs:      pipelineInputs,
		Outputs:     pipelineOutputs,
	}

	// mark the entire pipeline as compiled so it can't be compiled again
	p.compiled = true

	return pipelineDesc, nil
}
