package description

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unchartedsoftware/distil-compute/pipeline"
)

func TestBasicPipelineCompile(t *testing.T) {

	step0 := NewPipelineNode(createTestStep(0))
	step1 := NewPipelineNode(createTestStep(1))
	step2 := NewPipelineNode(createTestStep(2))

	step0.Add(step1)
	step1.Add(step2)

	desc, err := NewPipelineBuilder("test", "test pipeline", step0).Compile()
	assert.NotNil(t, desc)
	assert.NoError(t, err)

	steps := desc.GetSteps()
	assert.Equal(t, 3, len(steps))

	// validate step inputs
	assert.Equal(t, "inputs.0", steps[0].GetPrimitive().GetArguments()[stepInputsKey].GetContainer().GetData())
	testStep(t, 0, step0.step, steps)

	assert.Equal(t, "steps.0.produce", steps[1].GetPrimitive().GetArguments()[stepInputsKey].GetContainer().GetData())
	testStep(t, 1, step1.step, steps)

	assert.Equal(t, "steps.1.produce", steps[2].GetPrimitive().GetArguments()[stepInputsKey].GetContainer().GetData())
	testStep(t, 2, step2.step, steps)

	// validate outputs
	assert.Equal(t, 1, len(desc.GetOutputs()))
	assert.Equal(t, "steps.2.produce", desc.GetOutputs()[0].GetData())
}

func TestMultiInputPipelineCompile(t *testing.T) {

	step0 := NewPipelineNode(createTestStep(0))
	step1 := NewPipelineNode(createTestStep(1))
	step2 := NewPipelineNode(createTestStep(2))

	step0.Add(step2)

	desc, err := NewPipelineBuilder("test", "test pipeline", step0, step1).Compile()

	assert.NotNil(t, desc)
	assert.NoError(t, err)

	steps := desc.GetSteps()
	assert.Equal(t, 3, len(steps))

	// validate step inputs
	assert.Equal(t, "inputs.0", steps[0].GetPrimitive().GetArguments()[stepInputsKey].GetContainer().GetData())
	testStep(t, 0, step0.step, steps)

	assert.Equal(t, "inputs.1", steps[1].GetPrimitive().GetArguments()[stepInputsKey].GetContainer().GetData())
	testStep(t, 1, step1.step, steps)

	assert.Equal(t, "steps.0.produce", steps[2].GetPrimitive().GetArguments()[stepInputsKey].GetContainer().GetData())
	testStep(t, 2, step2.step, steps)

	// validate outputs
	assert.Equal(t, 2, len(desc.GetOutputs()))
	assert.Equal(t, "steps.1.produce", desc.GetOutputs()[0].GetData())
	assert.Equal(t, "steps.2.produce", desc.GetOutputs()[1].GetData())
}

func TestMultiInputNamedArgPipelineCompile(t *testing.T) {

	step0 := NewPipelineNode(createTestStepWithAll(0, []string{"produce"}, []string{"arg.0", "arg.1"}))
	step1 := NewPipelineNode(createTestStep(1))
	step2 := NewPipelineNode(createTestStep(2))

	step0.Add(step2)

	desc, err := NewPipelineBuilder("test", "test pipeline", step0, step1).Compile()

	assert.NotNil(t, desc)
	assert.NoError(t, err)

	steps := desc.GetSteps()
	assert.Equal(t, 3, len(steps))

	// validate step inputs
	assert.Equal(t, "inputs.0", steps[0].GetPrimitive().GetArguments()["arg.0"].GetContainer().GetData())
	assert.Equal(t, "inputs.1", steps[0].GetPrimitive().GetArguments()["arg.1"].GetContainer().GetData())
	testStep(t, 0, step0.step, steps)

	assert.Equal(t, "inputs.2", steps[1].GetPrimitive().GetArguments()[stepInputsKey].GetContainer().GetData())
	testStep(t, 1, step1.step, steps)

	assert.Equal(t, "steps.0.produce", steps[2].GetPrimitive().GetArguments()[stepInputsKey].GetContainer().GetData())
	testStep(t, 2, step2.step, steps)

	// validate outputs
	assert.Equal(t, 2, len(desc.GetOutputs()))
	assert.Equal(t, "steps.1.produce", desc.GetOutputs()[0].GetData())
	assert.Equal(t, "steps.2.produce", desc.GetOutputs()[1].GetData())
}

