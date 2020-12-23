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

import (
	"github.com/uncharted-distil/distil-compute/pipeline"
)

// InferenceStepData provides data for a pipeline description placeholder step,
// which marks the point at which a TA2 should be begin pipeline inference.
type InferenceStepData struct {
	inputRefs map[string]DataRef
	Inputs    []string
	Outputs   []string
}

// NewInferenceStepData creates a InferenceStepData instance with default values.
func NewInferenceStepData(arguments map[string]DataRef) *InferenceStepData {
	values := make([]string, len(arguments))
	i := 0
	for _, arg := range arguments {
		values[i] = arg.RefString()
	}
	return &InferenceStepData{
		Inputs:    values,
		Outputs:   []string{"produce"},
		inputRefs: arguments,
	}
}

// GetPrimitive returns nil since there is no primitive associated with a placeholder
// step.
func (s *InferenceStepData) GetPrimitive() *pipeline.Primitive {
	return nil
}

// GetArguments adapts the internal placeholder step argument type to the primitive
// step argument type.
func (s *InferenceStepData) GetArguments() map[string]DataRef {
	return s.inputRefs
}

// GetHyperparameters returns an empty map since inference steps don't
// take hyper parameters.
func (s *InferenceStepData) GetHyperparameters() map[string]interface{} {
	return map[string]interface{}{}
}

// GetOutputMethods returns a list of methods that will be called to generate
// primitive output.  These feed into downstream primitives.
func (s *InferenceStepData) GetOutputMethods() []string {
	return s.Outputs
}

// BuildDescriptionStep creates protobuf structures from a pipeline step
// definition.
func (s *InferenceStepData) BuildDescriptionStep() (*pipeline.PipelineDescriptionStep, error) {
	// generate arguments entries
	inputs := []*pipeline.StepInput{}
	for _, v := range s.Inputs {
		input := &pipeline.StepInput{
			Data: v,
		}
		inputs = append(inputs, input)
	}

	// list of methods that will generate output - order matters because the steps are
	// numbered
	outputs := []*pipeline.StepOutput{}
	for _, v := range s.Outputs {
		output := &pipeline.StepOutput{
			Id: v,
		}
		outputs = append(outputs, output)
	}

	// create the pipeline description structure
	step := &pipeline.PipelineDescriptionStep{
		Step: &pipeline.PipelineDescriptionStep_Placeholder{
			Placeholder: &pipeline.PlaceholderPipelineDescriptionStep{
				Inputs:  inputs,
				Outputs: outputs,
			},
		},
	}

	return step, nil
}
