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

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uncharted-distil/distil-compute/pipeline"
)

func TestBasicPipelineCompile(t *testing.T) {

	inputs := []string{"inputs"}
	outputs := []DataRef{&StepDataRef{2, "produce"}}

	step0 := createTestStep(0, []string{"produce"}, map[string]DataRef{"inputs": &PipelineDataRef{0}})
	step1 := createTestStep(1, []string{"produce"}, map[string]DataRef{"inputs": &StepDataRef{0, "produce"}})
	step2 := createTestStep(2, []string{"produce"}, map[string]DataRef{"inputs": &StepDataRef{1, "produce"}})
	pipelineSteps := []Step{step0, step1, step2}

	desc, err := NewPipelineBuilder("test", "test pipeline", inputs, outputs, pipelineSteps).Compile()

	assert.NotNil(t, desc)
	assert.NoError(t, err)

	steps := desc.GetSteps()
	assert.Equal(t, 3, len(steps))

	// validate step inputs
	assert.Equal(t, "inputs.0", steps[0].GetPrimitive().GetArguments()[stepInputsKey].GetContainer().GetData())
	testStep(t, 0, step0, steps)

	assert.Equal(t, "steps.0.produce", steps[1].GetPrimitive().GetArguments()[stepInputsKey].GetContainer().GetData())
	testStep(t, 1, step1, steps)

	assert.Equal(t, "steps.1.produce", steps[2].GetPrimitive().GetArguments()[stepInputsKey].GetContainer().GetData())
	testStep(t, 2, step2, steps)

	// validate outputs
	assert.Equal(t, 1, len(desc.GetOutputs()))
	assert.Equal(t, "steps.2.produce", desc.GetOutputs()[0].GetData())
}

func TestMultiInputPipelineCompile(t *testing.T) {

	inputs := []string{"inputs"}
	outputs := []DataRef{
		&StepDataRef{1, "produce"},
		&StepDataRef{2, "produce"},
	}

	step0 := createTestStep(0, []string{"produce"}, map[string]DataRef{"inputs": &PipelineDataRef{0}})
	step1 := createTestStep(1, []string{"produce"}, map[string]DataRef{"inputs": &PipelineDataRef{1}})
	step2 := createTestStep(2, []string{"produce"}, map[string]DataRef{"inputs": &StepDataRef{0, "produce"}})
	pipelineSteps := []Step{step0, step1, step2}

	desc, err := NewPipelineBuilder("test", "test pipeline", inputs, outputs, pipelineSteps).Compile()

	assert.NotNil(t, desc)
	assert.NoError(t, err)

	steps := desc.GetSteps()
	assert.Equal(t, 3, len(steps))

	// validate step inputs
	assert.Equal(t, "inputs.0", steps[0].GetPrimitive().GetArguments()[stepInputsKey].GetContainer().GetData())
	testStep(t, 0, step0, steps)

	assert.Equal(t, "inputs.1", steps[1].GetPrimitive().GetArguments()[stepInputsKey].GetContainer().GetData())
	testStep(t, 1, step1, steps)

	assert.Equal(t, "steps.0.produce", steps[2].GetPrimitive().GetArguments()[stepInputsKey].GetContainer().GetData())
	testStep(t, 2, step2, steps)

	// validate outputs
	assert.Equal(t, 2, len(desc.GetOutputs()))
	assert.Equal(t, "steps.1.produce", desc.GetOutputs()[0].GetData())
	assert.Equal(t, "steps.2.produce", desc.GetOutputs()[1].GetData())
}

func TestMultiInputNamedArgPipelineCompile(t *testing.T) {

	inputs := []string{"inputs"}
	outputs := []DataRef{
		&StepDataRef{1, "produce"},
		&StepDataRef{2, "produce"},
	}

	step0 := createTestStep(0, []string{"produce"}, map[string]DataRef{"arg.0": &PipelineDataRef{0}, "arg.1": &PipelineDataRef{1}})
	step1 := createTestStep(1, []string{"produce"}, map[string]DataRef{"inputs": &PipelineDataRef{2}})
	step2 := createTestStep(2, []string{"produce"}, map[string]DataRef{"inputs": &StepDataRef{0, "produce"}})
	pipelineSteps := []Step{step0, step1, step2}

	desc, err := NewPipelineBuilder("test", "test pipeline", inputs, outputs, pipelineSteps).Compile()

	assert.NotNil(t, desc)
	assert.NoError(t, err)

	steps := desc.GetSteps()
	assert.Equal(t, 3, len(steps))

	// validate step inputs
	assert.Equal(t, "inputs.0", steps[0].GetPrimitive().GetArguments()["arg.0"].GetContainer().GetData())
	assert.Equal(t, "inputs.1", steps[0].GetPrimitive().GetArguments()["arg.1"].GetContainer().GetData())
	testStep(t, 0, step0, steps)

	assert.Equal(t, "inputs.2", steps[1].GetPrimitive().GetArguments()[stepInputsKey].GetContainer().GetData())
	testStep(t, 1, step1, steps)

	assert.Equal(t, "steps.0.produce", steps[2].GetPrimitive().GetArguments()[stepInputsKey].GetContainer().GetData())
	testStep(t, 2, step2, steps)

	// validate outputs
	assert.Equal(t, 2, len(desc.GetOutputs()))
	assert.Equal(t, "steps.1.produce", desc.GetOutputs()[0].GetData())
	assert.Equal(t, "steps.2.produce", desc.GetOutputs()[1].GetData())
}