func TestBranchPipelineCompile(t *testing.T) {

	step0 := NewPipelineNode(createTestStepWithOutputs(0, []string{"produce.0", "produce.1"}))
	step1 := NewPipelineNode(createTestStep(1))
	step2 := NewPipelineNode(createTestStep(2))
	step0.Add(step1)
	step0.Add(step2)

	desc, err := NewPipelineBuilder("test", "test pipeline", step0).Compile()
	assert.NotNil(t, desc)
	assert.NoError(t, err)

	steps := desc.GetSteps()
	assert.Equal(t, 3, len(steps))

	// validate step inputs
	assert.Equal(t, "inputs.0", steps[0].GetPrimitive().GetArguments()[stepInputsKey].GetContainer().GetData())
	testStep(t, 0, step0.step, steps)

	assert.Equal(t, "steps.0.produce.0", steps[1].GetPrimitive().GetArguments()[stepInputsKey].GetContainer().GetData())
	testStep(t, 1, step1.step, steps)

	assert.Equal(t, "steps.0.produce.1", steps[2].GetPrimitive().GetArguments()[stepInputsKey].GetContainer().GetData())
	testStep(t, 2, step2.step, steps)

	// validate outputs
	assert.Equal(t, 2, len(desc.GetOutputs()))
	assert.Equal(t, "steps.1.produce", desc.GetOutputs()[0].GetData())
	assert.Equal(t, "steps.2.produce", desc.GetOutputs()[1].GetData())
}

func TestBasicNamedArgCompile(t *testing.T) {

	step0 := NewPipelineNode(createTestStepWithOutputs(0, []string{"produce.0", "produce.1"}))
	step1 := NewPipelineNode(createTestStepWithAll(1, []string{"produce"}, []string{"arg.0", "arg.1"}))
	step2 := NewPipelineNode(createTestStep(2))

	step0.Add(step1)
	step1.Add(step2)

	desc, err := NewPipelineBuilder("test", "test pipeline", step0).Compile()
	assert.NotNil(t, desc)
	assert.NoError(t, err)

	steps := desc.GetSteps()
	assert.Equal(t, 3, len(steps))

	// validate step inputs
	assert.Equal(t, "inputs.0", steps[0].GetPrimitive().GetArguments()[stepInputsKey].GetContainer().GetData())
	testStep(t, 0, step0.step, steps)

	assert.Equal(t, "steps.0.produce.0", steps[1].GetPrimitive().GetArguments()["arg.0"].GetContainer().GetData())
	assert.Equal(t, "steps.0.produce.1", steps[1].GetPrimitive().GetArguments()["arg.1"].GetContainer().GetData())
	testStep(t, 1, step1.step, steps)

	// validate outputs
	assert.Equal(t, 1, len(desc.GetOutputs()))
	assert.Equal(t, "steps.2.produce", desc.GetOutputs()[0].GetData())
}

func TestDiamondPipeline(t *testing.T) {

	step0 := NewPipelineNode(createTestStepWithOutputs(0, []string{"produce.0", "produce.1"}))
	step1 := NewPipelineNode(createTestStep(1))
	step2 := NewPipelineNode(createTestStepWithOutputs(2, []string{"produce.0", "produce.1"}))
	step3 := NewPipelineNode(createTestStepWithAll(3, []string{"produce"}, []string{"arg.0", "arg.1", "arg.2"}))

	step0.Add(step1)
	step0.Add(step2)
	step1.Add(step3)
	step2.Add(step3)

	desc, err := NewPipelineBuilder("test", "test pipeline", step0).Compile()
	assert.NotNil(t, desc)
	assert.NoError(t, err)

	steps := desc.GetSteps()
	assert.Equal(t, 4, len(steps))

	// validate step inputs
	assert.Equal(t, "inputs.0", steps[0].GetPrimitive().GetArguments()[stepInputsKey].GetContainer().GetData())
	testStep(t, 0, step0.step, steps)

	assert.Equal(t, "steps.0.produce.0", steps[1].GetPrimitive().GetArguments()[stepInputsKey].GetContainer().GetData())
	testStep(t, 1, step1.step, steps)

	assert.Equal(t, "steps.0.produce.1", steps[2].GetPrimitive().GetArguments()[stepInputsKey].GetContainer().GetData())
	testStep(t, 2, step2.step, steps)

	assert.Equal(t, "steps.1.produce", steps[3].GetPrimitive().GetArguments()["arg.0"].GetContainer().GetData())
	assert.Equal(t, "steps.2.produce.0", steps[3].GetPrimitive().GetArguments()["arg.1"].GetContainer().GetData())
	assert.Equal(t, "steps.2.produce.1", steps[3].GetPrimitive().GetArguments()["arg.2"].GetContainer().GetData())
	testStep(t, 3, step3.step, steps)

	// validate outputs
	assert.Equal(t, 1, len(desc.GetOutputs()))
	assert.Equal(t, "steps.3.produce", desc.GetOutputs()[0].GetData())
}

