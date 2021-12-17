//
//   Copyright © 2021 Uncharted Software Inc.
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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddPipelineBadInputs(t *testing.T) {
	inputsa := []string{"inputsa"}
	outputsa := []DataRef{&StepDataRef{2, "produce"}}
	inputsb := []string{"inputsb"}
	outputsb := []DataRef{&StepDataRef{2, "produce"}}

	step0a := createTestStep(0, []string{"produce"}, map[string]DataRef{"inputs": &PipelineDataRef{0}})
	step1a := createTestStep(1, []string{"produce"}, map[string]DataRef{"inputs": &StepDataRef{0, "produce"}})
	step2a := createTestStep(2, []string{"produce"}, map[string]DataRef{"inputs": &StepDataRef{1, "produce"}})
	pipelineStepsa := []Step{step0a, step1a, step2a}
	desca, err := NewPipelineBuilder("testa", "testa pipeline", inputsa, outputsa, pipelineStepsa).Compile()
	assert.NoError(t, err)
	pipa := &FullySpecifiedPipeline{
		Pipeline: desca,
		steps:    pipelineStepsa,
		name:     "testa",
		inputs:   inputsa,
		outputs:  outputsa,
	}

	step0b := createTestStep(0, []string{"produce"}, map[string]DataRef{"inputs": &PipelineDataRef{0}})
	step1b := createTestStep(1, []string{"produce"}, map[string]DataRef{"inputs": &StepDataRef{0, "produce"}})
	step2b := createTestStep(2, []string{"produce"}, map[string]DataRef{"inputs": &StepDataRef{1, "produce"}})
	pipelineStepsb := []Step{step0b, step1b, step2b}
	descb, err := NewPipelineBuilder("testb", "testb pipeline", inputsb, outputsb, pipelineStepsb).Compile()
	assert.NoError(t, err)
	pipb := &FullySpecifiedPipeline{
		Pipeline: descb,
		steps:    pipelineStepsb,
		name:     "testb",
		inputs:   inputsb,
		outputs:  outputsb,
	}

	err = pipa.AddPipeline(pipb)
	assert.Error(t, err)

	inputsb = []string{"inputsa", "inputsb"}
	pipb = &FullySpecifiedPipeline{
		Pipeline: descb,
		steps:    pipelineStepsb,
		name:     "testb",
		inputs:   inputsb,
		outputs:  outputsb,
	}

	err = pipa.AddPipeline(pipb)
	assert.Error(t, err)

	inputsa = []string{"inputsa", "inputsb", "inputsc"}
	pipa = &FullySpecifiedPipeline{
		Pipeline: desca,
		steps:    pipelineStepsa,
		name:     "testa",
		inputs:   inputsa,
		outputs:  outputsa,
	}

	err = pipa.AddPipeline(pipb)
	assert.Error(t, err)
}

func TestAddPipeline(t *testing.T) {
	inputs := []string{"inputs"}
	outputsa := []DataRef{&StepDataRef{2, "produce"}}
	outputsb := []DataRef{&StepDataRef{2, "produce"}}

	step0a := createTestStep(0, []string{"produce"}, map[string]DataRef{"inputs": &PipelineDataRef{0}})
	step1a := createTestStep(1, []string{"produce"}, map[string]DataRef{"inputs": &StepDataRef{0, "produce"}})
	step2a := createTestStep(2, []string{"produce"}, map[string]DataRef{"inputs": &StepDataRef{1, "produce"}})
	pipelineStepsa := []Step{step0a, step1a, step2a}
	desca, err := NewPipelineBuilder("testa", "testa pipeline", inputs, outputsa, pipelineStepsa).Compile()
	assert.NoError(t, err)
	pipa := &FullySpecifiedPipeline{
		Pipeline: desca,
		steps:    pipelineStepsa,
		name:     "testa",
		inputs:   inputs,
		outputs:  outputsa,
	}

	step0b := createTestStep(0, []string{"produce"}, map[string]DataRef{"inputs": &PipelineDataRef{0}})
	step1b := createTestStep(1, []string{"produce"}, map[string]DataRef{"inputs": &StepDataRef{0, "produce"}})
	step2b := createTestStep(2, []string{"produce"}, map[string]DataRef{"inputs": &StepDataRef{1, "produce"}})
	pipelineStepsb := []Step{step0b, step1b, step2b}
	descb, err := NewPipelineBuilder("testb", "testb pipeline", inputs, outputsb, pipelineStepsb).Compile()
	assert.NoError(t, err)
	pipb := &FullySpecifiedPipeline{
		Pipeline: descb,
		steps:    pipelineStepsb,
		name:     "testb",
		inputs:   inputs,
		outputs:  outputsb,
	}

	err = pipa.AddPipeline(pipb)
	assert.NoError(t, err)

	assert.Equal(t, 6, len(pipa.steps))
	assert.Equal(t, "steps.3.produce", pipa.steps[4].GetArguments()["inputs"].RefString())
	assert.Equal(t, "steps.4.produce", pipa.steps[5].GetArguments()["inputs"].RefString())

	assert.Equal(t, 2, len(pipa.outputs))
	assert.Equal(t, "steps.2.produce", pipa.outputs[0].RefString())
	assert.Equal(t, "steps.5.produce", pipa.outputs[1].RefString())
}