func TestBasicInferenceCompile(t *testing.T) {

	inputs := []string{"inputs"}
	outputs := []DataRef{&StepDataRef{2, "produce"}}

	step0 := createTestStep(0, []string{"produce"}, map[string]DataRef{"inputs": &PipelineDataRef{0}})
	step1 := createTestStep(1, []string{"produce"}, map[string]DataRef{"inputs": &StepDataRef{0, "produce"}})
	step2 := NewInferenceStepData(map[string]DataRef{"inputs": &StepDataRef{1, "produce"}})
	pipelineSteps := []Step{step0, step1, step2}

	desc, err := NewPipelineBuilder("test", "test pipeline", inputs, outputs, pipelineSteps).Compile()

	assert.NotNil(t, desc)
	assert.NoError(t, err)

	steps := desc.GetSteps()
	assert.Equal(t, 3, len(steps))

	// validate step inputs
	assert.Equal(t, "inputs.0", steps[0].GetPrimitive().GetArguments()[stepInputsKey].GetContainer().GetData())
	assert.Equal(t, "steps.0.produce", steps[1].GetPrimitive().GetArguments()[stepInputsKey].GetContainer().GetData())
	assert.Equal(t, "steps.1.produce", steps[2].GetPlaceholder().GetInputs()[0].GetData())
	testStep(t, 0, step0, steps)

	// validate outputs
	assert.Equal(t, 1, len(desc.GetOutputs()))
	assert.Equal(t, "steps.2.produce", desc.GetOutputs()[0].GetData())

	assert.Nil(t, steps[2].GetPrimitive().GetHyperparams())
	assert.Nil(t, steps[2].GetPrimitive().GetPrimitive())
}

func TestRecompileFailure(t *testing.T) {

	inputs := []string{"inputs"}
	outputs := []DataRef{&StepDataRef{2, "produce"}}

	step0 := createTestStep(0, []string{"produce"}, map[string]DataRef{"inputs": &PipelineDataRef{0}})
	step1 := createTestStep(1, []string{"produce"}, map[string]DataRef{"inputs": &StepDataRef{0, "produce"}})
	step2 := createTestStep(2, []string{"produce"}, map[string]DataRef{"inputs": &StepDataRef{1, "produce"}})
	pipelineSteps := []Step{step0, step1, step2}

	builder := NewPipelineBuilder("test", "test pipeline", inputs, outputs, pipelineSteps)

	desc, err := builder.Compile()
	desc, err = builder.Compile()

	assert.Nil(t, desc)
	assert.Error(t, err)
}

// TODO: not yet supported - add additional validation
//
// func TestMultiInferenceFailure(t *testing.T)
// func TestInferenceChildFailure(t *testing.T)
// func TestCycleFailure(t *testing.T)

func createLabels(counter int64) []string {
	return []string{fmt.Sprintf("alpha-%d", counter), fmt.Sprintf("bravo-%d", counter)}
}

func createTestStep(step int64, outputMethods []string, arguments map[string]DataRef) *StepData {
	labels := createLabels(step)
	return NewStepData(
		&pipeline.Primitive{
			Id:         fmt.Sprintf("0000-primtive-%d", step),
			Version:    "1.0.0",
			Name:       fmt.Sprintf("primitive-%d", step),
			PythonPath: fmt.Sprintf("d3m.primitives.distil.primitive.%d", step),
		},
		outputMethods,
		map[string]interface{}{
			"testString":         fmt.Sprintf("hyperparam-%d", step),
			"testBool":           step%2 == 0,
			"testInt":            step,
			"testFloat":          float64(step) + 0.5,
			"testStringArray":    labels,
			"testBoolArray":      []bool{step%2 == 0, step%2 != 0},
			"testIntArray":       []int64{step, step + 1},
			"testFloatArray":     []float64{float64(step) + 0.5, float64(step) + 1.5},
			"testIntMap":         map[string]int64{labels[0]: int64(step), labels[1]: int64(step + 1)},
			"testFloatMap":       map[string]float64{labels[0]: float64(step) + 0.5, labels[1]: float64(step) + 1.5},
			"testNestedIntArray": [][]int64{{step, step + 1}, {step + 2, step + 3}},
			"testNestedIntMap":   map[string][]int64{labels[0]: {step, step + 1}, labels[1]: {step + 2, step + 3}},
			"testPrimitiveRef":   &PrimitiveReference{int(step)},
		},
		arguments,
	)
}