func TestMergePipeline(t *testing.T) {

	step0 := NewPipelineNode(createTestStep(0))
	step1 := NewPipelineNode(createTestStep(1))
	step2 := NewPipelineNode(createTestStepWithAll(2, []string{"produce"}, []string{"arg.0", "arg.1"}))

	step0.Add(step2)
	step1.Add(step2)

	desc, err := NewPipelineBuilder("test", "test pipeline", step0, step1).Compile()
	assert.NotNil(t, desc)
	assert.NoError(t, err)

	steps := desc.GetSteps()
	assert.Equal(t, 3, len(steps))

	// validate step inputs
	assert.Equal(t, "inputs.0", steps[0].GetPrimitive().GetArguments()[stepInputsKey].GetContainer().GetData())
	testStep(t, 0, step0.step, steps)

	assert.Equal(t, "inputs.1", steps[1].GetPrimitive().GetArguments()[stepInputsKey].GetContainer().GetData())
	testStep(t, 1, step1.step, steps)

	assert.Equal(t, "steps.0.produce", steps[2].GetPrimitive().GetArguments()["arg.0"].GetContainer().GetData())
	assert.Equal(t, "steps.1.produce", steps[2].GetPrimitive().GetArguments()["arg.1"].GetContainer().GetData())
	testStep(t, 2, step2.step, steps)

	// validate outputs
	assert.Equal(t, 1, len(desc.GetOutputs()))
	assert.Equal(t, "steps.2.produce", desc.GetOutputs()[0].GetData())
}

func TestBasicInferenceCompile(t *testing.T) {

	step0 := NewPipelineNode(createTestStep(0))
	step1 := NewPipelineNode(createTestStep(1))
	step2 := NewPipelineNode(NewInferenceStepData())

	step0.Add(step1)
	step1.Add(step2)

	desc, err := NewPipelineBuilder("test", "test pipeline", step0).Compile()
	assert.NotNil(t, desc)
	assert.NoError(t, err)

	steps := desc.GetSteps()
	assert.Equal(t, 3, len(steps))

	// validate step inputs
	assert.Equal(t, "inputs.0", steps[0].GetPrimitive().GetArguments()[stepInputsKey].GetContainer().GetData())
	testStep(t, 0, step0.step, steps)

	// validate outputs
	assert.Equal(t, 1, len(desc.GetOutputs()))
	assert.Equal(t, "steps.2.produce", desc.GetOutputs()[0].GetData())

	assert.Nil(t, steps[2].GetPrimitive().GetHyperparams())
	assert.Nil(t, steps[2].GetPrimitive().GetPrimitive())
}

func TestRecompileFailure(t *testing.T) {
	step0 := NewPipelineNode(createTestStep(0))
	step1 := NewPipelineNode(createTestStep(1))
	step2 := NewPipelineNode(createTestStep(2))

	step0.Add(step1)
	step1.Add(step2)

	builder := NewPipelineBuilder("test", "test pipeline", step0)
	desc, err := builder.Compile()
	desc, err = builder.Compile()

	assert.Nil(t, desc)
	assert.Error(t, err)
}

func TestMultiInferenceFailure(t *testing.T) {
	step0 := NewPipelineNode(createTestStepWithOutputs(0, []string{"produce.0", "produce.1"}))
	step1 := NewPipelineNode(NewInferenceStepData())
	step2 := NewPipelineNode(NewInferenceStepData())

	step0.Add(step1)
	step1.Add(step2)

	desc, err := NewPipelineBuilder("test", "test pipeline", step0).Compile()

	assert.Error(t, err)
	assert.Nil(t, desc)
}

func TestInferenceChildFailure(t *testing.T) {
	step0 := NewPipelineNode(createTestStep(0))
	step1 := NewPipelineNode(NewInferenceStepData())
	step2 := NewPipelineNode(createTestStep(1))

	step0.Add(step1)
	step1.Add(step2)
	desc, err := NewPipelineBuilder("test", "test pipeline", step0).Compile()

	assert.Error(t, err)
	assert.Nil(t, desc)
}

func TestCycleFailure(t *testing.T) {
	step0 := NewPipelineNode(createTestStep(0))
	step1 := NewPipelineNode(createTestStep(1))
	step2 := NewPipelineNode(createTestStep(2))

	step0.Add(step1)
	step1.Add(step2)
	step2.Add(step0)
	desc, err := NewPipelineBuilder("test", "test pipeline", step0).Compile()

	assert.Error(t, err)
	assert.Nil(t, desc)
}

func createLabels(counter int64) []string {
	return []string{fmt.Sprintf("alpha-%d", counter), fmt.Sprintf("bravo-%d", counter)}
}

func createTestStepWithOutputs(step int64, outputMethods []string) *StepData {
	return createTestStepWithAll(step, outputMethods, nil)
}

func createTestStep(step int64) *StepData {
	return createTestStepWithOutputs(step, []string{"produce"})
}

func createTestStepWithAll(step int64, outputMethods []string, arguments []string) *StepData {
	labels := createLabels(step)
	return NewStepDataWithAll(
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

	assert.EqualValues(t, step.GetPrimitive(), steps[index].GetPrimitive().GetPrimitive())
}
