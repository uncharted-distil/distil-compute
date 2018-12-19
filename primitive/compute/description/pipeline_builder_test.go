package description

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicPipelineCompile(t *testing.T) {
	builder := NewPipelineBuilder("test", "test pipeline")

	step0 := NewPipelineNode(createTestStep(0))
	step1 := NewPipelineNode(createTestStep(1))
	step2 := NewPipelineNode(createTestStep(2))

	step0.Add(step1)
	step1.Add(step2)

	builder.AddSource(step0)
	desc, err := builder.Compile()
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

func TestMultioutputPipelineCompile(t *testing.T) {
	builder := NewPipelineBuilder("test", "test pipeline")

	step0 := NewPipelineNode(createTestStep(0))
	step1 := NewPipelineNode(createTestStep(1))
	step2 := NewPipelineNode(createTestStep(2))

	step0.Add(step1)
	step1.Add(step2)

	builder.AddSource(step0)
	desc, err := builder.Compile()
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

func TestDiamondPipeline(t *testing.T) {
	builder := NewPipelineBuilder("test", "test pipeline")

	step0 := NewPipelineNode(createTestStep(0))
	step1 := NewPipelineNode(createTestStep(1))
	step2 := NewPipelineNode(createTestStep(2))
	step3 := NewPipelineNode(createTestStep(3))

	step0.Add(step1)
	step0.Add(step2)
	step1.Add(step3)
	step2.Add(step3)

	builder.AddSource(step0)
	pipelineDesc, err := builder.Compile()
	assert.NotNil(t, pipelineDesc)
	assert.NoError(t, err)
}