func ConvertToStringArray(list *pipeline.ValueList) []string {
	arr := []string{}
	for _, v := range list.Items {
		arr = append(arr, v.GetString_())
	}
	return arr
}

func ConvertToBoolArray(list *pipeline.ValueList) []bool {
	arr := []bool{}
	for _, v := range list.Items {
		arr = append(arr, v.GetBool())
	}
	return arr
}

func ConvertToIntArray(list *pipeline.ValueList) []int64 {
	arr := []int64{}
	for _, v := range list.Items {
		arr = append(arr, v.GetInt64())
	}
	return arr
}

func ConvertToFloatArray(list *pipeline.ValueList) []float64 {
	arr := []float64{}
	for _, v := range list.Items {
		arr = append(arr, v.GetDouble())
	}
	return arr
}

func ConvertToIntMap(dict *pipeline.ValueDict) map[string]int64 {
	mp := map[string]int64{}
	for k, v := range dict.Items {
		mp[k] = v.GetInt64()
	}
	return mp
}

func ConvertToFloatMap(dict *pipeline.ValueDict) map[string]float64 {
	mp := map[string]float64{}
	for k, v := range dict.Items {
		mp[k] = v.GetDouble()
	}
	return mp
}

func ConvertToNestedIntArray(list *pipeline.ValueList) [][]int64 {
	arr := [][]int64{}
	for _, v := range list.Items {
		inner := []int64{}
		for _, w := range v.GetList().Items {
			inner = append(inner, w.GetInt64())
		}
		arr = append(arr, inner)
	}
	return arr
}

func ConvertToNestedIntMap(dict *pipeline.ValueDict) map[string][]int64 {
	mp := map[string][]int64{}
	for k, v := range dict.Items {
		inner := []int64{}
		for _, w := range v.GetList().Items {
			inner = append(inner, w.GetInt64())
		}
		mp[k] = inner
	}
	return mp
}

func testStep(t *testing.T, index int64, step Step, steps []*pipeline.PipelineDescriptionStep) {
	labels := createLabels(index)

	assert.Equal(t, fmt.Sprintf("hyperparam-%d", index),
		steps[index].GetPrimitive().GetHyperparams()["testString"].GetValue().GetData().GetRaw().GetString_())

	assert.Equal(t, int64(index), steps[index].GetPrimitive().GetHyperparams()["testInt"].GetValue().GetData().GetRaw().GetInt64())

	assert.Equal(t, index%2 == 0, steps[index].GetPrimitive().GetHyperparams()["testBool"].GetValue().GetData().GetRaw().GetBool())

	assert.Equal(t, float64(index)+0.5, steps[index].GetPrimitive().GetHyperparams()["testFloat"].GetValue().GetData().GetRaw().GetDouble())

	assert.Equal(t, labels,
		ConvertToStringArray(steps[index].GetPrimitive().GetHyperparams()["testStringArray"].GetValue().GetData().GetRaw().GetList()))

	assert.Equal(t, []int64{int64(index), int64(index) + 1},
		ConvertToIntArray(steps[index].GetPrimitive().GetHyperparams()["testIntArray"].GetValue().GetData().GetRaw().GetList()))

	assert.Equal(t, []float64{float64(index) + 0.5, float64(index) + 1.5},
		ConvertToFloatArray(steps[index].GetPrimitive().GetHyperparams()["testFloatArray"].GetValue().GetData().GetRaw().GetList()))

	assert.Equal(t, []bool{index%2 == 0, index%2 != 0},
		ConvertToBoolArray(steps[index].GetPrimitive().GetHyperparams()["testBoolArray"].GetValue().GetData().GetRaw().GetList()))

	assert.Equal(t, map[string]int64{labels[0]: int64(index), labels[1]: int64(index + 1)},
		ConvertToIntMap(steps[index].GetPrimitive().GetHyperparams()["testIntMap"].GetValue().GetData().GetRaw().GetDict()))

	assert.Equal(t, map[string]float64{labels[0]: float64(index) + 0.5, labels[1]: float64(index) + 1.5},
		ConvertToFloatMap(steps[index].GetPrimitive().GetHyperparams()["testFloatMap"].GetValue().GetData().GetRaw().GetDict()))

	assert.Equal(t, [][]int64{{index, index + 1}, {index + 2, index + 3}},
		ConvertToNestedIntArray(steps[index].GetPrimitive().GetHyperparams()["testNestedIntArray"].GetValue().GetData().GetRaw().GetList()))

	assert.Equal(t, map[string][]int64{labels[0]: {index, index + 1}, labels[1]: {index + 2, index + 3}},
		ConvertToNestedIntMap(steps[index].GetPrimitive().GetHyperparams()["testNestedIntMap"].GetValue().GetData().GetRaw().GetDict()))

	assert.Equal(t, int32(index), steps[index].GetPrimitive().GetHyperparams()["testPrimitiveRef"].GetPrimitive().GetData())

	assert.EqualValues(t, step.GetPrimitive(), steps[index].GetPrimitive().GetPrimitive())
}
